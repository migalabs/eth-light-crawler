package crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/migalabs/eth-light-crawler/pkg/config"
	"github.com/migalabs/eth-light-crawler/pkg/discv5"
	"github.com/migalabs/eth-light-crawler/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Crawler struct {
	ctx context.Context

	startT   time.Time
	duration time.Duration

	ethNode       *enode.LocalNode
	discv5Service *discv5.Discv5Service
}

func New(ctx context.Context, dbPath string, port int) (*Crawler, error) {
	// Generate a new PrivKey
	privK, err := utils.GenNewPrivKey()
	if err != nil {
		return nil, errors.Wrap(err, "error generating privkey")
	}

	// Init the DB
	enodeDB, err := enode.OpenDB(dbPath)
	if err != nil {
		return nil, err
	}
	// Generate a Enode with custom ENR
	ethNode := enode.NewLocalNode(enodeDB, privK)

	// define the Handler for when we discover a new ENR
	enrHandler := func(node *enode.Node) {
		// check if the node is valid
		err := node.ValidateComplete()
		if err != nil {
			log.Warnf("error validating the ENR - ", err.Error())
		}
		// extract the information from the enode
		id := node.ID()
		seq := node.Seq()
		ip := node.IP()
		udp := node.UDP()
		tcp := node.TCP()
		pubkey := node.Pubkey()
		record := node.Record()
		// eth2 := record.Load("")
		// attnets := record.Load("")

		fmt.Println(seq, id, pubkey, ip, udp, tcp, record)

	}

	// Generate the Discovery5 service
	discv5Serv, err := discv5.NewService(ctx, port, privK, ethNode, config.EthBootonodes, enrHandler)

	return &Crawler{
		ctx:           ctx,
		ethNode:       ethNode,
		discv5Service: discv5Serv,
	}, nil
}

func (c *Crawler) Run(duration time.Duration) error {
	// if duration has not been set, run until Crtl+C
	c.discv5Service.Run()
	// otherwise, run it for X time

	return nil
}
