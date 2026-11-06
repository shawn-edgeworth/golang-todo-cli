package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	todo "github.com/shawn-edgeworth/golang-todo-cli/todo"
)

func main() {
	args := os.Args
	var todos []*todo.Todo
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Todo CLI v1.0")

	if len(args) > 1 {
		exist, err := ensureTodosFileExist("todos")

		if !exist {
			todos = make([]*todo.Todo, 0)
		} else {
			todos, err = load("todos")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		value := ""
		if len(args) > 2 {
			value = args[2]
		}

		todos = executeCommand(args[1], value, todos)
		err = save("todos", todos)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	for {
		command, value, err := getCommand(reader)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if command == "exit" {
			os.Exit(0)
		}

		todos = executeCommand(command, value, todos)
	}
}

func executeCommand(command, value string, todos []*todo.Todo) []*todo.Todo {
	switch command {
	case "load":
		if len(value) == 0 {
			fmt.Println("Аргумент не найден")
			return todos
		}
		newTodos, err := load(value)
		if err != nil {
			fmt.Println(err)
		} else {
			todos = newTodos
		}
		return todos
	case "save":
		if len(value) == 0 {
			fmt.Println("Аргумент не найден")
			return todos
		}
		save(value, todos)
		return todos
	case "add":
		if len(value) == 0 {
			fmt.Println("Аргумент не найден")
			return todos
		}
		todos = append(todos, todo.NewTodo(value))
		return todos
	case "list":
		if len(todos) < 1 {
			fmt.Println("Список задач пуст")
			return todos
		}
		printTodos(todos)
		return todos
	case "done":
		index, err := parseInteger(value)
		if err != nil {
			fmt.Println(err)
			return todos
		}
		if index < 1 || index > len(todos) {
			fmt.Println("Такой задачи не существует")
			return todos
		}
		todos[index-1].Complete()
		return todos
	case "undone":
		index, err := parseInteger(value)
		if err != nil {
			fmt.Println(err)
			return todos
		}
		if index < 1 || index > len(todos) {
			fmt.Println("Такой задачи не существует")
			return todos
		}
		todos[index-1].Undone()
		return todos
	default:
		fmt.Printf("Неизвестная команда: %s\n", command)
		return todos
	}
}

func getCommand(reader *bufio.Reader) (string, string, error) {
	fmt.Print("Введи команду:")
	input, err := getInput(reader)
	if err != nil {
		return "", "", err
	}

	command, value, _ := strings.Cut(input, " ")
	command = strings.ToLower(command)
	return command, value, nil
}

func getInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	return input, nil
}

func parseInteger(value string) (int, error) {
	if len(value) <= 0 {
		return -1, errors.New("filename is empty")
	}
	index, err := strconv.Atoi(value)
	return index, err
}

func printTodos(todos []*todo.Todo) {
	fmt.Println()
	fmt.Println("Список задач:")
	for i, todo := range todos {
		completed := "[ ]"
		if todo.Done {
			completed = "[x]"
		}
		fmt.Printf("%d. %s %s\n", i+1, completed, todo.Value)
	}
}

func load(filename string) ([]*todo.Todo, error) {
	if len(filename) <= 0 {
		return nil, errors.New("filename is empty")
	}

	filenameWithExt := fmt.Sprintf("saves/%s.json", filename)
	b, err := os.ReadFile(filenameWithExt)
	if err != nil {
		return nil, err
	}

	todos := make([]*todo.Todo, 0)
	err = json.Unmarshal(b, &todos)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func ensureTodosFileExist(filename string) (bool, error) {
	if len(filename) <= 0 {
		return false, nil
	}

	filenameWithExt := fmt.Sprintf("saves/%s.json", filename)
	_, err := os.Stat(filenameWithExt)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return true, err
}

func save(filename string, todos []*todo.Todo) error {
	if len(filename) <= 0 {
		return errors.New("filename is empty")
	}

	b, err := json.Marshal(todos)
	if err != nil {
		return err
	}

	filenameWithExt := fmt.Sprintf("saves/%s.json", filename)
	err = os.WriteFile(filenameWithExt, b, 0644)

	return err
}
