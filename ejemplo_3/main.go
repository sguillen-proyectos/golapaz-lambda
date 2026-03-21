package main

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"path"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nfnt/resize"
)

func Handler(ctx context.Context, s3Event events.S3Event) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("error loading AWS config: %w", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		fmt.Printf("Processing: bucket=%s key=%s\n", bucket, key)

		result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err != nil {
			return fmt.Errorf("error downloading s3://%s/%s: %w", bucket, key, err)
		}
		defer result.Body.Close()

		img, err := png.Decode(result.Body)
		if err != nil {
			return fmt.Errorf("error decoding PNG: %w", err)
		}

		resizedImg := resize.Resize(200, 200, img, resize.Lanczos3)

		var buf bytes.Buffer
		if err := png.Encode(&buf, resizedImg); err != nil {
			return fmt.Errorf("error encoding resized PNG: %w", err)
		}

		filename := path.Base(key)
		destKey := fmt.Sprintf("resized/%s", filename)
		contentType := "image/png"
		_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      &bucket,
			Key:         &destKey,
			Body:        bytes.NewReader(buf.Bytes()),
			ContentType: &contentType,
		})
		if err != nil {
			return fmt.Errorf("error uploading resized image to s3://%s/%s: %w", bucket, destKey, err)
		}

		fmt.Printf("Resized image uploaded to s3://%s/%s\n", bucket, destKey)
	}

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
