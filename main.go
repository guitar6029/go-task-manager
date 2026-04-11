package main

// @title Task Manager API
// @version 1.0
// @description Task management API with JWT authentication
// @host localhost
// @BasePath /

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	_ "taskmanager/docs"
	"taskmanager/internal/api"

	dbpkg "taskmanager/internal/db"
	model "taskmanager/internal/model"
	redispkg "taskmanager/internal/redis"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// var commands = []string{"help", "q (quit)", "add <task>", "list <done | pending> [--limit=N] [--offset=N]", "next", "prev", "delete <id>", "done <id>"}
var commands = []string{"help", "q (quit)", "add <task>", "list"}

//var offset = 0
//var currentFilter = ""

//var currentLimit = 5

func main() {

	// redis
	redis := redispkg.NewClient()

	//test connection
	_, err := redis.Ping(redispkg.Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis!")

	env := os.Getenv("APP_ENV")
	switch env {
	case "local":
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Println("No .env file found")
		}
	case "dev":
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("No .env file found")
		}
	default:
		log.Println("Running with system env variables")
	}

	// default mode selection
	mode := "api"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	db, err := dbInit()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("error closing db: ", err)
		}
	}()

	switch mode {
	case "api":
		startAPI(db, redis)
	case "cli":
		startCLI(db)
	default:
		fmt.Println("Unknown Command")
	}

}

func dbInit() (*sql.DB, error) {
	//db init
	db, err := dbpkg.Init()
	if err != nil {
		return nil, fmt.Errorf("error initializing DB: %s", err)
	}

	return db, nil
}

func startAPI(db *sql.DB, rdb *redis.Client) {
	fmt.Println("Initializing API Program")
	api.Start(db, rdb)
}

func startCLI(db *sql.DB) {
	fmt.Println("Initializing CLI Program")
	fmt.Println("CLI mode (local dev tool)")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">")
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
		handler(db, parts[0], parts[1:])
	}
}

func handler(db *sql.DB, command string, args []string) {
	switch command {
	case "help":
		fmt.Println("Commands:")
		for _, c := range commands {
			fmt.Println(c)
		}
	// case "add":
	// 	var err error
	// 	title := strings.Join(args, " ")
	// 	id, err := dbpkg.CreateTask(db, title)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	fmt.Printf("Added %s to ID %d\n", title, id)

	// case "list":
	// 	offset = 0
	// 	listType, limitOverride, offsetOverride := parseListArgs(args)

	// 	currentFilter = listType

	// 	if limitOverride > 0 {
	// 		currentLimit = limitOverride
	// 	}
	// 	if offsetOverride > 0 {
	// 		offset = offsetOverride
	// 	}

	// 	tasks, err := service.GetTasks(db, currentFilter, currentLimit, offset)
	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 		return
	// 	}

	// 	if len(tasks) == 0 {
	// 		fmt.Println("No tasks")
	// 		return
	// 	}
	// 	ListTasks(tasks)
	// 	fmt.Printf("Showing %d - %d\n", offset+1, offset+len(tasks))

	default:
		fmt.Println("Unknown command")
	}
	// case "next":
	// 	offset += currentLimit

	// 	tasks, err := service.GetTasks(db, currentFilter, currentLimit, offset)
	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 		return
	// 	}

	// 	if len(tasks) == 0 {
	// 		fmt.Println("No more tasks")
	// 		offset -= currentLimit
	// 		return
	// 	}

	// 	ListTasks(tasks)
	// 	fmt.Printf("Showing %d - %d\n", offset+1, offset+len(tasks))
	// case "prev":

	// 	if offset == 0 {
	// 		fmt.Println("Already at first page")
	// 		return
	// 	}

	// 	offset -= currentLimit

	// 	tasks, err := service.GetTasks(db, currentFilter, currentLimit, offset)
	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 		return
	// 	}

	// 	ListTasks(tasks)
	// 	fmt.Printf("Showing %d - %d\n", offset+1, offset+len(tasks))

	// case "delete":
	// 	if len(args) == 0 {
	// 		fmt.Println("Missing task id")
	// 		return
	// 	}
	// 	taskID, err := getTaskID(args)
	// 	if err != nil {
	// 		fmt.Println("Error : ", err)
	// 		return
	// 	}
	// 	err = service.DeleteTask(db, taskID)
	// 	if err != nil {
	// 		fmt.Println("Task not found")

	// 	} else {
	// 		fmt.Println("Deleted task : ", taskID)
	// 	}
	// case "done":

	// 	taskID, err := getTaskID(args)
	// 	if err != nil {
	// 		fmt.Println("Error : ", err)
	// 		return
	// 	}
	// 	task, err := service.MarkTaskDone(db, taskID)
	// 	if err != nil {

	// 		fmt.Println("Task not found")
	// 	} else {
	// 		fmt.Printf("Task %s %d marked as done\n", task.Title, taskID)
	// 	}

}

// func getTaskID(args []string) (int, error) {
// 	if len(args) == 0 {
// 		return 0, fmt.Errorf("missing task id")
// 	}

// 	taskID, err := strconv.Atoi(args[0])
// 	if err != nil {
// 		return 0, fmt.Errorf("missing task id")
// 	}
// 	return taskID, nil
// }

func ListTasks(tasks []model.Task) {
	for _, t := range tasks {
		status := "❌"
		if t.Done {
			status = "✅"
		}
		fmt.Printf("%d: %s %s\n", t.ID, t.Title, status)
	}

}

// func parseListArgs(args []string) (string, int, int) {
// 	listType := ""
// 	limitOverride := 0
// 	offset := 0

// 	for _, arg := range args {
// 		if arg == "done" || arg == "pending" {
// 			listType = arg
// 		} else if strings.HasPrefix(arg, "--limit=") {
// 			val := strings.TrimPrefix(arg, "--limit=")
// 			if n, err := strconv.Atoi(val); err == nil {
// 				limitOverride = n
// 			}
// 		} else if strings.HasPrefix(arg, "--offset=") {
// 			val := strings.TrimPrefix(arg, "--offset=")
// 			if n, err := strconv.Atoi(val); err == nil {
// 				offset = n
// 			}
// 		}
// 	}

// 	return listType, limitOverride, offset
// }
