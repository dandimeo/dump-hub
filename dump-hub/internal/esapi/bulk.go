package esapi

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
	"bufio"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/x0e1f/dump-hub/internal/common"
	"github.com/x0e1f/dump-hub/internal/parser"
)

/*
BulkWorker - Elasticsearch BulkAPI Worker
*/
type BulkWorker struct {
	mutex     sync.Mutex
	chunkSize int
}

/*
newBulkWorker - Create BulkWorker
*/
func newBulkWorker(chunkSize int) *BulkWorker {
	bulkW := BulkWorker{
		chunkSize: chunkSize,
	}
	return &bulkW
}

/*
IndexFile - Process file for indexing (mutual)
*/
func (eClient *Client) IndexFile(parser *parser.Parser) {
	eClient.UpdateUploadStatus(
		parser.Checksum,
		common.Enqueued,
	)

	eClient.bulkw.mutex.Lock()
	log.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAA")
	eClient.UpdateUploadStatus(
		parser.Checksum,
		common.Processing,
	)

	file, err := os.Open(parser.Filepath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	quitChan := make(chan struct{})
	entryChan := make(chan *common.Entry)
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go indexRoutine(
			i,
			&wg,
			eClient,
			quitChan,
			entryChan,
		)
	}

	currentLine := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if currentLine < parser.Start {
			continue
		}

		entry := parser.ParseEntry(scanner.Text())
		if entry == nil {
			continue
		}
		entryChan <- entry

		currentLine++
	}

	close(quitChan)
	close(entryChan)
	wg.Wait()

	eClient.Refresh()
	eClient.UpdateUploadStatus(
		parser.Checksum,
		common.Complete,
	)

	eClient.bulkw.mutex.Unlock()
}
