package todos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Job interface {
	ExitChan() chan error
	Run(todos map[string]Todo) (map[string]Todo, error)
}

func ProcessJobs(jobs <-chan Job, db string) {
	todos := make(map[string]Todo, 0) // initial a null map
	var err error
	todos, err = loadTodoListfromFile(db)

	if err != nil {
		fmt.Println("error: ", err)
	}

	for {
		j := <-jobs

		// store the data to  todos map ( hash list )
		todosMod, err := j.Run(todos) //  transfer todo data to channel and check change or not in return

		if err == nil && todosMod != nil { // identify the change
			saveTodoListtoFile(db, todosMod) // save to file
		}

		j.ExitChan() <- err
	}
}

// loadTodoListfromFile
// read data from file and Unmarshal to todos
// input:     db string:  the full file name string
// return:    map[string]Todo
func loadTodoListfromFile(db string) (map[string]Todo, error) {

	todos := make(map[string]Todo, 0)   // initial a null map
	content, err := ioutil.ReadFile(db) // read from file

	if err == nil {
		if err = json.Unmarshal(content, &todos); err == nil {
			//
			// fmt.Println(todos)
		}
	} else {
		// todos is null
		// initial todos as sample data
	}

	return todos, nil

}

// saveTodoListtoFile
// save todos to file
// input:    db string: the full file name string
// input:    todosMod map[string]Todo : the map ( hash lish ) of todos
// return:   err error:    success return nil
func saveTodoListtoFile(db string, todosMod map[string]Todo) (err error) {
	if todosMod != nil { // identify the change
		b, err := json.Marshal(todosMod)
		if err == nil {
			err = ioutil.WriteFile(db, b, 0644) // save to file
		}
	}
	return nil
}

/**
func ProcessJobs(jobs <-chan Job, db string) {
	for {
		j := <-jobs

		todos := make(map[string]Todo, 0)         // initial a null map
		content, err := ioutil.ReadFile(db)       // read from file
		if err == nil {
			if err = json.Unmarshal(content, &todos); err == nil {   // store the data to  todos map ( hash list )
				todosMod, err := j.Run(todos)     //  transfer todo data to channel and check change or not in return

				if err == nil && todosMod != nil {  // identify the change
					b, err := json.Marshal(todosMod)
					if err == nil {
						err = ioutil.WriteFile(db, b, 0644)       // save to file
					}
				}
			}
		}
		j.ExitChan() <- err
	}
}
*/
