package tracer

import (
	"gin_websocket/lib/logger"
	"strings"
)

type Tracer struct {
	tags []string
}

func (t *Tracer) AddTag(tag Tag) {
	if tag.Key == TagHttpClientURL && tag.Value == "::1" {
		tag.Value = "127.0.0.1"
	}
	t.tags = append(t.tags, StringTag(tag))
}

func (t *Tracer) Finish() {
	if len(t.tags) > 0 {
		result := strings.Join(t.tags, " | ")
		logger.Runtime.Info(result)
	}
}
