package main

import (
	"context"
	"database/sql"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ruhancs/ms-wallet-go/internal/database"
	"github.com/ruhancs/ms-wallet-go/internal/event"
	"github.com/ruhancs/ms-wallet-go/internal/event/handler"
	createaccount "github.com/ruhancs/ms-wallet-go/internal/usecase/create_account"
	createclient "github.com/ruhancs/ms-wallet-go/internal/usecase/create_client"
	createtransaction "github.com/ruhancs/ms-wallet-go/internal/usecase/create_transaction"
	"github.com/ruhancs/ms-wallet-go/internal/web"
	"github.com/ruhancs/ms-wallet-go/internal/web/webserver"
	"github.com/ruhancs/ms-wallet-go/pkg/events"
	"github.com/ruhancs/ms-wallet-go/pkg/kafka"
	unitofwork "github.com/ruhancs/ms-wallet-go/pkg/unit_of_work"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//conexao com kafka
	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id": "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	//quando chamar o evento transactioncreated o handler sera executado
	eventDispatcher.Register("TransacionCreated",handler.NewTransactionCreatedKafaka(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()

	balanceUpdatedEvent := event.NewBalanceUpdated()
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))

	clientDb := database.NewClientDb(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	unitOfWork := unitofwork.NewUnitOfWork(ctx,db)
	
	//registrando os repositorios no unitofwork
	unitOfWork.Register("AccountDb", func (tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})
	unitOfWork.Register("TransactionDB", func (tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := createclient.NewCreateClientUseCase(clientDb)
	createAccountUseCase := createaccount.NewCreateAccountUseCase(accountDb,clientDb)
	createTransactionUseCase := createtransaction.NewCreateTransactionUseCase(unitOfWork,eventDispatcher,transactionCreatedEvent, balanceUpdatedEvent)

	webserver := webserver.NewWebServer(":8080")//criando o webserver que ira receber os handlers com os usecases
	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	//rotas
	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAcount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server runnig on port 8080")
	webserver.Start()
	
}