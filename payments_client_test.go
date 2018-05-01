package main

import (
	//"log"
	//"io/ioutil"
	//"encoding/json"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestUserAuthFast(t *testing.T) {
	credentials := &Credentials{"ShahzadI", "CDA8777865C7CC3C"}

	token, httpStatusCode := getAuthTokenFast(*credentials, authURL)
	if fasthttp.StatusOK != httpStatusCode {
		t.Errorf("User Authentication returned http StatusCode was incorrect, got: %d, want: %d.", httpStatusCode, http.StatusOK)
	}
	fmt.Println("Token = ", token)
}

func TestUserAuthFastBasic(t *testing.T) {
	credentials := &Credentials{"ShahzadI", "CDA8777865C7CC3C"}
	//getAuthTokenFast(*credentials, authURL)
	jsonBytes, err := json.Marshal(credentials)
	if err != nil {
		panic(err)
	}

	token, httpStatusCode := sendRequestGetResponseFast(jsonBytes, nil, "POST", authURL)

	// token, httpStatusCode := getAuthTokenFast(*credentials, authURL)
	if fasthttp.StatusOK != httpStatusCode {
		t.Errorf("User Authentication returned http StatusCode was incorrect, got: %d, want: %d.", httpStatusCode, http.StatusOK)
	}
	fmt.Println("Token = ", token)

}

func TestUserAuth(t *testing.T) {
	credentials := &Credentials{"ShahzadI", "CDA8777865C7CC3C"}

	token, httpStatusCode := getAuthToken(*credentials, authURL)
	if http.StatusOK != httpStatusCode {
		t.Errorf("User Authentication returned http StatusCode was incorrect, got: %d, want: %d.", httpStatusCode, http.StatusOK)
	}
	fmt.Println("Token = ", string(*token))

}
