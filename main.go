package main

// Loaded in auth/database.go with sql.Open
// goimports -local github.com/thevickypedia/filebrowser -w .
import (
	_ "github.com/mattn/go-sqlite3"

	"github.com/thevickypedia/filebrowser/v2/auth"
	"github.com/thevickypedia/filebrowser/v2/cmd"
)

func main() {
	auth.DataBase()
	cmd.Execute()
}
