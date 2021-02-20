package render

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type H map[string]interface{}

var StatusCtxKey = &contextKey{"Status"}

func Status(r *http.Request, status int) {
	*r = *r.WithContext(context.WithValue(r.Context(), StatusCtxKey, status))
}

func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if status, ok := r.Context().Value(StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write(buf.Bytes())
}

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "chi render context value " + k.name
}
