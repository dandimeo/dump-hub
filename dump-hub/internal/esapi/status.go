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
	"encoding/json"

	"github.com/olivere/elastic/v7"
	"github.com/x0e1f/dump-hub/internal/common"
)

/*
GetStatus - Get status documents (paginated)
*/
func (eClient *Client) GetStatus(from int, size int) (*common.StatusResult, error) {
	query := elastic.NewMatchAllQuery()
	sortQ := elastic.NewFieldSort("status")

	results, err := eClient.client.Search().
		Index("dump-hub-status").
		SortBy(sortQ).
		Query(query).
		From(from).
		Size(size).
		Do(eClient.ctx)
	if err != nil {
		return nil, err
	}

	statusData := common.StatusResult{}
	for _, hit := range results.Hits.Hits {
		status := common.Status{}
		err := json.Unmarshal(hit.Source, &status)
		if err != nil {
			return nil, err
		}

		statusData.Results = append(
			statusData.Results,
			status,
		)
	}
	statusData.Tot = int(results.Hits.TotalHits.Value)

	return &statusData, nil
}

func (eClient *Client) GetDocumentStatus(checkSum string) (*common.Status, error) {
	result, err := eClient.client.Get().
		Index("dump-hub-status").
		Id(checkSum).
		Do(eClient.ctx)
	if err != nil {
		return nil, err
	}

	status := common.Status{}
	err = json.Unmarshal(result.Source, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

/*
NewStatusDocument - New status document on dump-hub-status index
*/
func (eClient *Client) NewStatusDocument(h *common.Status, checkSum string) error {
	data, err := json.Marshal(h)
	if err != nil {
		return err
	}

	_, err = eClient.client.Index().
		Index("dump-hub-status").
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
UpdateUploadStatus - Update status field of an upload status document
*/
func (eClient *Client) UpdateUploadStatus(checkSum string, newStatus int) error {
	_, err := eClient.client.Update().
		Index("dump-hub-status").
		Id(checkSum).
		Doc(map[string]interface{}{"status": newStatus}).
		Refresh("true").
		Do(eClient.ctx)
	if err != nil {
		return err
	}

	return nil
}
