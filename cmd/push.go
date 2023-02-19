package cmd

import (
	"fmt"

	"github.com/silveiralexf/metricspusher/pkg/metrics"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

const (
	jobName = "metrics_pusher"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push metrics to a target Prometheus Pushgateway",
	Run: func(cmd *cobra.Command, args []string) {
		err := doPush()
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func doPush() error {
	metricLabels, err := setMetricLabels()
	if err != nil {
		return err
	}

	ct := metrics.CompletionTimestamp{}
	err = ct.NewGauge(metricLabels)
	if err != nil {
		return err
	}

	if cmdFlags.StartForSeconds != -1 {
		if cmdFlags.StartForSeconds < 5 {
			msg := fmt.Errorf("insufficient interval of %vs provided. minimum of 5 seconds is required", cmdFlags.StartForSeconds)
			return msg
		}

		err = ct.StartGaugeForSeconds(jobName, cmdFlags.URL, cmdFlags.StartForSeconds)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("pushed %v for service %v", cmdFlags.MetricName, cmdFlags.Service)
		logger.Info(msg)
		return nil
	}

	if cmdFlags.Start {
		err = ct.StartGauge(jobName, cmdFlags.URL)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("pushed %v for service %v", cmdFlags.MetricName, cmdFlags.Service)
		logger.Info(msg)
		return nil
	}

	if cmdFlags.Stop {
		err = ct.StopGauge(jobName, cmdFlags.URL)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("removed %v for service %v", cmdFlags.MetricName, cmdFlags.Service)
		logger.Info(msg)
		return nil
	}

	return fmt.Errorf("no valid action informed. '--help' for usage")
}

func setMetricLabels() (map[string]string, error) {
	hostname, ipAddress, err := metrics.GetLocalAddress()
	if err != nil {
		return nil, err
	}

	switch cmdFlags.MetricName {
	case "gaugetimestamp":
		metricLabels := map[string]string{
			"service":     cmdFlags.Service,
			"type":        cmdFlags.JobType,
			"stage":       cmdFlags.JobStage,
			"result":      cmdFlags.JobResult,
			"environment": cmdFlags.Environment,
			"source":      hostname,
			"instance":    ipAddress,
		}
		return removeEmpty(metricLabels), nil
	}

	return nil, fmt.Errorf("'%v' is not recognized. '--help' for usage", cmdFlags.MetricName)
}

func removeEmpty(m map[string]string) map[string]string {
	result := map[string]string{}
	for k, v := range m {
		if v != "" {
			result[k] = v
		}
	}
	return result
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.PersistentFlags().StringVar(
		&cmdFlags.URL,
		"url",
		"http://localhost:9091",
		"Prometheus Pushgateway target URL",
	)

	pushCmd.PersistentFlags().StringVar(
		&cmdFlags.Service,
		"service",
		"cicd_pipeline",
		"Adds 'service' label to help better identifying the job",
	)

	pushCmd.PersistentFlags().StringVar(
		&cmdFlags.Environment,
		"env",
		"dev",
		"Target environment to which the metric is related",
	)

	pushCmd.PersistentFlags().StringVar(
		&cmdFlags.MetricName,
		"metric",
		"gaugetimestamp",
		"Select the name of metric type to be pushed, use 'list' command for supported commands",
	)

	pushCmd.PersistentFlags().StringVar(
		&cmdFlags.JobType,
		"type",
		"",
		"Adds 'type' label for better identifying the service (optional)",
	)

	pushCmd.PersistentFlags().StringVar(
		&cmdFlags.JobStage,
		"stage",
		"",
		"Adds 'stage' label for helping identifying stage of long running jobs (optional, use with caution to avoid high cardinality)",
	)

	pushCmd.PersistentFlags().StringVar(
		&cmdFlags.JobResult,
		"result",
		"SUCCESS",
		"Adds 'result' label for helping identifying the status of the job or its stage (supported: 'SUCCESS' or 'FAILURE')",
	)

	pushCmd.PersistentFlags().IntVar(
		&cmdFlags.StartForSeconds,
		"for",
		-1,
		"Signals the startup of a job and trigger its removal after enoght time in seconds to be scrapped",
	)

	pushCmd.PersistentFlags().BoolVar(
		&cmdFlags.Start,
		"start",
		false,
		"Signals the startup of a job (required)",
	)

	pushCmd.PersistentFlags().BoolVar(
		&cmdFlags.Stop,
		"stop",
		false,
		"Signals the end of a job (required)",
	)
}
