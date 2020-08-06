package info

import "net/http"

func serveVersion(response http.ResponseWriter, _ *http.Request) {
	writeFile("version", response)
}
