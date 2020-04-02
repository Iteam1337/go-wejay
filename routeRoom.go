package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

type routeRoom struct {
	re *regexp.Regexp
}

func init() {
	router.room = routeRoom{regexp.MustCompile(`[a-z0-9-]`)}
}

func (route routeRoom) parsedRoomName(input []string) (out string) {
	if input == nil || strings.TrimSpace(input[0]) == "" {
		return
	}

	pre := strings.TrimSpace(input[0])
	pre = strings.ReplaceAll(pre, " ", "-")
	log.Println(pre)

	m := route.re.FindAllString(pre, -1)
	if m == nil {
		return
	}

	out = strings.Join(m, "")

	return
}

func (route *routeRoom) join(w http.ResponseWriter, r *http.Request) {
	room := route.parsedRoomName(r.URL.Query()["name"])

	if room == "" {
		http.Redirect(w, r, "//"+r.Host+"/", 307)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, e := w.Write([]byte(`{"path":"/room/` + room + `"}`)); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func (route *routeRoom) leave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, e := w.Write([]byte(`{"path":"/room/leave"}`)); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func (route *routeRoom) rooms(w http.ResponseWriter, r *http.Request) {

}
