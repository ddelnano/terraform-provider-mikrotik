package client

import (
	"reflect"
	"testing"
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

	expectedScript := Script{
		Name:                   scriptName,
		Owner:                  scriptOwner,
		Source:                 scriptSource,
		Policy:                 scriptPolicies,
		DontRequirePermissions: scriptDontReqPerms,
	}
	script, err := NewClient(GetConfigFromEnv()).CreateScript(
		scriptName,
		scriptOwner,
		scriptSource,
		scriptPolicies,
		scriptDontReqPerms,
	)

	if err != nil {
		t.Errorf("Error creating a script with: %v", err)
	}

	expectedScript.Id = script.Id

	defer c.DeleteScript(scriptName)
	if !reflect.DeepEqual(script, expectedScript) {
		t.Errorf("The script does not match what we expected. actual: %v expected: %v", script, expectedScript)
	}

	err = c.DeleteScript(scriptName)

	if err != nil {
		t.Errorf("Error deleting a script with: %v", err)
	}
}

func TestFindScriptOnNonExistantScript(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	name := "script-not-found"
	script, err := c.FindScript(name)

	if err != nil {
		t.Errorf("Failed to find script `%s` with error: %v", name, err)
	}

	if script.Name != "" {
		t.Errorf("Script should have a blank name")
	}
}
