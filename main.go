package main

// Loaded in auth/database.go with sql.Open
// goimports -local github.com/thevickypedia/filebrowser -w .
import (
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/thevickypedia/filebrowser/v2/auth"
	"github.com/thevickypedia/filebrowser/v2/cmd"
)

func main() {
	auth.DataBase()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
