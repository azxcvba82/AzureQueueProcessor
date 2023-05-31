package utils

import "os"

func GetSQLConnectString() string {
	return os.Getenv("SQLCONNECTSTRING")
}
