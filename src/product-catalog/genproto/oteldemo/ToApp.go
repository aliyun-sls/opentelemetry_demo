package oteldemo

import (
	"io/ioutil"
	"log"
	"net/http"
)

func Order() {
	//service := os.Getenv("SERVICE_NAME")
	urlgoapp := "http://order:8080/order"

	resp, err := http.Get(urlgoapp)
	if err != nil {
		log.Fatalf("Error calling endpoint: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	log.Printf("Response: %s", body)
}
