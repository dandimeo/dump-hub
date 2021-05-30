package common

import (
	"io"
	"log"
	"os"
)

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

/*
Banner - Dump Hub Cool Banner
*/
const Banner = `                          
   _                   _       _   
 _| |_ _ _____ ___ ___| |_ _ _| |_ 
| . | | |     | . |___|   | | | . |
|___|___|_|_|_|  _|   |_|_|___|___|
              |_|       
			             
`

/*
Host - API Host
Port - API Port
BaseAPI - API root folder
*/
const (
	Host    = "0.0.0.0"
	Port    = 8080
	BaseAPI = "/api/"
)

/*
EHost - Elasticsearch IP
EPort - Elasticsearch port
*/
const (
	EHost = "elasticsearch"
	EPort = 9200
)

/*
Error - Error while indexing file
Enqueued - File waiting to be indexed
Processing - File indexing in progress
Complete - File indexing complete
Deleting - Deleting entries
*/
const (
	Processing = 0
	Deleting   = 1
	Enqueued   = 2
	Error      = 3
	Complete   = 4
)

// Default GO init func to set log file redirection globally
func init() {
	filename := "/var/log/dump-hub/dump-hub.log" // Needs root permission to access file

	logFile, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0600,
	)

	if err == nil {
		defer logFile.Close()

		logWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(logWriter)
	} else {
		log.Printf("[ERROR] error creating logfile %s", err.Error())
		log.Println("logging to stdout only")
		log.SetOutput(os.Stdout)
	}
}