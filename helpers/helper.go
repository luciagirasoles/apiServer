package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Catch Return panic
func Catch(err error) {
	if err != nil {
		panic(err)
	}
}

// RespondwithError return error message
func RespondwithError(w http.ResponseWriter, code int, msg string) {
	fmt.Println("respond with error ", msg)
	RespondwithJSON(w, code, map[string]string{"message": msg})
}

// RespondwithJSON write json response format
func RespondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	fmt.Printf("payload %#v \n", payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
