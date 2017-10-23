package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func main() {
	cfgBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to read config: %s", err)
		os.Exit(1)
	}

	cfg := Config{}

	err = json.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to unmarshal config: %s", err)
		os.Exit(1)
	}

	cw := NewClientFromPath(cfg)
	dec := json.NewDecoder(os.Stdin)

	for {
		var ev Event

		if err := dec.Decode(&ev); err == io.EOF {
			fmt.Fprintf(os.Stderr, "Reached EOF")
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to unmarshal event: %s", err)
			break
		}

		if ev.Kind != EventKindHeartbeat || len(ev.InstanceID) == 0 {
			continue
		}

		err := submitEvent(cw, ev)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to put metric data: %s", err)
			continue
		}
	}
}

func submitEvent(cw *cloudwatch.CloudWatch, ev Event) error {
	input := &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("BOSH/HealthMonitor"),
	}

	for _, m := range ev.Metrics {
		ts := time.Unix(m.Timestamp, 0).UTC()

		val, err := strconv.ParseFloat(m.Value, 64)
		if err != nil {
			return err
		}

		metric := &cloudwatch.MetricDatum{
			Timestamp: &ts,

			MetricName: aws.String(m.Name),
			Value:      &val,

			Dimensions: []*cloudwatch.Dimension{
				// todo director
				&cloudwatch.Dimension{
					Name:  aws.String("deployment"),
					Value: aws.String(ev.Deployment),
				},
				&cloudwatch.Dimension{
					Name:  aws.String("instance"),
					Value: aws.String(ev.Job + "/" + ev.InstanceID),
				},
				&cloudwatch.Dimension{
					Name:  aws.String("instance_group"),
					Value: aws.String(ev.Job),
				},
				&cloudwatch.Dimension{
					Name:  aws.String("instance_id"),
					Value: aws.String(ev.InstanceID),
				},
				&cloudwatch.Dimension{
					Name:  aws.String("agent_id"),
					Value: aws.String(ev.AgentID),
				},
			},
		}

		input.MetricData = append(input.MetricData, metric)
	}

	_, err := cw.PutMetricData(input)

	return err
}
