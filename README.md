# Dump Hub <img src="dump-hub-fe/src/favicon.ico" width=22px/>

<img src="assets/header.png"/>

Self hosted search engine for data leaks and password dumps. Upload and parse multiple files, then quickly search through all stored items with the power of Elasticsearch.

**Disclaimer:** _This project does not include, and will never include, any data. Data must be uploaded by end users on their own instances of Dump Hub. I take no responsability for the nature of uploaded data._

Dump Hub currenlty supports dumps in csv/combo-list format, the parser is not strict, so if, for instance, one of the lines is not correctly formatted it will still be indexed as a single document and the parsing process will continue through the entire file.

## Docker Compose

Dump Hub can be executed by using docker-compose:

```
git clone https://github.com/x0e1f/dump-hub.git
cd dump-hub
docker-compose up --build
```

Dump Hub will bind port 443 on 0.0.0.0 by default.

**Warning:** _Do not expose Dump Hub on public networks! Please edit the **docker-compose.yml** file to fit your needs, evaluate to use your own SSL certificate and evaluate the usage of basic auth on nginx._

## Usage

To start using Dump Hub open a web browser on https://[dump-hub-ip].

**Warning:** _When you upload a file on Dump Hub, the entries will be parsed and indexed on Elasticsearch. You will be able to perform search queries on every field really quickly, but this comes at a cost: **bigger disk usage**. Keep that in mind if you need to work with a lot of data!_

The upload of a new file can be executed by following 2 steps:

<img src="assets/folder-line.svg" width=20px> **Multiple File Upload**

To upload files on Dump Hub use the upload page and select the desired files. Once one or more files are selected the upload process will begin. You will find a list of already uploaded files on the first section of the same page.

<img src="assets/plus-circle-line.svg" width=20px> **Analyze**

In order to add entries on Elasticsearch use the analyze page. From the analyze page you are able to select one of the already uploaded file. Select one file by clicking on it. Once one file is selected, a preview of the file content will be displayed on the box below. You can edit two settings by using the form above the preview box:

- <img src="assets/flag-line.svg" width=15px> **Starting Line:** The parsing will start from this line. When editing this value the preview will update accordingly.
- <img src="assets/filter-line.svg" width=15px> **Separator:** This is the separator character. This char will be used to split entries on line (just like a standard csv).

If the parser is correctly configured you will be able to see parsed items as columns in the table at the bottom of the page. From this table you can select which columns will be parsed and included in the final document (highlighted in green). Each of those fields will be indexed and fully searchable. When the desired result appears in the table you can start the analyze process by clicking on **Analyze File** button.

<img src="assets/block-line.svg" width=20px> **Data**

From the Data page you are able to view the list of:

- Entries in processing status (Files that are currently being analyzed and uploaded to Elasticsearch).
- Entries in deleting status (Files that are currently being deleted from Elasticsearch).
- Entries in pending status (Files that are waiting to be analyzed or deleted. Only one file can be analyzed or deleted at one time).
- Entries in error status.
- Entries in completed status.

From this page you are able to delete entries in completed or error status.

## License

The MIT License (MIT)

Copyright (c) 2021 Davide Pataracchia

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
