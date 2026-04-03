package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	dbpkg "taskmanager/db"
	model "taskmanager/model"
)

var commands = []string{"help", "q (quit)", "add <task>", "list <done | pending> [--limit=N] [--offset=N]", "next", "prev", "delete <id>", "done <id>"}
var offset = 0
var currentFilter = ""

var currentLimit = 5

func main() {

	db, err := dbpkg.Init()
	if err != nil {
		fmt.Println("Error initializing DB: ", err)
		return
	}
	defer db.Close()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if input == "q" {
			fmt.Println("Goodbye.")
			break
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		handler(db, command, args)

	}

}

func handler(db *sql.DB, command string, args []string) {
	switch command {
	case "add":
		var err error
		title := strings.Join(args, " ")
		id, err := dbpkg.CreateTask(db, title)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Added %s to ID %d\n", title, id)

	case "list":
		offset = 0
		listType, limitOverride, offsetOverride := parseListArgs(args)

		currentFilter = listType

		if limitOverride > 0 {
			currentLimit = limitOverride
		}
		if offsetOverride > 0 {
			offset = offsetOverride
		}

		tasks, err := dbpkg.GetTasks(db, currentFilter, currentLimit, offset)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks")
			return
		}
		ListTasks(tasks)
		fmt.Printf("Showing %d - %d\n", offset+1, offset+len(tasks))
	case "next":
		offset += currentLimit

		tasks, err := dbpkg.GetTasks(db, currentFilter, currentLimit, offset)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("No more tasks")
			offset -= currentLimit
			return
		}

		ListTasks(tasks)
		fmt.Printf("Showing %d - %d\n", offset+1, offset+len(tasks))
	case "prev":

		if offset == 0 {
			fmt.Println("Already at first page")
			return
		}

		offset -= currentLimit

		tasks, err := dbpkg.GetTasks(db, currentFilter, currentLimit, offset)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		ListTasks(tasks)
		fmt.Printf("Showing %d - %d\n", offset+1, offset+len(tasks))
	case "help":
		fmt.Println("Commands:")
		for _, c := range commands {
			fmt.Println(c)
		}
	case "delete":
		if len(args) == 0 {
			fmt.Println("Missing task id")
			return
		}
		taskID, err := getTaskID(args)
		if err != nil {
			fmt.Println("Error : ", err)
			return
		}
		err = dbpkg.DeleteTask(db, taskID)
		if err != nil {
			fmt.Println("Task not found")

		} else {
			fmt.Println("Deleted task : ", taskID)
		}
	case "done":

		taskID, err := getTaskID(args)
		if err != nil {
			fmt.Println("Error : ", err)
			return
		}
		err = dbpkg.UpdateTaskStatus(db, taskID, true)
		if err != nil {

			fmt.Println("Task not found")
		} else {
			fmt.Printf("Task %d marked as done\n", taskID)
		}

	default:
		fmt.Println("Unknown command")
	}
}

func getTaskID(args []string) (int, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("missing task id")
	}

	taskID, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("missing task id")
	}
	return taskID, nil
}

func ListTasks(tasks []model.Task) {
	for _, t := range tasks {
		status := "❌"
		if t.Done {
			status = "✅"
		}
		fmt.Printf("%d: %s %s\n", t.ID, t.Title, status)
	}

}

func parseListArgs(args []string) (string, int, int) {
	listType := ""
	limitOverride := 0
	offset := 0

	for _, arg := range args {
		if arg == "done" || arg == "pending" {
			listType = arg
		} else if strings.HasPrefix(arg, "--limit=") {
			val := strings.TrimPrefix(arg, "--limit=")
			if n, err := strconv.Atoi(val); err == nil {
				limitOverride = n
			}
		} else if strings.HasPrefix(arg, "--offset=") {
			val := strings.TrimPrefix(arg, "--offset=")
			if n, err := strconv.Atoi(val); err == nil {
				offset = n
			}
		}
	}

	return listType, limitOverride, offset
}
