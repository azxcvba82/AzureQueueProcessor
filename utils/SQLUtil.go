package utils

import (
	"context"
	"database/sql"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func SQLQueryClassic(sqlConnectionString string, sqlCommand string, args ...any) (r *sql.Rows, err error) {

	db, err := sql.Open("mysql", sqlConnectionString)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	rows, err := db.Query(sqlCommand, args...)
	if err != nil {
		return nil, err
	}
	return rows, err
}

func SQLQuery(model interface{}, sqlConnectionString string, sqlCommand string, args ...any) (err error) {

	db, err := sqlx.Open("mysql", sqlConnectionString)

	if err != nil {
		return err
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)

	if strings.Contains(reflect.ValueOf(model).Type().String(), "[]") {
		err = db.Select(model, sqlCommand, args...)
	} else {
		err = db.Get(model, sqlCommand, args...)
	}
	return err

}

func SQLExec(sqlConnectionString string, withTransaction bool, sqlCommand string, args ...any) (id int64, cnt int64, err error) {
	db, err := sql.Open("mysql", sqlConnectionString)
	defer db.Close()

	if err != nil {
		return -1, -1, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)

	var execResult sql.Result
	var insertId int64
	var rowsAffected int64
	if withTransaction == false {
		execResult, err = db.Exec(sqlCommand, args...)
	} else {
		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return -1, -1, err
		}
		defer tx.Rollback()
		execResult, err = tx.ExecContext(ctx, sqlCommand, args...)
		if err != nil {
			return -1, -1, err
		}

		insertId, err := execResult.LastInsertId()
		if err != nil {
			return -1, -1, err
		}

		rowsAffected, err := execResult.RowsAffected()
		if err != nil {
			return -1, -1, err
		}

		if err = tx.Commit(); err != nil {
			return insertId, rowsAffected, err
		}
	}
	if err != nil {
		return -1, -1, err
	}

	insertId, err = execResult.LastInsertId()
	if err != nil {
		return -1, -1, err
	}
	rowsAffected, err = execResult.RowsAffected()
	return insertId, rowsAffected, err
}
