package application

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/aja-silason/ledger/internal/infra/postgres"
	"github.com/aja-silason/ledger/internal/security"
	"github.com/google/uuid"
)

type WidrawCardLess struct {
	repo        *postgres.WithdrawRepository
	accountRepo *postgres.AccountRepository
}

func NewWidrawCardLess(
	repo *postgres.WithdrawRepository,
	accountRepo *postgres.AccountRepository,
) *WidrawCardLess {
	return &WidrawCardLess{
		repo:        repo,
		accountRepo: accountRepo,
	}
}

func (w *WidrawCardLess) DemandWidraw(ctx context.Context, input *WithDrawInput) (domain.SuccessMessage, error) {
	account, err := w.accountRepo.FindByID(input.AccountID)
	if err != nil {
		return nil, err
	}

	referenceCode, err := security.Generate8DigitCode()
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar código: %w", err)
	}

	code := strconv.Itoa(int(input.SecretCode))
	hashedCode, err := security.HashCode(code)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar hash do código: %w", err)
	}

	withdrawID := uuid.New()
	now := time.Now().UTC()
	durationMinutes := (60 * 8)
	withdraw := &domain.Withdraw{
		ID:        withdrawID,
		AccountID: account.ID,
		Amount:    input.Amount,
		Status:    domain.PENDING,
		Code:      referenceCode,
		CodeHash:  hashedCode,
		ExpiresAt: now.Add(time.Duration(durationMinutes) * time.Minute),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err = withdraw.ValidateAmountToWithDraw(); err != nil {
		return nil, err
	}

	if err = withdraw.ValidateCode(); err != nil {
		return nil, err
	}

	_, err = w.repo.Save(withdraw)
	if err != nil {
		return nil, err
	}

	return domain.NewSuccessMessage("Seu Código de Levantamento: " + referenceCode), nil
}
