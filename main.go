package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Task struct {
	id          int
	description string
	completed   bool
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
		if check := task.completed; check {
			checked = "X"
		}
		result += fmt.Sprintf("T%v\t[%s] %s\n", task.id, checked, task.description)
	}
	return result
}

func (tasks Tasks) retrieveByID(id int) *Task {
	for i, task := range tasks {
		if task.id == id {
			return &tasks[i]
		}
	}
	return nil
}

var counter int = 0

func main() {
	var tasks Tasks
	// check scanner.Err()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Todo>")
		// handle Scan() error later
		scanner.Scan()
		input := scanner.Text()
		command, options, _ := strings.Cut(input, " ")
		switch command {
		case "add":
			task := Task{id: getID(), description: options}
			tasks = append(tasks, task)
		case "list":
			fmt.Println(tasks)
		case "check":
			id, error := strconv.Atoi(options)
			if error == nil {
				task := tasks.retrieveByID(id)
				if task != nil {
					task.completed = true
					fmt.Println("Task checked")
				} else {
					fmt.Println("Task not found")
				}
			} else {
				fmt.Println("Invalid ID")
			}
		case "uncheck":
			id, error := strconv.Atoi(options)
			if error == nil {
				task := tasks.retrieveByID(id)
				if task != nil {
					task.completed = false
					fmt.Println("Task unchecked")
				} else {
					fmt.Println("Task not found")
				}
			} else {
				fmt.Println("Invalid ID")
			}
		case "edit":
			idStr, description, _ := strings.Cut(options, " ")
			id, error := strconv.Atoi(idStr)
			if error == nil {
				task := tasks.retrieveByID(id)
				if task != nil {
					task.description = description
					fmt.Println("Task edited successfully")
				} else {
					fmt.Println("Task not found")
				}
			} else {
				fmt.Println("Invalid ID")
			}
		case "delete":
			id, error := strconv.Atoi(options)
			if error == nil {
				task := tasks.retrieveByID(id)
				if task != nil {
					tasks = slices.DeleteFunc(tasks, func(task Task) bool { return task.id == id })
					fmt.Println("Task deleted successfully")
				} else {
					fmt.Println("Task not found")
				}
			} else {
				fmt.Println("Invalid ID")
			}
		case "help":
			help()
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Unknown command")
		}
	}
}
