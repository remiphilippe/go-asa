package goasa_test

import (
	"fmt"
	"testing"

	"github.com/golang/glog"
	goasa "github.com/remiphilippe/go-asa"
)

func TestNetworkObjectGroupCreate(t *testing.T) {
	var n *goasa.NetworkObject
	var g1, g2 *goasa.NetworkObjectGroup
	var all1, all2 []*goasa.NetworkObject
	var o []*goasa.NetworkObjectGroup

	limit := 2

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	for i := 1; i < limit+1; i++ {
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

		all1 = append(all1, n)
	}

	for i := 1; i < limit+1; i++ {
		ip := fmt.Sprintf("2.2.2.%d", i)
		n = new(goasa.NetworkObject)
		n.Host.Kind = "IPv4Address"
		n.Host.Value = ip
		n.Name = ip
		n.ObjectID = ip

		err = asa.CreateNetworkObject(n, goasa.DuplicateActionDoNothing)
		if err != nil {
			t.Errorf("error: %s\n", err)
		}

		all2 = append(all2, n)
	}

	g1 = new(goasa.NetworkObjectGroup)
	g1.Name = "g1"
	g1.ObjectID = "g1"
	for i := range all1 {
		g1.Members = append(g1.Members, all1[i].Reference())
	}

	err = asa.CreateNetworkObjectGroup(g1, goasa.DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	g2 = new(goasa.NetworkObjectGroup)
	g2.Name = "g2"
	g2.ObjectID = "g2"
	for i := range all1 {
		g2.Members = append(g2.Members, all2[i].Reference())
	}

	err = asa.CreateNetworkObjectGroup(g2, goasa.DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	o, err = asa.GetAllNetworkObjectGroups()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(o) != 2 {
		t.Errorf("error: not matching limit, expecting %d got %d\n", 2, len(o))
	}
}

func TestNetworkObjectGroupGetAll(t *testing.T) {
	// // disabling this test for now as it's a bit long
	// if true {
	// 	return
	// }
	var err error

	var g []*goasa.NetworkObjectGroup

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	g, err = asa.GetAllNetworkObjectGroups()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(g) != 2 {
		t.Errorf("error: too many results, expecting 2 got %d\n", len(g))
	}
}

func TestNetworkObjectGroupDelete(t *testing.T) {
	var err error
	var g []*goasa.NetworkObjectGroup

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	g, err = asa.GetAllNetworkObjectGroups()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	for i := range g {
		err = asa.DeleteNetworkObjectGroup(g[i])
		if err != nil {
			t.Errorf("error: %s\n", err)
		}

		for _, m := range g[i].Members {
			err = asa.DeleteNetworkObject(m)
			if err != nil {
				t.Errorf("error: %s\n", err)
			}
		}
	}
}

func TestCreateNetworkObjectGroupFromIPs(t *testing.T) {
	var err error

	asa, err := initTest()
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	ips1 := []string{
		"1.2.3.4",
		"5.6.7.8",
		"9.10.11.12",
	}

	g, err := asa.CreateNetworkObjectGroupFromIPs("testAutoCreate", ips1, goasa.DuplicateActionError)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(g.Members) != 3 {
		t.Errorf("Not enough objects... %d\n", len(g.Members))
	}

	err = asa.DeleteNetworkObjectGroup(g)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	for i := range g.Members {
		err = asa.DeleteNetworkObject(g.Members[i])
		if err != nil {
			t.Errorf("error: %s\n", err)
		}
	}

}
