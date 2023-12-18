package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	"fmt"
	"net"

	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/dbo"
)

func (dbo *DatabaseObject) CountBlocksByIPAddress(ipa net.IP) (n int, err error) {
	err = dbo.PreparedStatement(DbPsid_CountBlocksByIPAddress).QueryRow(ipa).Scan(&n)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) InsertBlock(ipa net.IP, durationSec uint) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_AddBlock).Exec(ipa, durationSec)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) IncreaseBlockDuration(ipa net.IP, deltaDurationSec uint) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_IncreaseBlockDuration).Exec(deltaDurationSec, ipa)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}
