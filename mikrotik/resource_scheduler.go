package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceScheduler() *schema.Resource {
	return &schema.Resource{
		Description: "Creates a Mikrotik scheduler.",

		CreateContext: resourceSchedulerCreate,
		ReadContext:   resourceSchedulerRead,
		UpdateContext: resourceSchedulerUpdate,
		DeleteContext: resourceSchedulerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the task.",
			},
			"on_event": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the script to execute. It must exist `/system script`.",
			},
			"start_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date of the first script execution.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time of the first script execution.",
			},
			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Interval between two script executions, if time interval is set to zero, the script is only executed at its start time, otherwise it is executed repeatedly at the time interval is specified.",
			},
		},
	}
}

func resourceSchedulerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sched := prepareScheduler(d)

	c := m.(*client.Mikrotik)

	scheduler, err := c.CreateScheduler(sched)
	if err != nil {
		return diag.FromErr(err)
	}

	return schedulerToData(scheduler, d)
}

func resourceSchedulerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	scheduler, err := c.FindScheduler(d.Id())
	if _, ok := err.(*client.NotFound); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return schedulerToData(scheduler, d)
}

func resourceSchedulerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	sched := prepareScheduler(d)
	sched.Id = d.Id()

	scheduler, err := c.UpdateScheduler(sched)
	if err != nil {
		return diag.FromErr(err)
	}

	return schedulerToData(scheduler, d)
}

func resourceSchedulerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Id()

	c := m.(*client.Mikrotik)

	err := c.DeleteScheduler(name)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func schedulerToData(s *client.Scheduler, d *schema.ResourceData) diag.Diagnostics {
	if s == nil {
		return diag.Errorf("scheduler was not found")
	}

	values := map[string]interface{}{
		"name":       s.Name,
		"on_event":   s.OnEvent,
		"start_time": s.StartTime,
		"start_date": s.StartDate,
		"interval":   s.Interval,
	}

	d.SetId(s.Name)

	var diags diag.Diagnostics

	for key, value := range values {
		if err := d.Set(key, value); err != nil {
			diags = append(diags, diag.Errorf("failed to set %s: %v", key, err)...)
		}
	}

	return diags
}

func prepareScheduler(d *schema.ResourceData) *client.Scheduler {
	return &client.Scheduler{
		Name:      d.Get("name").(string),
		OnEvent:   d.Get("on_event").(string),
		StartDate: d.Get("start_date").(string),
		StartTime: d.Get("start_time").(string),
		Interval:  types.MikrotikDuration(d.Get("interval").(int)),
	}
}
