package main

// aws s3 sync s3://seiscompml07 ./ --exclude "*"  --include "2015p*"

import (
	"github.com/GeoNet/haz/database"
	"github.com/GeoNet/haz/msg"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	db       database.DB
	spoolDir = os.Getenv("SC3_SPOOL_DIR")
)

func main() {
	var err error

	db, err = database.InitPG()
	if err != nil {
		log.Fatalf("ERROR: problem with DB config: %s", err)
	}
	defer db.Close()

	db.Check()

	files, err := ioutil.ReadDir(spoolDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	sc3ml := make(chan os.FileInfo)

	go func() {
		defer close(sc3ml)

		for _, fi := range files {
			if !strings.HasSuffix(fi.Name(), ".xml") {
				continue
			}
			sc3ml <- fi
		}
	}()

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			procSC3ML(sc3ml)

			wg.Done()
		}()
	}
	wg.Wait()
}

func procSC3ML(sc3ml <-chan os.FileInfo) {
	for fi := range sc3ml {
		log.Println(fi.Name())
		q := msg.ReadSC3ML07(spoolDir + "/" + fi.Name())
		if q.Err() != nil {
			log.Println("WARN ignoring errored SC3ML: " + fi.Name() + " " + q.Err().Error())
			continue
		}

		err := db.SaveQuake(q)
		if err != nil {
			log.Fatal(err)
		}
	}
}
