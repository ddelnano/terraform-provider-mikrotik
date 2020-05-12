package mikrotik

import (
	"errors"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceScheduler() *schema.Resource {
	return &schema.Resource{
		Create: resourceSchedulerCreate,
		Read:   resourceSchedulerRead,
		Update: resourceSchedulerUpdate,
		Delete: resourceSchedulerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"on_event": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"start_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"interval": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceSchedulerCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	onEvent := d.Get("on_event").(string)
	interval := d.Get("interval").(int)
	c := m.(client.Mikrotik)

	scheduler, err := c.CreateScheduler(
		name,
		onEvent,
		interval,
	)
	if err != nil {
		return err
	}

	schedulerToData(scheduler, d)
	return nil
}

func resourceSchedulerRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.Mikrotik)

	scheduler, err := c.FindScheduler(
		d.Id(),
	)

	if err != nil {
		return err
	}
	return schedulerToData(scheduler, d)
}

func resourceSchedulerUpdate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	onEvent := d.Get("on_event").(string)
	interval := d.Get("interval").(int)

	c := m.(client.Mikrotik)

	scheduler, err := c.UpdateScheduler(
		name,
		onEvent,
		interval,
	)

	if err != nil {
		return err
	}
	return schedulerToData(scheduler, d)
}

func resourceSchedulerDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Id()

	c := m.(client.Mikrotik)

	err := c.DeleteScheduler(name)

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func schedulerToData(s *client.Scheduler, d *schema.ResourceData) error {
	if s == nil {
		return errors.New("scheduler was not found")
	}
	d.SetId(s.Name)
	d.Set("name", s.Name)
	d.Set("on_event", s.OnEvent)
	d.Set("start_time", s.StartTime)
	d.Set("start_date", s.StartDate)
	d.Set("interval", s.Interval)
	return nil
}
