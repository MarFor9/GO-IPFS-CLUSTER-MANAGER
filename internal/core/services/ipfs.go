package services

import (
	"IPFS-CLUSTER-MANAGER/internal/core/config"
	"IPFS-CLUSTER-MANAGER/internal/core/domain"
	"IPFS-CLUSTER-MANAGER/internal/log"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type RuntimeConfig struct {
	IpfsNodeUrls    []string
	IpfsClusterUrls []string
}

type Ipfs struct {
	cfg        *config.Configuration
	runtimeCfg *RuntimeConfig
	mu         sync.RWMutex // Mutex for thread-safe access to runtimeCfg
}

func NewIpfs(cfg *config.Configuration) *Ipfs {
	return &Ipfs{
		cfg: cfg,
		runtimeCfg: &RuntimeConfig{
			IpfsNodeUrls:    []string{cfg.Ipfs0NodeUrl, cfg.Ipfs1NodeUrl, cfg.Ipfs2NodeUrl},
			IpfsClusterUrls: []string{cfg.Ipfs0ClusterUrl, cfg.Ipfs1ClusterUrl, cfg.Ipfs2ClusterUrl},
		},
	}
}

func (i *Ipfs) GetPins(ctx context.Context) (*[]domain.Pin, error) {
	log.Info(ctx, "[IPFS Cluster - GetPins] <- Enter")

	for _, baseURL := range i.runtimeCfg.IpfsClusterUrls {

		req, err := http.NewRequest("GET", baseURL+"/pins", nil)
		if err != nil {
			continue
		}

		response, err := http.DefaultClient.Do(req)
		if err != nil || response.StatusCode < http.StatusOK && response.StatusCode >= http.StatusMultipleChoices {
			log.Error(ctx, fmt.Sprintf("error reading file from cluster: %s", baseURL), err)
			continue
		}

		if response.StatusCode == http.StatusNoContent {
			return &[]domain.Pin{}, errors.New("no content")
		}

		body, err := io.ReadAll(response.Body)
		err = response.Body.Close()
		if err != nil {
			continue
		}

		var allPins []domain.Pin
		parts := strings.Split(string(body), "\n")
		for _, part := range parts {
			var pin domain.Pin
			if err := json.Unmarshal([]byte(part), &pin); err != nil {
				log.Error(ctx, "error unmarshalling part", err)
				continue
			}
			allPins = append(allPins, pin)
		}

		if len(allPins) == 0 {
			return nil, errors.New("no valid pins found")
		}

		log.Info(ctx, "[IPFS Cluster - GetPins] <- Leave")
		return &allPins, nil
	}
	return nil, errors.New("unable to retrieve file from any cluster nodes")
}

func (i *Ipfs) AddClusterNodePair(ctx context.Context, nodeUrl string, clusterUrl string) error {
	// Acquire a write lock to ensure thread-safe access to the runtime configuration
	i.mu.Lock()
	defer i.mu.Unlock()

	// Append the new node and cluster URL if it's not already in the list
	if !contains(i.runtimeCfg.IpfsNodeUrls, nodeUrl) && !contains(i.runtimeCfg.IpfsClusterUrls, clusterUrl) {
		i.runtimeCfg.IpfsNodeUrls = append(i.runtimeCfg.IpfsNodeUrls, nodeUrl)
		i.runtimeCfg.IpfsClusterUrls = append(i.runtimeCfg.IpfsClusterUrls, clusterUrl)
	}

	return nil
}

// Helper function to check if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (i *Ipfs) GetFile(ctx context.Context, cid string) ([]byte, error) {
	log.Info(ctx, "[IPFS Node - GetFile] <- Enter")

	for _, baseURL := range i.runtimeCfg.IpfsNodeUrls {

		req, err := http.NewRequest("POST", baseURL+"/api/v0/cat", nil)
		if err != nil {
			continue
		}

		q := req.URL.Query()
		q.Add("arg", cid)
		req.URL.RawQuery = q.Encode()

		response, err := http.DefaultClient.Do(req)
		if err != nil || response.StatusCode != http.StatusOK {
			log.Error(ctx, fmt.Sprintf("error reading file from cluster: %s", baseURL), err)
			continue
		}

		bodyBytes, err := io.ReadAll(response.Body)
		err = response.Body.Close()
		if err != nil {
			continue
		}

		log.Info(ctx, fmt.Sprintf("File was read from IPFS node: %s", baseURL))
		log.Info(ctx, "[IPFS Node - GetFile] <- Leave")
		return bodyBytes, nil
	}
	return nil, errors.New("unable to retrieve file from any cluster nodes")
}
func (i *Ipfs) GetStatus(ctx context.Context) domain.IPFSHealthCheckResponse {
	log.Info(ctx, "[IPFS Cluster/Node - GetStatus] <- Enter")
	response := domain.IPFSHealthCheckResponse{}

	// Assuming the lengths of cluster and node URLs are the same
	for idx := range i.runtimeCfg.IpfsClusterUrls {
		clusterURL := i.runtimeCfg.IpfsClusterUrls[idx]
		nodeURL := i.runtimeCfg.IpfsNodeUrls[idx]

		// Check cluster
		clusterAlive, clusterResponseTime := healthCheck(clusterURL+"/health", "GET")
		clusterStatus := domain.Down
		if clusterAlive {
			clusterStatus = domain.Alive
		}
		cluster := domain.IPFSStatus{
			Url:          clusterURL + "/health",
			Status:       clusterStatus,
			ResponseTime: clusterResponseTime.String(),
		}

		// Check node
		nodeAlive, nodeResponseTime := healthCheck(nodeURL+"/api/v0/version", "POST")
		nodeStatus := domain.Down
		if nodeAlive {
			nodeStatus = domain.Alive
		}
		node := domain.IPFSStatus{
			Url:          nodeURL + "/api/v0/version",
			Status:       nodeStatus,
			ResponseTime: nodeResponseTime.String(),
		}

		// Add pair to response
		response.Status = append(response.Status, domain.ClusterNodePairStatus{Cluster: cluster, Node: node})
	}

	log.Info(ctx, "[IPFS Cluster/Node - GetStatus] <- Leave")
	return response
}

func healthCheck(url string, method string) (bool, time.Duration) {
	startTime := time.Now()

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return false, time.Since(startTime)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(context.Background(), fmt.Sprintf("error checking health of %s", url), err)
		return false, time.Since(startTime)
	}
	defer response.Body.Close()

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return true, time.Since(startTime)
	}

	return false, time.Since(startTime)
}

func (i *Ipfs) AddFile(ctx context.Context, requestBody *multipart.Reader) (string, error) {
	log.Info(ctx, "[IPFS Cluster - AddFile] <- Enter")

	file, _, formDataContentType, err := readFile(ctx, requestBody)
	if err != nil {
		return "", err
	}

	var ipfsResponse domain.IPFSClusterAddResponse
	for _, baseURL := range i.runtimeCfg.IpfsClusterUrls {
		parsedURL, err := url.Parse(baseURL + "/add")
		if err != nil {
			continue
		}

		finalURL := parsedURL.String()

		req, err := http.NewRequest("POST", finalURL, &file)
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", formDataContentType)

		response, err := http.DefaultClient.Do(req)
		if err != nil || response.StatusCode != http.StatusOK {
			log.Error(ctx, fmt.Sprintf("error adding file to cluster: %s", baseURL), err)
			continue
		}
		defer response.Body.Close()

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Error(ctx, fmt.Sprintf("error reading response body from cluster: %s", baseURL), err)
			continue
		}

		err = json.Unmarshal(responseData, &ipfsResponse)
		if err != nil {
			log.Error(ctx, fmt.Sprintf("error unmarshalling response body from cluster: %s", baseURL), err)
			continue
		}

		// success break the loop
		break
	}
	if ipfsResponse.Cid == "" {
		return "", errors.New("unable to add file to any cluster nodes")
	}
	log.Info(ctx, "[IPFS Cluster - AddFile] <- Leave")
	return ipfsResponse.Cid, nil
}

func readFile(ctx context.Context, requestBody *multipart.Reader) (bytes.Buffer, string, string, error) {
	var file []byte
	var filename string
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	for {
		part, err := requestBody.NextPart()
		if err != nil {
			break
		}

		readAll, err := io.ReadAll(part)
		if err != nil {
			log.Error(ctx, "Error reading multipart form", "err", err)
			return bytes.Buffer{}, "", "", err
		}

		switch part.FormName() {
		case "file":
			if len(readAll) == 0 {
				return bytes.Buffer{}, "", "", errors.New("file is empty")
			}
			file = readAll
			filename = part.FileName()
		}
	}

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return bytes.Buffer{}, "", "", err
	}

	if _, err := part.Write(file); err != nil {
		return bytes.Buffer{}, "", "", err
	}

	err = writer.Close()
	if err != nil {
		return bytes.Buffer{}, "", "", err
	}

	return b, filename, writer.FormDataContentType(), nil
}
