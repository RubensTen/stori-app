package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"stori-app/helpers"
	"stori-app/repositories"
	"stori-app/service"
)

func HandleRequest(ctx context.Context, s3Event events.S3Event) (string, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Lambda Started")

	err := godotenv.Load(".env")
	if err != nil {
		logger.Warn(".env file not found", zap.String("error", err.Error()))
	}
	
	db := &helpers.DBConnection{}
	db, dbErr := db.GetInstanceDB()
	if dbErr != nil {
		logger.Error("Error to get instance DB", zap.String("error", dbErr.Error()))
		return "error to get the instance of db", dbErr
	} else {
		logger.Sugar().Info("DB connection succesfully")
	}

	defer func() {
		logger.Info("Executing Db.CloseConnection")
		if closeErr := db.CloseConnection(); closeErr != nil {
			logger.Error("Error to CloseConnection:", zap.String("error", closeErr.Error()))
		} else {
			logger.Info("Db connection closed succesfully")
		}
	}()

	data, err := helpers.GetCSVDataFromEvent(ctx, s3Event)
	if err != nil {
		return "error to read data from csv", err
	}

	repository := repositories.NewDBRepository(db)
	transactions := helpers.MapToTransactions(data)
	err = service.SendEmailsAndStore(transactions, repository)
	if err != nil {
		logger.Error(fmt.Sprintf("Error SendEmailsAndStore %s", err.Error()))
		return "error to send email or store account info", err
	}
	logger.Info(fmt.Sprintf("Process Successfully %s", time.DateTime))
	return "HandleRequest ending successfully", nil
}

func main() {
	lambda.Start(HandleRequest)
}
