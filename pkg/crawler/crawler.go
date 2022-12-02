package crawler

import (
	"context"
	log "github.com/sirupsen/logrus"
	
)

type Crawler struct {
	ctx context.Context

	startT time.Time
	duration time.Duration

	enode *
}


func New(ctx context.Context, port ) *Crawler {

	// Generate a new PrivKey

	// Generate a Enode with custom ENR

	// Init the DB

	// Generate the Discovery5 service
	
	return 
}

func (c *Crawler) Run(duration time.Duration) error {
	// if duration has not been set, run until Crtl+C

	// otherwise, run it for X time

	return nil
}