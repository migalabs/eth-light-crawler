package cmd

import (
	"github.com/libp2p/go-libp2p-kad-dht/crawler"
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
		&cli.StringFlag{
			Name:     "port",
			Usage:    "port number that we want to use/advertise in the Ethereum network",
			Value:    "9001",
			Required: true,
		},
	},
}

func RunDiscv5(ctx *cli.Context) error {
	// parse the configuration from the flags
	conf := config.DefaultConfig
	conf.Apply(ctx)

	// Create a new crawler
	crawlr := crawler.New()

	log.WithFields(log.Fields{
		"peerID":    "whatever the peerID is resulting from the Privkey",
		"IP":        conf.IP,
		"UDP":       conf.UDP,
		"TCP":       conf.TCP,
		"bootnodes": len(config.EthBootonodes),
		"log-info":  conf.LogLvl,
	}).Info("Starting discv node")

	// run the crawler for XX time

	return nil
}
