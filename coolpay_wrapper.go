package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/valyala/fasthttp"
)

const (
	authURL         = "https://coolpay.herokuapp.com/api/login"
	addRecipientURL = "https://coolpay.herokuapp.com/api/recipients"
	makePaymentURL  = "https://coolpay.herokuapp.com/api/payments"
	listPaymentsURL = "https://coolpay.herokuapp.com/api/payments"
)

type Token string

type CredentialsAPICP struct {
	Username string `json:"username"`
	Apikey   string `json:"apikey"`
}

type RecipientCP struct {
	Recipient RecipientInfoCP `json:"recipient"`
}

type RecipientInfoCP struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type PaymentCP struct {
	Payment PaymentInfoCP `json:"payment"`
}
type PaymentsCP struct {
	Payments []PaymentInfoCP `json:"payments"`
}

type PaymentInfoCP struct {
	Id          string  `json:"id,omitempty"`
	Amount      float32 `json:"amount,string"`
	Currency    string  `json:"currency"`
	RecipientId string  `json:"recipient_id"`
	Status      string  `json:"status,omitempty"`
}

/*
func main() {
	//curl -v # verbose
	//curl -H "Content-Type: application/json" -X POST -d '{"username":"ShahzadI","apikey":"CDA8777865C7CC3C"}' https://coolpay.herokuapp.com/api/login

	//curl -H "Content-Type: application/json" -H "Authorization: Bearer 12345.yourtoken.67890" -X POST -d '{"recipient":{"name":"Jake Farland"}}' https://coolpay.herokuapp.com/api/recipients
	//curl -H "Content-Type: application/json" -H "Authorization: Bearer 15d8180b-d3e4-41d1-b114-f215b6e540fb" -X POST -d '{"recipient":{"name":"Jake Farland"}}' https://coolpay.herokuapp.com/api/recipients

	//curl -H "Content-Type: application/json" -H "Authorization: Bearer 12345.yourtoken.67890" -X POST -d '{"payment":{"amount":10.5,"currency":"GBP","recipient_id":"previously.added.recipient.id"}}' https://coolpay.herokuapp.com/api/payments
	//curl -H "Content-Type: application/json" -H "Authorization: Bearer 15d8180b-d3e4-41d1-b114-f215b6e540fb" -X POST -d '{"payment":{"amount":10.5,"currency":"GBP","recipient_id":"6e7b146e-5957-11e6-8b77-86f30ca893d3"}}' https://coolpay.herokuapp.com/api/payments

	//curl -i -H "Content-Type: application/json" -H "Authorization: Bearer 12345.yourtoken.67890" https://coolpay.herokuapp.com/api/payments
	//curl -i -H "Accept: application/json" "server:5050/a/c/getName{"param0":"pradeep"}"

	credentials := &CredentialsAPICP{"ShahzadI", "CDA8777865C7CC3C"}
	token, httpStatusCode := getAuthToken(*credentials, authURL)

	if http.StatusOK == httpStatusCode {
		recipientInfo := &RecipientInfoCP{Name: "Janu Jarman"}
		returnedRecipientInfo, _ := addRecipient(*recipientInfo, token, addRecipientURL)

		paymentInfo := &PaymentInfoCP{Amount: 10.5, Currency: "GBP", RecipientId: returnedRecipientInfo.Id}
		returnedPaymentInfo, _ := makePaymentToRecipient(*paymentInfo, token, makePaymentURL)

		status, _ := verifyPaymentToRecipient(*returnedPaymentInfo, token, listPaymentsURL)
		if "paid" == status {
			fmt.Println("\nPayment VERIFIED")
		} else {
			fmt.Printf("\nPayment NOT Verified with status:%s\n", status)
		}
	}
}

*/

func verifyPaymentToRecipient(paymentInfoToVerify PaymentInfoCP, token *Token, url string) (string, int) {
	fmt.Printf("\nVerifying payment amount: %f ...\n", paymentInfoToVerify.Amount)

	resp := sendRequestGetResponse(nil, token, "GET", url)
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		log.Println("List all payments request did NOT return successfully")
		return "", resp.StatusCode
	}

	var returnedPayments PaymentsCP
	if err := json.NewDecoder(resp.Body).Decode(&returnedPayments); err != nil {
		log.Println(err)
	}

	status := verifyPayment(paymentInfoToVerify, returnedPayments.Payments)
	return status, resp.StatusCode
}

func verifyPayment(paymentInfoToVerify PaymentInfoCP, allPayments []PaymentInfoCP) string {
	for _, payment := range allPayments {
		if paymentInfoToVerify.Id == payment.Id {
			if paymentInfoToVerify.Amount == payment.Amount &&
				paymentInfoToVerify.Currency == payment.Currency && paymentInfoToVerify.RecipientId == payment.RecipientId {
				return payment.Status
			}
			return "Payment credentials DONOT match"
		}
	}
	return "NOT found"
}

func makePaymentToRecipientCP(paymentInfo PaymentInfoCP, token *Token, url string) (*PaymentInfoCP, int) {
	fmt.Printf("\nMaking payment: %f ...\n", paymentInfo.Amount)
	payment := &PaymentCP{paymentInfo}
	jsonBytes, err := json.Marshal(*payment)
	if err != nil {
		panic(err)
	}
	resp := sendRequestGetResponse(jsonBytes, token, "POST", url)
	defer resp.Body.Close()

	if http.StatusCreated != resp.StatusCode {
		log.Println("Make payment request did NOT return successfully")
		return nil, resp.StatusCode
	}

	var returnedPayment PaymentCP
	if err := json.NewDecoder(resp.Body).Decode(&returnedPayment); err != nil {
		log.Println(err)
	}
	return &returnedPayment.Payment, resp.StatusCode
}

func addRecipientCP(recipientInfo RecipientInfoCP, token *Token, url string) (*RecipientInfoCP, int) {
	fmt.Printf("\nAdding recipient: %s ...\n", recipientInfo.Name)
	recipient := &RecipientCP{recipientInfo}
	jsonBytes, err := json.Marshal(*recipient)
	if err != nil {
		panic(err)
	}
	resp := sendRequestGetResponse(jsonBytes, token, "POST", url)
	defer resp.Body.Close()

	if http.StatusCreated != resp.StatusCode {
		log.Println("Add recipient request did NOT return successfully")
		return nil, resp.StatusCode
	}
	var returnedRecipient RecipientCP
	if err := json.NewDecoder(resp.Body).Decode(&returnedRecipient); err != nil {
		log.Println(err)
	}
	return &returnedRecipient.Recipient, resp.StatusCode
}

func getAuthTokenCP(credentials CredentialsAPICP, url string) (*Token, int) {
	fmt.Printf("\nAuthenticating user: %s ...\n", credentials.Username)
	jsonBytes, err := json.Marshal(credentials)
	if err != nil {
		panic(err)
	}
	resp := sendRequestGetResponse(jsonBytes, nil, "POST", url)
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		log.Println("User Authentication request did NOT return successfully")
		return nil, resp.StatusCode
	}

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	// Fill the record with the data from the JSON
	var token Token

	data := make(map[string]string)
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println(err)
	}
	token = Token(data["token"])
	return &token, resp.StatusCode
}

func getAuthTokenCPFast(credentials CredentialsAPICP, url string) (authToken *Token, respStatusCode int) {
	fmt.Printf("\nFast Authenticating user: %s ...\n", credentials.Username)

	jsonBytes, err := json.Marshal(credentials)
	if err != nil {
		panic(err)
	}

	respData, respStatusCode := sendRequestGetResponseFast(jsonBytes, nil, "POST", url)
	fmt.Printf("StatusCode: %d\nrespData is: %s\n", respStatusCode, respData)
	if fasthttp.StatusOK != respStatusCode {
		log.Println("User Authentication request did NOT return successfully")
		return nil, respStatusCode
	}

	// Fill the record with the data from the JSON
	var token Token

	token = Token(respData["token"])
	fmt.Printf("authToken is: %s\n", string(token))
	return &token, respStatusCode
}

func sendRequestGetResponse(jsonBytes []byte, token *Token, method, url string) *http.Response {
	fmt.Println(string(jsonBytes))
	var req *http.Request
	var err error
	if nil == jsonBytes {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	}

	req.Header.Set("Content-Type", "application/json")
	if nil != token {
		req.Header.Set("Authorization", "Bearer "+string(*token))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

func sendRequestGetResponseFast(jsonBytes []byte, token *Token, method, url string) (data map[string]string, statusCode int) {
	fmt.Println(string(jsonBytes))
	data = make(map[string]string)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	if nil != jsonBytes {
		req.SetBody(jsonBytes)
	}
	// req.SetBodyString("p=q")

	req.Header.Set("Content-Type", "application/json")
	if nil != token {
		fmt.Println("must set the token")
		//req.Header.Set("Authorization", "Bearer "+string(*token))
		req.Header.Set("Authorization", "Bearer "+string(*token))
	}
	fmt.Printf("request Header is: %s\n", req.Header.String())
	fmt.Printf("request body is: %s\n", req.Body())
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		// println("Error:", err.Error())
		panic(err)
	} else {
		bodyBytes := resp.Body()
		fmt.Printf("recieved data is: %s\n", string(bodyBytes))
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			log.Println(err)
		}
	}
	fmt.Printf("recieved data unmarshalled is: %#v\n", data)

	return data, resp.StatusCode()

	// bodyBytes := resp.Body()
	// println(string(bodyBytes))
	// User-Agent: fasthttp
	// Body: p=q
}
