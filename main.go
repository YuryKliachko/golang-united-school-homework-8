package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type Arguments map[string]string

func validateArgs(operation string, fileName string, userId string, item string) string {
	var errMessage string
	if operation == "" {
		return "-operation flag has to be specified"
	}
	if fileName == "" {
		return "-fileName flag has to be specified"
	}

	if !IsValidOperation(operation) {
		return fmt.Sprintf("Operation %s not allowed!", operation)
	}
	if item == "" && operation == Add {
		return "-item flag has to be specified"
	}
	if userId == "" && (operation == FindById || operation == Remove) {
		return "-id flag has to be specified"
	}
	return errMessage
}

func readFile(fileName string) ([]byte, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	file.Close()
	return bytes, nil
}

func writeFile(fileName string, data []byte) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	file.Write(data)
	file.Close()
	return nil
}

func listUsers(fileName string, writer io.Writer) error {
	bytes, err := readFile(fileName)
	if err != nil {
		return err
	}
	writer.Write(bytes)
	return nil
}

func addUser(fileName string, writer io.Writer, item string) error {
	fmt.Println("Adding user to file")
	bytes, err := readFile(fileName)
	if err != nil {
		return err
	}

	var newUser User
	json.Unmarshal([]byte(item), &newUser)

	var users []User
	json.Unmarshal(bytes, &users)

	for _, user := range users {
		if newUser.Id == user.Id {
			errorMessage := fmt.Sprintf("Item with id %s already exists", newUser.Id)
			fmt.Println(errorMessage)
			writer.Write([]byte(errorMessage))
			return nil
		}
	}
	users = append(users, newUser)
	serialiezedUsers, err := json.Marshal(users)
	if err != nil {
		return err
	}
	writeFile(fileName, serialiezedUsers)
	return nil
}

func findById(fileName string, writer io.Writer, userId string) error {
	fmt.Println("Trying to find user by ID")
	bytes, err := readFile(fileName)
	if err != nil {
		return err
	}
	var users []User
	json.Unmarshal(bytes, &users)

	var foundUser []byte
	for _, user := range users {
		if user.Id == userId {
			serializedUser, err := json.Marshal(user)
			if err != nil {
				return err
			}
			foundUser = serializedUser
		}
	}
	writer.Write(foundUser)
	return nil
}

func removeUser(fileName string, writer io.Writer, userId string) error {
	fmt.Println("Trying to remove user by ID ", userId)
	bytes, err := readFile(fileName)
	if err != nil {
		return err
	}
	var users []User
	json.Unmarshal(bytes, &users)
	fmt.Println("Current users", users)
	var newUsers []User
	for index, user := range users {
		if user.Id == userId {
			newUsers = append(users[:index], users[index+1:]...)
			serializedUsers, err := json.Marshal(newUsers)
			if err != nil {
				return err
			}
			fmt.Println("New users", newUsers)
			writeFile(fileName, serializedUsers)
			return nil
		}
	}
	writer.Write([]byte(fmt.Sprintf("Item with id %s not found", userId)))
	return nil
}

func Perform(args Arguments, writer io.Writer) error {
	operation := args["operation"]
	fileName := args["fileName"]
	userId := args["id"]
	item := args["item"]

	validationError := validateArgs(operation, fileName, userId, item)
	if validationError != "" {
		return errors.New(validationError)
	}

	var operationError error
	switch operation {
	case List:
		operationError = listUsers(fileName, writer)
	case Add:
		operationError = addUser(fileName, writer, item)
	case FindById:
		operationError = findById(fileName, writer, userId)
	default:
		operationError = removeUser(fileName, writer, userId)
	}
	return operationError
}

func parseArgs() Arguments {
	var id = flag.String("id", "", "user identificator")
	var operation = flag.String("operation", "", "operation to perform")
	var item = flag.String("item", "", "item to save")
	var fileName = flag.String("fileName", "", "name of a file to get/put data from/to")
	flag.Parse()
	return Arguments{
		"id":        *id,
		"operation": *operation,
		"item":      *item,
		"fileName":  *fileName,
	}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
