package main

import "banner/internal/worker/rabbitmq"

func main() {
	w := rabbitmq.New()
	w.MustRun()
}
