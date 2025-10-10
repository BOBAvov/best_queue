package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

type StatusApi struct {
	Message string `json:"message"`
}

func TestApiIsStarting(t *testing.T) {
	res, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println(res.StatusCode)
		t.Error("status code is not 200, got", res.StatusCode)
	}

	ans := new(StatusApi)
	if err = json.NewDecoder(res.Body).Decode(&ans); err != nil {
		t.Fatal(err)
	}

	fmt.Println(ans.Message)
}
