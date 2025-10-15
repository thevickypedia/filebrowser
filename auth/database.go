package auth

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var db *sql.DB
var err error

var authDB = "auth.db"

var authErrTable = "auth_errors"
var authErrColumns = []string{"host TEXT", "block_until INTEGER"}

var tokenTracker = "token_tracker"
var tokenTrackerColumns = []string{"token TEXT"}

func initializeDatabase() {
	// Initialize the database connection and create the table in the init function
	db, err = makeDBConnection()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	err = createTable(authErrTable, authErrColumns)
	if err != nil {
		log.Fatalf("Failed to create table [%s]: %v", authErrTable, err)
	}

	err = createTable(tokenTracker, tokenTrackerColumns)
	if err != nil {
		log.Fatalf("Failed to create table [%s]: %v", tokenTracker, err)
	}
}

func makeDBConnection() (*sql.DB, error) {
	// Check if datastore path is valid
	absPath, err := filepath.Abs(authDB)
	if err != nil {
		return nil, err
	}

	// Ensure the directory exists
	dir := filepath.Dir(absPath)
	errTmp := os.MkdirAll(dir, 0755)
	if errTmp != nil {
		return nil, errTmp
	}

	// Open database connection
	dbTmp, errTmp := sql.Open("sqlite3", absPath)
	if errTmp != nil {
		return nil, errTmp
	}

	// Set connection timeout
	// db.SetConnMaxLifetime(0)
	// db.SetMaxOpenConns(1)
	// db.SetMaxIdleConns(1)

	return dbTmp, nil
}

func createTable(tableName string, columns []string) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, joinColumns(columns))
	_, err := db.Exec(query)
	return err
}

func joinColumns(columns []string) string {
	return strings.Join(columns, ", ")
}

func getForbiddenRecord(host string) (int64, error) {
	var blockUntil int64
	query := fmt.Sprintf("SELECT block_until FROM %s WHERE host = ?", authErrTable) //nolint:gosec
	err := db.QueryRow(query, host).Scan(&blockUntil)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // No record found
		}
		log.Printf("Warning: Failed to get blocked records for host [%s] from %s table - %s", host, authErrTable, err)
		return 0, err
	}
	return blockUntil, nil
}

func putForbiddenRecord(host string, blockUntil int64) {
	query := fmt.Sprintf("INSERT INTO %s (host, block_until) VALUES (?, ?)", authErrTable) //nolint:gosec
	_, err := db.Exec(query, host, blockUntil)
	if err != nil {
		log.Printf("Warning: Failed to put block_until [%d] for host [%s] in %s table - %s", blockUntil, host, authErrTable, err)
	}
}

func removeForbiddenRecord(host string) {
	query := fmt.Sprintf("DELETE FROM %s WHERE host = ?", authErrTable) //nolint:gosec
	_, err := db.Exec(query, host)
	if err != nil {
		log.Printf("Warning: Failed to remove host [%s] from %s table - %s", host, authErrTable, err)
	}
}

func GetAllowedJWT() []string {
	query := fmt.Sprintf("SELECT * FROM %s", tokenTracker) //nolint:gosec
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Warning: Failed to get allowed JWTs from %s table - %s", tokenTracker, err)
		return []string{}
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		rowErr := rows.Scan(&token)
		if rowErr != nil {
			log.Printf("Warning: Failed to scan token from %s table - %s", tokenTracker, rowErr)
			continue
		}
		tokens = append(tokens, token)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Warning: Error occurred during rows iteration from %s table - %s", tokenTracker, err)
	}

	return tokens
}

func PutAllowedJWT(token string) error {
	query := fmt.Sprintf("INSERT INTO %s (token) VALUES (?)", tokenTracker) //nolint:gosec
	_, err := db.Exec(query, token)
	if err != nil {
		log.Printf("Warning: Failed to put token in %s: %v", tokenTracker, err)
	}
	return err
}

func RemoveAllowedJWT(token string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE token = ?", tokenTracker) //nolint:gosec
	_, err := db.Exec(query, token)
	if err != nil {
		log.Printf("Warning: Failed to remove token from %s: %v", tokenTracker, err)
	}
	return err
}
