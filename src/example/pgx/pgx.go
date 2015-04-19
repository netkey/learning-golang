package main

import (
	"example/pgx/task"
	"fmt"
	log "github.com/inconshreveable/log15"
	//"github.com/jackc/cli"
	//"github.com/jackc/pgx"
	"github.com/vaughan0/go-ini"

	"os"
	"path/filepath"
	"strconv"
)

//*********************************************************************************

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

func loadConfig(path string) (ini.File, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("Invalid config path: %v", err)
	}

	file, err := ini.LoadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to load config file: %v", err)
	}

	return file, nil
}

//*********************************************************************************

func main() {

	srvlog := log.New("module", "postgres/pgx")
	srvlog.Info("Program starting", "args", os.Args)

	var err error
	// init postgres DB connection
	dbhost := "localhost"
	dbuser := "postgres"
	dbpassword := "postgres"
	dbname := "tsingcloud"

	srvlog.Warn("database info", log.Ctx{"dbhost": dbhost, "dbuser": dbuser})

	//var taskPgdb *task.PostgresDB

	//taskPgdb =&task.PostgresDB{}

	taskPgdb := new(task.PostgresDB)

	taskPgdb.InitDb(dbhost, dbuser, dbpassword, dbname)

	defer taskPgdb.Pool.Close()

	if len(os.Args) == 1 {
		printHelp()
		srvlog.Warn("print help")
		os.Exit(0)
	}

	switch os.Args[1] {

	case "test":
		err = taskPgdb.Transfer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "query error: %v\n", err)
			srvlog.Warn("query error", "error", err)
			os.Exit(1)
		}

	case "list":
		err = taskPgdb.ListTasks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to list tasks: %v\n", err)
			srvlog.Warn("Unable to list tasks", "error", err)
			os.Exit(1)
		}

	case "add":
		err = taskPgdb.AddTask(os.Args[2])
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
		err = taskPgdb.UpdateTask(int32(n), os.Args[3])
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
		err = taskPgdb.RemoveTask(int32(n))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to remove task: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "Invalid command")
		srvlog.Warn("print help")
		printHelp()
		os.Exit(1)
	}
}
