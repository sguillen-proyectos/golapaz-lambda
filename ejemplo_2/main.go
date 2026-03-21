package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		fmt.Printf("Event Source: %s\n", record.EventSource)
		fmt.Printf("Event Name: %s\n", record.EventName)
		fmt.Printf("Bucket: %s\n", record.S3.Bucket.Name)
		fmt.Printf("Object Key: %s\n", record.S3.Object.Key)
		fmt.Printf("Object Size: %d\n", record.S3.Object.Size)
	}

	rawJSON, err := json.MarshalIndent(s3Event, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing event: %w", err)
	}
	fmt.Printf("Full Event:\n%s\n", string(rawJSON))

	return nil
}

func main() {
	runLocally := os.Getenv("RUN_LOCALLY")

	if runLocally != "" {
		fmt.Println("Running locally - no Lambda runtime detected")
		return
	}
	lambda.Start(Handler)
}
