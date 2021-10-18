package main

import(
	"fmt"
	"log"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)
func init_user_db() {
	// AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAZAZDZ6YVHU4TTHLL",
			"Z21ebf95ZIYGkyp43WpCcAvYiPYApVaBpsimpk7i",
			"", // a token will be created when the session it's used.
		  ),
	}))
	
	// A DynamoDB client
	svc := dynamodb.New(sess)

	// Create a table
	tableName := "Users"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("username"),
				AttributeType: aws.String("S"),
			},
			
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("username"),
				KeyType:       aws.String("HASH"),
			},
			
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	fmt.Println("Created the Users table", tableName)

}

func init_post_db() {
	// AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		// Endpoint: aws.String("http://3.87.3.254"),
		// EndPoint: aws.String("https://dynamodb.us-east-1.amazonaws.com"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAZAZDZ6YVHU4TTHLL",
			"Z21ebf95ZIYGkyp43WpCcAvYiPYApVaBpsimpk7i",
			"", // a token will be created when the session it's used.
		  ),
	}))
	
	// A DynamoDB client
	svc := dynamodb.New(sess)

	// Create a table
	tableName := "Posts"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("user"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("url"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("user"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("url"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	fmt.Println("Created the Posts table", tableName)
}
func uploadToDB(p Post) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAZAZDZ6YVHU4TTHLL",
			"Z21ebf95ZIYGkyp43WpCcAvYiPYApVaBpsimpk7i",
			"", // a token will be created when the session it's used.
		  ),
	}))
	
	// A DynamoDB client
	svc := dynamodb.New(sess)

	postAMAP, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		panic("Cannot marshal post into AttributeValue map")
	}
	// create the api params
	// fmt.Println("marshalled struct: %+v", postAMAP)
	params := &dynamodb.PutItemInput{
		TableName: aws.String("Posts"),
		Item:      postAMAP,
	}
	// put the item
	_, err = svc.PutItem(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		return
	}
	fmt.Println(params.Item)
	// print the response data
	fmt.Println("Successly save %+v 's post to DB", p.Url)
	
}
func deleteFromDB(user string, url string, curTable string) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAZAZDZ6YVHU4TTHLL",
			"Z21ebf95ZIYGkyp43WpCcAvYiPYApVaBpsimpk7i",
			"", // a token will be created when the session it's used.
		  ),
	}))
	
	// A DynamoDB client
	svc := dynamodb.New(sess)

	// create the api params
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(curTable),
		Key: map[string]*dynamodb.AttributeValue{
			"user": {
				S: aws.String(user),
			},
			"url": {
				S: aws.String(url),
			},
		},
	}
	// delete the item
	resp, err := svc.DeleteItem(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		return
	}

	// print the response data
	fmt.Println("Successfully delete %+v from DB", user)
	fmt.Println(resp)
}

func readApostFromDB(user string, url string, curTable string) (*Post, error){
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAZAZDZ6YVHU4TTHLL",
			"Z21ebf95ZIYGkyp43WpCcAvYiPYApVaBpsimpk7i",
			"", // a token will be created when the session it's used.
		  ),
	}))
	
	// A DynamoDB client
	svc := dynamodb.New(sess)

	// create the api params
	params := &dynamodb.GetItemInput{
		TableName: aws.String(curTable),
		Key: map[string]*dynamodb.AttributeValue{
			"user": {
				S: aws.String(user),
			},
			"url": {
				S: aws.String(url),
			},
		},
	}
	
	// read the item
	resp, err := svc.GetItem(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
	}
	// dump the response data
	fmt.Println(resp.Item)
	var p Post
	// unmarshal the dynamodb attribute values into a custom struct
	err = dynamodbattribute.UnmarshalMap(resp.Item, &p)
	// print the response data
	fmt.Printf("Unmarshaled Post = %+v\n", p)
	return &p, nil
} 

func readFromDB(user string, curTable string) ([]Post, error){
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAZAZDZ6YVHU4TTHLL",
			"Z21ebf95ZIYGkyp43WpCcAvYiPYApVaBpsimpk7i",
			"", // a token will be created when the session it's used.
		  ),
	}))
	
	// A DynamoDB client
	svc := dynamodb.New(sess)
	
	// create the api params
	params := &dynamodb.QueryInput{
		TableName:              aws.String("Posts"),
		KeyConditionExpression: aws.String("#user = :uname"),
		ExpressionAttributeNames: map[string]*string{
			"#user": aws.String("user"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":uname": {
				S: aws.String(user),
			},
		},
		
	}

	// read the item
	resp, err := svc.Query(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
	}

	// dump the response data
	fmt.Println(resp)

	// Unmarshal the slice of dynamodb attribute values
	// into a slice of custom structs
	var posts []Post
	err = dynamodbattribute.UnmarshalListOfMaps(resp.Items, &posts)
	return posts, nil
	
} 
func queryHandler(page *dynamodb.QueryOutput, lastPage bool) bool {
	// dump the response data
	//fmt.Println(page)

	// Unmarshal the slice of dynamodb attribute values
	// into a slice of custom structs
	var posts []Post
	err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &posts)
	if err != nil {
		// print the error and continue receiving pages
		fmt.Printf("\nCould not unmarshal AWS data: err = %v\n", err)
		return true
	}

	// print the response data
	for _, m := range posts {
		fmt.Printf("Post: '%s' (%s)\n", m.User, m.Url)
	}

	// if not done receiving all of the pages
	if lastPage == false {
		fmt.Printf("\n*** NOT DONE RECEIVING PAGES ***\n\n")
	} else {
		fmt.Printf("\n*** RECEIVED LAST PAGE ***\n\n")
	}

	// continue receiving pages (can be used to limit the number of pages)
	return true
}
func updateFromDB(p Post) (*Post, error){
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAZAZDZ6YVHU4TTHLL",
			"Z21ebf95ZIYGkyp43WpCcAvYiPYApVaBpsimpk7i",
			"", // a token will be created when the session it's used.
		  ),
	}))
	
	// A DynamoDB client
	svc := dynamodb.New(sess)
	
	// query
	user := p.User
	url := p.Url

	current_time := time.Now()
	// update values
	firstname := p.Firstname
	lastname := p.Lastname
	description := p.Description
	updatetime := current_time.Format("2006-01-02 15:04:05")
	// create the api params
	params := &dynamodb.UpdateItemInput{
		TableName: aws.String("Posts"),
		Key: map[string]*dynamodb.AttributeValue{
			"user": {
				S: aws.String(user),
			},
			"url": {
				S: aws.String(url),
			},
		},
		UpdateExpression: aws.String("set firstname=:r, lastname=:p, description=:a, updatetime=:e"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {S: aws.String(firstname)},
			":p": {S: aws.String(lastname)},
			":a": {S: aws.String(description)},
			":e": {S: aws.String(updatetime)},
			//":a": {SS: aws.StringSlice(actors)},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueAllNew),
	}
	// update the item
	resp, err := svc.UpdateItem(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
	}

	// unmarshal the dynamodb attribute values into a custom struct
	var cp Post
	err = dynamodbattribute.UnmarshalMap(resp.Attributes, &cp)

	// print the response data
	fmt.Printf("Updated Post = %+v\n", cp)

	return &cp, nil
}

func saveUserToDB(uname string, upass string) error{
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAZAZDZ6YVHU4TTHLL",
			"Z21ebf95ZIYGkyp43WpCcAvYiPYApVaBpsimpk7i",
			"", // a token will be created when the session it's used.
		  ),
	}))
	
	// A DynamoDB client
	svc := dynamodb.New(sess)
	var u User
	u.Username = uname
	u.Password = upass
	userAMAP, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		panic("Cannot marshal post into AttributeValue map")
	}
	// create the api params
	fmt.Println("marshalled struct: %+v", userAMAP)
	params := &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item:      userAMAP,
	}
	// put the item
	_, err = svc.PutItem(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		return err
	}
	// print the response data
	fmt.Println("Successly save an user to DB")
	return nil
}
func readUserFromDB(uname string, curTable string) (string, error){
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAZAZDZ6YVHU4TTHLL",
			"Z21ebf95ZIYGkyp43WpCcAvYiPYApVaBpsimpk7i",
			"", // a token will be created when the session it's used.
		  ),
	}))
	
	// A DynamoDB client
	svc := dynamodb.New(sess)

	// create the api params
	params := &dynamodb.GetItemInput{
		TableName: aws.String(curTable),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(uname),
			},
		},
	}
	// read the item
	resp, err := svc.GetItem(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		return "can read from DB", err
	}
	// unmarshal the dynamodb attribute values into a custom struct
	var user User
	err = dynamodbattribute.UnmarshalMap(resp.Item, &user)
	// print the response data
	fmt.Printf("Unmarshaled Post = %+v\n", user)
	return user.Password, nil
}