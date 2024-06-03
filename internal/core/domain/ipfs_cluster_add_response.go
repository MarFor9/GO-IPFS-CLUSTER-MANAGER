package domain

type IPFSClusterAddResponse struct {
	Name        string   `json:"name"`
	Cid         string   `json:"cid"`
	Size        int      `json:"size"`
	Allocations []string `json:"allocations"`
}
