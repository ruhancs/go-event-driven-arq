package web

import (
	"encoding/json"
	"net/http"

	createaccount "github.com/ruhancs/ms-wallet-go/internal/usecase/create_account"
)

type WebAccountHandler struct {
	AccountUseCase createaccount.CreateAccountUseCase
}

func NewWebAccountHandler(accountUseCase createaccount.CreateAccountUseCase) *WebAccountHandler {
	return &WebAccountHandler{
		AccountUseCase: accountUseCase,
	}
}

func (handler *WebAccountHandler) CreateAcount(res http.ResponseWriter, req *http.Request) {
	var dto createaccount.CreateAccountInputDto
	//pegar dados do request e parcear para json e inseri no dto
	err := json.NewDecoder(req.Body).Decode(&dto)
	if  err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	output,err := handler.AccountUseCase.Execute(dto)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	
	//encodar o output
	err = json.NewEncoder(res).Encode(output)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	res.WriteHeader(http.StatusCreated)
}