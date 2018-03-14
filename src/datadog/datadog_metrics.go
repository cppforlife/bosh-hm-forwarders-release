package main

import (
	"fmt"
	"strconv"

	dgclient "datadog/client"
)

type DatadogMetrics struct {
	client dgclient.Client
	event  Event
}

func (a DatadogMetrics) Send() error {
	metrics := []dgclient.Metric{}

	tags := a.tags()

	var parseErr error

	for _, metric := range a.event.Metrics {
		metricValFloat, err := strconv.ParseFloat(metric.Value, 64)
		if err != nil {
			parseErr = err
			continue
		}

		// todo metric.Tags??

		metrics = append(metrics, dgclient.Metric{
			Metric: "bosh.healthmonitor." + metric.Name,
			Points: []dgclient.MetricDataPoint{{float64(metric.Timestamp), metricValFloat}},
			Tags:   tags,
		})
	}

	var postErr error

	if len(metrics) > 0 {
		postErr = a.client.PostMetrics(metrics)
	}

	if postErr != nil {
		return fmt.Errorf("posting metrics: %s", postErr)
	}
	if parseErr != nil {
		return fmt.Errorf("parsing metric value: %s", parseErr)
	}

	return nil
}

func (a DatadogMetrics) tags() []string {
	tags := []string{
		"job:" + a.event.Job,
		"index:" + a.event.Index,
		"id:" + a.event.InstanceID,
		"deployment:" + a.event.Deployment,
		"agent:" + a.event.AgentID,
	}
	for _, team := range a.event.Teams {
		tags = append(tags, "team:"+team)
	}
	return tags
}
