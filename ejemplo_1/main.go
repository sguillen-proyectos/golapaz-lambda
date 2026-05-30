package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler() (string, error) {
	fmt.Printf("Handler 1\n")
	return "handler 1", nil
}

func Handler2(ctx context.Context, event map[string]string) error {
	fmt.Printf("Handler 2\n")
	for k, v := range event {
		fmt.Printf("%v = %v\n", k, v)
	}
	return nil
}

func main() {
	runLocally := os.Getenv("RUN_LOCALLY")
	fmt.Println("Hola comunidad Go!")

	if runLocally != "" {
		Handler()
		return
	}
	lambda.Start(Handler2)
}
