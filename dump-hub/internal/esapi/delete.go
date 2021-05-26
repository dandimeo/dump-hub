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
	"log"

	"github.com/olivere/elastic/v7"
	"github.com/x0e1f/dump-hub/internal/common"
)

/*
DeleteEntries - Delete entries associated to a file (checkSum)
*/
func (eClient *Client) DeleteEntries(checkSum string) {
	for {
		if !eClient.isBusy() {
			break
		}
	}
	eClient.setBusy(true)

	eClient.UpdateUploadStatus(
		checkSum,
		common.Deleting,
	)

	matchQ := elastic.NewMatchQuery(
		"origin_id",
		checkSum,
	)
	query := elastic.
		NewBoolQuery().
		Must(matchQ)

	scroll := eClient.client.Scroll().
		Index("dump-hub").
		Query(query).
		Size(1)

	chunk := []string{}
	for {
		result, err := scroll.Do(eClient.ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			eClient.UpdateUploadStatus(
				checkSum,
				common.Error,
			)
			log.Println(err)
		}
		for _, hit := range result.Hits.Hits {
			if len(chunk) > eClient.bulkw.cSize {
				err := eClient.BulkDelete(chunk)
				if err != nil {
					eClient.UpdateUploadStatus(
						checkSum,
						common.Error,
					)
					log.Println(err)
					return
				}
				chunk = []string{}
			}
			chunk = append(chunk, hit.Id)
		}

		if len(chunk) > 0 {
			err := eClient.BulkDelete(chunk)
			if err != nil {
				eClient.UpdateUploadStatus(
					checkSum,
					common.Error,
				)
				log.Println(err)
				return
			}
		}
	}

	eClient.Refresh()
	_, err := eClient.client.Delete().
		Index("dump-hub-status").
		Id(checkSum).
		Do(eClient.ctx)
	if err != nil {
		log.Println(err)
	}

	eClient.setBusy(false)
}

/*
BulkDelete - Delete entries with BulkAPI
*/
func (eClient *Client) BulkDelete(chunk []string) error {
	bulkRequest := eClient.client.Bulk()

	for _, id := range chunk {
		req := elastic.NewBulkDeleteRequest().
			Index("dump-hub").
			Id(id)

		bulkRequest = bulkRequest.Add(req)
	}

	_, err := bulkRequest.
		Do(eClient.ctx)
	if err != nil {
		return err
	}

	return nil
}
