package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nsqio/go-nsq"
)

type Booking struct {
	Code        string
	Username    string
	Destination string
}

type BookingReq struct {
	Username    string `json:"username"`
	Destination string `json:"destination"`
}

type Response struct {
	Message string
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/booking", RequestBooking).Methods("POST")

	log.Fatal(http.ListenAndServe(":5000", router))

}

func RequestBooking(w http.ResponseWriter, r *http.Request) {
	log.Println("Jalan nih")
	var bookingReq BookingReq
	_ = json.NewDecoder(r.Body).Decode(&bookingReq)

	code := StringWithCharset(5, charset)

	booking := Booking{code, bookingReq.Username, bookingReq.Destination}

	SendMessage(booking)

	response := Response{"Booking Success, Your Booking Code : " + booking.Code}
	json.NewEncoder(w).Encode(response)
}

func SendMessage(booking Booking) {
	config := nsq.NewConfig()
	// p, err := nsq.NewProducer("127.0.0.1:4150", config)
	p, err := nsq.NewProducer("68.183.237.182:4150", config)

	if err != nil {
		failOnError(err, "Failed to create producer")
	}
	msg, err := json.Marshal(booking)

	err = p.Publish("baru_NSQ", []byte(string(msg)))
	if err != nil {
		failOnError(err, "Failed to publish a message")
	}
	log.Println("cek jalan atau tidak")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
