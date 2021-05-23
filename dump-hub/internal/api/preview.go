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
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/x0e1f/dump-hub/internal/common"
)

func previewFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var previewReq common.PreviewReq

		err := json.NewDecoder(r.Body).Decode(&previewReq)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		fileName := common.EncodeFilename(previewReq.FileName)
		filePath := filepath.Join(uploadFolder, fileName)
		previewData, err := readPreview(
			filePath,
			previewReq.Start,
		)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		preview := common.PreviewResult{
			Preview: *previewData,
		}
		response, err := json.Marshal(preview)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func readPreview(filePath string, start int) (*[]string, error) {
	previewData := []string{}

	if start < 0 {
		start = 0
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	for scanner.Scan() {
		line := scanner.Text()
		if currentLine >= (start + previewSize) {
			break
		}

		if currentLine >= start {
			previewData = append(previewData, line)
		}
		currentLine++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &previewData, nil
}