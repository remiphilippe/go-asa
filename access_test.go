package goasa_test

import (
	"testing"

	"github.com/remiphilippe/go-asa"
)

func TestAccesskObjectGetAll(t *testing.T) {
	var err error

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	_, err = asa.GetAllGlobalAccessObjects()
	if err == nil {
		t.Errorf("error: we were expecting 404 (no results)\n")
	}
}

func TestCreateAccessRule(t *testing.T) {
	var err error
	var o []*goasa.AccessRule
	var r *goasa.AccessRule

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	ip := new(goasa.ReferenceObject)
	ip.Kind = "NetworkProtocol"
	ip.Value = "ip"

	r = new(goasa.AccessRule)
	r.SourceAddress = asa.SystemObjectAny()
	r.DestinationAddress = asa.SystemObjectAny()
	r.SourceService = ip
	r.DestinationService = ip
	r.Permit = true
	r.Active = true

	err = asa.CreateAccessRule("global", r, goasa.DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	o, err = asa.GetAllGlobalAccessObjects()
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	if len(o) != 1 {
		t.Errorf("error: expecting one result\n")
	}

	err = asa.DeleteAccessRule("global", o[0])
	if err != nil {
		t.Errorf("error: %s\n", err)
	}
}
