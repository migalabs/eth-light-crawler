package db

import (
	"encoding/hex"

	gcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/migalabs/eth-light-crawler/pkg/discv5"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (d *DBClient) dropEnrTable() error {
	log.Debugf("droping enrs table in the db")

	_, err := d.psqlPool.Exec(d.ctx, `
		DROP TABLE enrs;
	`)
	return err

}

func (d *DBClient) initEnrDatabase() error {
	log.Debugf("initializing enrs table in the db")

	// try create the table in the DB
	_, err := d.psqlPool.Exec(
		d.ctx, `
		CREATE TABLE IF NOT EXISTS enrs(
			id SERIAL,
			timestamp BIGINT NOT NULL,
			node_id TEXT NOT NULL,
			seq BIGINT NOT NULL,
			ip TEXT NOT NULL,
			tcp INT,
			udp INT,
			pubkey TEXT NOT NULL,
			fork_digest TEXT,
			next_fork_version TEXT,
			attnets TEXT, 
			attnets_number INT,

			PRIMARY KEY(node_id)	
		);
		`,
	)
	if err != nil {
		return errors.Wrap(err, "unable to create table enrs in the db")
	}

	return nil
}

// Insert ENR in the DB
// insert into the db if new one, update the data if the ENR has a higher Seq number
func (d *DBClient) InsertEnr(enr *discv5.EnrNode) error {
	log.Debug("inserting enr in the db")

	pubBytes := gcrypto.FromECDSAPub(enr.Pubkey)
	pubkey := hex.EncodeToString(pubBytes)

	_, err := d.psqlPool.Exec(
		d.ctx, `
			INSERT INTO enrs(
				timestamp,
				node_id,
				seq,
				ip,
				tcp,
				udp,
				pubkey,
				fork_digest,
				next_fork_version,
				attnets,
				attnets_number)
			VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)	
		`,
		enr.Timestamp.Unix(),
		enr.ID.String(),
		enr.Seq,
		enr.IP,
		enr.TCP,
		enr.UDP,
		pubkey,
		enr.Eth2Data.ForkDigest.String(),
		enr.Eth2Data.NextForkVersion.String(),
		hex.EncodeToString(enr.Attnets.Raw[:]),
		enr.Attnets.NetNumber,
	)
	if err != nil {
		log.Error(enr)
		return errors.Wrap(err, "unable to insert enr")
	}

	return nil
}

// Update ENR in the DB
// insert into the db if new one, update the data if the ENR has a higher Seq number
func (d *DBClient) UpdateEnr(enr *discv5.EnrNode) error {
	log.Debug("inserting enr in the db")

	pubBytes := gcrypto.FromECDSAPub(enr.Pubkey)
	pubkey := hex.EncodeToString(pubBytes)

	_, err := d.psqlPool.Exec(
		d.ctx, `
			UPDATE enrs SET
				timestamp=$2,
				seq=$3,
				ip=$4,
				tcp=$5,
				udp=$6,
				pubkey=$7,
				fork_digest=$8,
				next_fork_version=$9,
				attnets=$10,
				attnets_number=$11
			WHERE node_id=$1
		`,
		enr.ID.String(),
		enr.Timestamp.Unix(),
		enr.Seq,
		enr.IP,
		enr.TCP,
		enr.UDP,
		pubkey,
		enr.Eth2Data.ForkDigest.String(),
		enr.Eth2Data.NextForkVersion.String(),
		hex.EncodeToString(enr.Attnets.Raw[:]),
		enr.Attnets.NetNumber,
	)
	if err != nil {
		log.Error(enr)
		return errors.Wrap(err, "unable to update enr")
	}
	return nil
}
