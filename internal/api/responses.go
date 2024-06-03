package api

import "IPFS-CLUSTER-MANAGER/internal/core/domain"

func toApiIpfsHealthCheckResponse(response domain.IPFSHealthCheckResponse) IPFSHealthCheckResponse {
	var apiResponse IPFSHealthCheckResponse

	for _, pair := range response.Status {
		clusterStatus := IPFSStatus{
			Url:          pair.Cluster.Url,
			Status:       Status(pair.Cluster.Status),
			ResponseTime: pair.Cluster.ResponseTime,
		}
		nodeStatus := IPFSStatus{
			Url:          pair.Node.Url,
			Status:       Status(pair.Node.Status),
			ResponseTime: pair.Node.ResponseTime,
		}

		apiPair := ClusterNodePairStatus{
			Cluster: clusterStatus,
			Node:    nodeStatus,
		}

		apiResponse.Status = append(apiResponse.Status, apiPair)
	}

	return apiResponse
}

func toApiPins(pins *[]domain.Pin) []Pin {
	var apiPins []Pin

	for _, p := range *pins {
		apiPin := Pin{
			Cid:         p.CID,
			Name:        p.Name,
			Allocations: p.Allocations,
			Created:     p.Created,
			Origins:     p.Origins,
			Metadata:    &p.Metadata,
			PeerMap:     make(map[string]PeerStatus),
		}

		for peerId, domainPeerStatus := range p.PeerMap {
			peerStatus := PeerStatus{
				AttemptCount:      domainPeerStatus.AttemptCount,
				Error:             domainPeerStatus.Error,
				IpfsPeerAddresses: domainPeerStatus.IPFSPeerAddresses,
				IpfsPeerId:        domainPeerStatus.IPFSPeerID,
				Peername:          domainPeerStatus.PeerName,
				PriorityPin:       domainPeerStatus.PriorityPin,
				Status:            domainPeerStatus.Status,
				Timestamp:         domainPeerStatus.Timestamp,
			}
			apiPin.PeerMap[peerId] = peerStatus
		}

		apiPins = append(apiPins, apiPin)
	}

	return apiPins
}
