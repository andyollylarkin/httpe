package main

import (
	"fmt"
	"net/http"

	"github.com/andyollylarkin/httpe"
)

func main() {
	m := httpe.NewErrorMessageRaw(httpe.Code("BAD"), 1, http.StatusInternalServerError)

	fmt.Println(m.Error())
}
