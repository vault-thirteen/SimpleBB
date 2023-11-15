package ul

import (
	"database/sql/driver"
	"testing"

	"github.com/vault-thirteen/tester"
)

func Test_NewFromArray(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1.
	ul, err = NewFromArray(nil)
	aTest.MustBeNoError(err)
	x := UidList(nil)
	aTest.MustBeEqual(ul, &x)

	// Test #2.
	ul, err = NewFromArray([]uint{})
	aTest.MustBeNoError(err)
	x = UidList([]uint{})
	aTest.MustBeEqual(ul, &x)

	// Test #3.
	ul, err = NewFromArray([]uint{1, 2, 3})
	aTest.MustBeNoError(err)
	x = UidList([]uint{1, 2, 3})
	aTest.MustBeEqual(ul, &x)
}

func Test_Size(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1.
	ul = nil
	aTest.MustBeEqual(ul.Size(), 0)

	// Test #2.
	ul, err = NewFromArray([]uint{})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.Size(), 0)

	// Test #3.
	ul, err = NewFromArray([]uint{1, 2, 3})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.Size(), 3)
}

func Test_CheckIntegrity(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1.
	ul = nil
	err = ul.CheckIntegrity()
	aTest.MustBeNoError(err)

	// Test #2.
	tmp := UidList([]uint{1, 2, 3})
	ul = &tmp
	err = ul.CheckIntegrity()
	aTest.MustBeNoError(err)

	// Test #3.
	tmp = UidList([]uint{1, 2, 3, 2})
	ul = &tmp
	err = ul.CheckIntegrity()
	aTest.MustBeAnError(err)
}

func Test_AddItem(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Tests.
	ul = New()
	err = ul.AddItem(1)
	aTest.MustBeNoError(err)
	err = ul.AddItem(2)
	aTest.MustBeNoError(err)
	err = ul.AddItem(3)
	aTest.MustBeNoError(err)
	err = ul.AddItem(2)
	aTest.MustBeAnError(err)
}

func Test_RemoveItem(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1.
	ul = New()
	err = ul.RemoveItem(1)
	aTest.MustBeAnError(err)

	// Test #2. Non-existent item.
	ul, err = NewFromArray([]uint{1, 2, 3, 4})
	aTest.MustBeNoError(err)
	err = ul.RemoveItem(5)
	aTest.MustBeAnError(err)

	// Test #3. First item.
	ul, err = NewFromArray([]uint{1, 2, 3})
	aTest.MustBeNoError(err)
	err = ul.RemoveItem(1)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(*ul, UidList([]uint{2, 3}))

	// Test #4. Middle item.
	ul, err = NewFromArray([]uint{1, 2, 3})
	aTest.MustBeNoError(err)
	err = ul.RemoveItem(2)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(*ul, UidList([]uint{1, 3}))

	// Test #3.
	ul, err = NewFromArray([]uint{1, 2, 3})
	aTest.MustBeNoError(err)
	err = ul.RemoveItem(3)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(*ul, UidList([]uint{1, 2}))
}

func Test_Scan(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1.
	ul = nil
	err = ul.Scan(false)
	aTest.MustBeAnError(err)

	// Test #2.
	ul = New()
	err = ul.Scan(nil)
	aTest.MustBeNoError(err)

	// Test #3.
	ul = New()
	err = ul.Scan("string")
	aTest.MustBeAnError(err)

	// Test #4.
	ul = New()
	err = ul.Scan([]byte("[1,2,3]"))
	aTest.MustBeNoError(err)
	tmp := UidList([]uint{1, 2, 3})
	aTest.MustBeEqual(ul, &tmp)

	// Test #5.
	ul = New()
	err = ul.Scan([]byte("[1,2,3,NaN]"))
	aTest.MustBeAnError(err)

	// Test #6.
	ul = New()
	err = ul.Scan([]byte(""))
	aTest.MustBeAnError(err)

	// Test #7.
	ul = New()
	err = ul.Scan([]byte("[]"))
	aTest.MustBeNoError(err)
}

func Test_Value(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var dv driver.Value
	var err error

	// Test #1.
	ul = nil
	dv, err = ul.Value()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(dv, nil)

	// Test #2.
	ul, err = NewFromArray([]uint{1, 2, 3})
	aTest.MustBeNoError(err)
	dv, err = ul.Value()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(dv, driver.Value([]byte("[1,2,3]")))
}

func Test_ValuesString(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error
	var vs string

	// Test #1.
	ul, err = NewFromArray(nil)
	aTest.MustBeNoError(err)
	vs, err = ul.ValuesString()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(vs, "")

	// Test #2.
	ul, err = NewFromArray([]uint{})
	aTest.MustBeNoError(err)
	vs, err = ul.ValuesString()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(vs, "")

	// Test #3.
	ul, err = NewFromArray([]uint{1})
	aTest.MustBeNoError(err)
	vs, err = ul.ValuesString()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(vs, "1")

	// Test #4.
	ul, err = NewFromArray([]uint{1, 2, 3})
	aTest.MustBeNoError(err)
	vs, err = ul.ValuesString()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(vs, "1,2,3")
}

func Test_OnPage(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1.
	ul, err = NewFromArray(nil)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.OnPage(0, 1), UidList([]uint(nil)))

	// Test #2.
	ul, err = NewFromArray(nil)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.OnPage(1, 1), UidList([]uint(nil)))

	// Test #3.
	ul, err = NewFromArray([]uint{})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.OnPage(1, 1), UidList([]uint(nil)))

	// Test #4.
	ul, err = NewFromArray([]uint{1, 2, 3, 4, 5})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.OnPage(1, 5), UidList([]uint{1, 2, 3, 4, 5}))

	// Test #5.
	ul, err = NewFromArray([]uint{1, 2, 3})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.OnPage(1, 5), UidList([]uint{1, 2, 3}))

	// Test #6.
	ul, err = NewFromArray([]uint{1, 2, 3, 4, 5, 6, 7})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.OnPage(1, 5), UidList([]uint{1, 2, 3, 4, 5}))

	// Test #7.
	ul, err = NewFromArray([]uint{1, 2, 3, 4, 5})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.OnPage(2, 5), UidList([]uint(nil)))

	// Test #8.
	ul, err = NewFromArray([]uint{1, 2, 3, 4, 5, 6})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.OnPage(2, 5), UidList([]uint{6}))

	// Test #9.
	ul, err = NewFromArray([]uint{1, 2, 3, 4, 5, 6, 7})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.OnPage(2, 5), UidList([]uint{6, 7}))

	// Test #10.
	ul, err = NewFromArray([]uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.OnPage(2, 5), UidList([]uint{6, 7, 8, 9, 10}))
}
