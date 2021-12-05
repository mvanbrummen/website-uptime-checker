package main

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type widget struct {
	UserID int       `dynamo:",hash"`  // Hash key, a.k.a. partition key
	Time   time.Time `dynamo:",range"` // Range key, a.k.a. sort key

	Msg string `dynamo:"Message"` // Change name in the database
}

func main() {
	awsRegion := os.Getenv("AWS_REGION")
	dynamdoEndpoint := os.Getenv("AWS_DYNAMO_ENDPOINT")

	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{
		Region:   aws.String(awsRegion),
		Endpoint: aws.String(dynamdoEndpoint),
	})

	db.CreateTable("Widgets", &widget{}).Run()

	table := db.Table("Widgets")

	// put item
	w := widget{UserID: 613, Time: time.Now(), Msg: "hello"}
	err := table.Put(w).Run()

	if err != nil {
		panic(err)
	}

	// get the same item
	var result widget
	err = table.Get("UserID", w.UserID).
		Range("Time", dynamo.Equal, w.Time).
		One(&result)

	log.Println(result)

}
