package routes

import (
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

			c.JSON(200, websites)
		})
	}

}
