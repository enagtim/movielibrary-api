package main

import (
	"net/http"
)

func main() {
	http.ListenAndServe(":8002", nil)
}
