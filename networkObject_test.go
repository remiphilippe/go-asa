package goasa_test

import (
	"fmt"
	"testing"

	"github.com/remiphilippe/go-asa"
)

func TestNetworkObjectGetAll(t *testing.T) {
	// disabling this test for now as it's a bit long
	// if true {
	// 	return
	// }

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	var err error
	var n *goasa.NetworkObject
	var o []*goasa.NetworkObject
	var all []*goasa.NetworkObject

	limit := 202

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	o, err = asa.GetAllNetworkObjects()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(o) > 0 {
		t.Errorf("error: too many results, expecting 0 got %d\n", len(o))
	}

	for i := 0; i < limit; i++ {
		ip := fmt.Sprintf("1.1.1.%d", i)
		n = new(goasa.NetworkObject)
		n.Host.Kind = "IPv4Address"
		n.Host.Value = ip
		n.Name = ip
		n.ObjectID = ip

		err = asa.CreateNetworkObject(n, goasa.DuplicateActionDoNothing)
		if err != nil {
			t.Errorf("error: %s\n", err)
		}

		all = append(all, n)
	}

	o, err = asa.GetAllNetworkObjects()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(o) != limit {
		t.Errorf("error: not matching limit, expecting %d got %d\n", limit, len(o))
	}

	o, err = asa.GetNetworkObjects(10)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(o) != 10 {
		t.Errorf("error: not matching limit, expecting %d got %d\n", 10, len(o))
	}

	//TODO: this one fails, need to fix it, second limit in too high
	//I0824 11:46:39.460953   24091 networkObject.go:78] Upperbound: 110, Offset: 0, limit: 100
	//I0824 11:46:39.587332   24091 networkObject.go:78] Upperbound: 110, Offset: 100, limit: 100

	// o, err = asa.GetNetworkObjects(110)
	// if err != nil {
	// 	t.Errorf("error: %s\n", err)
	// 	return
	// }

	// if len(o) != 110 {
	// 	t.Errorf("error: not matching limit, expecting %d got %d\n", 110, len(o))
	// }

	for i := range all {
		err = asa.DeleteNetworkObject(all[i])
		if err != nil {
			t.Errorf("error: %s\n", err)
		}
	}
}

func TestNetworkObjectCreate(t *testing.T) {
	var err error

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	n := new(goasa.NetworkObject)
	n.Host.Kind = "IPv4Address"
	n.Host.Value = "1.1.1.1"
	n.Name = "1.1.1.1"
	n.ObjectID = "1.1.1.1"

	err = asa.CreateNetworkObject(n, goasa.DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	o, err := asa.GetNetworkObjectByID(n.ObjectID)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	if o.Name != n.Name {
		t.Errorf("error: no matching name, expecting %s got %s\n", n.Name, o.Name)
	}

	if o.ObjectID != n.ObjectID {
		t.Errorf("error: no matching objectid, expecting %s got %s\n", n.ObjectID, o.ObjectID)
	}

	err = asa.DeleteNetworkObject(n)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}
}

func TestNetworkObjectDuplicate(t *testing.T) {
	var err error

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	n := new(goasa.NetworkObject)
	n.Host.Kind = "IPv4Address"
	n.Host.Value = "1.1.1.1"
	n.Name = "1.1.1.1"
	n.ObjectID = "1.1.1.1"

	err = asa.CreateNetworkObject(n, goasa.DuplicateActionError)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	err = asa.CreateNetworkObject(n, goasa.DuplicateActionError)
	if err == nil {
		t.Errorf("error: we should have had an error here....\n")
	}

	err = asa.DeleteNetworkObject(n)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}
}

func TestCreateNetworkObjectFromIPs(t *testing.T) {
	var err error

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	ips1 := []string{
		"1.2.3.4",
		"5.6.7.8",
		"9.10.11.12",
	}

	ips2 := []string{
		"1.2.3.4",
		"5.6.7.8",
		"9.10.11.12",
		"13.14.15.16",
	}

	ns, err := asa.CreateNetworkObjectsFromIPs(ips1)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(ns) != 3 {
		t.Errorf("we should have at least 3 members, have %d\n", len(ns))
	}

	ns2, err := asa.CreateNetworkObjectsFromIPs(ips2)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(ns2) != 4 {
		t.Errorf("we should have at least 4 members, have %d\n", len(ns2))
	}

	found := false
	for i := range ns2 {
		if ns2[i].Host.Value == "13.14.15.16" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("didn't find 13.14.15.16 in array\n")
	}

	for i := range ns2 {
		err = asa.DeleteNetworkObject(ns2[i])
		if err != nil {
			t.Errorf("error: %s\n", err)
		}
	}

}
