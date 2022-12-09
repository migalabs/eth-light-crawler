package discv5

import (
	"crypto/ecdsa"
	"encoding/binary"
	"math/bits"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/migalabs/armiarma/src/utils"
	"github.com/protolambda/zrnt/eth2/beacon/common"
)

type EnrNode struct {
	Timestamp time.Time
	ID        enode.ID
	IP        net.IP
	Seq       uint64
	UDP       int
	TCP       int
	Pubkey    *ecdsa.PublicKey
	Eth2Data  utils.Eth2ENREntry
	Attnets   utils.AttnetsENREntry
}

func NewEnrNode(nodeID enode.ID) *EnrNode {

	return &EnrNode{
		Timestamp: time.Now(),
	}
}

func (e *EnrNode) ParseEth2Data() (*common.Eth2Data, error) {
	return e.Eth2Data.Eth2Data()
}

type Attnets struct {
	Raw       common.AttnetBits
	NetNumber int
}

func (e *EnrNode) ParseAttnets() (*Attnets, error) {

	// count the number of bits in the Attnets
	bits := CountBits(e.Attnets)

	var attnets common.AttnetBits
	err := attnets.UnmarshalText(e.Attnets)
	if err != nil {
		return nil, err
	}

	return &Attnets{
		Raw:       attnets,
		NetNumber: bits,
	}, nil
}

func CountBits(byteArr []byte) int {
	rawInt := binary.BigEndian.Uint64(byteArr)
	return bits.OnesCount64(rawInt)
}
