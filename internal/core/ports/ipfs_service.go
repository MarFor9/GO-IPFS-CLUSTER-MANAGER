package ports

import (
	"IPFS-CLUSTER-MANAGER/internal/core/domain"
	"context"
	"mime/multipart"
)

type IpfsService interface {
	GetStatus(ctx context.Context) domain.IPFSHealthCheckResponse
	AddFile(ctx context.Context, requestBody *multipart.Reader) (string, error)
	GetFile(ctx context.Context, cid string) ([]byte, error)
	AddClusterNodePair(ctx context.Context, nodeUrl string, clusterUrl string) error
	GetPins(ctx context.Context) (*[]domain.Pin, error)
}
