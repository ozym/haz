// sqs handles messages from AWS SQS via channels.
package sqs

import (
	"github.com/GeoNet/haz/msg"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"time"
)

var (
	MaxNumberOfMessages = 1
	VisibilityTimeout   = 600
	WaitTimeSeconds     = 20
)

var (
	retry     = time.Duration(30) * time.Second
	accessKey = os.Getenv("SQS_ACCESS_KEY")
	secretKey = os.Getenv("SQS_SECRET_KEY")
	queueName = os.Getenv("SQS_QUEUE_NAME")
	awsRegion = os.Getenv("AWS_REGION")
)

// InitRx handles receiving and deleting messages from AWS SQS via a pair of channels.
// Messages from SQS can be received from the read chan.  Receipt handles for messages that should be deleted from SQS
// can be sent to the write chan.
//
// Messages are read from the SQS queue s.MaxNumberOfMessages at a time.  They remain
// invisible or inflight on the SQS queue until either they are deleted by sending the ReceiptHandle to the write
// chan or they become visible again after s.VisibilityTimeout seconds when they are redelivered.
// Messages are delivered by SQS at least once.  Applications should handle receiving message duplicates.
//
// The chans block for slow consumers.
func InitRx() (<-chan msg.Raw, chan<- string, error) {
	var rx = make(chan msg.Raw)
	var dx = make(chan string)

	cred := credentials.NewStaticCredentials(accessKey, secretKey, "")
	_, err := cred.Get()
	if err != nil {
		log.Fatal("Get credential error (did you put SQS in config?):", err)
		return rx, dx, err
	}
	sess := session.New(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: cred,
	})
	svc := sqs.New(sess)
	var q *sqs.GetQueueUrlOutput

	for {
		params := &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		}
		q, err = svc.GetQueueUrl(params)
		if err != nil {
			log.Printf("WARN - problem getting SQS queue: %s", err)
			log.Println("WARN - sleeping and trying to get SQS queue again")
			time.Sleep(retry)
			continue
		}

		break
	}

	go receive(*q.QueueUrl, svc, rx)
	go delete(*q.QueueUrl, svc, dx)

	return rx, dx, nil
}

func receive(qUrl string, svc *sqs.SQS, rx chan msg.Raw) {
	for {
		param := &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(qUrl),
			MaxNumberOfMessages: aws.Int64(int64(MaxNumberOfMessages)),
			VisibilityTimeout:   aws.Int64(int64(VisibilityTimeout)),
			WaitTimeSeconds:     aws.Int64(int64(WaitTimeSeconds)),
		}

		r, err := svc.ReceiveMessage(param)
		if err != nil {
			log.Println("WARN - problem receiving messages from SQS, sleeping, continuing.")
			time.Sleep(retry)
			continue
		}

		if (len(r.Messages)) > 0 {
			for _, raw := range r.Messages {
				m := msg.Raw{
					Body:          *raw.Body,
					ReceiptHandle: *raw.ReceiptHandle,
				}
				rx <- m
			}
		}
	}
}

func delete(qUrl string, svc *sqs.SQS, dx chan string) {
	for {
		m := <-dx

		params := &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(qUrl), // Required
			ReceiptHandle: aws.String(m),    // Required
		}

		_, err := svc.DeleteMessage(params)

		if err != nil {
			log.Println("WARN - problem deleting messages from SQS, continuing.")
		}

	}
}
