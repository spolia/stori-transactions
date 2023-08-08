package account_test

import (
	"context"
	"testing"
	"time"

	"github.com/spolia/stori-transactions/internal/account"
	"github.com/spolia/stori-transactions/internal/account/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_When_SaveAndNotifyMovements_ReturnsOk_Then_OK(t *testing.T) {
	// Given
	movements := []account.Movements{
		{ID: "1", Date: time.Date(2023, time.May, 12, time.Now().Hour(), time.Now().Minute(), 0, 0, time.UTC), Amount: 100, Type: "debit"},
		{ID: "2", Date: time.Date(2023, time.May, 13, time.Now().Hour(), time.Now().Minute(), 0, 0, time.UTC), Amount: 200, Type: "credit"},
		{ID: "3", Date: time.Date(2023, time.April, 12, time.Now().Hour(), time.Now().Minute(), 0, 0, time.UTC), Amount: 300, Type: "debit"}}

	// When
	var repositoryMock repositoryMock
	repositoryMock.On("SaveMovements").Return(nil).Once()
	repositoryMock.On("ExistUser").Return(true)
	repositoryMock.On("GetEmail").Return("jsmith@gmail.com", nil).Once()

	var smtClientMock smtClientMock
	smtClientMock.On("SendEmail").Return(nil).Once()

	serviceAccount := account.New(&smtClientMock, &repositoryMock)

	// Then
	err := serviceAccount.SaveAndNotifyMovements(context.Background(), movements, "jsmith")
	require.NoError(t, err)
}

type repositoryMock struct {
	mock.Mock
}

type smtClientMock struct {
	mock.Mock
}

func (r *repositoryMock) SaveMovements(ctx context.Context, movements []repository.Movements, alias string) error {
	args := r.Called()
	return args.Error(0)
}

func (r *repositoryMock) ExistUser(ctx context.Context, alias string) bool {
	args := r.Called()
	return args.Bool(0)
}

func (r *repositoryMock) GetEmail(ctx context.Context, alias string) (string, error) {
	args := r.Called()
	return args.String(0), args.Error(1)
}

func (r *repositoryMock) SaveUser(ctx context.Context, user repository.User) error {
	args := r.Called()
	return args.Error(0)
}

func (s *smtClientMock) SendEmail(ctx context.Context, email string, movementsByMonth map[string]int, totalBalance float64, avgDebitByMonth, avgCreditByMonth map[string]float64) error {
	args := s.Called()
	return args.Error(0)
}
