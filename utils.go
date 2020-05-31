/**
@author anilasimbilen
**/
package main

import (
	"encoding/json"
	"net/http"
)

type RejectProto struct {
	statusCode int
	status     func(code int) RejectProto
	send       func()
}
type ResolveProto struct {
	statusCode int
	status     func(code int) ResolveProto
	send       func()
	body       interface{}
}

func reject(w http.ResponseWriter, err error) RejectProto {
	var res RejectProto

	if initType := w.Header().Get("content-type"); initType == "" {
		w.Header().Add("content-type", "application/json")
	}
	__send := func() {
		w.WriteHeader(res.statusCode)
		w.Write([]byte(`{"message": ` + err.Error() + ` }`))
	}
	__status := func(code int) RejectProto {
		res.statusCode = code
		return res
	}

	res = RejectProto{
		status:     __status,
		send:       __send,
		statusCode: http.StatusInternalServerError,
	}
	return res
}

func resolve(w http.ResponseWriter, body interface{}) ResolveProto {
	var res ResolveProto
	if initType := w.Header().Get("content-type"); initType == "" {
		w.Header().Add("content-type", "application/json")
	}
	__send := func() {
		w.WriteHeader(res.statusCode)
		json.NewEncoder(w).Encode(res.body)
	}
	__status := func(code int) ResolveProto {
		res.statusCode = code
		return res
	}
	res = ResolveProto{
		status:     __status,
		send:       __send,
		statusCode: http.StatusAccepted,
		body:       body,
	}
	return res
}
