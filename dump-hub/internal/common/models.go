package common

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
Entry :: Entry document
*/
type Entry struct {
	Origin   string   `json:"origin"`
	OriginID string   `json:"origin_id"`
	Data     []string `json:"data"`
}

/*
Status :: Dump Hub Status document
*/
type Status struct {
	Date     string `json:"date"`
	Filename string `json:"filename"`
	Checksum string `json:"checksum"`
	Status   int    `json:"status"`
}

/*
File :: File in uploads folder
*/
type File struct {
	FileName string `json:"filename"`
	Size     int64  `json:"size"`
}

/*
StatusData :: Dump Hub Status API Response
*/
type StatusData struct {
	Results []Status `json:"results"`
	Tot     int      `json:"tot"`
}

/*
SearchReq :: Search API request
*/
type SearchReq struct {
	Query string `json:"query"`
	Page  int    `json:"page"`
}

/*
SearchResult :: Search API response
*/
type SearchResult struct {
	Results []Entry `json:"results"`
	Tot     int     `json:"tot"`
}

/*
StatusReq :: Status API request
*/
type StatusReq struct {
	Page int `json:"page"`
}

/*
FilesResult :: Files API response
*/
type FilesResult struct {
	Dir   string `json:"dir"`
	Files []File `json:"files"`
}

/*
DeleteReq :: Delete API request
*/
type DeleteReq struct {
	Checksum string `json:"checksum"`
}

/*
PreviewReq :: Preview API request
*/
type PreviewReq struct {
	FileName string `json:"filename"`
	Start    int    `json:"start"`
}

/*
PreviewResult :: Preview API response
*/
type PreviewResult struct {
	Preview []string `json:"preview"`
}

/*
AnalyzeReq :: Analyze API request
*/
type AnalyzeReq struct {
	Filename string `json:"filename"`
	Pattern  string `json:"pattern"`
	Columns  []int  `json:"columns"`
}
