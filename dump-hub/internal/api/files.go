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
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/x0e1f/dump-hub/internal/common"
	"github.com/x0e1f/dump-hub/internal/elastic"
)

/*
files :: Get files in upload folder (GET)
*/
func files(eClient *elastic.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var results common.FilesResult

		files, err := readUploadFolder()
		if err != nil {
			log.Println(err)
		}
		results.Files = files

		response, err := json.Marshal(results)
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

/*
deleteFile :: Delete file in upload folder (DELETE)
*/
func deleteFile(eClient *elastic.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		data, err := base64.
			StdEncoding.
			DecodeString(id)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			log.Println(err)
			return
		}
		fileName := string(data)

		fileName = encodeFilename(fileName)
		filePath := filepath.Join(uploadFolder, fileName)
		err = os.Remove(filePath)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
}
