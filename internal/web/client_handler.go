package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	createclient "github.com/ruhancs/ms-wallet-go/internal/usecase/create_client"
)

type WebClientHandler struct {
	CreatClientUseCase createclient.CreateClientUseCase
}

func NewWebClientHandler(createClientUseCase createclient.CreateClientUseCase) *WebClientHandler {
	return &WebClientHandler{
		CreatClientUseCase: createClientUseCase,
	}
}

func (handle *WebClientHandler) CreateClient(res http.ResponseWriter, req *http.Request) {
	fmt.Println("ENTROU NO HANDLER")
	var dto createclient.CreateClientInputDto
	//pegar dados do request e parcear para json e inseri no dto
	err := json.NewDecoder(req.Body).Decode(&dto)
	if  err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	output,err := handle.CreatClientUseCase.Execute(dto)
	fmt.Println(err)
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