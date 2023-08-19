package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	createtransaction "github.com/ruhancs/ms-wallet-go/internal/usecase/create_transaction"
)

type WebTransactionHandler struct {
	CreateTransactionUseCase createtransaction.CreateTransactionUseCase
}

func NewWebTransactionHandler(createTransactionUseCase createtransaction.CreateTransactionUseCase) *WebTransactionHandler {
	return &WebTransactionHandler{
		CreateTransactionUseCase: createTransactionUseCase,
	}
}

func (handler *WebTransactionHandler) CreateTransaction(res http.ResponseWriter, req *http.Request) {
	var dto createtransaction.CreateTransactionInputDto
	//pegar dados do request e parcear para json e inseri no dto
	err := json.NewDecoder(req.Body).Decode(&dto)
	if  err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := req.Context()//pegar o contexto da request

	output,err := handler.CreateTransactionUseCase.Execute(ctx,dto)
	fmt.Println(err)
	if err != nil {
		fmt.Println("ERROR Usecase")
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(err.Error()))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	
	//encodar o output
	err = json.NewEncoder(res).Encode(output)
	if err != nil {
		fmt.Println("ERROR CONVERSAO JSON")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	res.WriteHeader(http.StatusCreated)
}