package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
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

	probesDao := db.NewProbesDao(dynamoDB)

	r := gin.Default()

	routes.RegisterRoutes(r, probesDao)

	// schedule the website checks
	client := resty.New()

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for range ticker.C {
			fmt.Println("Tick...")
			for _, w := range probesDao.GetWebsites() {
				url := strings.Replace(w.SK, db.WebsiteSKPrefix, "", 1)
				log.Println(url)

				resp, err := client.R().
					EnableTrace().
					Get(url)

				if err != nil {
					log.Println(err)
					continue
				}

				log.Println("Storing probe")
				p, err := probesDao.PutWebsiteProbe(url, fmt.Sprintf("%s", resp), resp.StatusCode(), int(resp.Time().Milliseconds()))

				log.Println(" probe is")
				log.Println(p)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	port := os.Getenv("PORT")

	if port != "" {
		endless.ListenAndServe(port, r)
	} else {
		endless.ListenAndServe(":8080", r)
	}
}
