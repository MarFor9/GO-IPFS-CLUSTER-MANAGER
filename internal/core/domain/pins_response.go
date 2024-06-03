package domain

import "time"

type Pin struct {
	CID         string                 `json:"cid"`
	Name        string                 `json:"name"`
	Allocations []string               `json:"allocations"`
	Origins     []string               `json:"origins"`
	Created     time.Time              `json:"created"`
	Metadata    map[string]interface{} `json:"metadata"`
	PeerMap     map[string]PeerStatus  `json:"peer_map"`
}

type PeerStatus struct {
	PeerName          string    `json:"peername"`
	IPFSPeerID        string    `json:"ipfs_peer_id"`
	IPFSPeerAddresses []string  `json:"ipfs_peer_addresses"`
	Status            string    `json:"status"`
	Timestamp         time.Time `json:"timestamp"`
	Error             *string   `json:"error"`
	AttemptCount      int       `json:"attempt_count"`
	PriorityPin       bool      `json:"priority_pin"`
}
