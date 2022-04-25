package access

import (
	"errors"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
)

const (
	// Prefix uses in this page
	ENV_PREFIX = "GAF_MYSQL"
	env_uname  = ENV_PREFIX + "_USERNAME"
	env_passwd = ENV_PREFIX + "_PASSWORD"
	env_proto  = ENV_PREFIX + "_PROTOCOL"
	env_addr   = ENV_PREFIX + "_ADDRESS"
	env_dbname = ENV_PREFIX + "_DATABASE_NAME"
)

func getMySQLConfigFromEnv() (*mysql.Config, error) {
	uname := os.Getenv(env_uname)
	passwd := os.Getenv(env_passwd)
	proto := strings.ToLower(os.Getenv(env_proto))
	addr := os.Getenv(env_addr)
	dbname := os.Getenv(env_dbname)

	if len(uname) == 0 || len(passwd) == 0 || len(addr) == 0 || len(dbname) == 0 {
		return nil, errors.New("some mysql configuration missed")
	} else if len(proto) == 0 {
		proto = "tcp"
	}

	return &mysql.Config{
		User:                 uname,
		Passwd:               passwd,
		Addr:                 addr,
		DBName:               dbname,
		Net:                  proto,
		AllowNativePasswords: true,
	}, nil
}
