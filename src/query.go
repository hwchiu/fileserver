package fileserver

import (
	"net/url"
)

type QueryUrl struct {
	Url url.Values `json:"url"`
}

func New(url url.Values) *QueryUrl {
	return &QueryUrl{
		Url: url,
	}
}

func (query *QueryUrl) Str(key string) (string, bool) {
	values := query.Url[key]

	if len(values) == 0 {
		return "", false
	}

	return values[0], true
}
