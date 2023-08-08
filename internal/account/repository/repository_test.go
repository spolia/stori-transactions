package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/spolia/stori-transactions/internal/account/repository"
	"github.com/stretchr/testify/require"
)

func TestSaveMovement_When_SaveIsOK_Returns_ok(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repo := repository.New(db)
	defer db.Close()
	movements := []repository.Movements{
		{ID: "1", Date: "2021-01-01", Amount: 100, Type: "debit"},
		{ID: "2", Date: "2021-01-02", Amount: 200, Type: "credit"},
		{ID: "3", Date: "2021-01-03", Amount: 300, Type: "debit"}}
	var alias = "spolia"

	// When
	mock.ExpectBegin()
	query1 := "INSERT INTO movements (date, type_movement, amount, alias) VALUES (?, ?, ?, ?);"
	mock.ExpectPrepare(query1)
	mock.ExpectExec(query1).WithArgs(movements[0].Date, movements[0].Type, movements[0].Amount, alias).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(query1).WithArgs(movements[1].Date, movements[1].Type, movements[1].Amount, alias).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(query1).WithArgs(movements[2].Date, movements[2].Type, movements[2].Amount, alias).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	// then
	err = repo.SaveMovements(context.Background(), movements, alias)
	require.NoError(t, err)
}
