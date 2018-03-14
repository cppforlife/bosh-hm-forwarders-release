package main

import (
	dgclient "datadog/client"
)

type DatadogEvent struct {
	client dgclient.Client
	event  Event
}

func (a DatadogEvent) Send() error {
	err := a.client.PostEvent(dgclient.Event{
		Title: a.event.Title,
		Text:  a.event.Summary,
		Time:  a.event.CreatedAt,

		Priority:  a.priority(),
		AlertType: a.alertType(),

		Tags: a.tags(),
	})

	return err
}

func (a DatadogEvent) priority() string {
	switch SeverityToString(a.event.Severity) {
	case "alert", "critical", "error":
		return "normal"
	default:
		return "low"
	}
}

func (a DatadogEvent) alertType() string {
	switch SeverityToString(a.event.Severity) {
	case "alert", "critical", "error":
		return "error"
	default:
		return "warning"
	}
}

func (a DatadogEvent) tags() []string {
	return []string{
		"source:" + a.event.Source,
		"deployment:" + a.event.Deployment,
	}
}
