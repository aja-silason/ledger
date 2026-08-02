package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/aja-silason/ledger/internal/application"
	"github.com/aja-silason/ledger/internal/infra/http"
	"github.com/aja-silason/ledger/internal/infra/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

func main() {
	_ = uuid.New()
	_ = decimal.NewFromInt(100)

	dbURL := "postgres://ledger:ledger@localhost:5001/ledger?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Erro na base de dados: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Falha na autenticação ou conexão com o PostgreSQL: %v", err)
	}
	fmt.Println("Conexão com o PostgreSQL estabelecida com sucesso!")

	// Injeções do Banco de dados
	accountRepository := postgres.NewPostgresSQLAccountRepository(db)
	balanceRepository := postgres.NewPostgresSQLBalanceRepository(db)
	transactionsRepository := postgres.NewPostgresSQLTransactionRepository(db)
	entriesRepository := postgres.NewPostgresSQLEntriesRepository(db)
	accountGateway := postgres.NewPostgresSQLAccountGateway(db)
	withDrawRepository := postgres.NewPostgresSQLWithdrawRepository(db)

	// Injeções dos Services
	createAccount := application.NewCreateAccountService(accountRepository, balanceRepository)
	getAccount := application.NewGetAccountServiceFinder(accountGateway)
	depositIn := application.NewDepositIn(balanceRepository, accountRepository, transactionsRepository)
	transferMoney := application.NewTransferMoney(balanceRepository, accountRepository, transactionsRepository, entriesRepository)

	withdrawCardless := application.NewWidrawCardLess(withDrawRepository, accountRepository)
	cancelWithdrawCardless := application.NewCancelWithDrawCardLess(withDrawRepository)
	drawedWithdrawCardless := application.NewDrawedCardLess(withDrawRepository, transactionsRepository, entriesRepository, balanceRepository)

	accountController := http.NewAccountController(createAccount, getAccount, depositIn, transferMoney, withdrawCardless, cancelWithdrawCardless, drawedWithdrawCardless)

	controllers := &http.Controllers{
		AccountController: accountController,
	}

	r := http.SetupRoutes(controllers)

	fmt.Println("Servidor iniciado na porta 8080")
	r.Run(":8080")
}
