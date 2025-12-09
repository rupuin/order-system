package main

import (
	"log"
	"net/http"
	"svc-order/async"
	"svc-order/handlers"
	"svc-order/persistence"

	"github.com/gorilla/mux"
)

func main() {
	log.Printf("Starting svc-order...")

	orderProducer := async.NewProducer(async.GetBrokers(), "order_events")
	defer orderProducer.Close()

	repo := persistence.NewRepository()
	orderHandler := handlers.NewOrderHandler(orderProducer, repo)

	router := mux.NewRouter()

	router.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")

	log.Printf("svc-order ready on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
