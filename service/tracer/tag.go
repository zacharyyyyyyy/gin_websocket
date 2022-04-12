package tracer

import "fmt"

const (
	TagHttpMethod      = "http.method"
	TagHttpStatusCode  = "http.status_code"
	TagHttpURL         = "http.url"
	TagHttpRawURL      = "http.raw_url"
	TagHttpClientURL   = "http.client_url"
	TagHttpParam       = "http.param"
	TagHttpContentType = "http.content_type"
)

type Tag struct {
	Key   string
	Value string
}

func StringTag(tag Tag) string {
	return fmt.Sprintf("%s:%s", tag.Key, tag.Value)
}
