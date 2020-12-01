package services

import (
	"context"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type ProdService struct {
}

func (c *ProdService) GetProdStock(ctx context.Context, req *ProdRequest) (*ProdResponse, error) {
	
	return &ProdResponse{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		ProdStock:     req.ProdId + 1,
	}, nil
}
