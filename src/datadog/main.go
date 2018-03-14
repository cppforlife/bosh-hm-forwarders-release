package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	dgclient "datadog/client"
)

func main() {
	cfgBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "[datadog] error: failed to read config: %s\n", err)
		os.Exit(1)
	}

	cfg := Config{}

	err = json.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[datadog] error: failed to unmarshal config: %s\n", err)
		os.Exit(1)
	}

	err = cfg.Validate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[datadog] error: failed to validate config: %s\n", err)
		os.Exit(1)
	}

	dec := json.NewDecoder(os.Stdin)
	client := dgclient.NewHTTPClient(cfg.APIKey)

	fmt.Fprintf(os.Stderr, "[datadog] starting\n")

	for {
		var ev Event

		if err := dec.Decode(&ev); err == io.EOF {
			fmt.Fprintf(os.Stderr, "[datadog] Reached EOF\n")
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "[datadog] error: failed to unmarshal event: %s\n", err)
			break
		}

		switch ev.Kind {
		case EventKindHeartbeat:
			if len(ev.InstanceID) != 0 {
				err := DatadogMetrics{client, ev}.Send()
				if err != nil {
					fmt.Fprintf(os.Stderr, "[datadog] error: failed to send heartbeat metrics: %s\n", err)
				}
			}

		case EventKindAlert:
			err := DatadogEvent{client, ev}.Send()
			if err != nil {
				fmt.Fprintf(os.Stderr, "[datadog] error: failed to send alert event: %s\n", err)
			}
		}
	}

	fmt.Fprintf(os.Stderr, "[datadog] exiting\n")
}
