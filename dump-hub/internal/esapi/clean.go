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
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/olivere/elastic/v7"
	"github.com/x0e1f/dump-hub/internal/common"
)

/*
cleanTmp - Clean tmp folder
*/
func cleanTmp() {
	log.Println("Cleaning tmp folder...")

	dir, err := ioutil.ReadDir("/tmp")
	if err != nil {
		log.Println(err)
	}
	for _, d := range dir {
		os.Remove(path.Join("/tmp", d.Name()))
	}
}

/*
cleanStatus - Clean unprocessed files and update status
*/
func (eClient *Client) cleanStatus() {
	log.Println("Cleaning status of unprocessed files...")

	matchQ := elastic.NewMatchQuery(
		"status",
		common.Processing,
	)
	query := elastic.
		NewBoolQuery().
		Must(matchQ)

	scroll := eClient.client.Scroll().
		Index("dump-hub-status").
		Query(query).
		Size(1)

	for {
		result, err := scroll.Do(eClient.ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}

		for _, hit := range result.Hits.Hits {
			err = eClient.UpdateUploadStatus(hit.Id, common.Error)
			if err != nil {
				log.Println(err)
			}
		}
	}

	matchQ = elastic.NewMatchQuery(
		"status",
		common.Deleting,
	)
	query = elastic.
		NewBoolQuery().
		Must(matchQ)

	scroll = eClient.client.Scroll().
		Index("dump-hub-status").
		Query(query).
		Size(1)

	for {
		result, err := scroll.Do(eClient.ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}

		for _, hit := range result.Hits.Hits {
			err = eClient.UpdateUploadStatus(hit.Id, common.Error)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
