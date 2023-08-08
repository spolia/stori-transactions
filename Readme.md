# Stori-Transactions

![technology Go](https://img.shields.io/badge/technology-go-blue.svg)

## Overview

This project implements a system that processes a file from a mounted directory that contain a list of debit and credit transactions of an account; then send summary information to the user in the form of an email. 
The summary email contains information on the total balance of the account, the number of
transactions grouped by month, and the average credit and average debit amounts grouped by month.

## App structure

This application was developed in `golang` , `gorilla/mux` , `package oriented design` as project structure and `mysql` as database.

The project is composed of two main components:

- `users` : manage users creations.
- `movements`: manage the user account movements like debit and credit transactions. 
- [check](/migrations/init.sql) the defined sql scheme .

## Endpoints

- `POST /movements/notify` : save user movements and notify through email
- `POST /users` : User registration. Users with the same alias are not allowed.

## How to run this project

1. Make sure you have already installed both Docker Engine and Docker Compose in the last version.
2. In order to send a mail you need to set your email and the password of this. Check [this guide](https://www.getmailbird.com/gmail-app-password/) first in order to generate a Gmail App Password from your account, this will allow you to use an encrypted password instead of the real one. 
3. Set your email and your encrypted password in the environment variable file (.env file)
4. Type `make build` to build the docker compose and then `make up` to up the compose.


## Examples 
```
curl --location --request POST 'localhost:8080/movements/notify' \
--form 'transactions=@"/stori-transactions/transactions.csv"' \
--form 'alias="jsmith"'
```

```
 curl --location --request POST 'localhost:8080/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "alias": "jsmith",
    "firstname": "Jane",
    "lastname": "Smith",
    "email": "jsmith@gmail.com",
    "password": "password2"

}' 
```

# Test

- Type `make test` to run the unit tests.
- Two different csv files are provided. 