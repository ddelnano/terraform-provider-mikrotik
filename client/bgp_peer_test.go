package client

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

// required BGP Peer fields
var remoteAs int = 65533
var remoteAddress string = "172.21.16.0"

var peerTTL string = "default"
var addressFamilies string = "ip"
var defaultOriginate string = "never"
var holdTime string = "3m"
var nextHopChoice string = "default"

func TestAddBgpPeerAndDeleteBgpPeer(t *testing.T) {
	SkipLegacyBgpIfUnsupported(t)
	c := NewClient(GetConfigFromEnv())

	instanceName := "peer-test"
	bgpPeerName := "test-peer"

	_, err := c.AddBgpInstance(&BgpInstance{Name: instanceName, As: 65530, RouterID: "172.16.0.254"})
	if err != nil {
		t.Fatalf("unable to create BGP instance used for testing: %v", err)
	}
	defer func(c *Mikrotik, name string) {
		_ = c.DeleteBgpInstance(name)
	}(c, instanceName)

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
	SkipLegacyBgpIfUnsupported(t)
	c := NewClient(GetConfigFromEnv())

	instanceName := "peer-update-test"
	bgpPeerName := "test-peer-update"

	_, err := c.AddBgpInstance(&BgpInstance{Name: instanceName, As: 65530, RouterID: "172.16.1.254"})
	if err != nil {
		t.Fatalf("unable to create BGP instance used for testing: %v", err)
	}
	defer func(c *Mikrotik, name string) {
		_ = c.DeleteBgpInstance(name)
	}(c, instanceName)

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
	SkipLegacyBgpIfUnsupported(t)
	c := NewClient(GetConfigFromEnv())

	name := "bgp peer does not exist"
	_, err := c.FindBgpPeer(name)

	require.Truef(t, IsNotFoundError(err),
		"Expecting to receive NotFound error for bgp peer %q.", name)
}
