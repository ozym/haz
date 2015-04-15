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
// Messages from SQS can be received from the read chan.  Messages that should be deleted from SQS
// can be sent to the write chan.
//
// decode should decode []byte to a msg.MessageTx type with the underlying concrete implementation
// that you expect from the SQS queue.  Messages that cause a decode error are not sent to the read chan.
//
// Set deleteErr true to delete any messages from the queue that cause decode to throw an error.  false to not
// delete (which would make sense for use with a redrive queue).
//
// The chan block for slow consumers.
//
// A type assertion will be needed to receive messages of the concrete type you require e.g.,
//
//   rx, dx, err := sqs.InitRx(s, msg.HazDecode, false)
//   for {
//		m := <-rx
//
//		h, ok := m.(msg.Haz)
//		if ok {
//                         // do something with h
//				dx <- m // delete m from SQS
//			}
//		} else {
//			log.Println("received non msg.Haz message, deleting.")
//			dx <- m
//		}
//	}
func InitRx(s *cfg.SQS, decode func([]byte, string) (msg.MessageTx, error), deleteErr bool) (<-chan msg.MessageTx, chan<- msg.MessageTx, error) {
	var rx = make(chan msg.MessageTx)
	var dx = make(chan msg.MessageTx)

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

	go receive(q, rp, rx, decode, deleteErr)
	go delete(q, dx)

	return rx, dx, nil
}

func receive(q *sqs.Queue, receiveParams map[string]string, rx chan msg.MessageTx, decode func([]byte, string) (msg.MessageTx, error), deleteErr bool) {
	for {
		r, err := q.ReceiveMessageWithParameters(receiveParams)
		if err != nil {
			log.Println("WARN - problem receiving messages from SQS, sleeping, continuing.")
			time.Sleep(retry)
			continue
		}

		if (len(r.Messages)) > 0 {
			for _, raw := range r.Messages {
				m, err := decode([]byte(raw.Body), raw.ReceiptHandle)
				if err != nil && deleteErr {
					log.Println("WARN - error decoding SQS message: " + err.Error())
					_, err = q.DeleteMessage(&raw)
					if err != nil {
						log.Println("WARN - problem deleting messages from SQS, continuing.")
					}
				} else {
					rx <- m
				}
			}
		}
	}
}

func delete(q *sqs.Queue, dx chan msg.MessageTx) {
	for {
		m := <-dx
		d := sqs.Message{
			ReceiptHandle: m.ReceiptHandle(),
		}

		_, err := q.DeleteMessage(&d)
		if err != nil {
			log.Println("WARN - problem deleting messages from SQS, continuing.")
		}
	}
}
