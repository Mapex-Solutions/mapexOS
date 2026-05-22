package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	mockmqtt "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mock_servers/mqtt"
	mocknats "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mock_servers/nats"
	mockrabbitmq "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mock_servers/rabbitmq"
	mocksmtp "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mock_servers/smtp"
)

// Mock server ports — must match config.sh
const (
	httpPort     = 9999
	mqttPort     = 1884
	natsPort     = 4333
	rabbitmqPort = 5673
	smtpPort     = 2525
)

func main() {
	var cleanups []func()

	// MQTT
	cleanupMQTT, err := mockmqtt.ForBenchmark(mqttPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "MQTT: %v\n", err)
		os.Exit(1)
	}
	cleanups = append(cleanups, cleanupMQTT)
	fmt.Printf("  MQTT    :%d\n", mqttPort)

	// NATS (executor mock — NOT the JetStream used for message consumption)
	cleanupNATS, err := mocknats.ForBenchmark(natsPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NATS: %v\n", err)
		os.Exit(1)
	}
	cleanups = append(cleanups, cleanupNATS)
	fmt.Printf("  NATS    :%d\n", natsPort)

	// RabbitMQ (AMQP 0-9-1)
	cleanupRabbitMQ, err := mockrabbitmq.ForBenchmark(rabbitmqPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "RabbitMQ: %v\n", err)
		os.Exit(1)
	}
	cleanups = append(cleanups, cleanupRabbitMQ)
	fmt.Printf("  RabbitMQ:%d\n", rabbitmqPort)

	// SMTP
	cleanupSMTP, err := mocksmtp.ForBenchmark(smtpPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SMTP: %v\n", err)
		os.Exit(1)
	}
	cleanups = append(cleanups, cleanupSMTP)
	fmt.Printf("  SMTP    :%d\n", smtpPort)

	// HTTP (inline — handles webhooks for http, teams, slack executors)
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			io.Copy(io.Discard, r.Body)
			r.Body.Close()
			w.WriteHeader(200)
		})
		if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP: %v\n", err)
			os.Exit(1)
		}
	}()
	fmt.Printf("  HTTP    :%d\n", httpPort)

	fmt.Println("  All mock servers running. Kill to stop.")

	// Wait for signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	for _, cleanup := range cleanups {
		cleanup()
	}
}
