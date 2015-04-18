package main

import (
	"fmt"
	//"github.com/jackc/pgx"
	log "github.com/inconshreveable/log15"
	"os"
	"strconv"
)

func checkError(err error) error {
	if err != nil {
		panic(err)
	}
	return nil
}

func printHelp() {
	fmt.Print(`Todo pgx demo
Usage:
  todo list
  todo add task
  todo update task_num item
  todo remove task_num
Example:
  todo add 'Learn Go'
  todo list
`)
}

//*********************************************************************************

func main() {
	var err error
	// init postgres DB connection

	//var config pgx.ConnConfig

	dbhost := "localhost"
	dbuser := "postgres"
	dbpassword := "postgres"
	dbname := "tsingcloud"

	var pgdb PostgresDB

	pgdb.InitConfig(dbhost, dbuser, dbpassword, dbname)

	pgdb.InitConnection()

	defer pgdb.Pool.Close()

	/**
	conn, err = pgx.Connect(extractConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	*/
	if len(os.Args) == 1 {
		printHelp()
		os.Exit(0)
	}

	switch os.Args[1] {

	case "test":
		err = pgdb.Transfer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "query error: %v\n", err)
			log.Crit("query error", "error", err)
			os.Exit(1)
		}

	case "list":
		err = pgdb.listTasks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to list tasks: %v\n", err)
			log.Crit("Unable to list tasks", "error", err)
			os.Exit(1)
		}

	case "add":
		err = pgdb.addTask(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to add task: %v\n", err)
			os.Exit(1)
		}

	case "update":
		n, err := strconv.ParseInt(os.Args[2], 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable convert task_num into int32: %v\n", err)
			os.Exit(1)
		}
		err = pgdb.updateTask(int32(n), os.Args[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to update task: %v\n", err)
			os.Exit(1)
		}

	case "remove":
		n, err := strconv.ParseInt(os.Args[2], 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable convert task_num into int32: %v\n", err)
			os.Exit(1)
		}
		err = pgdb.removeTask(int32(n))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to remove task: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "Invalid command")
		printHelp()
		os.Exit(1)
	}
}
