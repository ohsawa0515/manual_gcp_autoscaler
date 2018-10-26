package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

var (
	mig, projectID, region, zone string
	min, max                     int64
)

// AutoscalerIface -
type AutoscalerIface interface {
	MergeNum(int64, int64)
	UpdateSize(context.Context) error
}

// RegionAutoscalerClient -
type RegionAutoscalerClient struct {
	Autoscaler *compute.Autoscaler
	Svc        *compute.Service
	project    string
	region     string
}

// ZoneAutoscalerClient -
type ZoneAutoscalerClient struct {
	Autoscaler *compute.Autoscaler
	Svc        *compute.Service
	project    string
	zone       string
}

// NewAutoscalerClient -
func NewAutoscalerClient(ctx context.Context, service *compute.Service, project, region, zone, mig string) (AutoscalerIface, error) {
	if len(region) > 0 {
		resp, err := service.RegionAutoscalers.Get(project, region, mig).Context(ctx).Do()
		if err != nil {
			return nil, err
		}
		return &RegionAutoscalerClient{
			Autoscaler: resp,
			Svc:        service,
			project:    project,
			region:     region,
		}, nil
	} else if len(zone) > 0 {
		resp, err := service.Autoscalers.Get(project, zone, mig).Context(ctx).Do()
		if err != nil {
			return nil, err
		}
		return &ZoneAutoscalerClient{
			Autoscaler: resp,
			Svc:        service,
			project:    project,
			zone:       zone,
		}, nil
	}
	return nil, errors.New("Please specify either option `region` or `zone`")
}

// MergeNum sets sizes of min and max if the number of min and max is not the default value(0).
func (c *RegionAutoscalerClient) MergeNum(max, min int64) {
	// default value
	if min != 0 && min > 0 {
		c.Autoscaler.AutoscalingPolicy.MinNumReplicas = min
	}
	if max != 0 {
		c.Autoscaler.AutoscalingPolicy.MaxNumReplicas = max
	}
}

// MergeNum sets sizes of min and max if the number of min and max is not the default value(0).
func (c *ZoneAutoscalerClient) MergeNum(max, min int64) {
	// default value
	if min != 0 && min > 0 {
		c.Autoscaler.AutoscalingPolicy.MinNumReplicas = min
	}
	if max != 0 {
		c.Autoscaler.AutoscalingPolicy.MaxNumReplicas = max
	}
}

// UpdateSize updates the number of Max and Min in the managed instance group.
func (c *RegionAutoscalerClient) UpdateSize(ctx context.Context) error {
	if _, err := c.Svc.RegionAutoscalers.Update(c.project, c.region, c.Autoscaler).Context(ctx).Do(); err != nil {
		return err
	}
	return nil
}

// UpdateSize updates the number of Max and Min in the managed instance group.
func (c *ZoneAutoscalerClient) UpdateSize(ctx context.Context) error {
	if _, err := c.Svc.Autoscalers.Update(c.project, c.zone, c.Autoscaler).Context(ctx).Do(); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.StringVar(&projectID, "project", "", "Project ID.")
	flag.StringVar(&region, "region", "", "Region of the managed regional instance group. e.g. asia-northeast1")
	flag.StringVar(&zone, "zone", "", "Zone of the managed instance group. e.g. asia-northeast1-a")
	flag.StringVar(&mig, "mig", "", "Name of managed instance group.")
	flag.Int64Var(&min, "min", 0, "The minimum number of replicas that the autoscaler can scale down to. This cannot be less than 0.")
	flag.Int64Var(&max, "max", 0, "The maximum number of instances that the autoscaler can scale up to.")
	flag.Parse()

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		log.Fatal(err)
	}
	computeService, err := compute.New(client)
	if err != nil {
		log.Fatal(err)
	}

	autoscaler, err := NewAutoscalerClient(ctx, computeService, projectID, region, zone, mig)
	if err != nil {
		log.Fatal(err)
	}
	autoscaler.MergeNum(max, min)

	if err := autoscaler.UpdateSize(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Update is completed.")
	os.Exit(0)
}
