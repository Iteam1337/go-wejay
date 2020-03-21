package cookie

import (
	"encoding/base64"
	"net/http"
)

// GetID …
func GetID(r *http.Request) (id string, err error) {
	var idChars []byte
	cookie, err := r.Cookie("user")
	if err != nil {
		return
	}

	idChars, err = base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return
	}

	id = string(idChars)
	return
}

// GetIDORreturn403 …
func GetIDORreturn403(w http.ResponseWriter, r *http.Request) (id string, err error) {
	var idChars []byte
	cookie, err := r.Cookie("user")
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		return
	}

	idChars, err = base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		http.Error(w, "Couldn't parse cookie", http.StatusBadRequest)
	}

	id = string(idChars)
	return
}

// CreateAndSet …
func CreateAndSet(w http.ResponseWriter, r *http.Request) (id string) {
	id = r.FormValue("state")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	cookie := http.Cookie{
		Name: "user",
		Path: "/",
	}

	cookie.Value = base64.StdEncoding.EncodeToString([]byte(id))

	http.SetCookie(w, &cookie)

	return
}
