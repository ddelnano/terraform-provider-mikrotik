package client

import (
	"reflect"
	"testing"
)

var bgpName string = "test-bgp"
var as int = 65533
var updatedAs int = 65534
var clientToClientReflection bool = true
var clusterID string = "172.21.16.1"
var noClusterID string = ""
var bgpComment string = "test-comment"
var confederation int = 8
var updatedConfederation int = 5
var confederationPeers string = ""
var disabled bool = false
var ignoreAsPathLen bool = false
var outFilter string = ""
var redistributeConnected bool = false
var redistributeOspf bool = false
var redistributeOtherBgp bool = false
var redistributeRip bool = false
var redistributeStatic bool = false
var routerID string = "172.21.16.2"
var routingTable string = ""

func TestAddBgpInstanceAndDeleteBgpInstance(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	expectedBgpInstance := &BgpInstance{
		Name:                     bgpName,
		As:                       as,
		ClientToClientReflection: clientToClientReflection,
		IgnoreAsPathLen:          ignoreAsPathLen,
		OutFilter:                outFilter,
		RedistributeConnected:    redistributeConnected,
		RedistributeOspf:         redistributeOspf,
		RedistributeOtherBgp:     redistributeOtherBgp,
		RedistributeRip:          redistributeRip,
		RedistributeStatic:       redistributeStatic,
		RouterID:                 routerID,
		RoutingTable:             routingTable,
	}
	bgpInstance, err := c.AddBgpInstance(expectedBgpInstance)
	if err != nil {
		t.Fatalf("Error creating a bpg instance with: %v", err)
	}

	expectedBgpInstance.ID = bgpInstance.ID

	if !reflect.DeepEqual(bgpInstance, expectedBgpInstance) {
		t.Errorf("The bgp instance does not match what we expected. actual: %v expected: %v", bgpInstance, expectedBgpInstance)
	}

	err = c.DeleteBgpInstance(bgpInstance.Name)

	if err != nil {
		t.Errorf("Error deleting bgp instance with: %v", err)
	}
}

func TestAddAndUpdateBgpInstanceWithOptionalFieldsAndDeleteBgpInstance(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	expectedBgpInstance := &BgpInstance{
		Name:                     bgpName,
		As:                       as,
		ClientToClientReflection: clientToClientReflection,
		Comment:                  bgpComment,
		ConfederationPeers:       confederationPeers,
		Disabled:                 disabled,
		IgnoreAsPathLen:          ignoreAsPathLen,
		OutFilter:                outFilter,
		RedistributeConnected:    redistributeConnected,
		RedistributeOspf:         redistributeOspf,
		RedistributeOtherBgp:     redistributeOtherBgp,
		RedistributeRip:          redistributeRip,
		RedistributeStatic:       redistributeStatic,
		RouterID:                 routerID,
		RoutingTable:             routingTable,
		ClusterID:                clusterID,
		Confederation:            confederation,
	}
	bgpInstance, err := c.AddBgpInstance(expectedBgpInstance)
	if err != nil {
		t.Fatalf("Error creating a bpg instance with: %v", err)
	}

	expectedBgpInstance.ID = bgpInstance.ID

	if !reflect.DeepEqual(bgpInstance, expectedBgpInstance) {
		t.Errorf("The bgp instance does not match what we expected. actual: %v expected: %v", bgpInstance, expectedBgpInstance)
	}

	// update fields
	expectedBgpInstance.Confederation = updatedConfederation
	expectedBgpInstance.As = updatedAs

	bgpInstance, err = c.UpdateBgpInstance(expectedBgpInstance)

	if !reflect.DeepEqual(bgpInstance, expectedBgpInstance) {
		t.Errorf("The bgp instance does not match what we expected. actual: %v expected: %v", bgpInstance, expectedBgpInstance)
	}

	err = c.DeleteBgpInstance(bgpInstance.Name)

	if err != nil {
		t.Errorf("Error deleting bgp instance with: %v", err)
	}
}

func TestFindBgpInstance_onNonExistantBgpInstance(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	name := "bgp instance does not exist"
	_, err := c.FindBgpInstance(name)

	if _, ok := err.(*NotFound); !ok {
		t.Errorf("Expecting to receive NotFound error for bgp instance `%s`, instead error was nil.", name)
	}
}
