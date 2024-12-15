package password

import "testing"

func Test_HashPassword(t *testing.T) {
	world, err := HashPassword("Qaz@1234")
	if err != nil {
		t.Error(err)
	}
	t.Log(world)
}

func Test_CheckPasswordHash(t *testing.T) {
	world, err := HashPassword("Qaz@1234")
	if err != nil {
		t.Error(err)
	}
	hash := CheckPasswordHash("Qaz@1234", world)
	if !hash {
		t.Error("hash failed")
	}
}
