package utils

import (
	"errors"
	"net/http"
)

// ParseRequest …
func ParseRequest(state string, r *http.Request) (code string, err error) {
	values := r.URL.Query()

	if e := values.Get("error"); e != "" {
		err = errors.New("spotify: auth failed - " + e)
		return
	}

	if code = values.Get("code"); code == "" {
		err = errors.New("spotify: didn't get access code")
		return
	}

	if actualState := values.Get("state"); actualState != state {
		err = errors.New("spotify: redirect state parameter doesn't match")
		return
	}

	return
}
