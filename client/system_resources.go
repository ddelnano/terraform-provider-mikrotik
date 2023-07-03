package client

import (
	"log"

	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
)

type SystemResources struct {
	Uptime  types.MikrotikDuration `mikrotik:"uptime"`
	Version string                 `mikrotik:"version"`
}

func (d *SystemResources) ActionToCommand(action Action) string {
	return map[Action]string{
		Find: "/system/resource/print",
	}[action]
}

func (client Mikrotik) GetSystemResources() (*SystemResources, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}
	sysResources := &SystemResources{}
	cmd := Marshal(sysResources.ActionToCommand(Find), sysResources)

	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	err = Unmarshal(*r, sysResources)
	return sysResources, err
}
