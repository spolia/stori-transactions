package account

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spolia/stori-transactions/internal/account/repository"
)

type Movements struct {
	ID     string
	Date   time.Time
	Amount float64
	Type   string
}

type User struct {
	Alias     string `json:"alias" validate:"required"`
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

var ErrorUserAlreadyExist = errors.New("already exist")

type Repository interface {
	SaveMovements(ctx context.Context, movements []repository.Movements, alias string) error
	ExistUser(ctx context.Context, alias string) bool
	GetEmail(ctx context.Context, alias string) (string, error)
	SaveUser(ctx context.Context, user repository.User) error
}

type SMTPClient interface {
	SendEmail(ctx context.Context, email string, movementsByMonth map[string]int, totalBalance float64, avgDebitByMonth, avgCreditByMonth map[string]float64) error
}

type Service struct {
	smtpClient SMTPClient
	repository Repository
}

func New(smtClient SMTPClient, repo Repository) *Service {
	return &Service{
		smtpClient: smtClient,
		repository: repo,
	}
}

// SaveAndNotifyMovements saves the movements and sends an email with the movements by month, total balance and average amounts by month
func (s *Service) SaveAndNotifyMovements(ctx context.Context, movements []Movements, alias string) error {
	//first check if the user exist
	if !s.repository.ExistUser(ctx, alias) {
		return fmt.Errorf("user not found %s", alias)
	}

	var movementsByMonth = make(map[string]int, 12)
	var avgDebitByMonth, avgCreditByMonth = make(map[string]float64, 12), make(map[string]float64, 12)
	var totalBalance float64
	var repositoryTrx = make([]repository.Movements, 0, len(movements))

	for _, t := range movements {
		// Update total balance
		// Update number of movements and average amounts by month
		if t.Type == "credit" {
			totalBalance += t.Amount
			avgCreditByMonth[t.Date.Month().String()] += t.Amount
		} else if t.Type == "debit" {
			totalBalance -= t.Amount
			avgDebitByMonth[t.Date.Month().String()] += t.Amount
		}

		movementsByMonth[t.Date.Month().String()]++
		repositoryTrx = append(repositoryTrx, repository.Movements{ID: t.ID, Date: t.Date.String(), Amount: t.Amount, Type: t.Type})
	}

	// Calculate average amounts by month
	for month := range movementsByMonth {
		avgCreditByMonth[month] = avgCreditByMonth[month] / float64(movementsByMonth[month])
		avgDebitByMonth[month] = avgDebitByMonth[month] / float64(movementsByMonth[month])
	}

	// Save movements
	if err := s.repository.SaveMovements(ctx, repositoryTrx, alias); err != nil {
		return err
	}

	var email string
	var err error
	// Get the user email
	if email, err = s.repository.GetEmail(ctx, alias); err != nil {
		return err
	}
	// Send email
	if err = s.smtpClient.SendEmail(ctx, email, movementsByMonth, totalBalance, avgDebitByMonth, avgCreditByMonth); err != nil {
		return err
	}

	return nil
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, user User) error {
	return s.repository.SaveUser(ctx, repository.User{
		Alias:     user.Alias,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	})
}
