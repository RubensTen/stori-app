package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"stori-app/helpers"
	"stori-app/models"
	"stori-app/repositories"
)

type TransactionAmmount struct {
	Ammount         float64
	TransactionDate string
}

type AccountInfo struct {
	Email               string
	ClientName          string
	MonthlyTransactions map[string]int
	TotalBalance        float64
	DebitAmounts        []TransactionAmmount
	CreditAmounts       []TransactionAmmount
}

const (
	debitAccount  = "DEBIT"
	creditAccount = "CREDIT"
)

func calculateAverage(values []TransactionAmmount) float64 {
	var sum float64
	for _, value := range values {
		sum += float64(value.Ammount)
	}
	if len(values) == 0 {
		return 0.0
	}
	return sum / float64(len(values))
}

// transactions with plus sign are credit transactions
func isCreditTransaction(transaction float64) bool {
	if transaction >= 0 {
		return true
	}
	return false
}

func mapTransactionsToAccount(transactions []helpers.Transactions) map[string]AccountInfo {
	accountsInfo := make(map[string]AccountInfo)
	for _, transaction := range transactions {
		month, err := helpers.GetMonthName(strings.Split(transaction.Date, "/")[0])
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		entry, exist := accountsInfo[transaction.Email]
		// check if the user not exist
		if !exist {
			entry = AccountInfo{
				ClientName:   transaction.Name,
				Email:        transaction.Email,
				TotalBalance: transaction.Transaction,
			}
			// map for months and count
			monthMap := make(map[string]int)
			monthMap[month] = 1
			entry.MonthlyTransactions = monthMap
		} else {
			balance := entry.TotalBalance + transaction.Transaction
			newBalance, err := helpers.FixFloatPrecision(balance)
			if err != nil {
				log.Printf("Can't fix balance Precision, setting balance without parsed")
				entry.TotalBalance = balance
			} else {
				entry.TotalBalance = newBalance
			}
			// check if the month not exist
			if _, exist := entry.MonthlyTransactions[month]; !exist {
				entry.MonthlyTransactions[month] = 1
			} else {
				entry.MonthlyTransactions[month] += 1
			}
		}

		if isCreditTransaction(transaction.Transaction) {
			entry.CreditAmounts = append(entry.CreditAmounts, TransactionAmmount{
				Ammount:         transaction.Transaction,
				TransactionDate: transaction.Date,
			})
		} else {
			entry.DebitAmounts = append(entry.DebitAmounts, TransactionAmmount{
				Ammount:         transaction.Transaction,
				TransactionDate: transaction.Date,
			})
		}

		accountsInfo[transaction.Email] = entry
	}
	return accountsInfo
}

func saveAccountWithTransactions(client models.ClientModel, accountType string, balance float64, transactions []TransactionAmmount, repository *repositories.DBRepository) error {
	account, err := repository.CreateOrUpdateAccount(models.AccountModel{
		ClientID: client.Id,
		Type:     accountType,
		Balance:  balance,
	})
	if err != nil {
		return errors.New(fmt.Sprintf("error to create %s account for client :%s cause: %s", accountType, client.Email, err.Error()))
	}

	var transactionsModel []models.TransactionModel
	for _, transaction := range transactions {
		transactionsModel = append(transactionsModel, models.TransactionModel{
			AccountID:       account.Id,
			Ammount:         transaction.Ammount,
			TransactionDate: transaction.TransactionDate,
		})
	}
	err = repository.CreateTransaction(transactionsModel)
	if err != nil {
		return errors.New(fmt.Sprintf("error to create many %s transactions for client :%s cause: %s", accountType, client.Email, err.Error()))
	}
	return nil
}

func saveData(repository *repositories.DBRepository, accountInfo AccountInfo) error {
	client, err := repository.CreateClient(accountInfo.ClientName, accountInfo.Email)
	if err != nil {
		return err
	}

	// create credit account
	if len(accountInfo.CreditAmounts) > 0 {
		fmt.Printf("creating credit account for %s", accountInfo.Email)
		err := saveAccountWithTransactions(*client, creditAccount, accountInfo.TotalBalance, accountInfo.CreditAmounts, repository)
		if err != nil {
			return err
		}
	}

	// create debit account
	if len(accountInfo.DebitAmounts) > 0 {
		fmt.Printf("creating debit account for %s", accountInfo.Email)
		err := saveAccountWithTransactions(*client, debitAccount, accountInfo.TotalBalance, accountInfo.DebitAmounts, repository)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendEmailsAndStore(transactions []helpers.Transactions, repository *repositories.DBRepository) error {
	if os.Getenv("SENDER_EMAIL") == "" {
		return errors.New("SENDER_EMAIL not found, it's required")
	}
	accountsInfo := mapTransactionsToAccount(transactions)
	for _, request := range accountsInfo {
		log.Printf("Email in progress to: %s", request.Email)
		bodyEmail := helpers.GenerateSummaryTemplate(helpers.SummaryTemplate{
			ClientName:          request.ClientName,
			TotalBalance:        request.TotalBalance,
			MonthlyTransactions: request.MonthlyTransactions,
			AverageDebitAmount:  calculateAverage(request.DebitAmounts),
			AverageCreditAmount: calculateAverage(request.CreditAmounts),
		})
		err := helpers.SendEmail(helpers.EmailRequest{
			Sender:  os.Getenv("SENDER_EMAIL"),
			To:      request.Email,
			Subject: "Stori - Account Balance and Transactions Report",
			Body:    bodyEmail,
		})
		if err != nil {
			log.Printf("error to send email: %s", err.Error())
		}
		err = saveData(repository, request)
		if err != nil {
			log.Printf("error to saveData: %s", err.Error())
		}
	}
	return nil
}
