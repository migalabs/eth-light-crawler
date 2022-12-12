package cmd

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/migalabs/armiarma/src/discovery/dv5"
	"github.com/migalabs/armiarma/src/info"
	"github.com/migalabs/eth-light-crawler/pkg/config"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var Discovery5 = &cli.Command{
	Name:   "discv5",
	Usage:  "crawl Ethereum's public DHT thought the Discovery 5.1 protocol",
	Action: RunDiscv5,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "verbosity of the logs that will be displayed [debug,warn,info,error]",
			EnvVars: []string{"IPFS_CID_HOARDER_LOGLEVEL"},
			Value:   "info",
		},
	},
}

var (
	nodeKey *ecdsa.PrivateKey
)

func GenNewPrivKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(gcrypto.S256(), rand.Reader)
}

type LocalNode struct {
	ctx       context.Context
	LocalNode enode.LocalNode
	info_data *info.Eth2InfoData
}

func RunDiscv5(ctx *cli.Context) error {
	// all the magic goes here
	// Ethereum compatible PrivateKey
	// check: https://github.com/migalabs/armiarma/blob/ca3d2f6adea364fc7f38bdabda912b5541bb4154/src/utils/keys.go#L52
	// Private Key
	nodeKey, _ := GenNewPrivKey()

	log.WithFields(log.Fields{
		"peerID":    "whatever the peerID is resulting from the Privkey",
		"IP":        config.DefaultIP,
		"port":      config.DefaultPort,
		"bootnodes": len(config.EthBootonodes),
	}).Info("Starting discv node")

	// Ethereum node
	// check: https://github.com/ethereum/go-ethereum/blob/c2e0abce2eedc1ba2a1b32c46fd07ef18a25354a/p2p/enode/localnode.go#L70
	db, _ := enode.OpenDB("")
	ln := enode.NewLocalNode(db, nodeKey)

	// Discovery5 service
	// check: https://github.com/migalabs/armiarma/blob/ca3d2f6adea364fc7f38bdabda912b5541bb4154/src/discovery/dv5/dv5_service.go#L58
	// Bootnodes are in pkg/config/bootnodes
	c := ctx.Context
	infObj, _ := info.InitEth2(c)
	dv5, _ := dv5.NewDiscovery(
		ctx.Context,
		ln, // LocalNode?
		nodeKey,
		Bootonodes, // pkg/config/bootnodes?
		infObj.ForkDigest,
		9006)

	return nil
}
