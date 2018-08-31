package goasa_test

import (
	"fmt"
	"testing"

	goasa "github.com/remiphilippe/go-asa"
)

func TestNetworkServiceGroupCreate(t *testing.T) {
	var n *goasa.NetworkService
	var g1, g2, g3 *goasa.NetworkServiceGroup
	var all1, all2 []*goasa.NetworkService
	var o []*goasa.NetworkServiceGroup

	limit := 2

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	for i := 1; i < limit+1; i++ {
		name := fmt.Sprintf("tcp-%d", i+100)
		port := fmt.Sprintf("tcp/%d", i+100)
		n, err = asa.CreateTCPService(name, port, goasa.DuplicateActionDoNothing)
		if err != nil {
			t.Errorf("error: %s\n", err)
		}

		all1 = append(all1, n)
	}

	for i := 1; i < limit+1; i++ {
		name := fmt.Sprintf("udp-%d", i+200)
		port := fmt.Sprintf("udp/%d", i+200)
		n, err = asa.CreateUDPService(name, port, goasa.DuplicateActionDoNothing)
		if err != nil {
			t.Errorf("error: %s\n", err)
		}

		all2 = append(all2, n)
	}

	g1 = new(goasa.NetworkServiceGroup)
	g1.Name = "g1"
	g1.Kind = "object#NetworkServiceGroup"
	for i := range all1 {
		g1.Members = append(g1.Members, all1[i].Reference())
	}

	err = asa.CreateNetworkServiceGroup(g1, goasa.DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	g2 = new(goasa.NetworkServiceGroup)
	g2.Name = "g2"
	g2.Kind = "object#NetworkServiceGroup"
	for i := range all2 {
		g2.Members = append(g2.Members, all2[i].Reference())
	}

	err = asa.CreateNetworkServiceGroup(g2, goasa.DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	g3 = new(goasa.NetworkServiceGroup)
	g3.Name = "g3"
	g3.Kind = "object#NetworkServiceGroup"
	for i := range all1 {
		g3.Members = append(g3.Members, all1[i].Reference())
	}
	for i := range all2 {
		g3.Members = append(g3.Members, all2[i].Reference())
	}

	err = asa.CreateNetworkServiceGroup(g3, goasa.DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	o, err = asa.GetAllNetworkServiceGroups()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(o) != 3 {
		t.Errorf("error: not matching limit, expecting %d got %d\n", 3, len(o))
	}

	err = asa.DeleteNetworkServiceGroup(g1)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	err = asa.DeleteNetworkServiceGroup(g2)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	err = asa.DeleteNetworkServiceGroup(g3)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	for i := range all1 {
		err = asa.DeleteNetworkService(all1[i])
		if err != nil {
			t.Errorf("error: %s\n", err)
		}
	}

	for i := range all2 {
		err = asa.DeleteNetworkService(all2[i])
		if err != nil {
			t.Errorf("error: %s\n", err)
		}
	}

}

func TestNetworkServiceGroupCreateFromMap(t *testing.T) {
	var err error

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	m := make(map[string][]int)
	m["udp"] = []int{10, 11, 12}
	m["tcp"] = []int{10, 11, 12}

	g, err := asa.CreateNetworkServiceGroupFromMap("mapBasedSG", m, goasa.DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	err = asa.DeleteNetworkServiceGroup(g)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	for i := range g.Members {
		err = asa.DeleteNetworkService(g.Members[i])
		if err != nil {
			t.Errorf("error: %s\n", err)
		}
	}
}
