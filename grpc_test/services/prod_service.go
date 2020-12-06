package services

import (
	"context"
	"google.golang.org/protobuf/runtime/protoimpl"
	"math/rand"
)

type ProdService struct {
}

func (c *ProdService) GetProdStock(ctx context.Context, req *ProdRequest) (*ProdResponse, error) {
	return &ProdResponse{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		ProdStock:     rand.Int31n(1000),
	}, nil
}
