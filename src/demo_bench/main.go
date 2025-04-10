package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Config struct {
	TargetURL        string
	Rate             int
	Workers          int
	ProbabilityTable map[string]int
}

type TestScenario struct {
	Name   string
	Weight int
	Exec   func() vegeta.Target
}

var (
	cfg        Config
	products   = []string{"0PUK6V6EV0", "1YMWWN1N4O", "2ZYFJ3GM2N", "66VCHSJNUP", "6E92ZMYYFZ"}
	categories = []string{"binoculars", "telescopes", "accessories", "assembly"}
	people     []map[string]interface{}
)

func getEnvI(key string, defaultValue int) int {
	rateStr := os.Getenv(key)

	if rateStr != "" {
		if rate, err := strconv.Atoi(rateStr); err != nil {
			log.Printf("Invalid REQUEST_RATE value: %v. Using default value: %d", err, defaultValue)
			return defaultValue
		} else {
			return rate
		}
	}

	return defaultValue
}

func getEnvStr(key string, defaultValue string) string {
	rateStr := os.Getenv(key)

	if rateStr != "" {
		return rateStr
	}

	return defaultValue
}

func init() {
	flag.StringVar(&cfg.TargetURL, "target", getEnvStr("REQUEST_URL", "http://gateway:8085"), "Target base URL")
	flag.IntVar(&cfg.Rate, "rate", getEnvI("REQUEST_RATE", 1000), "Requests per ms")
	flag.IntVar(&cfg.Workers, "workers", getEnvI("REQUEST_WORKER_NUM", 10), "Number of workers")

	cfg.ProbabilityTable = map[string]int{
		"index":           1,
		"browse_product":  5,
		"recommendations": 5,
		"ads":             5,
		"view_cart":       5,
		"add_to_cart":     5,
		"checkout":        3,
		"checkout_multi":  2,
	}
}

func main() {
	loadTestData()
	attacker := vegeta.NewAttacker(
		vegeta.Client(&http.Client{}),
		vegeta.Workers(uint64(cfg.Workers)),
	)

	targeter := vegeta.NewStaticTargeter(generateTargets()...)
	metrics := &vegeta.Metrics{}

	for res := range attacker.Attack(targeter, vegeta.Rate{Freq: cfg.Rate, Per: time.Millisecond}, 0, "Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	vegeta.NewTextReporter(metrics).Report(os.Stdout)
}

func loadTestData() {
	file, err := os.Open("people.json")
	if err != nil {
		log.Fatalf("Error loading test data: %v", err)
	}
	defer file.Close()
	json.NewDecoder(file).Decode(&people)
}

func generateTargets() []vegeta.Target {
	scenarios := []TestScenario{
		{"index", cfg.ProbabilityTable["index"], index},
		{"browse_product", cfg.ProbabilityTable["browse_product"], browseProduct},
		{"recommendations", cfg.ProbabilityTable["recommendations"], getRecommendations},
		{"ads", cfg.ProbabilityTable["ads"], getAds},
		{"view_cart", cfg.ProbabilityTable["view_cart"], viewCart},
		{"add_to_cart", cfg.ProbabilityTable["add_to_cart"], addToCart},
		{"checkout", cfg.ProbabilityTable["checkout"], checkout},
		{"checkout_multi", cfg.ProbabilityTable["checkout_multi"], checkoutMulti},
	}

	var targets []vegeta.Target
	for _, s := range scenarios {
		for i := 0; i < s.Weight; i++ {
			targets = append(targets, s.Exec())
		}
	}
	return targets
}

func index() vegeta.Target {
	return vegeta.Target{
		Method: "GET",
		URL:    fmt.Sprintf("%s/", cfg.TargetURL),
	}
}

func browseProduct() vegeta.Target {
	return vegeta.Target{
		Method: "GET",
		URL:    fmt.Sprintf("%s/api/products/%s", cfg.TargetURL, randomProduct()),
	}
}

func getRecommendations() vegeta.Target {
	return vegeta.Target{
		Method: "GET",
		URL:    fmt.Sprintf("%s/api/recommendations?productIds=%s", cfg.TargetURL, randomProduct()),
	}
}

func getAds() vegeta.Target {
	return vegeta.Target{
		Method: "GET",
		URL:    fmt.Sprintf("%s/api/data/?contextKeys=%s", cfg.TargetURL, randomCategory()),
	}
}

func viewCart() vegeta.Target {
	return vegeta.Target{
		Method: "GET",
		URL:    fmt.Sprintf("%s/api/cart", cfg.TargetURL),
	}
}

func addToCart() vegeta.Target {
	user := uuid.NewString()
	product := randomProduct()
	return vegeta.Target{
		Method: "POST",
		URL:    fmt.Sprintf("%s/api/cart", cfg.TargetURL),
		Body:   []byte(fmt.Sprintf(`{"item":{"productId":"%s","quantity":%d},"userId":"%s"}`, product, rand.Intn(5)+1, user)),
	}
}

func checkout() vegeta.Target {
	user := uuid.NewString()
	return vegeta.Target{
		Method: "POST",
		URL:    fmt.Sprintf("%s/api/checkout", cfg.TargetURL),
		Body:   generateCheckoutPayload(user, 1),
	}
}

func checkoutMulti() vegeta.Target {
	user := uuid.NewString()
	return vegeta.Target{
		Method: "POST",
		URL:    fmt.Sprintf("%s/api/checkout", cfg.TargetURL),
		Body:   generateCheckoutPayload(user, rand.Intn(3)+2),
	}
}

func randomProduct() string {
	return products[rand.Intn(len(products))]
}

func randomCategory() string {
	return categories[rand.Intn(len(categories))]
}

func generateCheckoutPayload(user string, items int) []byte {
	person := people[rand.Intn(len(people))]
	person["userId"] = user
	jsonData, _ := json.Marshal(person)
	return jsonData
}
