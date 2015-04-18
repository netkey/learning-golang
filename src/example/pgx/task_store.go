package main

import (
	"fmt"
	log "github.com/inconshreveable/log15"
	"github.com/jackc/pgx"
	"os"
	//"strconv"
)

type PostgresDB struct {
	Pool       *pgx.ConnPool
	poolConfig pgx.ConnPoolConfig
}

type PostgresTx struct {
	tx *pgx.Tx
}

func (pgdb *PostgresDB) afterConnect(conn *pgx.Conn) (err error) {

	_, err = conn.Prepare("getTask", `
    select id,description from tasks where id=$1
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("listTask", `
    select * from tasks
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("addTask", `
    insert into tasks(description) values( $1 )
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("updateTask", `
    update tasks
      set description=$2
      where id=$1
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("deleteTask", `
    delete from tasks where id=$1
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("transfer", `select * from transfer('Bob','Mary',14.00)`)
	if err != nil {
		return
	}

	// There technically is a small race condition in doing an upsert with a CTE
	// where one of two simultaneous requests to the shortened URL would fail
	// with a unique index violation. As the point of this demo is pgx usage and
	// not how to perfectly upsert in PostgreSQL it is deemed acceptable.
	_, err = conn.Prepare("putTask", `
    with upsert as (
      update tasks
      set description=$2
      where id=$1
      returning *
    )
    insert into tasks(id, description)
    select $1, $2 where not exists(select 1 from upsert)
  `)
	return
}

func (pgdb *PostgresDB) InitConfig(dbhost, dbuser, dbpassword, dbname string) error {

	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     dbhost,
			User:     dbuser,
			Password: dbpassword,
			Database: dbname,
			//Logger:   log.New("module", "pgx"),
		},
		MaxConnections: 5,
		AfterConnect:   pgdb.afterConnect,
	}

	pgdb.poolConfig = connPoolConfig
	return nil
}

func (pgdb *PostgresDB) InitConnection() error {
	//var pool *pgx.ConnPool
	var err error

	pgdb.Pool, err = pgx.NewConnPool(pgdb.poolConfig)
	if err != nil {
		log.Crit("Unable to create connection pool", "error", err)
		os.Exit(1)
	}

	log.Crit("database connect sueecss")
	return nil
}

func (pgdb *PostgresDB) Transfer() error {
	rows, _ := pgdb.Pool.Query("transfer") // limit 4 offset 2")

	for rows.Next() {

		var transfer string
		err := rows.Scan(&transfer)
		if err != nil {
			return err
		}
		fmt.Printf("select * from transfer('Bob','Mary',14.00) return: %s\n", transfer)
	}

	return rows.Err()
}

func (pgdb *PostgresDB) listTasks() error {
	rows, _ := pgdb.Pool.Query("listTask") // limit 4 offset 2")

	for rows.Next() {
		var id int32
		var description string
		err := rows.Scan(&id, &description)
		if err != nil {
			return err
		}
		fmt.Printf("%d. %s\n", id, description)
	}

	return rows.Err()
}

func (pgdb *PostgresDB) addTask(description string) error {

	length := len(description)
	fmt.Println("length of description is: ", length)

	if length > 0 {

		tx, err := pgdb.Pool.Begin()
		checkError(err)
		// Rollback is safe to call even if the tx is already closed, so if
		// the tx commits successfully, this is a no-op
		defer tx.Rollback()

		_, err = pgdb.Pool.Exec("addTask", description)
		checkError(err)
		err = tx.Commit()
		checkError(err)

	} else {
		fmt.Println(" description is null")
	}

	return nil
}

func (pgdb *PostgresDB) updateTask(itemNum int32, description string) error {

	tx, err := pgdb.Pool.Begin()
	checkError(err)
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback()

	_, err = pgdb.Pool.Exec("updateTask", itemNum, description)

	checkError(err)
	err = tx.Commit()

	return checkError(err)

}

func (pgdb *PostgresDB) removeTask(itemNum int32) error {

	_, err1 := pgdb.Pool.Exec("deleteTask", itemNum)
	return err1

}
