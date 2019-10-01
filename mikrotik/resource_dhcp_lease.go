package mikrotik

import (
	"fmt"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLease() *schema.Resource {
	return &schema.Resource{
		Create: resourceLeaseCreate,
		Read:   resourceLeaseRead,
		Update: resourceLeaseUpdate,
		Delete: resourceLeaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"macaddress": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceLeaseCreate(d *schema.ResourceData, m interface{}) error {
	address := d.Get("address").(string)
	macaddress := d.Get("macaddress").(string)

	c := m.(client.Mikrotik)

	r, err := c.AddDhcpLease(address, macaddress)
	if err != nil {
		return err
	}

	// If API is successful we should only get a single sentence and list back like so
	// 2019/02/28 20:13:15 !done @ [{`ret` `*14`}]
	var id string
	for _, reply := range r.Re {
		for _, item := range reply.List {
			if item.Key == ".id" {
				id = item.Value
			}
		}
	}

	lease := &client.DhcpLease{
		Id:         id,
		Address:    address,
		MacAddress: macaddress,
	}
	leaseToData(lease, d)
	return nil
}

func resourceLeaseRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	lease, err := c.FindDhcpLease(d.Id())

	// TODO: Ignoring this error can cause all resources to think they
	// need to be created. We should more appropriately handle this. The
	// error where the DHCP lease is not found is not actually an error and
	// needs to be disambiguated from real failures
	if err != nil {
		d.SetId("")
		return nil
	}

	// FIXME
	if lease == nil {
		d.SetId("")
		return nil
	}

	leaseToData(lease, d)
	return nil
}

func resourceLeaseUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	macaddress := d.Get("macaddress").(string)
	address := d.Id()

	lease, err := c.FindDhcpLease(address)

	if err != nil {
		return err
	}

	fmt.Printf("[DEBUG] About to update dhcp lease with %v", lease)
	err = c.UpdateDhcpLease(lease.Id, address, macaddress)

	if err != nil {
		return err
	}

	// TODO: the c.UpdateDhcpLease call should return a
	// new DhcpLease instead of mutating the current one.
	lease.MacAddress = macaddress
	leaseToData(lease, d)
	return nil
}

func resourceLeaseDelete(d *schema.ResourceData, m interface{}) error {
	address := d.Id()

	c := m.(client.Mikrotik)

	lease, err := c.FindDhcpLease(address)

	if err != nil {
		return err
	}
	err = c.DeleteDhcpLease(lease.Id)

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func leaseToData(lease *client.DhcpLease, d *schema.ResourceData) error {
	d.SetId(lease.Address)
	d.Set("numerical_id", lease.Id)
	d.Set("address", lease.Address)
	d.Set("macaddress", lease.MacAddress)
	return nil
}
