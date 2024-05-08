package helpers

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getS3Client(ctx context.Context) (*s3.Client, error) {
	//sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithClientLogMode(aws.LogRequestWithBody | aws.LogResponseWithBody))
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load default config: %s", err)
		return nil, err
	}
	return s3.NewFromConfig(sdkConfig), nil
}

func GetBucketObject(ctx context.Context, s3Event events.S3Event) (*s3.GetObjectOutput, error) {
	var object *s3.GetObjectOutput
	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.URLDecodedKey
		client, err := getS3Client(ctx)
		if err != nil {
			return nil, err
		}
		object, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err != nil {
			log.Printf("error getting object %s/%s: %s", bucket, key, err)
			return nil, err
		}
	}
	return object, nil
}
