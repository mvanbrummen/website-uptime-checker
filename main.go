package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
	"github.com/mvanbrummen/website-uptime-probe/pkg/routes"
)

func main() {
	awsRegion := os.Getenv("AWS_REGION")
	dynamdoEndpoint := os.Getenv("AWS_DYNAMO_ENDPOINT")

	sess := session.Must(session.NewSession())
	_ = dynamo.New(sess, &aws.Config{
		Region:   aws.String(awsRegion),
		Endpoint: aws.String(dynamdoEndpoint),
	})

	r := gin.Default()

	routes.RegisterRoutes(r)

	r.Run(":8080")
}
