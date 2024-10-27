package base

import (
	"fmt"
	"reflect"

	num "github.com/vault-thirteen/auxie/number"
)

const (
	ErrF_CanNotScanSourceAsInt = "source can not be scanned as int: %s"
)

func ScanSrcAsInt(src any) (i int, err error) {
	switch src.(type) {
	case int64:
		i = int(src.(int64))

	case int:
		i = src.(int)

	case []byte:
		i, err = num.ParseInt(string(src.([]byte)))
		if err != nil {
			return 0, err
		}

	case string:
		i, err = num.ParseInt(src.(string))
		if err != nil {
			return 0, err
		}

	default:
		return 0, fmt.Errorf(ErrF_CanNotScanSourceAsInt, reflect.TypeOf(src).String())
	}

	return i, nil
}
