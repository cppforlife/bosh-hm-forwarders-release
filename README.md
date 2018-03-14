# bosh-hm-forwarders-release

Release contains following BOSH HealthMonitor forwarders:

- [AWS CloudWatch](docs/cloudwatch.md)
- [DataDog](docs/datadog.md)

## Example usage

Standalone BOSH HM (on bosh-lite):

```
$ bosh -n -d bosh-hm deploy manifests/hm.yml -o manifests/datadog.yml \
  -v datadog_api_key=... \
  -l ~/workspace/deployments/vbox/creds.yml \
  -v director_address=10.254.50.4 \
  -v nats_address=10.254.50.4
```
