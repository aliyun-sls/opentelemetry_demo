package oteldemo

import (
	"io/ioutil"
	"log"
	"net/http"
)

func ToApp01() {
	//service := os.Getenv("SERVICE_NAME")
	url := "http://goapp01-svc-hqtwd.default.svc.cluster.local:8080/getapp01"

	resp, err := http.Get(url)
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
