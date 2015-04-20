// sns is for sending messages to AWS SNS topics.
package sns

import (
	"fmt"
	"github.com/AdRoll/goamz/aws" // Have to use the AdRoll version can't use a fork with GeoNet/goamz/sns
	"github.com/GeoNet/cfg"
	"github.com/GeoNet/goamz/sns"
	"github.com/GeoNet/msg"
	"log"
	"time"
)

var retry = time.Duration(30) * time.Second

type SNS struct {
	c *cfg.SNS
	s *sns.SNS
}

func Init(c *cfg.SNS) (SNS, error) {
	sn, err := sns.New(aws.Auth{AccessKey: c.AccessKey, SecretKey: c.SecretKey},
		aws.GetRegion(c.AWSRegion))

	s := SNS{
		c: c,
		s: sn,
	}

	return s, err
}

// Publish publishes m to the AWS SNS topic cfg.SNS.TopicArn.
// Errors on publish or an empty MessageId in the publish response are an error.
// If an error is encountered then publishing is attempted retries more times with
// a pause of 30s between each attempt.  retries can be 0 to attempt publishing only once.
func (s *SNS) Publish(m msg.Raw, retries int) (err error) {
	var r *sns.PublishResponse
	c := 0
	for {
		r, err = s.s.Publish(&sns.PublishOptions{TopicArn: s.c.TopicArn, Message: m.Body, Subject: m.Subject})
		if err == nil {
			if r.MessageId == "" {
				err = fmt.Errorf("Empty MessageID from send to SNS ")
			}

		}
		if err == nil {
			break
		}
		if c >= retries {
			break
		}
		c++

		log.Println("WARN " + err.Error())
		log.Println("WARN error sending to SNS.  Sleeping and trying again")

		time.Sleep(retry)
	}
	return err
}
