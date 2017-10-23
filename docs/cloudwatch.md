## bosh-hm-forwarder-cloudwatch job

Only forwards heartbeats.

Tested with following IAM policy:

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [ "cloudwatch:Put*" ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
```
