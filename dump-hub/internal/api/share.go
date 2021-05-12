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
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/x0e1f/dump-hub/internal/common"
)

const uploadFolder = "/opt/uploads/"
const maxFileSize = 10000 * 1000000
const chunkSize = 1000
const pageSize = 20

/*
computeChecksum :: compute file checksum
*/
func computeChecksum(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

/*
readUploadFolder :: Return the list of files in the uploads folder
*/
func readUploadFolder() ([]common.File, error) {
	var files = []common.File{}

	fileInfo, err := ioutil.ReadDir(uploadFolder)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		/* Ignore hidden files */
		if file.Name()[0] == '.' {
			continue
		}

		/* File size in MB */
		fileSize := file.Size() / 1000000

		/* Retrieve original filename */
		fileName, err := decodeFilename(file.Name())
		if err != nil {
			log.Println(err)
			continue
		}

		uFile := common.File{
			FileName: fileName,
			Size:     fileSize,
		}
		files = append(files, uFile)
	}

	return files, nil
}

/*
encodeFilename :: Encode filename to base64
*/
func encodeFilename(fileName string) string {
	/* Remove file extension */
	fileName = strings.TrimSuffix(
		fileName,
		filepath.Ext(fileName),
	)

	/* Encode to base64 */
	fileName = base64.StdEncoding.
		EncodeToString([]byte(fileName))

	return fileName
}

/*
decodeFilename :: Decode filename from base64
*/
func decodeFilename(fileName string) (string, error) {
	/* Decode from base64 */
	data, err := base64.
		StdEncoding.
		DecodeString(fileName)
	if err != nil {
		return "", err
	}
	fileName = string(data)

	return fileName, nil
}
