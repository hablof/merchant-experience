package gateway

import (
	"net/http"
	"path/filepath"
)

type myHandler struct {
	filename  string
	status200 bool
}

func newTestHandler(filename string, status200 bool) http.Handler {
	mh := myHandler{
		filename:  filename,
		status200: status200,
	}

	sm := http.NewServeMux()
	sm.HandleFunc("/", mh.File)

	return sm
}

func (h *myHandler) File(w http.ResponseWriter, r *http.Request) {

	s := filepath.Join("test", h.filename)

	if h.status200 {
		http.ServeFile(w, r, s)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Status 500"))
	}
}
