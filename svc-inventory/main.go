package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"svc-inventory/async"
	"svc-inventory/handlers"
	"sync"
	"syscall"
)

func main() {
	log.Printf("Starting svc-inventory...")

	ctx, cancel := setupShudown()
	defer cancel()

	var wg sync.WaitGroup

	startOrderConsumer(ctx, &wg)

	waitForShudownSignal()
	log.Println("shudown signal received, stopping consumer")

	cancel()

	wg.Wait()
}

func setupShudown() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func waitForShudownSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

func startOrderConsumer(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		consumer := async.NewConsumer(
			async.GetBrokers(),
			"order_events",
			"inventory-group",
		)
		defer consumer.Close()

		log.Println("Starting Order consumer...")
		consumer.ProcessMessages(ctx, handlers.HandleOrderCreated)
		log.Println("Order Consumer stopped")
	}()

}
