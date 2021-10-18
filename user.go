package main

import (
    "fmt"
)

const (
    USER_INDEX = "Users"
)

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func checkUser(user *User)(bool, error) {
	serachReult, err := readUserFromDB(user.Username, USER_INDEX)
	if err != nil {
		println("got failuer reading form DB")
		return false, err
	}
	if user.Password == serachReult {
		fmt.Printf("Login as %s\n", user.Username)
		return true, nil
	}
    return false, nil
}
func addUser(user *User) (bool, error) {
    searchResult, err := readUserFromDB(user.Username, USER_INDEX)
    if err != nil {
        return false, err
    }

    if len(searchResult) > 0 {
        return false, nil
    }

    err = saveUserToDB(user.Username, user.Password)
    if err != nil {
        return false, err
    }
    fmt.Printf("User is added: %s\n", user.Username)
    return true, nil
}
