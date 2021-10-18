package main

type Post struct {
	User        string    `json:"user" dynamodbav"user"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	Description string    `json:"description"`
	Uploadtime  string 	  `json:"uploadtime"`
	Updatetime  string    `json:"updatetime"`
	Url         string    `json:"url"`
	Type        string    `json:"type"`
}
