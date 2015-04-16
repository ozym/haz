// sqs handles messages from AWS SQS via channels.
package sqs

import (
	"github.com/GeoNet/cfg"
	"github.com/GeoNet/goamz/sqs"
	"github.com/GeoNet/msg"
	"log"
	"strconv"
	"time"
)

var retry = time.Duration(30) * time.Second

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
func InitRx(s *cfg.SQS) (<-chan msg.Raw, chan<- string, error) {
	var rx = make(chan msg.Raw)
	var dx = make(chan string)

	a, err := sqs.NewFrom(s.AccessKey, s.SecretKey, s.AWSRegion)
	if err != nil {
		return rx, dx, err
	}

	var q *sqs.Queue
	for {
		q, err = a.GetQueue(s.QueueName)
		if err != nil {
			log.Printf("WARN - problem getting SQS queue: %s", err)
			log.Println("WARN - sleeping and trying to get SQS queue again")
			time.Sleep(retry)
			continue
		}
		break
	}

	var rp = map[string]string{
		"MaxNumberOfMessages": strconv.Itoa(s.MaxNumberOfMessages),
		"VisibilityTimeout":   strconv.Itoa(s.VisibilityTimeout),
		"WaitTimeSeconds":     strconv.Itoa(s.WaitTimeSeconds),
	}

	go receive(q, rp, rx)
	go delete(q, dx)

	return rx, dx, nil
}

func receive(q *sqs.Queue, receiveParams map[string]string, rx chan msg.Raw) {
	for {
		r, err := q.ReceiveMessageWithParameters(receiveParams)
		if err != nil {
			log.Println("WARN - problem receiving messages from SQS, sleeping, continuing.")
			time.Sleep(retry)
			continue
		}

		if (len(r.Messages)) > 0 {
			for _, raw := range r.Messages {
				m := msg.Raw{
					Body:          raw.Body,
					ReceiptHandle: raw.ReceiptHandle,
				}
				rx <- m
			}
		}
	}
}

func delete(q *sqs.Queue, dx chan string) {
	for {
		m := <-dx
		d := sqs.Message{
			ReceiptHandle: m,
		}

		_, err := q.DeleteMessage(&d)
		if err != nil {
			log.Println("WARN - problem deleting messages from SQS, continuing.")
		}
	}
}
