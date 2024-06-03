package domain

type Status string

const (
	Alive Status = "Alive"
	Down  Status = "Down"
)

type IPFSHealthCheckResponse struct {
	Status []ClusterNodePairStatus `json:"status"`
}

type ClusterNodePairStatus struct {
	Cluster IPFSStatus `json:"cluster"`
	Node    IPFSStatus `json:"node"`
}

type IPFSStatus struct {
	Url          string `json:"url"`
	Status       Status `json:"status"`
	ResponseTime string `json:"responseTime"`
}
