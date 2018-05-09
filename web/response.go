package web

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

type Response interface {
	WriteTo(http.ResponseWriter) error
	GetStatusCode() int
}

type ResponseWithLength interface {
	Response
	GetLength() int
}

type NormalResponse struct {
	statusCode int
	body       []byte
	header     http.Header
}

var (
	Forbidden   = Empty(403)
	NotFound    = Empty(404)
	ServerError = Empty(500)
)

func (r *NormalResponse) WriteTo(w http.ResponseWriter) (err error) {
	headers := w.Header()
	headers.Set("Content-Type", strconv.Itoa(r.GetLength()))
	for k, v := range r.header {
		headers[k] = v
	}
	w.WriteHeader(r.statusCode)

	if r.body != nil {
		_, err = w.Write(r.body)
	}
	return
}

func (r *NormalResponse) GetStatusCode() int {
	return r.statusCode
}

func (r *NormalResponse) GetLength() int {
	return len(r.body)
}

func (r *NormalResponse) Header(key string, value string) *NormalResponse {
	r.header.Add(key, value)
	return r
}

func (r *NormalResponse) Type(contentType string) *NormalResponse {
	return r.Header("Content-Type", contentType)
}

func Empty(statusCode int) *NormalResponse {
	return &NormalResponse{
		statusCode: statusCode,
		body:       nil,
		header:     http.Header{},
	}
}

func Redirect(statusCode int, location string) *NormalResponse {
	return Empty(statusCode).Header("Location", location)
}

func Error(message string, err error) *NormalResponse {
	stack := string(debug.Stack())

	log.Println(message, err)
	for _, line := range strings.Split(stack, "\n") {
		log.Println(line)
	}

	return ServerError
}

func Data(statusCode int, body []byte) *NormalResponse {
	return &NormalResponse{
		statusCode: statusCode,
		body:       body,
		header:     http.Header{},
	}
}

func Text(statusCode int, body string) *NormalResponse {
	return Data(statusCode, []byte(body))
}

func Json(statusCode int, body interface{}) *NormalResponse {
	b, err := json.Marshal(body)
	if err != nil {
		return Error("json marshal failed", err)
	}
	return Data(statusCode, b).Type("application/json")
}
