package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/spolia/stori-transactions/internal/account"
	"github.com/spolia/stori-transactions/internal/account/repository"
)

var validate = validator.New()

// create user
func createUser(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userRequest account.User
		if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validate.Struct(userRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userRequest.Alias = strings.ToLower(userRequest.Alias)
		err := service.CreateUser(r.Context(), userRequest)
		if err != nil {
			if err == repository.ErrorAlreadyExist {
				http.Error(w, "alias already exist", http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("ok")
	}
}

// save and notify movements from a csv file
func saveAndNotifyMovements(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.FormValue("alias")
		if alias == "" {
			http.Error(w, "alias is required", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("transactions")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		out, err := os.Create("/tmp/uploadedfile.csv")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer out.Close()

		// write the content from POST to the file
		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var movements []account.Movements
		movements, err = readFile("/tmp/uploadedfile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		if err = service.SaveAndNotifyMovements(r.Context(), movements, alias); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode("ok")
	}

}

func readFile(fileName string) ([]account.Movements, error) {
	// Open the CSV file
	file, err := os.Open(fmt.Sprintf("%s.csv", fileName))
	if err != nil {
		fmt.Println("opening file: ", err)
		return nil, err
	}

	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the CSV records one by one
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading records: %v", err)
	}

	// Slice to hold all transactions
	transactions := make([]account.Movements, 0, len(records))
	//check columns names
	checkColumnNames(records[0])
	records = records[1:]

	// Process each record and convert amount to float
	for _, record := range records {
		// parse date string to time.Time
		date, err := parseDate(record[1])
		if err != nil {
			return nil, fmt.Errorf("parsing date:%s %v", date, err)
		}

		amount, trxType, err := parseAmountandType(record[2])
		if err != nil {
			return nil, fmt.Errorf("parsing amount: %v", err)
		}

		transactions = append(transactions, account.Movements{
			ID:     record[0],
			Date:   date,
			Amount: amount,
			Type:   trxType,
		})
	}

	return transactions, nil
}

// parseDate parse the date string to time.Time
func parseDate(dateStr string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

// checkColumnNames check if the column names are correct
func checkColumnNames(columNames []string) error {
	if columNames[0] != "id" {
		return fmt.Errorf("column Name 0 must be id")
	}
	if columNames[1] != "date" {
		return fmt.Errorf("column Name 1 must be date")
	}
	if columNames[2] != "amount" {
		return fmt.Errorf("column Name 2 must be amount")
	}
	return nil
}

// parseAmountandType parse the amount and type from the amount string
func parseAmountandType(amountStr string) (float64, string, error) {
	var trxType = "credit"
	// Get the '+' or '-' sign from the amount string
	sign := amountStr[:1][0]

	// Check if the sign is '+' or '-'
	if sign != '+' && sign != '-' {
		return 0, "", fmt.Errorf("sign in the amount value is required, %v", amountStr)
	}

	// Check if the sign is '-' then the transaction type is debit
	if sign == '-' {
		trxType = "debit"
	}

	// Convert the amount string to float64
	amount, err := strconv.ParseFloat(amountStr[1:], 64)
	if err != nil {
		return 0, "", fmt.Errorf("amount value is not the right type, %v", err)
	}

	return amount, trxType, nil
}
