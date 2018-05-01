package main

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Welcome!\n")
}

func hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
}

func authenticateUser(ctx *fasthttp.RequestCtx) {
	token, httpStatusCode := sendAuthUserReq(ctx, authURL)
	if fasthttp.StatusOK != httpStatusCode {
		fmt.Errorf("User Authentication returned http StatusCode was incorrect, got: %d.", httpStatusCode)
		ctx.SetStatusCode(httpStatusCode)
	} else {
		returnAuthUserResp(ctx, token)
	}
}

/*
func addRecipient(ctx *fasthttp.RequestCtx) {
	token, httpStatusCode := sendAddRecipientReq(ctx, token, authURL)
	if fasthttp.StatusOK != httpStatusCode {
		fmt.Errorf("User Authentication returned http StatusCode was incorrect, got: %d.", httpStatusCode)
		ctx.SetStatusCode(httpStatusCode)
	} else {
		returnAuthUserResp(ctx, token)
	}
}
*/
func sendAuthUserReq(ctx *fasthttp.RequestCtx, authURL string) (authToken *Token, respStatusCode int) {
	var recReqCredentials UserCredentials
	postArgs := ctx.PostArgs()
	fmt.Fprintf(ctx, "login, %s!\n", ctx.PostBody())
	// with "Content-Type: application/x-www-form-urlencoded" -X POST -d "username=ShahzadI&password=CDA8777865C7CC3C"
	fmt.Printf("postArgs are: %+v\n", postArgs)
	fmt.Printf("postArgs username is: %s\n", string(postArgs.Peek("username")))
	fmt.Printf("postArgs password is: %s\n", string(postArgs.Peek("password")))

	/*
		// with "Content-Type: application/json" -X POST -d '{"username":"ShahzadI","password":"CDA8777865C7CC3C"}'
		if err := json.Unmarshal(ctx.PostBody(), &recReqCredentials); err != nil {
			log.Println(err)
		}
	*/
	// fmt.Fprintf(ctx, "login, %s!\n", postArgs)
	recReqCredentials.Username = string(postArgs.Peek("username"))
	recReqCredentials.Password = string(postArgs.Peek("password"))
	fmt.Printf("username is: %s\n", recReqCredentials.Username)
	fmt.Printf("password is: %s\n", recReqCredentials.Password)

	senReqCredentials := &CredentialsAPICP{"ShahzadI", "CDA8777865C7CC3C"}
	return getAuthTokenCPFast(*senReqCredentials, authURL)
}

func returnAuthUserResp(ctx *fasthttp.RequestCtx, authToken *Token) {
	fmt.Println("Token = ", authToken)
	mapData := map[string]string{"token": string(*authToken)}
	respJSONBytes, _ := json.Marshal(mapData)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(respJSONBytes)
}

/*
func sendAddRecipientReq(recipientInfo RecipientInfoCP, token *Token, url string) (*RecipientInfoCP, int) {
	recipientInfo := &RecipientInfoCP{Name: "Janu Jarman"}
	returnedRecipientInfo, _ := addRecipient(*recipientInfo, token, addRecipientURL)

}*/
