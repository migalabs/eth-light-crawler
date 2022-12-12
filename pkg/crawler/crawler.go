package crawler

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/migalabs/armiarma/src/utils"
	"github.com/migalabs/eth-light-crawler/pkg/config"
	"github.com/migalabs/eth-light-crawler/pkg/db"
	"github.com/migalabs/eth-light-crawler/pkg/discv5"
	ut "github.com/migalabs/eth-light-crawler/pkg/utils"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	log "github.com/sirupsen/logrus"
)

var (
	EmptyBytes error = errors.New("array of bytes is empty")
)

type Crawler struct {
	ctx context.Context

	startT   time.Time
	duration time.Duration

	ethNode       *enode.LocalNode
	discv5Service *discv5.Discv5Service
	dbClient      *db.DBClient

	enrCache map[enode.ID]int64
}

func New(ctx context.Context, dbEndpoint string, dbPath string, port int, resetDB bool) (*Crawler, error) {
	// Generate a new PrivKey
	privK, err := ut.GenNewPrivKey()
	if err != nil {
		return nil, errors.Wrap(err, "error generating privkey")
	}

	// Init the ethereum peerstore
	enodeDB, err := enode.OpenDB(dbPath)
	if err != nil {
		return nil, err
	}

	// Create a new
	sqlDB, err := db.NewDBClient(ctx, dbEndpoint, true, resetDB)
	if err != nil {
		return nil, err
	}

	// Generate a Enode with custom ENR
	ethNode := enode.NewLocalNode(enodeDB, privK)

	// generate the cache of node_ids > seq numbers
	enrCache := make(map[enode.ID]uint64)

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

		// Retrieve the Fork Digest and the attestnets
		eth2Data, ok, err := utils.ParseNodeEth2Data(*node)
		if !ok {
			eth2Data = new(common.Eth2Data)
		} else {
			if err != nil {
				log.Error(errors.Wrap(err, "eth2 data parsing error"))
				// eth2Data = new(common.Eth2Data)
			}
		}

		attnets, ok, err := discv5.ParseAttnets(*node)
		if !ok {
			attnets = new(discv5.Attnets)
		} else {
			if err != nil {
				log.Error(errors.Wrap(err, "attnets parsing err"))
				// attnets = new(discv5.Attnets)
			}
		}
		// create a new ENR node
		enrNode := discv5.NewEnrNode(id)

		// add all the fields from the CL network
		enrNode.Seq = seq
		enrNode.IP = ip
		enrNode.TCP = tcp
		enrNode.UDP = udp
		enrNode.Pubkey = pubkey
		enrNode.Eth2Data = eth2Data
		enrNode.Attnets = attnets

		log.WithFields(log.Fields{
			"node_id":     id,
			"ip":          ip,
			"udp":         udp,
			"tcp":         tcp,
			"fork_digest": eth2Data.ForkDigest,
			"fork_epoch":  eth2Data.NextForkEpoch,
			"attnets":     hex.EncodeToString(attnets.Raw[:]),
			"att_number":  attnets.NetNumber,
			"enr":         node.String(),
		}).Info("Eth node found")

		// decide whether we need to insert or update an existing
		prevSeq, ok := enrCache[enrNode.ID]
		if !ok { // Insert not previously tracked enr
			sqlDB.InsertIntoDB(enrNode)
		} else if enrNode.Seq > prevSeq { // Update the the data of the given Node
			sqlDB.UpdateInDB(enrNode)
		}
		enrCache[enrNode.ID] = enrNode.Seq
	}

	// Generate the Discovery5 service
	discv5Serv, err := discv5.NewService(ctx, port, privK, ethNode, config.EthBootonodes, enrHandler)
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate the discv5 service")
	}

	return &Crawler{
		ctx:           ctx,
		ethNode:       ethNode,
		dbClient:      sqlDB,
		discv5Service: discv5Serv,
	}, nil
}

func (c *Crawler) Run(duration time.Duration) error {
	// if duration has not been set, run until Crtl+C
	c.discv5Service.Run()
	// otherwise, run it for X time

	return nil
}

func (c *Crawler) ID() string {
	return c.ethNode.ID().String()
}
