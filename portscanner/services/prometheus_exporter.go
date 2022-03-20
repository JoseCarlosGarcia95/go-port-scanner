package services

import (
	"net/http"
	"strconv"
	"time"
)

// IsPrometheusExporter checks if a node exporter is running on a given port.
func IsPrometheusExporter(host string, port uint32) bool {
	client := http.Client{
		Timeout: time.Second,
	}

	req, err := client.Get("http://" + host + ":" + strconv.Itoa(int(port)) + "/metrics")

	if err != nil {
		return false
	}

	// if the request is successful, the node exporter is running
	if req.StatusCode == 200 || req.StatusCode == 401 {
		return true
	}

	return false
}
