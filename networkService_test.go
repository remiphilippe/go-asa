package goasa_test

import (
	"testing"

	"github.com/remiphilippe/go-asa"
)

func TestTCPServiceCreate(t *testing.T) {
	var err error
	var n *goasa.NetworkService

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	n, err = asa.CreateTCPService("testTCP", "tcp/80", goasa.DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if n.Reference().Kind != "objectRef#TcpUdpServiceObj" {
		t.Errorf("error: reference should be objectRef#TcpUdpServiceObj not %s\n", n.Reference().Kind)
	}

	err = asa.DeleteNetworkService(n)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

}

func TestUDPServiceCreate(t *testing.T) {
	var err error
	var n *goasa.NetworkService

	asa, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	n, err = asa.CreateUDPService("testUDP", "udp/53", goasa.DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if n.Reference().Kind != "objectRef#TcpUdpServiceObj" {
		t.Errorf("error: reference should be objectRef#TcpUdpServiceObj not %s\n", n.Reference().Kind)
	}

	err = asa.DeleteNetworkService(n)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

}
