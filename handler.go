package main

import (
    "fmt"
    "net/http"
	"time"
	"encoding/json"
	"github.com/pborman/uuid"
	"regexp"
	"strings"

    jwt "github.com/form3tech-oss/jwt-go"
)
var mySigningKey = []byte("secret")

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse from body of request to get a json object.
	fmt.Println("Received one upload request")
    w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	if r.Method == "OPTIONS" {
        return
    }
	user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]

	post := Post{
		User: username.(string),
		Firstname: r.PostFormValue("firstname"),
		Lastname: r.PostFormValue("lastname"),
		Description: r.PostFormValue("description"),
	}
	current_time := time.Now()
	post.Uploadtime = current_time.Format("2006-01-02 15:04:05")
    post.Updatetime = current_time.Format("2006-01-02 15:04:05")
	file, header, err := r.FormFile("mediafile")
	if err != nil {
        http.Error(w, "Media file is not available", http.StatusBadRequest)
        fmt.Printf("Media file is not available %v\n", err)
        return
    }
	id := uuid.New()
	fileaddress, err := UploadImage(file, id, header);
	if err != nil {
		panic(err)
	}
	post.Url = fileaddress
    post.Type = "image"
	uploadToDB(post)
	w.WriteHeader(http.StatusOK)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one delete request")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	
	if r.Method == "OPTIONS" {
        return
    }
	post := Post{}
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		panic(err)
	}
	user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]

	post.User = username.(string)
	deleteFromDB(post.User, post.Url, "Posts")
	pidArr := strings.Split(post.Url, "com/")
	err = DeleteImage(pidArr[1])
	if err != nil {
		fmt.Println("Delete from S3 fail")
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}
func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get request")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	
	if r.Method == "OPTIONS" {
        return
    }
    // uname := r.URL.Query().Get("user")
	var post Post
	// err := json.NewDecoder(r.Body).Decode(&post)
    // if err != nil {
    //     http.Error(w, err.Error(), http.StatusBadRequest)
    //     panic(err)
    // }
	user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]

	post.User = username.(string)
	
	fmt.Println("user is " + post.User)
        res, err := readFromDB(post.User, "Posts")
        if err != nil {
            panic(err)
            fmt.Println("something wrong with DB")
        }
        js ,err := json.Marshal(res)
        if err != nil {
            panic(err)
        }
        w.WriteHeader(http.StatusOK)
	    w.Write(js)
        
    // else{
	//     res, err := readApostFromDB(post.User, post.Url, "Posts")
    //     if err != nil {
    //         panic(err)
    //         fmt.Println("something wrong with DB")
    //     }
    //     js ,err := json.Marshal(res)
    //     if err != nil {
    //         panic(err)
    //     }
    //     w.WriteHeader(http.StatusOK)
	//     w.Write(js)
    // }
	
	
}
func updateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one update request")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	
	if r.Method == "OPTIONS" {
        return
    }
	post := Post{}
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		panic(err)
	}
	user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]

	post.User = username.(string)
	res , err := updateFromDB(post)
	if err != nil {
		panic(err)
	}
	js ,err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signin request")
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
        return
    }

    //  Get User information from client
    decoder := json.NewDecoder(r.Body)
    var user User
    if err := decoder.Decode(&user); err != nil {
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return
    }

    exists, err := checkUser(&user)
    if err != nil {
        http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to read user from Elasticsearch %v\n", err)
        return
    }

    if !exists {
        http.Error(w, "User doesn't exists or wrong password", http.StatusUnauthorized)
        fmt.Printf("User doesn't exists or wrong password\n")
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 48).Unix(),
    })

    tokenString, err := token.SignedString(mySigningKey)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        fmt.Printf("Failed to generate token %v\n", err)
        return
    }

    w.Write([]byte(tokenString))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one signup request")
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == "OPTIONS" {
        return
    }

    decoder := json.NewDecoder(r.Body)
    var user User
    if err := decoder.Decode(&user); err != nil {
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return
    }

    if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        fmt.Printf("Invalid username or password\n")
        return
    }

    success, err := addUser(&user)
    if err != nil {
        http.Error(w, "Failed to save user to Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to save user to Elasticsearch %v\n", err)
        return
    }

    if !success {
        http.Error(w, "User already exists", http.StatusBadRequest)
        fmt.Println("User already exists")
        return
    }
    fmt.Printf("User added successfully: %s.\n", user.Username)
	w.WriteHeader(http.StatusOK)

}
