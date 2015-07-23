package main

import (
	"flag"
	"fmt"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/iam"
	"github.com/awslabs/aws-sdk-go/service/sns"
	"github.com/awslabs/aws-sdk-go/service/sqs"
	"log"
	"os"
	"time"
)

var queueNames = []string{
	`haz-web-consumer-api`, // for api.geonet.org.nz in AWS
	`haz-duty-consumer`,
	`haz-eqnews-consumer`,
	`haz-pim-consumer`,
	`haz-twitter-consumer`,
	`haz-ua-consumer`,
	`haz-db-consumer-api`,
}

var (
	sqsSvc = sqs.New(&aws.Config{})
	snsSvc = sns.New(&aws.Config{})
	iamSvc = iam.New(&aws.Config{})
)

var (
	prefix, rx, tx string
	createKeys     bool
)

func init() {
	if os.Getenv("AWS_REGION") == "" {
		log.Fatal("Please set env var AWS_REGION.  Use AWS_REGION=ap-southeast-2 for Sydney.")
	}

	flag.StringVar(&prefix, "prefix", "", "Required. Prefix for all the resources.")
	flag.BoolVar(&createKeys, "create-keys", false, "Create keys for the users.")

	flag.Parse()

	if prefix == "" {
		log.Fatal("please specify prefix.")
	}

	for i, _ := range queueNames {
		queueNames[i] = prefix + "-" + queueNames[i]
	}

	rx = prefix + "-haz-rx"
	tx = prefix + "-haz-tx"
}

func main() {
	rxUserArn, err := makeUser(rx)
	if err != nil {
		log.Panic(err)
	}

	txUserArn, err := makeUser(tx)
	if err != nil {
		log.Panic(err)
	}

	if createKeys {
		err = makeKeys(rx)
		if err != nil {
			log.Panic(err)
		}
		err = makeKeys(tx)
		if err != nil {
			log.Panic(err)
		}
	}

	hazArn, err := makeHazTopic(txUserArn)
	if err != nil {
		log.Panic(err)
	}

	dlqArn, err := makeDLQ()
	if err != nil {
		log.Panic(err)
	}

	for _, q := range queueNames {
		err := makeHazQueue(q, dlqArn, hazArn, rxUserArn)
		if err != nil {
			log.Panic(err)
		}
	}

	fmt.Printf("SNS Haz topic send user: %s\n", tx)
	fmt.Printf("SQS Haz queue receive user: %s\n", rx)

	fmt.Printf("\nHaz SNS topic arn: %s\n", hazArn)
	fmt.Println("\nHaz SNS topic subscriptions (please make sure all Haz sqs queues are subscribed):")
	// Assuming less than 100 queues/subscriptions
	sub, err := snsSvc.ListSubscriptionsByTopic(&sns.ListSubscriptionsByTopicInput{TopicARN: &hazArn})
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range sub.Subscriptions {
		a, err := snsSvc.GetSubscriptionAttributes(&sns.GetSubscriptionAttributesInput{SubscriptionARN: s.SubscriptionARN})
		if err != nil {
			log.Fatal(err)
		}
		am := *a.Attributes
		fmt.Printf("%s RawMessageDelivery=%s\n", *s.Endpoint, *am["RawMessageDelivery"])
	}

	fmt.Printf("\nSQS Haz message queues:\n")
	for _, q := range queueNames {
		fmt.Println(q)
	}
}

func makeKeys(name string) error {
	k, err := iamSvc.CreateAccessKey(&iam.CreateAccessKeyInput{UserName: &name})
	if err != nil {
		return err
	}
	fmt.Println("Access keys.  Note these now.  The secret can't be retrieved again.")
	fmt.Printf("User: %s ID: %s Secret: %s\n", name, *k.AccessKey.AccessKeyID, *k.AccessKey.SecretAccessKey)

	return nil
}

func makeUser(name string) (arn string, err error) {
	ru, err := iamSvc.GetUser(&iam.GetUserInput{UserName: &name})
	if err != nil { // most likely cause of error is the user not existing in which case try to create the user.
		_, err = iamSvc.CreateUser(&iam.CreateUserInput{UserName: &name})
		if err != nil {
			return
		}

		log.Println("Made user, waiting for it to propagate.")
		time.Sleep(time.Duration(30) * time.Second)

		ru, err = iamSvc.GetUser(&iam.GetUserInput{UserName: &name})
		if err != nil {
			return
		}
	}

	u := *ru.User
	arn = *u.ARN

	return
}

func makeHazQueue(name, dlqArn, hazTopicArn, rxUserArn string) error {
	l, err := sqsSvc.ListQueues(&sqs.ListQueuesInput{
		QueueNamePrefix: &name,
	})
	if err != nil {
		return err
	}

	switch len(l.QueueURLs) {
	case 1:
		return nil
	case 0:
		r, err := sqsSvc.CreateQueue(
			&sqs.CreateQueueInput{
				QueueName: &name,
				Attributes: &map[string]*string{
					`MessageRetentionPeriod`:        aws.String(`1209600`),
					`ReceiveMessageWaitTimeSeconds`: aws.String(`20`),
					`VisibilityTimeout`:             aws.String(`20`),
					`RedrivePolicy`:                 aws.String(`{"maxReceiveCount":"3", "deadLetterTargetArn":"` + dlqArn + `"}`),
					`Policy`: aws.String(`{
						"Version": "2012-10-17",
						"Id": "haz-sqs-policy",
						"Statement": [
						{
							"Sid": "allow-send-haz",
							"Effect": "Allow",
							"Principal": "*",
							"Action": ["sqs:SendMessage"],
							"Resource": "*",
							"Condition": {
								"StringEquals" : {
									"aws:SourceArn":"` + hazTopicArn + `"
								}
							}
							},
							{
								"Sid": "allow-read-haz",
								"Effect": "Allow",
								"Principal": {"AWS": "` + rxUserArn + `"},
								"Action": [
								"sqs:DeleteMessage",
								"sqs:ReceiveMessage",
								"sqs:GetQueueUrl"
								],
								"Resource": "*"
							}
							]
							}`),
				},
			})
		if err != nil {
			return err
		}

		// find queue ARN

		at, err := sqsSvc.GetQueueAttributes(
			&sqs.GetQueueAttributesInput{
				QueueURL: r.QueueURL,
				AttributeNames: []*string{
					aws.String(`QueueArn`),
				},
			})
		if err != nil {
			return err
		}

		atr := *at.Attributes
		arn := *atr[`QueueArn`]
		if arn == "" {
			return fmt.Errorf("Empty queue arn")

		}

		// Subscribe queue to haz topic

		sub, err := snsSvc.Subscribe(
			&sns.SubscribeInput{
				Endpoint: &arn,
				Protocol: aws.String(`sqs`),
				TopicARN: &hazTopicArn,
			})
		if err != nil {
			return err
		}

		// Make the subscription raw

		_, err = snsSvc.SetSubscriptionAttributes(
			&sns.SetSubscriptionAttributesInput{
				AttributeName:   aws.String(`RawMessageDelivery`),
				AttributeValue:  aws.String(`true`),
				SubscriptionARN: sub.SubscriptionARN,
			})
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("Found more than one queue with name prefix %s.  Can't recover.", name)
}

func makeDLQ() (string, error) {
	name := prefix + "-haz-dlq"

	l, err := sqsSvc.ListQueues(&sqs.ListQueuesInput{
		QueueNamePrefix: &name,
	})
	if err != nil {
		return "", err
	}

	switch len(l.QueueURLs) {
	case 1:
		// / lookup arn
		at, err := sqsSvc.GetQueueAttributes(
			&sqs.GetQueueAttributesInput{
				QueueURL: l.QueueURLs[0],
				AttributeNames: []*string{
					aws.String(`QueueArn`),
				},
			})
		if err != nil {
			return "", err
		}

		atr := *at.Attributes
		arn := *atr[`QueueArn`]
		if arn == "" {
			return "", fmt.Errorf("Empty queue arn")

		}
		return arn, nil
	case 0:
		var arn string
		r, err := sqsSvc.CreateQueue(
			&sqs.CreateQueueInput{
				QueueName: &name,
				Attributes: &map[string]*string{
					`MessageRetentionPeriod`: aws.String(`1209600`),
					`Policy`:                 aws.String(`{"Version": "2012-10-17", "Id": "dlq-default-policy"}`),
				},
			})
		if err != nil {
			return arn, err
		}

		at, err := sqsSvc.GetQueueAttributes(
			&sqs.GetQueueAttributesInput{
				QueueURL: r.QueueURL,
				AttributeNames: []*string{
					aws.String(`QueueArn`),
				},
			})
		if err != nil {
			return arn, err
		}

		atr := *at.Attributes
		arn = *atr[`QueueArn`]
		if arn == "" {
			err = fmt.Errorf("Empty queue arn")
			return arn, err
		}

		return arn, err
	}

	return "", fmt.Errorf("Found more than one queue with name prefix %s.  Can't recover.", name)
}

func makeHazTopic(txUserArn string) (arn string, err error) {
	resp, err := snsSvc.CreateTopic(
		&sns.CreateTopicInput{
			Name: aws.String(prefix + "-haz-topic"),
		})
	if err != nil {
		return
	}

	_, err = snsSvc.SetTopicAttributes(
		&sns.SetTopicAttributesInput{
			TopicARN:      resp.TopicARN,
			AttributeName: aws.String(`Policy`),
			AttributeValue: aws.String(`{
						"Version":"2012-10-17",
						"Id":"haz-topic-access",
						"Statement" :[
						{
							"Sid":"allow-send-haz",
							"Effect":"Allow",           
							"Principal" :{
								"AWS":"` + txUserArn + `"
								},
								"Action":["sns:GetTopicAttributes","sns:Publish"],
								"Resource":"*"
							}
							]
							}`),
		})
	arn = *resp.TopicARN
	return
}
