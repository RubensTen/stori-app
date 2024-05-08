package repositories

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"stori-app/helpers"
	"stori-app/models"
)

type DBRepository struct {
	dbInstance *helpers.DBConnection
}

func NewDBRepository(dbInstance *helpers.DBConnection) *DBRepository {
	return &DBRepository{dbInstance}
}

const (
	clientCollection      = "client"
	accountCollection     = "account"
	transactionCollection = "transaction"
)

func (repo *DBRepository) GetClientByEmail(email string) (*models.ClientModel, error) {
	var result models.ClientModel
	filter := bson.D{{Key: "email", Value: email}}
	if repo.dbInstance == nil || repo.dbInstance.Db == nil {
		fmt.Print("error repository instance is nil, throw a panic error")
	}
	err := repo.dbInstance.Db.Collection(clientCollection).FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (repo *DBRepository) CreateClient(name, email string) (*models.ClientModel, error) {
	// Verificar si el cliente ya existe.
	clientExist, err := repo.GetClientByEmail(email)
	if err != nil {
		return nil, err
	}

	if clientExist != nil {
		return clientExist, nil // Si ya existe, retornamos el existente.
	}

	client := &models.ClientModel{
		Name:  name,
		Email: email,
	}

	result, err := repo.dbInstance.Db.Collection(clientCollection).InsertOne(context.TODO(), client)
	if err != nil {
		return nil, err
	}

	if result != nil {
		clientId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, errors.New("error to parse client id to primitive object")
		}
		client.Id = clientId
	} else {
		return nil, errors.New("error client id not provided from insertone")
	}

	return client, nil
}

func (repo *DBRepository) GetAccountByIdAndType(clientId primitive.ObjectID, accountType string) (*models.AccountModel, error) {
	query := bson.D{{Key: "client_id", Value: clientId}, {Key: "type", Value: accountType}}
	var account models.AccountModel
	err := repo.dbInstance.Db.Collection(accountCollection).FindOne(context.TODO(), query).Decode(&account)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (repo *DBRepository) CreateOrUpdateAccount(account models.AccountModel) (*models.AccountModel, error) {
	accountExist, err := repo.GetAccountByIdAndType(account.ClientID, account.Type)
	if err != nil {
		return nil, err
	}

	// if the account exist the only update the balance
	if accountExist != nil {
		query := bson.D{{Key: "client_id", Value: account.ClientID}, {Key: "type", Value: account.Type}}
		update := bson.M{
			"$set": bson.M{
				"balance": account.Balance,
			},
		}
		_, err := repo.dbInstance.Db.Collection(accountCollection).UpdateOne(context.TODO(), query, update)
		if err != nil {
			return nil, err
		}
		accountExist.Balance = account.Balance
		return accountExist, nil
	}

	result, err := repo.dbInstance.Db.Collection(accountCollection).InsertOne(context.TODO(), account)
	if err != nil {
		return nil, err
	}

	if result != nil {
		accountId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, errors.New("error to parse accountId id to primitive object")
		}
		account.Id = accountId
	} else {
		return nil, errors.New("error accountId id not provided from insertone")
	}

	return &account, nil
}

func (repo *DBRepository) CreateTransaction(transactions []models.TransactionModel) error {
	var newTransactions []interface{}
	for _, transaction := range transactions {
		newTransactions = append(newTransactions, transaction)
	}
	_, err := repo.dbInstance.Db.Collection(transactionCollection).InsertMany(context.TODO(), newTransactions)
	if err != nil {
		return err
	}
	return nil
}
