package auth

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var db *sql.DB
var err error

func initializeDatabase() {
	// Initialize the database connection and create the table in the init function
	datastore := "auth.db"
	db, err = makeDBConnection(datastore)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create the table "auth_errors"
	// fixthis: unable to use the columns created - no such columns
	err = createTable("auth_errors", []string{"host TEXT", "block_until INTEGER"})
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func makeDBConnection(datastore string) (*sql.DB, error) {
	// Check if datastore path is valid
	absPath, err := filepath.Abs(datastore)
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
	return fmt.Sprintf("%s", columns)
}

func getRecord(host string) (int64, error) {
	var blockUntil int64
	query := "SELECT block_until FROM auth_errors WHERE host = ?"
	err := db.QueryRow(query, host).Scan(&blockUntil)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // No record found
		}
		log.Printf("Failed to get blocked records for host [%s] from auth_errors table - %s", host, err)
		return 0, err
	}
	return blockUntil, nil
}

func putRecord(host string, blockUntil int64) {
	query := "INSERT INTO auth_errors (host, block_until) VALUES (?, ?)"
	_, err := db.Exec(query, host, blockUntil)
	log.Printf("Failed to put block_until [%d] for host [%s] in auth_errors table - %s", blockUntil, host, err)
}

func removeRecord(host string) {
	query := "DELETE FROM auth_errors WHERE host = ?"
	_, err := db.Exec(query, host)
	log.Printf("Failed to remove host [%s] from auth_errors table - %s", host, err)
}
