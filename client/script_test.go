package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var scriptSource string = ":put testing"
var scriptName string = "testing"
var scriptOwner string = "owner"
var scriptPolicies []string = []string{
	"ftp",
	"reboot",
	"read",
	"write",
	"policy",
	"test",
	"password",
	"sniff",
	"sensitive",
	"romon",
}
var scriptDontReqPerms = true

func TestCreateScriptAndDeleteScript(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	expectedScript := &Script{
		Name:                   scriptName,
		Owner:                  scriptOwner,
		Source:                 scriptSource,
		Policy:                 scriptPolicies,
		DontRequirePermissions: scriptDontReqPerms,
	}
	script, err := NewClient(GetConfigFromEnv()).
		AddScript(&Script{
			Name:                   scriptName,
			Owner:                  scriptOwner,
			Source:                 scriptSource,
			Policy:                 scriptPolicies,
			DontRequirePermissions: scriptDontReqPerms,
		},
		)
	require.NoError(t, err)

	expectedScript.Id = script.Id

	defer c.DeleteScript(scriptName)
	require.Equal(t, expectedScript, script)

	err = c.DeleteScript(scriptName)

	require.NoError(t, err)
}

func TestFindScript_onNonExistantScript(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	name := "script-not-found"
	_, err := c.FindScript(name)

	if !IsNotFoundError(err) {
		t.Errorf("client should have received error indicating the following script `%s` was not found. Instead error was %v", name, err)
	}
}
