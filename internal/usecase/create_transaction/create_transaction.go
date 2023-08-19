package createtransaction

import (
	"context"

	"github.com/ruhancs/ms-wallet-go/internal/entity"
	"github.com/ruhancs/ms-wallet-go/internal/gateway"
	"github.com/ruhancs/ms-wallet-go/pkg/events"
	unitofwork "github.com/ruhancs/ms-wallet-go/pkg/unit_of_work"
)

type CreateTransactionInputDto struct {
	AccountIdFrom string `json:"account_id_from"`
	AccountIdTo string	`json:"account_id_to"`
	Amount float64	`json:"amount"`
}

type CreateTransactionOutputDto struct {
	TransactionId string `json:"transaction_id"` 
	AccountIdFrom string `json:"account_id_from"`
	AccountIdTo string `json:"account_id_to"` 
	Amount float64 `json:"amount"`
}

type BalanceUpdatedOutputDto struct {
	AccountIdFrom string `json:"account_id_from"`
	AccountIdTo string `json:"account_id_to"`
	BalanceAccountIdFrom float64 `json:"balance_account_from"`
	BalanceAccountIdTo float64 `json:"balance_account_to"`
}

type CreateTransactionUseCase struct {
	UnitOfWork unitofwork.UnitOfWorkInterface // para utilizar o repositorio(gateway) de transaction e account
	EventDispatcher events.EventDispatcherInterface
	TransactionCreated events.EventInterface
	BalanceUpdated events.EventInterface
}

func NewCreateTransactionUseCase(
	unitOfWork unitofwork.UnitOfWorkInterface,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
	balanceUpdated events.EventInterface,
	) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		UnitOfWork: unitOfWork,
		EventDispatcher: eventDispatcher,
		TransactionCreated: transactionCreated,
		BalanceUpdated: balanceUpdated,
	}
}

//contexto é utilizado para controlar as operacoes ex: se uma operacao demorar xx cancelar a operacao
func (usecase *CreateTransactionUseCase) Execute(context context.Context,input CreateTransactionInputDto) (*CreateTransactionOutputDto, error) {
	output := &CreateTransactionOutputDto{}
	balanceUpdateOutput := &BalanceUpdatedOutputDto{}

	//executa todas operacoes no db de uma vez
	// se algo der errado nna func o rollback sera feito em todas operacoes
	err := usecase.UnitOfWork.Do(context, func(_ *unitofwork.UnitOfWork) error {
		accountRepository := usecase.getAccountRepository(context)//pegar o repositorio de account
		transactionRepository := usecase.getTransactionepository(context)

		accountFrom,err := accountRepository.FindById(input.AccountIdFrom)
		if err != nil {
			return err
		}
		accountTo,err := accountRepository.FindById(input.AccountIdTo)
		if err != nil {
			return err
		}
		transaction,err := entity.NewTransaction(accountFrom,accountTo,input.Amount)
		if err != nil {
			return err
		}
		
		err = accountRepository.UpdateBalance(accountFrom)
		if err != nil {
			return err
		}
		err = accountRepository.UpdateBalance(accountTo)
		if err != nil {
			return err
		}
		
		err = transactionRepository.Create(transaction)
		if err != nil {
			return err
		}

		output.TransactionId = transaction.Id
		output.AccountIdFrom = input.AccountIdFrom
		output.AccountIdTo = input.AccountIdTo
		output.Amount = input.Amount

		balanceUpdateOutput.AccountIdFrom = input.AccountIdFrom
		balanceUpdateOutput.AccountIdTo = input.AccountIdTo
		balanceUpdateOutput.BalanceAccountIdFrom = accountFrom.Balance
		balanceUpdateOutput.BalanceAccountIdTo = accountTo.Balance

		return nil
	}) 
	if err != nil {
		return nil, err
	}

	usecase.TransactionCreated.SetPayload(output)
	usecase.EventDispatcher.Dispatch(usecase.TransactionCreated)//dispara o evento transactioncreated

	usecase.BalanceUpdated.SetPayload(balanceUpdateOutput)//setar a msg do payload
	usecase.EventDispatcher.Dispatch(usecase.BalanceUpdated)//dispara o evento transactioncreated

	return output,nil

}

//chamar a accountGateway atraves do unitofwork
func (usecase *CreateTransactionUseCase) getAccountRepository(context context.Context) gateway.AccountGateway {
	repository,err := usecase.UnitOfWork.GetRepository(context, "AccountDb") //pegar o repositorio de account
	if err != nil {
		panic(err)
	}
	return repository.(gateway.AccountGateway)
}

//chamar a accountGateway atraves do unitofwork
func (usecase *CreateTransactionUseCase) getTransactionepository(context context.Context) gateway.TransactionGateway {
	repository,err := usecase.UnitOfWork.GetRepository(context, "TransactionDB") //pegar o repositorio de account
	if err != nil {
		panic(err)
	}
	return repository.(gateway.TransactionGateway)
}