package internal

import (
	"net/http"
	"path/filepath"
	"strings"
)

func newTestServerHandler() http.Handler {
	fs := fileServer{}

	sm := http.NewServeMux()
	sm.HandleFunc("/", fs.GetTable)

	return sm
}

type fileServer struct {
}

func (fs *fileServer) GetTable(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(strings.TrimPrefix(r.URL.EscapedPath(), "/"), "/")

	http.ServeFile(w, r, filepath.Join(path...))
}
