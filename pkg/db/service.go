package db

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/migalabs/eth-light-crawler/pkg/discv5"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type dbAction int8

const (
	insertItem dbAction = iota
	updateItem
	deleteItem
)

const (
	bufferSize    = 2048
	maxPersisters = 1
)

type DBClient struct {
	// Control Variables
	ctx context.Context
	m   sync.RWMutex

	// Pgx Postgres variables
	loginStr string
	psqlPool *pgxpool.Pool

	persistC   chan *PersistableItem
	persisters int
	persistWG  *sync.WaitGroup
	doneC      chan struct{}
}

func NewDBClient(
	ctx context.Context,
	loginStr string,
	initialized bool,
	reset bool) (*DBClient, error) {

	logEntry := logrus.WithField("module", "db-client")
	logEntry.WithFields(logrus.Fields{"endpoint": loginStr}).Debug("attempt connection to DB")

	// check if the login string has enough len
	if len(loginStr) == 0 {
		return nil, errors.New("empty db-endpoint provided")
	}

	// try connecting to the DB from the given logingStr
	pPool, err := pgxpool.Connect(ctx, loginStr)
	if err != nil {
		return nil, err
	}

	// check if the connection is successful
	err = pPool.Ping(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to ping db")
	}

	// generate all the necessary/control channels
	persistC := make(chan *PersistableItem, bufferSize)
	var persistWG sync.WaitGroup

	// compose the DBClient
	dbClient := &DBClient{
		ctx:       ctx,
		loginStr:  loginStr,
		psqlPool:  pPool,
		persistC:  persistC,
		persistWG: &persistWG,
		doneC:     make(chan struct{}),
	}

	// initialize all the tables
	if initialized {
		err = dbClient.initTables(reset)
		if err != nil {
			return nil, errors.Wrap(err, "unable to initialize the SQL tables at "+loginStr)
		}
	}

	// run the db persisters
	go dbClient.spawnPersisters()

	return dbClient, nil
}

func (c *DBClient) initTables(resetTables bool) error {
	// initialize all the necessary tables to perform the crawl

	var err error

	// drop Enr Table if requested
	if resetTables {
		err = c.dropEnrTable()
		if err != nil {
			return err
		}
	}

	// init Enr table
	err = c.initEnrDatabase()
	if err != nil {
		return err
	}

	return nil
}

func (c *DBClient) spawnPersisters() {
	// spaw as many persisters as defined in `maxPersisters`
	for persister := 1; persister <= maxPersisters; persister++ {
		c.launchPersister(persister)
		c.persistWG.Add(1)
		c.persisters++
	}
	logrus.Debugf("spawned total of %d db persister", c.persisters)
}

func (c *DBClient) launchPersister(persisterID int) error {
	logEntry := logrus.WithFields(logrus.Fields{"persisterID": persisterID})

	go func() {
		defer c.persistWG.Done()
		var sdFlag bool = false

		logEntry.Info("inititalizing persister")
		// check with higher priority if the main-ctx died
		for {
			// check if we need to close the persiter (and if the channel still has stuff to read)
			if (len(c.persistC) == 0) && sdFlag {
				logEntry.Info("signal to close the persister detected and there is nothing to read, closing persister")
				return
			}

			select {
			case obj := <-c.persistC: // persist any kind of item
				switch obj.Action {
				case insertItem:
					switch obj.Item.(type) {
					case (*discv5.EnrNode):
						enr := obj.Item.(*discv5.EnrNode)
						logrus.Debugf("updating enr for node %s", enr.ID)
						err := c.InsertEnr(enr)
						if err != nil {
							logEntry.Error(err)
						}
					default:
						logEntry.Error("unrecognized type of object received to insert into DB", obj)
					}
				case updateItem:
					switch obj.Item.(type) {
					case (*discv5.EnrNode):
						enr := obj.Item.(*discv5.EnrNode)
						logrus.Debugf("udpating enr for node %s", enr.ID)
						err := c.UpdateEnr(enr)
						if err != nil {
							logEntry.Error(err)
						}
					default:
						logEntry.Error("unrecognized type of object received to update into DB", obj)
					}
				case deleteItem:
					logEntry.Info("Delete operation still not supported")
				default:
					logEntry.Info("unable to understand operation", obj.Action)
				}
			case <-c.ctx.Done(): // check if the context of the tool died
				logEntry.Info("context died, clossing persister")
				return

			case <-c.doneC:
				sdFlag = true
			}
		}
	}()

	return nil

}

func (c *DBClient) Close() {
	// notify the persisters to finish
	for i := 0; i < c.persisters; i++ {
		c.doneC <- struct{}{}
	}

	c.persistWG.Wait()

	// close all the exisiting channels
	close(c.doneC)
	close(c.persistC)

	// close safelly the connection with PSQL
	c.psqlPool.Close()
}

type PersistableItem struct {
	Action dbAction
	Item   interface{}
}

func newPersistable(item interface{}, action dbAction) *PersistableItem {
	return &PersistableItem{
		Action: action,
		Item:   item,
	}
}

func (c *DBClient) InsertIntoDB(persItem interface{}) {
	item := newPersistable(persItem, insertItem)
	c.persistC <- item
}

func (c *DBClient) UpdateInDB(persItem interface{}) {
	item := newPersistable(persItem, updateItem)
	c.persistC <- item
}
