package core

import (
	"net/http"

	"github.com/buger/jsonparser"
)

func HttpStatusCode498CheckWriter(w http.ResponseWriter) http.ResponseWriter {
	return &httpStatusCode498CheckWriter{
		writer: w,
	}
}

type httpStatusCode498CheckWriter struct {
	writer http.ResponseWriter
}

func (h *httpStatusCode498CheckWriter) Write(p []byte) (n int, err error) {
	var has498StatusCode bool

	jsonparser.ArrayEach(p, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		statusCode, _ := jsonparser.GetInt(value, "extensions", "statusCode")
		if statusCode == 498 {
			has498StatusCode = true
		}

	}, "errors")
	if has498StatusCode {
		h.writer.WriteHeader(498)
	}
	return h.writer.Write(p)
}

func (h *httpStatusCode498CheckWriter) Header() http.Header {
	return h.writer.Header()
}

func (h *httpStatusCode498CheckWriter) WriteHeader(statusCode int) {
	h.writer.WriteHeader(statusCode)
}

