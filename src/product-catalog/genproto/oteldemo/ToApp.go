package oteldemo

import (
	"io/ioutil"
	"log"
	"net/http"
)

func ToApp01() {
	//service := os.Getenv("SERVICE_NAME")
	urlgoapp := "http://goapp:8080/getapp01"
	urlrefuse := "http://refuse:8080/getapp01"

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

	resp, err = http.Get(urlrefuse)
	if err != nil {
		log.Fatalf("Error calling endpoint: %v", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	log.Printf("Response: %s", body)
}
