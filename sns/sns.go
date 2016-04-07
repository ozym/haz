// sns is for sending messages to AWS SNS topics.
package sns

import (
	"fmt"
	"github.com/GeoNet/haz/msg"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
	"os"
	"time"
)

var (
	retry     = time.Duration(30) * time.Second
	accessKey = os.Getenv("SNS_ACCESS_KEY")
	secretKey = os.Getenv("SNS_SECRET_KEY")
	topicArn  = os.Getenv("SNS_TOPIC_ARN")
	awsRegion = os.Getenv("AWS_REGION")
)

type SNS struct {
	s *sns.SNS
}

func Init() (SNS, error) {
	cred := credentials.NewStaticCredentials(accessKey, secretKey, "")
	_, err := cred.Get()
	if err != nil {
		log.Fatal("Get credential error (did you put SNS in config?):", err)
		return SNS{}, err
	}
	sess := session.New(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: cred,
	})

	s := SNS{
		s: sns.New(sess),
	}

	return s, err
}

// Publish publishes m to the AWS SNS topic cfg.SNS.TopicArn.
// Errors on publish or an empty MessageId in the publish response are an error.
// If an error is encountered then publishing is attempted retries more times with
// a pause of 30s between each attempt.  retries can be 0 to attempt publishing only once.
func (s *SNS) Publish(m msg.Raw, retries int) (err error) {
	c := 0
	for {
		params := &sns.PublishInput{
			Message:  aws.String(m.Body),
			TopicArn: aws.String(topicArn),
		}

		if m.Subject != "" {
			params.Subject = aws.String(m.Subject)
		}

		r, err := s.s.Publish(params)
		if err == nil {
			if *r.MessageId == "" {
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
