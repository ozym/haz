# haz-aws-messaging

Makes AWS resources for haz messaging.  Uses godep (`godep go build`).

Usage:

Export an env var for the AWS region you are using.  Set the AWS credentials profile as required.    
http://blogs.aws.amazon.com/security/post/Tx3D6U6WSFGOK2H/A-New-and-Standardized-Way-to-Manage-Credentials-in-the-AWS-SDKs

*NOTE: The identity in the given credential must has permission to create "IAM user", "SQS", and "SNS".*

```
export AWS_REGION=ap-southeast-2
export AWS_PROFILE=production
```

```
haz-aws-messaging --prefix=foo [--create-keys]
```

`prefix` is used as a prefix for all resources.  Use this to keep resources unique in an AWS account.  For dev your user name would be a good choice.
`--create-keys` is used to create keys for the tx and rx users.  A user can only have 2 sets of keys at any time.  After initial creation of keys use the IAM console for further key management.  Be sure to store the keys somewhere - they cannot be retrieved again.  Run as required to see the config.

Creates the following resources

* An SNS topic for sending haz messages to.
* A list of SQS queues that have raw message subscriptions to the haz topic.
* A dead letter queue that is used by each SQS with a redrive policy.
* A user with permission to send to the SNS topic.
* A user with permission to read from the SQS queues.

SQS queues can be added at anytime.  

Cloudwatch alarms and alerting topics should be created separately depending on current monitoring needs.  In general any client for a web application will be monitored via http probes.  Alerting clients will tend to need their SQS queues monitoring.
