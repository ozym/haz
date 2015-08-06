// +build devtest

// `go test -tags devtest` to test
package main

import (
	"github.com/GeoNet/goamz/aws"
	"github.com/GeoNet/goamz/s3"
	"log"
	"os"
	"os/exec"
	"testing"
)

func setup() {
	auth := aws.Auth{AccessKey: config.S3.AccessKey, SecretKey: config.S3.SecretKey}
	s = s3.New(auth, aws.Region{Name: config.S3.AWSRegion, S3Endpoint: config.S3.EndPoint})

	config.SeiscompmlS3.InDir += "_test"
	config.SeiscompmlS3.OutDir += "_test"
	config.SeiscompmlS3.UnprocessedDir += "_test"

	// create in/out/unprocessed if not exists
	if _, err := os.Stat(config.SeiscompmlS3.InDir); err != nil {
		log.Println("Creating in dir ", config.SeiscompmlS3.InDir)
		err = os.MkdirAll(config.SeiscompmlS3.InDir, 0755)
		if err != nil {
			log.Fatalln("Cannot make dir", err.Error())
		}
	}

	if _, err := os.Stat(config.SeiscompmlS3.OutDir); err != nil {
		log.Println("Creating out dir ", config.SeiscompmlS3.InDir)
		err = os.MkdirAll(config.SeiscompmlS3.OutDir, 0755)
		if err != nil {
			log.Fatalln("Cannot make dir", err.Error())
		}
	}

	if _, err := os.Stat(config.SeiscompmlS3.UnprocessedDir); err != nil {
		log.Println("Creating unprocessed dir ", config.SeiscompmlS3.InDir)
		err = os.MkdirAll(config.SeiscompmlS3.UnprocessedDir, 0755)
		if err != nil {
			log.Fatalln("Cannot make dir", err.Error())
		}
	}
}

func teardown() {
	// this cleans up all test folders and data
	// NOTE: The directories' full path might not been deleted because we created directories by MkdirAll
	os.RemoveAll(config.SeiscompmlS3.InDir)
	os.RemoveAll(config.SeiscompmlS3.OutDir)
	os.RemoveAll(config.SeiscompmlS3.UnprocessedDir)

	// xml stored to s3 will not be deleted.
}

func TestSciscompmlS3(t *testing.T) {
	setup()
	defer teardown()

	b := s.Bucket(config.S3.Bucket)

	cmd := exec.Command("cp", "./testdata/test2023p087463.xml", config.SeiscompmlS3.InDir)
	err := cmd.Run()
	if err != nil {
		t.Errorf("Copy test data failed:", err.Error())
	}

	checkAndProcess()

	h, err := b.GetResponse("2023p087463.xml")
	if err != nil {
		t.Error("Failed to get from S3 bucket:", err.Error())
	}

	agencyid, modificationtime := getMeta(h)
	if agencyid != "WEL(GNS_Primary)" || modificationtime != "2023-02-02T08:41:39.970373Z" {
		t.Error("Remote metadata does not match.")
	}

	err = b.Del("2023p087463.xml")
	if err != nil {
		log.Println("WARN: failed to delete test data in S3:", err.Error())
	}
}
