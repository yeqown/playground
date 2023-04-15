package basic_test

import (
	"testing"
	"unsafe"
)

type User struct {
	B0 float64
	B1 float64
	B2 float64
	B3 float64
	B4 float64
	B5 float64
}

func addMoney(u *User, moneyType uint, money float64) {
	//switch {
	//case moneyType == 0:
	//	u.B0 += money
	//case moneyType == 1:
	//	u.B1 += money
	//case moneyType == 2:
	//	u.B2 += money
	//case moneyType == 3:
	//	u.B3 += money
	//case moneyType == 4:
	//	u.B4 += money
	//case moneyType == 5:
	//	u.B5 += money
	//}

	p := unsafe.Pointer(uintptr(unsafe.Pointer(&u.B0)) + uintptr(moneyType)*unsafe.Sizeof(money))
	*(*float64)(p) += money
}

func Test_unsafe(t *testing.T) {
	u := new(User)

	addMoney(u, 0, 0.01)
	t.Logf("%+v", u)
	addMoney(u, 1, 1.00)
	t.Logf("%+v", u)
	addMoney(u, 2, 2.00)
	t.Logf("%+v", u)
	addMoney(u, 3, 3.00)
	t.Logf("%+v", u)
	addMoney(u, 4, 4.00)
	t.Logf("%+v", u)
	addMoney(u, 5, 5.00)
	t.Logf("%+v", u)

}
