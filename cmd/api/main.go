package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/aja-silason/ledger/internal/application"
	"github.com/aja-silason/ledger/internal/infra/http"
	"github.com/aja-silason/ledger/internal/infra/postgres"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Ledger Financial API
// @version         1.0
// @description     API de Livro Razão (Ledger) para controle financeiro com Contabilidade de Dupla Entrada e Idempotência.
// @host            localhost:8080
//
//	@BasePath        /
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
		log.Fatalf("❌ Falha na autenticação ou conexão com o PostgreSQL: %v", err)
	}
	fmt.Println("✅ Conexão com o PostgreSQL estabelecida com sucesso!")

	r := gin.Default()

	// ledgerRepo := application.NewCreateAccountService()
	postgresRepo := postgres.NewPostgresSQLAccountRepository(db)
	ledgerCreateService := application.NewCreateAccountService(postgresRepo)
	ledgerController := http.NewLedgerController(nil, ledgerCreateService)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/accounts", ledgerController.CreateAccountTransaction)
	r.POST("/create-transation", ledgerController.CreateTransaction)

	fmt.Println("🚀 Servidor iniciado na porta 8080")
	r.Run(":8080")
}
