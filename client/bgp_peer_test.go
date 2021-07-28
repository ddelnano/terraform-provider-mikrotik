package client

import (
	"reflect"
	"testing"
)

// required BGP Peer fields
var bgpPeerName string = "test-peer"
var remoteAs int = 65533
var remoteAddress string = "172.21.16.0"
var instanceName string = "test"

var peerTTL string = "default"
var addressFamilies string = "ip"
var defaultOriginate string = "never"
var holdTime string = "3m"
var nextHopChoice string = "default"

func TestAddBgpPeerAndDeleteBgpPeer(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	expectedBgpPeer := &BgpPeer{
		Name:             bgpPeerName,
		Instance:         instanceName,
		RemoteAs:         remoteAs,
		RemoteAddress:    remoteAddress,
		TTL:              peerTTL,
		AddressFamilies:  addressFamilies,
		DefaultOriginate: defaultOriginate,
		HoldTime:         holdTime,
		NexthopChoice:    nextHopChoice,
	}
	bgpPeer, err := c.AddBgpPeer(expectedBgpPeer)
	if err != nil {
		t.Fatalf("Error creating a bpg peer with: %v", err)
	}

	expectedBgpPeer.ID = bgpPeer.ID

	if !reflect.DeepEqual(bgpPeer, expectedBgpPeer) {
		t.Errorf("The bgp peer does not match what we expected. actual: %v expected: %v", bgpPeer, expectedBgpPeer)
	}

	err = c.DeleteBgpPeer(bgpPeer.Name)

	if err != nil {
		t.Errorf("Error deleting bgp peer with: %v", err)
	}
}

func TestAddAndUpdateBgpPeerWithOptionalFieldsAndDeleteBgpPeer(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	expectedBgpPeer := &BgpPeer{
		Name:             bgpPeerName,
		Instance:         instanceName,
		RemoteAs:         remoteAs,
		RemoteAddress:    remoteAddress,
		TTL:              peerTTL,
		AddressFamilies:  addressFamilies,
		DefaultOriginate: defaultOriginate,
		HoldTime:         holdTime,
		NexthopChoice:    nextHopChoice,
	}
	bgpPeer, err := c.AddBgpPeer(expectedBgpPeer)
	if err != nil {
		t.Fatalf("Error creating a bpg peer with: %v", err)
	}

	expectedBgpPeer.ID = bgpPeer.ID

	if !reflect.DeepEqual(bgpPeer, expectedBgpPeer) {
		t.Errorf("The bgp peer does not match what we expected. actual: %v expected: %v", bgpPeer, expectedBgpPeer)
	}

	// update fields
	expectedBgpPeer.UpdateSource = "172.21.16.1"
	expectedBgpPeer.TCPMd5Key = "test-key-name"
	expectedBgpPeer.RemotePort = 65000
	expectedBgpPeer.OutFilter = "test out filter"
	expectedBgpPeer.MaxPrefixRestartTime = "infinity"
	expectedBgpPeer.MaxPrefixLimit = 20
	expectedBgpPeer.KeepAliveTime = "30m"
	expectedBgpPeer.InFilter = "test in filter"
	expectedBgpPeer.Comment = "test comment"
	expectedBgpPeer.CiscoVplsNlriLenFmt = "bits"
	expectedBgpPeer.AllowAsIn = 0

	bgpPeer, err = c.UpdateBgpPeer(expectedBgpPeer)

	if !reflect.DeepEqual(bgpPeer, expectedBgpPeer) {
		t.Errorf("The bgp peer does not match what we expected. actual: %v expected: %v", bgpPeer, expectedBgpPeer)
	}

	err = c.DeleteBgpPeer(bgpPeer.Name)

	if err != nil {
		t.Errorf("Error deleting bgp peer with: %v", err)
	}
}

func TestFindBgpPeer_onNonExistantBgpPeer(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	name := "bgp peer does not exist"
	_, err := c.FindBgpPeer(name)

	if _, ok := err.(*NotFound); !ok {
		t.Errorf("Expecting to receive NotFound error for bgp peer `%s`, instead error was nil.", name)
	}
}
