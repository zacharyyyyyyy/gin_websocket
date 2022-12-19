package tracer

import (
	"net/http"
	"strings"

	"gin_websocket/lib/logger"
)

type Tracer struct {
	tags       []string
	resultCode int
}

func (t *Tracer) AddTag(tag Tag) {
	if tag.Key == TagHttpClientURL && tag.Value == "::1" {
		tag.Value = "127.0.0.1"
	}
	t.tags = append(t.tags, StringTag(tag))
}
func (t *Tracer) AddResultCode(code int) {
	t.resultCode = code
}

func (t *Tracer) Finish() {
	if len(t.tags) > 0 {
		result := strings.Join(t.tags, " | ")
		if t.resultCode == http.StatusOK {
			logger.Runtime.Info(result)
		} else {
			logger.Runtime.Error(result)
		}

	}
}
