package goasa_test

import "testing"

func TestSave(t *testing.T) {
	var err error

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	err = asa.Save()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}
}
