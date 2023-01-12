package client

import (
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFirewallFilter_customChain(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	rule := &FirewallFilterRule{
		Chain:           "mychain",
		Comment:         "Test rule",
		DestPort:        "1001",
		ConnectionState: types.MikrotikList{"new"},
		Protocol:        "tcp",
	}

	createdRule, err := c.AddFirewallFilterRule(rule)
	require.NoError(t, err)

	defer func(id string) {
		assert.NoError(t, c.DeleteFirewallFilterRule(id))
	}(createdRule.Id)

	rule.Id = createdRule.Id

	foundRule, err := c.FindFirewallFilterRule(createdRule.Id)
	require.NoError(t, err)
	assert.Equal(t, rule, foundRule)
}

func TestFirewallFilter_builtinChain(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	rule := &FirewallFilterRule{
		Chain:           "filter",
		Comment:         "Test rule for builtin chain",
		DestPort:        "1001",
		ConnectionState: types.MikrotikList{"established", "related"},
		Protocol:        "tcp",
	}

	createdRule, err := c.AddFirewallFilterRule(rule)
	require.NoError(t, err)

	defer func(id string) {
		assert.NoError(t, c.DeleteFirewallFilterRule(id))
	}(createdRule.Id)

	rule.Id = createdRule.Id

	foundRule, err := c.FindFirewallFilterRule(rule.Id)
	require.NoError(t, err)
	assert.Equal(t, rule, foundRule)

	rule.Protocol = "udp"
	rule.Comment = "Updated protocol"
	_, err = c.UpdateFirewallFilterRule(rule)
	require.NoError(t, err)

	foundRule, err = c.FindFirewallFilterRule(rule.Id)
	require.NoError(t, err)
	assert.Equal(t, rule, foundRule)
}
