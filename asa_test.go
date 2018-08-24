package goasa_test

import (
	"flag"
	"os"
	"testing"

	"github.com/golang/glog"
	goasa "github.com/remiphilippe/go-asa"
)

func initTest() (*goasa.ASA, error) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Lookup("v").Value.Set("2")

	params := make(map[string]string)
	params["username"] = os.Getenv("ASA_USER")
	params["password"] = os.Getenv("ASA_PASSWORD")
	params["debug"] = "true"
	params["insecure"] = "true"

	asa, err := goasa.NewASA(os.Getenv("ASA_HOST"), params)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return nil, err
	}

	return asa, nil
}

func TestLoginFail(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Lookup("v").Value.Set("2")

	params := make(map[string]string)
	params["username"] = "bob"
	params["password"] = "sponge"
	params["debug"] = "true"
	params["insecure"] = "true"

	asa, err := goasa.NewASA(os.Getenv("ASA_HOST"), params)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	_, err = asa.Get("monitoring/device/components/version", nil)
	if err == nil {
		t.Errorf("error: we should have failed here\n")
		return
	}
}

func TestLoginOK(t *testing.T) {
	var err error

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	_, err = asa.Get("monitoring/device/components/version", nil)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}
}

func TestError(t *testing.T) {
	var err error

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	_, err = asa.Get("toto", nil)
	if err == nil {
		t.Errorf("error: we should have failed here\n")
		return
	}

	if _, ok := err.(goasa.ASAError); !ok {
		t.Errorf("error: we don't have an ASAError, WTF?\n")
		return
	}
}
