package discv5

import (
	"context"
)

type Discv5Service struct {
	ctx context.Context

	enode
}

func NewService(ctx context.Context, enode *eth.enode) (*Discv5Service, error) {

	return &Discv5Service{
		ctx,
		enode,
	}
}
