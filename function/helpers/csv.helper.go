package helpers

import (
	"context"
	"encoding/csv"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func GetCSVDataFromEvent(ctx context.Context, s3Event events.S3Event) ([][]string, error) {
	s3Object, err := GetBucketObject(ctx, s3Event)
	if err != nil {
		log.Printf("error getting object %s", err.Error())
		return nil, err
	}
	defer s3Object.Body.Close()
	// read content object
	content := csv.NewReader(s3Object.Body)
	content.FieldsPerRecord = -1
	data, err := content.ReadAll()
	if err != nil {
		log.Printf("error getting data from s3object %s", err.Error())
		return nil, err
	}
	return data, nil
}
