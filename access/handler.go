package access

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"

	_ "github.com/rk0cc-xyz/gaf/storage"
)

var handler_instance *MySQLHandler = nil

// Implemented DatabaseFieldHandler in MySQL environment.
type MySQLHandler struct {
	sql_conf mysql.Config
	db       *sql.DB
}

// Get current instance of the handler.
func GetMySQLHandlerInstance() (*MySQLHandler, error) {
	if handler_instance == nil {
		conf, conferr := getMySQLConfigFromEnv()
		if conferr != nil {
			return nil, conferr
		}

		handler_instance = &MySQLHandler{
			sql_conf: *conf,
			db:       nil,
		}
	}

	return handler_instance, nil
}

// Open the connection between SQL and this package.
//
// If opened already, it just return same instance of the database.
func (msh MySQLHandler) OpenSQL() (*sql.DB, error) {
	db_closed := false

	if msh.db == nil || msh.db.Ping() != nil {
		db_closed = true
	}

	if db_closed {
		odb, odberr := sql.Open("mysql", msh.sql_conf.FormatDSN())
		if odberr != nil {
			return nil, odberr
		}
		msh.db = odb
	}

	return msh.db, nil
}

// Close current SQL connection.
func (msh MySQLHandler) CloseCurrentSQL() error {
	if msh.db != nil {
		return msh.db.Close()
	}

	return nil
}

func (msh MySQLHandler) WriteToDB(page int64, content []byte, updatedAt string) error {
	cdb, cdberr := msh.OpenSQL()
	if cdberr != nil {
		return cdberr
	}

	stmt, stmterr := cdb.Prepare("INSERT INTO `REPOSITORY_CONTENT` (`PAGE`, `CONTENT`, `UPDATED_AT`) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE `CONTENT` = ?, `UPDATED_AT` = ?")
	if stmterr != nil {
		return stmterr
	}

	if _, stmtexecerr := stmt.Exec(page, content, updatedAt, content, updatedAt); stmtexecerr != nil {
		return stmtexecerr
	}

	stmt.Close()

	return nil
}

func (msh MySQLHandler) ReadFromDB(page int64) ([]byte, *string, error) {
	cdb, cdberr := msh.OpenSQL()
	if cdberr != nil {
		return nil, nil, cdberr
	}

	var content []byte
	var updatedAt string

	qerr := cdb.QueryRow("SELECT `CONTENT`, `UPDATED_AT` FROM `REPOSITORY_CONTENT` WHERE `PAGE` = ?", page).Scan(&content, &updatedAt)

	if qerr != nil {
		return nil, nil, qerr
	}

	return content, &updatedAt, nil
}

// Receive current maximum page of the context.
func (msh MySQLHandler) GetMaxPage() (*int64, error) {
	cdb, cdberr := msh.OpenSQL()
	if cdberr != nil {
		return nil, cdberr
	}

	var maxPage int64

	qerr := cdb.QueryRow("SELECT `PAGE` FROM `REPOSITORY_CONTENT` ORDER BY `PAGE` DESC LIMIT 1").Scan(&maxPage)

	if qerr != nil {
		return nil, qerr
	}

	return &maxPage, nil
}

// Clear redundent page by giving the maximum page from fetch.
func (msh MySQLHandler) ClearExtraPages(fetchMaxPage int64) error {
	cdb, cdberr := msh.OpenSQL()
	if cdberr != nil {
		return cdberr
	}

	_, cerr := cdb.Exec("DELETE FROM `REPOSITORY_CONTENT` WHERE `PAGE` > ?", fetchMaxPage)
	if cerr != nil {
		return cerr
	}

	return nil
}
