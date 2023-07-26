package gateway

import (
	"errors"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGateway_Table(t *testing.T) {

	tests := []struct {
		name          string
		fileToServe   string
		respStatus200 bool
		wantErr       error
	}{
		{
			name:          "txt file",
			fileToServe:   "test.txt",
			respStatus200: true,
			wantErr:       nil,
		},
		{
			name:          "xlsx file",
			fileToServe:   "example_file.xlsx",
			respStatus200: true,
			wantErr:       nil,
		},
		{
			name:          "failed to fetch resource",
			fileToServe:   "test.txt",
			respStatus200: false,
			wantErr:       errors.New("failed to fetch resource"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(newTestHandler(tt.fileToServe, tt.respStatus200))

			g := NewGateway()

			r, err := g.Table(server.URL)

			assert.Equal(t, tt.wantErr, err, "method error")
			if err != nil {
				t.SkipNow()
			}

			bytesReadFromMethod, err := io.ReadAll(r)
			if err != nil {
				assert.FailNow(t, err.Error())
			}

			f, err := os.Open(filepath.Join("test", tt.fileToServe))
			if err != nil {
				assert.FailNow(t, err.Error())
			}

			bytesReadFromFile, err := io.ReadAll(f)
			if err != nil {
				assert.FailNow(t, err.Error())
			}

			assert.Equal(t, bytesReadFromFile, bytesReadFromMethod, "content")
		})
	}
}
