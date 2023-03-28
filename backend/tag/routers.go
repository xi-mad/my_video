package tag

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xi-mad/my_video/commom"
	"gorm.io/gorm"
)

func Register(router *gin.RouterGroup) {
	router.GET("/list", ListTag)
	router.POST("/create", CreateTag)
	router.PUT("/update", UpdateTag)
	router.DELETE("/delete", DeleteTag)
	router.GET("/options", Options)
}

func ListTag(c *gin.Context) {
	var model ListTagModel
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
	tag := make([]Tag, 0)
	err := commom.DB.Model(&Tag{}).Scopes(condition...).Find(&tag).Error
	c.JSON(200, commom.CommonResultAuto(tag, err))
}

func CreateTag(c *gin.Context) {
	var model CreateTagModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	c.JSON(200, commom.CommonResultAuto(nil, createTag(model)))
}

func UpdateTag(c *gin.Context) {
	var model UpdateTagModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	tag := Tag{
		ID:    model.ID,
		Name:  model.Name,
		Order: model.Order,
	}
	err := commom.DB.Updates(&tag).Error
	c.JSON(200, commom.CommonResultAuto(tag, err))
}

func DeleteTag(c *gin.Context) {
	var model DeleteTagModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	c.JSON(200, commom.CommonResultAuto(nil, deleteTag(model)))
}

func Options(c *gin.Context) {
	var tags []Tag
	err := commom.DB.Model(&Tag{}).Select("id, name").Find(&tags).Error
	options := make([]commom.SelectOption, 0)
	for _, tag := range tags {
		options = append(options, commom.SelectOption{
			Value: tag.ID,
			Label: tag.Name,
		})
	}
	c.JSON(200, commom.CommonResultAuto(options, err))
}
