package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mvanbrummen/website-uptime-probe/pkg/db"
	"github.com/mvanbrummen/website-uptime-probe/pkg/types"
)

func RegisterRoutes(r *gin.Engine, probesDao *db.ProbesDao) {

	v1 := r.Group("/v1")
	{
		v1.GET("/health", HealthRoute)
		v1.GET("/websites", func(c *gin.Context) {
			w := probesDao.GetWebsites()

			websites := types.MapWebsites(w)

			c.JSON(http.StatusOK, websites)
		})
		v1.PUT("/websites", func(c *gin.Context) {
			var req types.CreateWebsiteRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			website, err := probesDao.PutWebsite(req.URL, req.ProbeScheduleMinutes)

			if err != nil {
				internalServerError(c, err)
				return
			}

			c.JSON(http.StatusCreated, types.MapWebsite(website))
		})
	}

}

func internalServerError(c *gin.Context, err error) {
	log.Println(err)
	c.JSON(500, gin.H{
		"error": "Internal server error",
	})
}
