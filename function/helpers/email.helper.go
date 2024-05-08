package helpers

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type EmailRequest struct {
	Sender  string
	To      string
	Subject string
	Body    string
	Charset string
}

func SendEmail(request EmailRequest) error {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("failed to load default config: %s", err)
		return err
	}

	if request.Charset == "" {
		request.Charset = "UTF-8"
	}

	dest := &types.Destination{
		ToAddresses: []string{request.To},
	}
	msg := &types.Message{
		Body: &types.Body{
			Html: &types.Content{
				Data:    aws.String(request.Body),
				Charset: aws.String(request.Charset),
			},
		},
		Subject: &types.Content{
			Charset: aws.String(request.Charset),
			Data:    aws.String(request.Subject),
		},
	}
	sei := &ses.SendEmailInput{Destination: dest, Message: msg, Source: &request.Sender}
	client := ses.NewFromConfig(sdkConfig)
	output, err := client.SendEmail(context.TODO(), sei)
	if err != nil {
		log.Printf("SendEmail error: %s", err)
		return err
	}

	log.Printf("Sucessfully Email Sent to: %s with messageid: %s", request.To, *output.MessageId)
	log.Println()

	return nil
}
