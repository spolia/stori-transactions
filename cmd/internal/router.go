package internal

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spolia/stori-transactions/internal/account"
)

type Service interface {
	CreateUser(ctx context.Context, user account.User) error
	SaveAndNotifyMovements(ctx context.Context, movements []account.Movements, alias string) error
}

func API(r *mux.Router, service Service) {
	r.HandleFunc("/movements/notify", saveAndNotifyMovements(service)).Methods(http.MethodPost)
	r.HandleFunc("/users", createUser(service)).Methods(http.MethodPost)
}
