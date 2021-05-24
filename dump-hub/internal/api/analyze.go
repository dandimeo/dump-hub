package api

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
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/x0e1f/dump-hub/internal/common"
	"github.com/x0e1f/dump-hub/internal/elastic"
	"github.com/x0e1f/dump-hub/internal/parser"
)

func analyze(eClient *elastic.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var analyzeReq common.AnalyzeReq

		err := json.NewDecoder(r.Body).Decode(&analyzeReq)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if len(analyzeReq.Columns) < 1 {
			log.Println("Invalid columns value")
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		fileName := common.EncodeFilename(analyzeReq.Filename)
		originPath := filepath.Join(uploadFolder, fileName)
		if _, err := os.Stat(originPath); os.IsNotExist(err) {
			log.Println("File does not exist")
			http.Error(w, "", http.StatusNotFound)
			return
		}

		go analyzeFile(
			eClient,
			analyzeReq,
			fileName,
		)

		w.WriteHeader(http.StatusOK)
	}
}

func analyzeFile(eClient *elastic.Client, analyzeReq common.AnalyzeReq, fileName string) {
	filePath, err := moveTemp(fileName)
	if err != nil {
		log.Println(err)
		return
	}

	checkSum, err := common.ComputeChecksum(filePath)
	if err != nil {
		log.Println(err)
		return
	}

	originalFilename, _ := common.DecodeFilename(fileName)
	date := time.Now().Format("2006-01-02 15:04:05")
	status := common.Status{
		Date:     date,
		Filename: originalFilename,
		Checksum: checkSum,
		Status:   0,
	}
	err = eClient.NewStatusDocument(&status, checkSum)
	if err != nil {
		log.Println(err)
		return
	}

	parser, err := parser.New(
		analyzeReq.Pattern,
		analyzeReq.Columns,
		analyzeReq.Filename,
		filePath,
	)
	if err != nil {
		log.Println("Unable to create parser object")
		return
	}

	processEntry(eClient, parser)
}

func moveTemp(fileName string) (string, error) {
	originPath := filepath.Join(uploadFolder, fileName)
	hiddenPath := filepath.Join(uploadFolder, "."+fileName)
	filePath := filepath.Join("/tmp/", fileName)

	os.Rename(
		originPath,
		hiddenPath,
	)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	originFile, err := os.Open(hiddenPath)
	if err != nil {
		return "", err
	}
	defer originFile.Close()

	_, err = io.Copy(file, originFile)
	if err != nil {
		return "", err
	}
	originFile.Close()

	err = os.Remove(hiddenPath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func processEntry(e *elastic.Client, parser *parser.Parser) {
	file, err := os.Open(parser.Filepath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	/* Start routines */
	var wg sync.WaitGroup
	quitChan := make(chan struct{})
	entryChan := make(chan *common.Entry)
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go uploader(i, &wg, e, quitChan, entryChan)
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

	/* Refresh elastic index */
	e.Refresh()

	/* Update status (Complete)*/
	e.UpdateUploadStatus(
		parser.Checksum,
		1,
	)
}

/*
uploader :: Upload entries to elastic
*/
func uploader(id int, wg *sync.WaitGroup, e *elastic.Client, quitChan <-chan struct{}, entryChan <-chan *common.Entry) {
	wg.Add(1)
	run := true
	chunk := []*common.Entry{}

	for run {
		/* Chunk size reached */
		if len(chunk) >= elastic.ChunkSize {
			err := e.BulkInsert(chunk)
			if err != nil {
				log.Println(err)
			}
			chunk = []*common.Entry{}
		}

		select {
		case <-quitChan:
			run = false
		case entry := <-entryChan:
			if entry == nil {
				continue
			}
			chunk = append(chunk, entry)
		}
	}

	/* If there is still data, upload chunk */
	if len(chunk) > 0 {
		err := e.BulkInsert(chunk)
		if err != nil {
			log.Println(err)
		}
	}

	wg.Done()
}
