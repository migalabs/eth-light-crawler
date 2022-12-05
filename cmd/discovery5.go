package cmd

import (
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/migalabs/armiarma/src/utils"
	"github.com/migalabs/eth-light-crawler/pkg/config"
	"net"

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

type Eth2InfoData struct {
	IP            net.IP
	TcpPort       int
	UdpPort       int
	UserAgent     string
	TopicArray    []string
	Network       string
	ForkDigest    string
	DbEndpoint    string
	Eth2endpoint  string
	LogLevel      string
	PrivateKey    *crypto.Secp256k1PrivateKey
	BootNodesFile string
	OutputPath    string
}

func RunDiscv5(ctx *cli.Context) error {
	// all the magic goes here

	// Ethereum compatible PrivateKey
	// check: https://github.com/migalabs/armiarma/blob/ca3d2f6adea364fc7f38bdabda912b5541bb4154/src/utils/keys.go#L52
	i := Eth2InfoData{}
	// Private Key
	i.PrivateKey = utils.GeneratePrivKey()

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
