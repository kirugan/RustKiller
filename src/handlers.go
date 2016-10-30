package main

import (
	"net/http"
	"fmt"
)

func simpleHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("FUCK!")
}
