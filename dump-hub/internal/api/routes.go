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
	"net/http"

	"github.com/gorilla/mux"
)

/*
defineRoutes :: Define API routes and handlers
*/
func (engine *Engine) defineRoutes() {
	router := mux.NewRouter().StrictSlash(true)

	router.
		Path(engine.baseAPI + "upload").
		Methods(http.MethodPost).
		HandlerFunc(upload(engine.eClient))

	router.
		Path(engine.baseAPI + "history").
		Methods(http.MethodPost).
		HandlerFunc(getHistory(engine.eClient))

	router.
		Path(engine.baseAPI + "search").
		Methods(http.MethodPost).
		HandlerFunc(search(engine.eClient))

	router.
		Path(engine.baseAPI + "delete").
		Methods(http.MethodPost).
		HandlerFunc(delete(engine.eClient))

	router.
		Path(engine.baseAPI + "files").
		Methods(http.MethodGet).
		HandlerFunc(files(engine.eClient))

	router.
		Path(engine.baseAPI + "files/{id}").
		Methods(http.MethodDelete).
		HandlerFunc(deleteFile(engine.eClient))

	engine.router = router
}
