package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	"net"

	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/dbo"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

func (dbo *DatabaseObject) CountBlocksByIPAddress(ipa net.IP) (n cmb.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountBlocksByIPAddress).QueryRow(ipa)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) InsertBlock(ipa net.IP, durationSec cmb.Count) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_AddBlock).Exec(ipa, durationSec)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) IncreaseBlockDuration(ipa net.IP, deltaDurationSec cmb.Count) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_IncreaseBlockDuration).Exec(deltaDurationSec, ipa)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}
