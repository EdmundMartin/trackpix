package server

import (
	"net/http"
	"time"
)

// TrackingServer contains settings for tracking pixel endpoint
type TrackingServer struct {
	Headers    []string
	PixelBytes []byte
	Results    chan *TrackingData
}

// TrackingData represents
type TrackingData struct {
	Host     string
	Protocol string
	Remote   string
	Datetime time.Time
	Headers  map[string]string
}

func extractHeaders(keys []string, headers http.Header) map[string]string {
	trHeaders := make(map[string]string)
	for _, h := range keys {
		val, ok := headers[h]
		if ok {
			trHeaders[h] = val[0]
		}
	}
	return trHeaders
}

func servePixel(w http.ResponseWriter, img []byte) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(img)
}

// Pixel responds with a 1x1 pixel after extracting information from request
func (t *TrackingServer) Pixel(w http.ResponseWriter, r *http.Request) {
	td := &TrackingData{Host: r.Host, Protocol: r.Proto, Remote: r.RemoteAddr, Datetime: time.Now()}
	td.Headers = extractHeaders(t.Headers, r.Header)
	t.Results <- td
	servePixel(w, t.PixelBytes)
}
