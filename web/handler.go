package web

import "net/http"

type Handler interface {
	ServeHTTP(req *http.Request) Response
}

type HandlerFunc func(*http.Request) Response

func (fn HandlerFunc) ServeHTTP(req *http.Request) Response { return fn(req) }

func Wrap(handler Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		res := safeHandlerCall(handler, req)
		if res != nil {
			res.WriteTo(w)
		}
	})
}

func WrapFunc(fn func(*http.Request) Response) http.Handler {
	return Wrap(HandlerFunc(fn))
}

func safeHandlerCall(handler Handler, req *http.Request) (res Response) {
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			res = Error("panic", err)
		}
	}()

	return handler.ServeHTTP(req)
}
