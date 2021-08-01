package mikrotik

import (
	"context"
	"errors"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceScheduler() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSchedulerCreate,
		ReadContext:   resourceSchedulerRead,
		UpdateContext: resourceSchedulerUpdate,
		DeleteContext: resourceSchedulerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"on_event": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceSchedulerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sched := prepareScheduler(d)

	c := m.(client.Mikrotik)

	scheduler, err := c.CreateScheduler(sched)
	if err != nil {
		return diag.FromErr(err)
	}

	err = schedulerToData(scheduler, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSchedulerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Mikrotik)

	scheduler, err := c.FindScheduler(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = schedulerToData(scheduler, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSchedulerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Mikrotik)

	sched := prepareScheduler(d)
	sched.Id = d.Id()

	scheduler, err := c.UpdateScheduler(sched)
	if err != nil {
		return diag.FromErr(err)
	}

	err = schedulerToData(scheduler, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSchedulerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Id()

	c := m.(client.Mikrotik)

	err := c.DeleteScheduler(name)

	if err != nil {
		return diag.FromErr(err)
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

func prepareScheduler(d *schema.ResourceData) *client.Scheduler {
	scheduler := new(client.Scheduler)

	scheduler.Name = d.Get("name").(string)
	scheduler.OnEvent = d.Get("on_event").(string)
	scheduler.StartDate = d.Get("start_date").(string)
	scheduler.StartTime = d.Get("start_time").(string)
	scheduler.Interval = d.Get("interval").(int)

	return scheduler
}
