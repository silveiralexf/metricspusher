/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cmdFlags = Flags{}
	rootCmd  = &cobra.Command{
		Use:   "metricspusher",
		Short: "Generic Prometheus Pushgateway client for ephemeral jobs and applications",
	}
	logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.InfoLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}
)

type Flags struct {
	URL             string
	Service         string
	Environment     string
	MetricName      string
	JobStage        string
	JobResult       string
	JobType         string
	StartForSeconds int
	Start           bool
	Stop            bool
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
