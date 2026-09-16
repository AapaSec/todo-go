package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"encoding/json/jsontext"
	"encoding/json/v2"
)

type Task struct {
	Id          int
	Description string
	Completed   bool
}

type Tasks []Task

func help() {
	fmt.Println("Available commands: add, list, check, uncheck, edit, delete, help, exit")
}

func getID() int {
	counter++
	return counter
}

func (tasks Tasks) String() string {
	var result string
	if tasks == nil {
		return "Empty tasks list"
	}
	for _, task := range tasks {
		checked := " "
		if check := task.Completed; check {
			checked = "X"
		}
		result += fmt.Sprintf("T%v\t[%s] %s\n", task.Id, checked, task.Description)
	}
	return result
}

func (tasks Tasks) retrieveByID(id int) *Task {
	for i, task := range tasks {
		if task.Id == id {
			return &tasks[i]
		}
	}
	return nil
}

func (t *Tasks) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if k := dec.PeekKind(); k != jsontext.KindBeginArray {
		return &json.SemanticError{JSONKind: k}
	}
	if _, err := dec.ReadToken(); err != nil {
		return err
	}
	var task Task
	for dec.PeekKind() != jsontext.KindEndArray {
		if err := json.UnmarshalDecode(dec, &task); err != nil {
			return err
		}
		*t = append(*t, task)
	}
	if _, err := dec.ReadToken(); err != nil {
		return err
	}

	return nil
}

var counter int = 0

func main() {
	var tasks Tasks
	// check scanner.Err()
	scanner := bufio.NewScanner(os.Stdin)
	file, err := os.OpenFile("tasks.json", os.O_RDWR, 777)
	if err != nil {
		panic(err)
	}

	// TEST
	stat, _ := file.Stat()
	testJSON2 := make([]byte, stat.Size())
	n, err := file.Read(testJSON2)
	//testTask1 := Task{Id: -1, Description: "TEST"}
	//testTask2 := Task{Id: -2, Description: "TEST"}
	//testTask3 := Task{Id: -3, Description: "TEST"}
	//testTask4 := Task{Id: -4, Description: "TEST"}
	//tasks = append(tasks, testTask1, testTask2, testTask3, testTask4)

	fmt.Println(string(testJSON2))
	if err := json.Unmarshal(testJSON2, &tasks); err != nil {
		fmt.Println("AA")
		panic(err)
	}

	for {
		fmt.Print("Todo>")
		// handle Scan() error later
		scanner.Scan()
		input := scanner.Text()
		command, options, _ := strings.Cut(input, " ")
		// add save file option
		switch command {
		case "add":
			task := Task{Id: getID(), Description: options}
			tasks = append(tasks, task)
		case "list":
			fmt.Println(tasks)
		case "check":
			id, err := strconv.Atoi(options)
			if err == nil {
				task := tasks.retrieveByID(id)
				if task != nil {
					task.Completed = true
					fmt.Println("Task checked")
				} else {
					fmt.Println("Task not found")
				}
			} else {
				fmt.Println("Invalid ID")
			}
		case "uncheck":
			id, err := strconv.Atoi(options)
			if err == nil {
				task := tasks.retrieveByID(id)
				if task != nil {
					task.Completed = false
					fmt.Println("Task unchecked")
				} else {
					fmt.Println("Task not found")
				}
			} else {
				fmt.Println("Invalid ID")
			}
		case "edit":
			idStr, description, _ := strings.Cut(options, " ")
			id, err := strconv.Atoi(idStr)
			if err == nil {
				task := tasks.retrieveByID(id)
				if task != nil {
					task.Description = description
					fmt.Println("Task edited successfully")
				} else {
					fmt.Println("Task not found")
				}
			} else {
				fmt.Println("Invalid ID")
			}
		case "delete":
			id, err := strconv.Atoi(options)
			if err == nil {
				task := tasks.retrieveByID(id)
				if task != nil {
					tasks = slices.DeleteFunc(tasks, func(task Task) bool { return task.Id == id })
					fmt.Println("Task deleted successfully")
				} else {
					fmt.Println("Task not found")
				}
			} else {
				fmt.Println("Invalid ID")
			}
		case "help":
			help()
		case "save":

			// unwanted chracters added
			tasksJSON, _ := json.Marshal(tasks)
			_, err := os.Create("tasks.json")
			_, err = file.Write(tasksJSON)
			if err != nil {
				fmt.Println(n)
				panic(err)
			}
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Unknown command")
		}
	}
}
