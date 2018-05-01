package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const PORT = 8086

/*
func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // send data to client side
}
*/

func main() {
	listenPort := ":" + strconv.Itoa(PORT)
	fmt.Printf("listening on port: %s\n", listenPort)

	router := fasthttprouter.New()
	router.GET("/", index)
	router.GET("/hello/:name", hello)
	router.POST("/api/login", authenticateUser)
	//router.POST("/api/recipient", addRecipient)
	log.Fatal(fasthttp.ListenAndServe(listenPort, router.Handler))
	/*
		http.HandleFunc("/", sayHello)
		http.HandleFunc("/hello", sayhelloName) // set router
		if err := http.ListenAndServe(listenPort, nil); err != nil {
			panic(err)
		}
	*/
	/*
		credentials := &Credentials{"ShahzadI", "CDA8777865C7CC3C"}
		token, httpStatusCode := getAuthToken(*credentials, authURL)

		if http.StatusOK == httpStatusCode {
			recipientInfo := &RecipientInfo{Name: "Janu Jarman"}
			returnedRecipientInfo, _ := addRecipient(*recipientInfo, token, addRecipientURL)

			paymentInfo := &PaymentInfo{Amount: 10.5, Currency: "GBP", RecipientId: returnedRecipientInfo.Id}
			returnedPaymentInfo, _ := makePaymentToRecipient(*paymentInfo, token, makePaymentURL)

			status, _ := verifyPaymentToRecipient(*returnedPaymentInfo, token, listPaymentsURL)
			if "paid" == status {
				fmt.Println("\nPayment VERIFIED")
			} else {
				fmt.Printf("\nPayment NOT Verified with status:%s\n", status)
			}
		}
	*/
}
