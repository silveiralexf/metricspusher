package metrics

import (
	"fmt"
	"net"
	"os"
)

func GetLocalAddress() (hostname, ip string, err error) {
	hostname, err = os.Hostname()
	if err != nil {
		return "", "", err
	}

	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", "", err
	}

	for _, addr := range addr {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return hostname, ipnet.IP.String(), nil
			}
		}
	}
	return "", "", fmt.Errorf("failed to discover local address")
}

func validateMetricLabels(metricLabels map[string]string) error {
	if metricLabels == nil {
		return fmt.Errorf("empty payload provided")
	}

	for k, v := range metricLabels {
		if k == "" && v == "" {
			return fmt.Errorf("failed on parsing labels: ['%v'='%v']", k, v)
		}
	}

	if metricLabels["result"] != "SUCCESS" && metricLabels["result"] != "FAILURE" {
		return fmt.Errorf("result '%v' not supported, try 'SUCCESS' and 'FAILURE'", metricLabels["result"])
	}
	return nil
}
