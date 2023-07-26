package gateway

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type Gateway struct {
	hc http.Client
}

func NewGateway() *Gateway {
	c := http.Client{Timeout: 10 * time.Second}

	return &Gateway{
		hc: c,
	}
}

func (g *Gateway) Table(url string) (io.Reader, error) {

	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := g.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("response body close error: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to fetch resource: %s", resp.Status)

		return nil, errors.New("failed to fetch resource")
	}

	// buf := make([]byte, resp.ContentLength)
	buf, err := io.ReadAll(resp.Body)
	switch {
	case err == io.EOF: // всё хорошо

	case err != nil:
		return nil, err
	}

	return bytes.NewBuffer(buf), nil
}
