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
	"encoding/json"
	"io"
	"log"

	"github.com/olivere/elastic/v7"
	"github.com/x0e1f/dump-hub/internal/common"
)

/*
CreateIndex :: Create elasticsearch index if not exists
*/
func (eClient *Client) CreateIndex(index string, mapping string) error {
	exists, err := eClient.client.IndexExists(index).Do(eClient.ctx)
	if err != nil {
		return err
	}

	if !exists {
		_, err := eClient.client.
			CreateIndex(index).
			Body(mapping).
			Do(eClient.ctx)
		if err != nil {
			return err
		}
		log.Printf("Created elasticsearch index: %s", index)
	}

	return nil
}

/*
Search :: Search entries using simple query string API
*/
func (eClient *Client) Search(queryString string, from int, size int) (*common.SearchResult, error) {
	var results *elastic.SearchResult
	var err error

	if len(queryString) < 1 || queryString == "*" {
		query := elastic.NewMatchAllQuery()
		results, err = eClient.client.Search().
			Index("dump-hub").
			Query(query).
			From(from).
			Size(size).
			Do(eClient.ctx)
		if err != nil {
			return nil, err
		}
	} else {
		query := elastic.
			NewMultiMatchQuery(
				queryString,
				"_all",
			).
			Type("match_phrase")
		results, err = eClient.client.Search().
			Index("dump-hub").
			Query(query).
			From(from).
			Size(size).
			Do(eClient.ctx)
		if err != nil {
			return nil, err
		}
	}

	/* Populate search results */
	searchResult := common.SearchResult{}
	for _, hit := range results.Hits.Hits {
		entry := common.Entry{}
		err := json.Unmarshal(hit.Source, &entry)
		if err != nil {
			log.Println(err)
			break
		}

		searchResult.Results = append(
			searchResult.Results,
			entry,
		)
	}
	searchResult.Tot = int(results.Hits.TotalHits.Value)

	return &searchResult, nil
}

/*
BulkInsert :: Index entries with BulkAPI
*/
func (eClient *Client) BulkInsert(e []*common.Entry) error {
	bulkRequest := eClient.client.Bulk()

	for _, entry := range e {
		req := elastic.NewBulkIndexRequest().
			OpType("index").
			Index("dump-hub").
			Doc(entry)

		bulkRequest = bulkRequest.Add(req)
	}

	_, err := bulkRequest.
		Do(eClient.ctx)
	if err != nil {
		return err
	}

	return nil
}

/*
BulkDelete :: Delete entries with BulkAPI
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

/*
IsAlreadyUploaded :: Check if file is already uploaded (by checksum)
*/
func (eClient *Client) IsAlreadyUploaded(checkSum string) (bool, error) {
	exists, err := eClient.client.Exists().
		Index("dump-hub-history").
		Id(checkSum).
		Do(eClient.ctx)
	if err != nil {
		return false, err
	}

	return exists, nil
}

/*
NewStatusDocument :: New status document on dump-hub-uploads index
*/
func (eClient *Client) NewStatusDocument(h *common.Status, checkSum string) error {
	data, err := json.Marshal(h)
	if err != nil {
		return err
	}

	_, err = eClient.client.Index().
		Index("dump-hub-uploads").
		BodyString(string(data)).
		Id(checkSum).
		Refresh("true").
		Do(eClient.ctx)
	if err != nil {
		return err
	}

	return nil
}

/*
UpdateUploadStatus :: Update status field of an history element
*/
func (eClient *Client) UpdateUploadStatus(checkSum string, newStatus int) error {
	_, err := eClient.client.Update().
		Index("dump-hub-uploads").
		Id(checkSum).
		Doc(map[string]interface{}{"status": newStatus}).
		Do(eClient.ctx)
	if err != nil {
		return err
	}

	return nil
}

/*
GetStatus :: Get status documents
*/
func (eClient *Client) GetStatus(from int, size int) (*common.StatusData, error) {
	query := elastic.NewMatchAllQuery()

	results, err := eClient.client.Search().
		Index("dump-hub-uploads").
		Query(query).
		From(from).
		Size(size).
		Do(eClient.ctx)
	if err != nil {
		return nil, err
	}

	statusData := common.StatusData{}
	for _, hit := range results.Hits.Hits {
		history := common.Status{}
		err := json.Unmarshal(hit.Source, &history)
		if err != nil {
			log.Println(err)
			break
		}

		statusData.Results = append(
			statusData.Results,
			history,
		)
	}
	statusData.Tot = int(results.Hits.TotalHits.Value)

	return &statusData, nil
}

/*
DeleteEntries :: Delete entries associated to a file (checkSum)
*/
func (eClient *Client) DeleteEntries(checkSum string) {
	eClient.UpdateUploadStatus(checkSum, 2)

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
			log.Println(err)
		}
		for _, hit := range result.Hits.Hits {
			if len(chunk) > ChunkSize {
				err := eClient.BulkDelete(chunk)
				if err != nil {
					log.Println(err)
				}
				chunk = []string{}
			}

			chunk = append(chunk, hit.Id)
		}

		/* If there is still data, delete chunk */
		if len(chunk) > 0 {
			err := eClient.BulkDelete(chunk)
			if err != nil {
				log.Println(err)
			}
		}
	}

	/* Refresh elastic index */
	eClient.Refresh()

	/* Delete status document */
	_, err := eClient.client.Delete().
		Index("dump-hub-uploads").
		Id(checkSum).
		Do(eClient.ctx)
	if err != nil {
		log.Println(err)
	}
}

/*
Refresh :: Refresh dump-hub index
*/
func (eClient *Client) Refresh() error {
	_, err := eClient.client.Refresh().
		Index("dump-hub").
		Do(eClient.ctx)
	if err != nil {
		return err
	}

	return nil
}
