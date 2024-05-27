package internal

import (
	"io"
	"net/http"
)

func GetIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is rando by sahithyandev")
}
