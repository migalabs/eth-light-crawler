package config

import (
	cli "github.com/urfave/cli/v2"
)

type Config struct {
	IP     string
	UDP    string
	TCP    string
	LogLvl string
}

var DefaultConfig Config = Config{
	IP:     "0.0.0.0",
	UDP:    "9001",
	TCP:    "9001",
	LogLvl: "info",
}

func (c *Config) Apply(ctx *cli.Context) {
	// read log-level from the ctx
	if ctx.IsSet("log-level") {
		c.LogLvl = ctx.String("log-level")
	}
	// read the port from the ctx
	if ctx.IsSet("port") {
		port := ctx.String("port")
		c.UDP = port
		c.TCP = port
	}
	// more args?
}
