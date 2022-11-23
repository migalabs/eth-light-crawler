package cmd
package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"strings"
	
	"github.com/migalabs/eth-light-crawler/pkg/config"
	"github.com/ethereum/go-ethereum/p2p/enode"

	gcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p-core/crypto"

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

func RunDiscv5(ctx *cli.Context) error {
	var (
		nodeKey *ecdsa.PrivateKey
	)
	// all the magic goes here

	// Ethereum compatible PrivateKey
	// check: https://github.com/migalabs/armiarma/blob/ca3d2f6adea364fc7f38bdabda912b5541bb4154/src/utils/keys.go#L52
	utils.GeneratePrivKey()

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
	return nil
}
