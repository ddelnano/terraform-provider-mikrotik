package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var scriptSource string = ":put testing"
var scriptName string = "testing"
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
	_, owner, _, _, _, _ := GetConfigFromEnv()

	expectedScript := &Script{
		Name:                   scriptName,
		Owner:                  owner,
		Source:                 scriptSource,
		Policy:                 scriptPolicies,
		DontRequirePermissions: scriptDontReqPerms,
	}
	script, err := NewClient(GetConfigFromEnv()).
		AddScript(&Script{
			Name:                   scriptName,
			Source:                 scriptSource,
			Policy:                 scriptPolicies,
			DontRequirePermissions: scriptDontReqPerms,
		},
		)
	require.NoError(t, err)

	expectedScript.Id = script.Id

	defer func() {
		if err := c.DeleteScript(scriptName); err != nil {
			assert.True(t, IsNotFoundError(err), "the only acceptable error is NotFound")
		}
	}()

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
