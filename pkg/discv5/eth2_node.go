package discv5

import (
	"crypto/ecdsa"
	"net"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
)

type EnrNode struct {
	ID              enode.ID
	IP              net.IP
	Seq             int64
	UDP             int
	TCP             int
	PubKey          *ecdsa.PublicKey
	Record          *enr.Record
	ForkDigest      [4]byte
	NextForkVersion [4]byte
	NextForkEpoch   uint64
}

func ParseEnrRecords(records *enr.Record) (forkDigest [4]byte, ForkVersion [4]byte, ForkEpoch [4]byte) {

	return
}
