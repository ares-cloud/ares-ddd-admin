package ipcity

import "testing"

func Test_GetGetLocation(t *testing.T) {
	du, err := GetGetLocationBaiDu("125.76.174.38")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(du)
}
