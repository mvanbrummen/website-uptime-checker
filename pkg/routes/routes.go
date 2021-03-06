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
		v1.DELETE("/websites", func(c *gin.Context) {
			var req types.DeleteWebsiteRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if probesDao.GetWebsite(req.URL) == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Website not found"})
				return
			}

			err := probesDao.DeleteWebsite(req.URL)

			if err != nil {
				internalServerError(c, err)
				return
			}

			c.JSON(http.StatusNoContent, nil)
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

		v1.GET("/websites/probes", func(c *gin.Context) {
			url := c.Query("url")

			probes := types.MapProbes(probesDao.GetWebsiteProbes(url))

			c.JSON(http.StatusOK, probes)
		})
	}

}

func internalServerError(c *gin.Context, err error) {
	log.Println(err)
	c.JSON(500, gin.H{
		"error": "Internal server error",
	})
}
