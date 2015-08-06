// haz-ua-consumer listens to an AWS SQS queue for Haz JSON messages and
// generate tags for push message subscribers, then send it to UA to push message out.
package main

import (
	"github.com/AdRoll/goamz/aws"
	"github.com/GeoNet/cfg"
	"github.com/GeoNet/goamz/s3"
	"github.com/GeoNet/haz/msg"
	"github.com/GeoNet/log/logentries"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//go:generate configer seiscompml-to-s3.json
var (
	config = cfg.Load()
	s      *s3.S3
)

func init() {
	logentries.Init(config.Logentries.Token)
	msg.InitLibrato(config.Librato.User, config.Librato.Key, config.Librato.Source)
	config.SeiscompmlS3.CheckInterval = 60
}

// main kicks off SeisComPML processing and HeartBeat generation.
func main() {
	auth := aws.Auth{AccessKey: config.S3.AccessKey, SecretKey: config.S3.SecretKey}
	s = s3.New(auth, aws.Region{Name: config.S3.AWSRegion, S3Endpoint: config.S3.EndPoint})
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

	sc3ml()
}

// sc3 for processing SeisComPML files.
type sc3 struct {
	msg.Quake
	f string
}

// sc3ml polls the spool dir for SeisComPML files and runs Process on them.
func sc3ml() {
	inter := time.Duration(config.SeiscompmlS3.CheckInterval) * time.Second
	for {
		checkAndProcess()
		time.Sleep(inter)
	}
}

func checkAndProcess() {
	files, err := ioutil.ReadDir(config.SeiscompmlS3.InDir)
	if err != nil {
		log.Printf("WARN: %s\n", err.Error())
	}

	for _, fi := range files {
		log.Println(fi.Name())
		if strings.HasSuffix(fi.Name(), ".xml") {
			s := &sc3{f: config.SeiscompmlS3.InDir + "/" + fi.Name()}

			var d string

			if msg.Process(s) {
				d = config.SeiscompmlS3.OutDir + "/" + fi.Name()
				log.Printf("Process %s successed.\n", fi.Name())
			} else {
				d = config.SeiscompmlS3.UnprocessedDir + "/" + fi.Name()
				log.Printf("Process failed. Will move %s to unprocessed.\n", fi.Name())
			}

			// remove destination file anyway
			if _, err := os.Stat(d); err == nil {
				if err := os.Remove(d); err != nil {
					log.Printf("%s : %s\n", err.Error())
				}
			}
			// move to destination folder
			err := os.Rename(s.f, d)
			if err != nil {
				log.Printf("%s : %s\n", err.Error())
			}
		}
	}
}

// Process processes SeisComPML files.  Converts them to a msg.Quake, checks the quality, sends them
// to an AWS S3 bucket.
// Returns true if sent, false if unprocessed
func (sc *sc3) Process() bool {
	sc.Quake = msg.ReadSC3ML07(sc.f)
	if sc.Err() != nil {
		log.Printf("Read error for file %s : %s\n", sc.f, sc.Err())
		return false
	}

	sc.RxLog()
	// check if this need to be upload
	upload := true
	b := s.Bucket(config.S3.Bucket)

	if sc.Quake.AgencyID != "" && !sc.Quake.ModificationTime.IsZero() {
		// Check if we're updating a existing publicID...
		// So try to retrieve metadata from S3 and see if this is newer.
		h, err := b.GetResponse(sc.Quake.PublicID + ".xml")
		if err == nil {
			agencyid, modificationtime := getMeta(h)
			// log.Println("remote ", agencyid, modificationtime)
			if sc.Quake.AgencyID == agencyid {
				if sc.Quake.ModificationTime.Format(time.RFC3339Nano) <= modificationtime {
					upload = false
					log.Println("Skip upload ", sc.f)
				}
			}
		} else {
			// this is new, or the file in S3 doesn't contain valid metadata
		}
	}

	if upload {
		bt, err := ioutil.ReadFile(sc.f)
		if err != nil {
			log.Printf("Failed to read %s while trying to upload it: %s", sc.f, err.Error())
			return false
		}
		err = b.Put(sc.Quake.PublicID+".xml", bt, "application/xml", s3.Private, s3.Options{Meta: makeMeta(sc.Quake)})
		if err != nil {
			log.Println("Failed to store into S3 bucket:", err.Error())
			return false
		}
		log.Println("Upload to S3 complete for file ", sc.f)
		sc.TxLog()
		return true
	}
	return false
}

func makeMeta(q msg.Quake) (m map[string][]string) {
	m = make(map[string][]string)
	m["agencyid"] = []string{q.AgencyID}
	m["modificationtime"] = []string{q.ModificationTime.Format(time.RFC3339Nano)}
	return
}

func getMeta(h *http.Response) (agencyid string, modificationTime string) {
	agencyid = h.Header.Get("x-amz-meta-agencyid")
	modificationTime = h.Header.Get("x-amz-meta-modificationtime")
	return
}
