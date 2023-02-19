package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

// CompletionTimestamp records the timestamp of the last successful execution and result of a given job
type CompletionTimestamp struct {
	Name        string
	Description string
	TargetURL   string
	Labels      map[string]string
	Gauge       prometheus.Gauge
}

func (ct *CompletionTimestamp) NewGauge(metricLabels map[string]string) error {
	err := ct.setupGauge(metricLabels)
	if err != nil {
		return err
	}

	ct.Gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        ct.Name,
		Help:        ct.Description,
		ConstLabels: metricLabels,
	},
	)
	return nil
}

func (ct *CompletionTimestamp) setupGauge(metricLabels map[string]string) error {
	ct.Name = "metrics_pusher_completion_timestamp_seconds"
	ct.Description = "The timestamp of the last successful execution of a given job"

	if metricLabels == nil {
		return fmt.Errorf("failed on parsing labels provided")
	}

	err := validateMetricLabels(metricLabels)
	if err != nil {
		return err
	}

	ct.Labels = metricLabels
	return nil
}

func (ct *CompletionTimestamp) StartGaugeForSeconds(jobName, targetURL string, interval int) error {
	err := ct.StartGauge(jobName, targetURL)
	if ct.Gauge == nil {
		return err
	}

	time.Sleep(time.Duration(interval) * time.Second)

	err = ct.StopGauge(jobName, targetURL)
	if ct.Gauge == nil {
		return err
	}
	return nil
}

func (ct *CompletionTimestamp) StartGauge(jobName, targetURL string) error {
	if ct.Gauge == nil {
		return fmt.Errorf("gauge not initialized")
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(30)*time.Second)
	defer cancel()

	ct.Gauge.SetToCurrentTime()

	p := push.New(targetURL, jobName)

	err := p.Collector(ct.Gauge).PushContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to push '%v': [%v]", ct.Name, err)
	}

	return nil
}

func (ct *CompletionTimestamp) StopGauge(jobName, targetURL string) error {
	if ct.Gauge == nil {
		return fmt.Errorf("gauge not initialized")
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(30)*time.Second)
	defer cancel()

	p := push.New(targetURL, jobName)
	err := p.AddContext(ctx)
	if err != nil {
		return fmt.Errorf("failed adding context for stop action '%v': [%v]", ct.Name, err)
	}

	err = p.Delete()
	if err != nil {
		return fmt.Errorf("failed stop action on '%v': [%v]", ct.Name, err)
	}

	return nil
}
