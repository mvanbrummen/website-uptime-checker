package main

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
	"github.com/mvanbrummen/website-uptime-probe/pkg/db"
	"github.com/mvanbrummen/website-uptime-probe/pkg/routes"
)

func main() {
	awsRegion := os.Getenv("AWS_REGION")
	dynamdoEndpoint := os.Getenv("AWS_DYNAMO_ENDPOINT")

	sess := session.Must(session.NewSession())
	dynamoDB := dynamo.New(sess, &aws.Config{
		Region:   aws.String(awsRegion),
		Endpoint: aws.String(dynamdoEndpoint),
	})

	// test data
	dynamoDB.CreateTable("WebsiteProbes", &db.WebsiteProbes{}).Run()

	t := dynamoDB.Table("WebsiteProbes")

	err := t.Put(db.WebsiteProbes{
		PK: "WEBSITE#",
		SK: "WEBSITE#https://google.com",

		WebsiteAttributes: db.WebsiteAttributes{
			ProbeScheduleMinutes: 1,
			CreatedDate:          time.Now(),
			Active:               true,
		},
	}).Run()

	if err != nil {
		panic(err)
	}
	// end test data

	probesDao := db.NewProbesDao(dynamoDB)

	r := gin.Default()

	routes.RegisterRoutes(r, probesDao)

	r.Run(":8080")
}
