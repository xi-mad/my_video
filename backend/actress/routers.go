package actress

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xi-mad/my_video/commom"
	"gorm.io/gorm"
)

func Register(router *gin.RouterGroup) {
	router.GET("/list", ListActress)
	router.POST("/create", CreateActress)
	router.PUT("/update", UpdateActress)
	router.DELETE("/delete", DeleteActress)
	router.GET("/options", Options)
}

func ListActress(c *gin.Context) {
	var model ListActressModel
	if err := c.ShouldBindQuery(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}

	var condition []func(db *gorm.DB) *gorm.DB
	if model.Name != "" {
		condition = append(condition, func(db *gorm.DB) *gorm.DB {
			return db.Where("name like ?", fmt.Sprintf("%%%s%%", model.Name))
		})
	}
	actress := make([]Actress, 0)
	err := commom.DB.Model(&Actress{}).Scopes(condition...).Find(&actress).Error
	c.JSON(200, commom.CommonResultAuto(actress, err))
}

func CreateActress(c *gin.Context) {
	var model CreateActressModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	c.JSON(200, commom.CommonResultAuto(nil, createActress(model)))
}

func UpdateActress(c *gin.Context) {
	var model UpdateActressModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	actress := Actress{
		ID:    model.ID,
		Name:  model.Name,
		Order: model.Order,
	}
	err := commom.DB.Updates(&actress).Error
	c.JSON(200, commom.CommonResultAuto(actress, err))
}

func DeleteActress(c *gin.Context) {
	var model DeleteActressModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	c.JSON(200, commom.CommonResultAuto(nil, deleteActress(model)))
}

func Options(c *gin.Context) {
	var actress []Actress
	err := commom.DB.Model(&Actress{}).Select("id, name").Find(&actress).Error
	options := make([]commom.SelectOption, 0)
	for _, act := range actress {
		options = append(options, commom.SelectOption{
			Value: act.ID,
			Label: act.Name,
		})
	}
	c.JSON(200, commom.CommonResultAuto(options, err))
}
