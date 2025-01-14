package elastic

/*
The MIT License (MIT)
Copyright (c) 2021 Davide Pataracchia
Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:
The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.
*/

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
)

/*
Client :: Elasticsearch client object
*/
type Client struct {
	client *elastic.Client
	ctx    context.Context
	ip     string
	port   int
}

/*
New :: New client for Elasticsearch API
*/
func New(ip string, port int) *Client {
	e := &Client{
		ip:   ip,
		port: port,
	}

	conn := "http://" + e.ip + ":" + strconv.Itoa(e.port)
	log.Println("Waiting for elasticsearch node...")
	client, err := elastic.NewClient(elastic.SetURL(conn))
	for err != nil {
		client, err = elastic.NewClient(
			elastic.SetURL(conn),
			elastic.SetHealthcheckTimeoutStartup(30*time.Second),
			elastic.SetSniff(false),
		)
	}
	log.Println("Connected to elasticsearch!")
	e.client = client
	e.ctx = context.Background()

	err = e.CreateIndex("dump-hub", entryMapping)
	if err != nil {
		log.Fatal(err)
	}
	err = e.CreateIndex("dump-hub-history", historyMapping)
	if err != nil {
		log.Fatal(err)
	}
	e.waitGreen()

	var wg sync.WaitGroup
	wg.Add(2)

	go e.cleanHistory(&wg)
	go cleanTmp(&wg)
	wg.Wait()

	return e
}
