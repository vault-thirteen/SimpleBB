package dbo

import (
	"database/sql"
	"fmt"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

func CheckRowsAffected(sqlResult sql.Result, expectedValue cmb.Count) (err error) {
	var ra int64
	ra, err = sqlResult.RowsAffected()
	if err != nil {
		return err
	}

	if cmb.Count(ra) != expectedValue {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func GetLastInsertedId(sqlResult sql.Result) (lastInsertedId cmb.Id, err error) {
	var x int64
	x, err = sqlResult.LastInsertId()
	if err != nil {
		return LastInsertedIdOnError, err
	}

	return cmb.Id(x), nil
}

func CheckRowsAffectedAndGetLastInsertedId(sqlResult sql.Result, expectedValue cmb.Count) (lastInsertedId cmb.Id, err error) {
	err = CheckRowsAffected(sqlResult, expectedValue)
	if err != nil {
		return LastInsertedIdOnError, err
	}

	return GetLastInsertedId(sqlResult)
}
