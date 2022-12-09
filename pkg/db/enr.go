package db

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"

	"github.com/migalabs/eth-light-crawler/pkg/discv5"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

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
			fork_epoch INT,
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

	eth2Data, err := enr.ParseEth2Data()
	if err != nil {
		return errors.New(fmt.Sprintf("unable to parse the Eth2Data for enr %s", enr))
	}
	attnets, err := enr.ParseAttnets()
	if err != nil {
		return errors.New(fmt.Sprintf("unable to parse the attestation networks for enr %s", enr))
	}

	pubBytes, _ := x509.MarshalPKIXPublicKey(enr.Pubkey)
	pubkey := hex.EncodeToString(pubBytes)

	_, err = d.psqlPool.Exec(
		d.ctx, `
			INSERT INTO enrs(
				timestampt,
				node_id,
				seq,
				ip,
				tcp,
				udp,
				pubkey,
				fork_digest,
				next_fork_version,
				fork_epoch,
				attnets,
				attnets_number)
			VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)	
		`,
		enr.Timestamp.Unix(),
		enr.ID.String(),
		enr.Seq,
		enr.IP,
		enr.TCP,
		enr.UDP,
		pubkey,
		eth2Data.ForkDigest.String(),
		eth2Data.NextForkVersion.String(),
		eth2Data.NextForkEpoch,
		attnets.Raw.String(),
		attnets.NetNumber,
	)
	if err != nil {
		return err
	}

	return nil
}

// Update ENR in the DB
// insert into the db if new one, update the data if the ENR has a higher Seq number
func (d *DBClient) UpdateEnr(enr *discv5.EnrNode) error {
	log.Debug("inserting enr in the db")

	eth2Data, err := enr.ParseEth2Data()
	if err != nil {
		return errors.New(fmt.Sprintf("unable to parse the Eth2Data for enr", enr))
	}
	attnets, err := enr.ParseAttnets()
	if err != nil {
		return errors.New(fmt.Sprintf("unable to parse the attestation networks for enr", enr))
	}
	pubBytes, _ := x509.MarshalPKIXPublicKey(enr.Pubkey)
	pubkey := hex.EncodeToString(pubBytes)

	_, err = d.psqlPool.Exec(
		d.ctx, `
			UPDATE enrs SET(
				timestampt=$2,
				seq=$3,
				ip=$4,
				tcp=$5,
				udp=$6,
				pubkey=$7,
				fork_digest=$8,
				next_fork_version=$9,
				fork_epoch=$10,
				attnets=$11,
				attnets_number=$12)
			WHERE node_id=$1 and seq < $3	
		`,
		enr.ID.String(),
		enr.Timestamp.Unix(),
		enr.Seq,
		enr.IP,
		enr.TCP,
		enr.UDP,
		pubkey,
		eth2Data.ForkDigest.String(),
		eth2Data.NextForkVersion.String(),
		eth2Data.NextForkEpoch,
		attnets.Raw.String(),
		attnets.NetNumber,
	)
	if err != nil {
		return err
	}
	return nil
}
