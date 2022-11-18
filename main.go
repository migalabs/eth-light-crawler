package main

import (
	"context"
	"fmt"
	"os"

	"github.com/migalabs/eth-light-crawler/cmd"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var (
	CliName    = "ethereum light crawler"
	CliVersion = "v0.0.1"
)

func main() {
	fmt.Println(CliName, CliVersion)

	lightCrawler := cli.App{
		Name:      CliName,
		Usage:     "This tool uses Ethereum's peer-discovery protocol to measure the size of the Ethereum network (testnets included)",
		UsageText: "eth-light-crawler [subcommands] [arguments]",
		Commands: []*cli.Command{
			cmd.Discovery5,
		},
	}

	err := lightCrawler.RunContext(context.Background(), os.Args)
	if err != nil {
		log.Errorf("error running %s - %s", CliName, err.Error())
		os.Exit(1)
	}
}
