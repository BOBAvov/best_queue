package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sso/models"
	"testing"
)

func TestSignUpAndIn(t *testing.T) {
	username := "test"
	password := "test"
	tg_nick := "@test"
	group := "ИУ7-12Б"

	jsonData, err := json.Marshal(models.RegisterUser{
		Username: username,
		Password: password,
		TgNick:   tg_nick,
		Group:    group,
	})
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post("http://127.0.0.1:8080/auth/sign-up", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println(res.StatusCode)
		t.Errorf("Response should be StatusOK")
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

	err = res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	jsonData, err = json.Marshal(models.AuthUser{
		TgNick:   tg_nick,
		Password: password,
	})

	res, err = http.Post("http://127.0.0.1:8080/auth/sign-in", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Response should be StatusOK")
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
	err = res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

}
