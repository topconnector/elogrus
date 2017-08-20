// Copyright 2017 Mobile Data Books, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// dep init
import (
	"context"
	"fmt"
	"github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
	"gopkg.in/sohlich/elogrus.v2"
	"net/http"
	"os"
	"time"
)

type viewHandler_helloHandler struct {
	// Logger is the log.Logger instance used to log messages with the Logger middleware
	Logger *logrus.Logger
}

func (v *viewHandler_helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	t1 := time.Now()

	t1 = time.Now()
	t01 := t1.Sub(t0)
	tStr1 := fmt.Sprintf("Hello World from Go in minimal Docker container (4.28MB) - tc-helloworld-go-ws-logging-elasticsearch - v.1.0, it took %v to run", t01)
	//fmt.Printf("\n%s\n", tStr1)
	v.Logger.Info(tStr1)

	fmt.Fprintln(w, tStr1)

}

// http://elasticsearch-logging:9200
func main() {

	envNAMESPACE := os.Getenv("NAMESPACE")
	envNODE_NAME := os.Getenv("NODE_NAME")
	envELASTICSEARCH_URL := os.Getenv("ELASTICSEARCH_URL")

	log := logrus.New()
	log.Println("NAMESPACE:" + envNAMESPACE)
	log.Println("NODE_NAME:" + envNODE_NAME)
	log.Println("elastic.SetURL:" + envELASTICSEARCH_URL)

	// elastic.SetSniff(false)
	// client, err := elastic.NewClient(elastic.SetURL(envELASTICSEARCH_URL)) // http://10.96.185.105:9200 10.0.1.73:30138  http://10.32.0.8:9200
	// if err != nil {
	// 	log.Panic(err)
	// }

	// Obtain a client for an Elasticsearch cluster of one node,
	// running on envELASTICSEARCH_URL. Do not run the sniffer.
	// Set the healthcheck interval to 10s. When requests fail,
	// retry 5 times. Print error messages to os.Stderr and informational
	// messages to os.Stdout.
	client, err := elastic.NewClient(
		elastic.SetURL(envELASTICSEARCH_URL),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
	//elastic.SetRetrier(NewCustomRetrier()),
	//elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
	//elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags))
	)
	if err != nil {
		log.Panic(err)
	}

	// logrus.PanicLevel,
	// logrus.FatalLevel,
	// logrus.ErrorLevel,
	// logrus.WarnLevel,
	// logrus.InfoLevel,
	// logrus.DebugLevel,

	hook, err := elogrus.NewElasticHook(client, "localhost", logrus.DebugLevel, "tc-helloworld-go-ws-logging-elasticsearch-log")
	if err != nil {
		log.Panic(err)
	}
	log.Hooks.Add(hook)

	for index := 0; index < 100; index++ {
		logrus.Infof("Test msg %d", time.Now().Unix())
	}

	time.Sleep(5 * time.Second)

	termQuery := elastic.NewTermQuery("Host", "localhost")
	searchResult, err := client.Search().
		Index("tc-helloworld-go-ws-logging-elasticsearch-log").
		Query(termQuery).
		Do(context.TODO())

	if searchResult.Hits.TotalHits != 100 {
		//t.Error("Not all logs pushed to elastic")
		//t.FailNow()
		err := "Not all logs pushed to elastic"
		log.Panic(err)
	}

	//

	http.Handle("/", &viewHandler_helloHandler{
		Logger: log,
	})
	log.Println("tc-helloworld-go-ws-logging-elasticsearch: started, serving at 1010")

	//fmt.Println("Started, serving at 1010")
	err = http.ListenAndServe(":1010", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
