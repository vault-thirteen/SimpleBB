package ul

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	ErrDestinationIsNotInitialized = "destination is not initialized"
	ErrFUnsupportedDataType        = "unsupported data type: %s"
	ErrFDuplicateUid               = "duplicated uid: %v"
	ErrFUidIsNotFound              = "uid is not found: %v"
)

// UidList is a list unique identifiers.
//
// The main purpose of this list is to store a chronological order of all added
// items. An identifier is an unsigned integer number. The order of items in
// the list is important. New items are added to the end of the list, deleted
// items shift existing items.
type UidList []uint

func New() (ul *UidList) {
	return new(UidList)
}

func NewFromArray(uids []uint) (ul *UidList, err error) {
	tmp := UidList(uids)
	ul = &tmp

	err = ul.CheckIntegrity()
	if err != nil {
		return nil, err
	}

	return ul, nil
}

// CheckIntegrity verifies integrity of the list.
// All items must be unique to pass the check.
func (ul *UidList) CheckIntegrity() (err error) {
	if ul == nil {
		return nil
	}

	m := make(map[uint]bool)
	var isDuplicate bool

	for _, uid := range *ul {
		_, isDuplicate = m[uid]
		if isDuplicate {
			return fmt.Errorf(ErrFDuplicateUid, uid)
		}

		m[uid] = true
	}

	return nil
}

func (ul *UidList) Size() (n int) {
	if ul == nil {
		return 0
	} else {
		return len(*ul)
	}
}

// AddItem add a new identifier to the end of the list.
func (ul *UidList) AddItem(uid uint) (err error) {
	// Check for uniqueness.
	for _, x := range *ul {
		if x == uid {
			return fmt.Errorf(ErrFDuplicateUid, uid)
		}
	}

	*ul = append(*ul, uid)

	return nil
}

// RemoveItem deletes an identifier from the list shifting its items.
func (ul *UidList) RemoveItem(uid uint) (err error) {
	// Does the item really exist ?
	ok := false
	var pos int // Index of the removed item.
	var x uint
	for pos, x = range *ul {
		if x == uid {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf(ErrFUidIsNotFound, uid)
	}

	lastIndex := len(*ul) - 1

	if pos == lastIndex {
		[]uint(*ul)[lastIndex] = 0
		*ul = (*ul)[:lastIndex]
	} else {
		copy((*ul)[pos:], (*ul)[pos+1:])
		[]uint(*ul)[lastIndex] = 0
		*ul = (*ul)[:lastIndex]
	}

	return nil
}

// Scan method provides compatibility with SQL JSON data type.
func (ul *UidList) Scan(src any) (err error) {
	if ul == nil {
		return errors.New(ErrDestinationIsNotInitialized)
	}

	switch src.(type) {
	case []byte:
		{
			data := new(UidList)

			err = json.Unmarshal(src.([]byte), data)
			if err != nil {
				return err
			}

			if data != nil {
				*ul = *data
			}

			return nil
		}
	case nil:
		return nil
	default:
		return fmt.Errorf(ErrFUnsupportedDataType, reflect.TypeOf(src).String())
	}
}

// Value method provides compatibility with SQL JSON data type.
func (ul *UidList) Value() (dv driver.Value, err error) {
	if ul == nil {
		return nil, nil
	}

	var buf []byte
	buf, err = json.Marshal(ul)
	if err != nil {
		return nil, err
	}

	return driver.Value(buf), nil
}

func (ul *UidList) ValuesString() (values string, err error) {
	if ul == nil {
		return "", nil
	}

	if len(*ul) == 0 {
		return "", nil
	}

	var sb = strings.Builder{}
	iLast := len(*ul) - 1
	for i, uid := range *ul {
		if i < iLast {
			_, err = sb.WriteString(strconv.FormatUint(uint64(uid), 10) + ",")
		} else {
			_, err = sb.WriteString(strconv.FormatUint(uint64(uid), 10))
		}
		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

func (ul *UidList) OnPage(pageNumber uint, pageSize uint) (ulop UidList) {
	if pageNumber < 1 {
		return nil
	}

	if ul == nil {
		return nil
	}

	if *ul == nil {
		return nil
	}

	if len(*ul) == 0 {
		return nil
	}

	// Last index in array.
	iLast := uint(len(*ul) - 1)
	if iLast < 0 {
		return nil
	}

	// Left index of a virtual page.
	ipL := pageSize * (pageNumber - 1)
	if iLast < ipL {
		return nil
	}

	// Right index of a virtual page.
	ipR := ipL + pageSize - 1
	if iLast < ipR {
		return (*ul)[ipL : iLast+1]
	} else {
		return (*ul)[ipL : ipR+1]
	}
}
