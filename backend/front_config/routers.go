package front_config

import (
	"github.com/gin-gonic/gin"
	"github.com/xi-mad/my_video/common"
)

func Register(r *gin.RouterGroup) {
	r.GET("/", List)
}

func List(c *gin.Context) {
	c.JSON(200, common.CommonResultSuccess(Config{
		ObjectHeight:     common.DefaultConfig.App.ObjectHeight,
		CollectionHeight: common.DefaultConfig.App.CollectionHeight,
	}))
}
