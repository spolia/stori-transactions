package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/spolia/stori-transactions/cmd/internal"
	"github.com/spolia/stori-transactions/internal/account"
	"github.com/spolia/stori-transactions/internal/account/mail"
	"github.com/spolia/stori-transactions/internal/account/repository"
)

func main() {
	log.Println("starting")

	db, err := sql.Open("mysql", "tester:secret@tcp(db:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	// SMTP configuration
	sender := os.Getenv("EMAIL")
	password := os.Getenv("EncryptedPassword")
	smtpServer := "smtp.gmail.com"
	smtpPort := 587
	log.Println("sender: ", sender)
	log.Println("password: ", password)

	mailClientConf := mail.Configuration{
		Sender:     sender,
		Password:   password,
		SmtpServer: smtpServer,
		SmtpPort:   smtpPort,
	}

	accountService := account.New(mail.New(mailClientConf), repository.New(db))
	router := mux.NewRouter()
	internal.API(router, accountService)
	// localhost:8080
	http.ListenAndServe(":8080", router)
}
