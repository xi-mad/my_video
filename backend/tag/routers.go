package tag

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/xi-mad/my_video/common"
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
		c.JSON(200, common.CommonResultFailed(err))
		return
	}

	var condition []func(db *gorm.DB) *gorm.DB
	if model.Name != "" {
		condition = append(condition, func(db *gorm.DB) *gorm.DB {
			return db.Where("name like ?", fmt.Sprintf("%%%s%%", model.Name))
		})
	}
	tag := make([]Tag, 0)
	err := common.DB.Model(&Tag{}).Scopes(condition...).Find(&tag).Error
	c.JSON(200, common.CommonResultAuto(tag, err))
}

func CreateTag(c *gin.Context) {
	var model CreateTagModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	c.JSON(200, common.CommonResultAuto(nil, createTag(model)))
}

func UpdateTag(c *gin.Context) {
	var model UpdateTagModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	tag := Tag{
		ID:    model.ID,
		Name:  model.Name,
		Order: model.Order,
	}
	err := common.DB.Updates(&tag).Error
	c.JSON(200, common.CommonResultAuto(tag, err))
}

func DeleteTag(c *gin.Context) {
	var model DeleteTagModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	c.JSON(200, common.CommonResultAuto(nil, deleteTag(model)))
}

func Options(c *gin.Context) {
	var tags []Tag
	if err := common.DB.Model(&Tag{}).Select("id, name").Find(&tags).Error; err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	c.JSON(200, common.CommonResultSuccess(lo.Map(tags, func(item Tag, index int) common.SelectOption {
		return common.SelectOption{
			Value: item.ID,
			Label: item.Name,
		}
	})))
}
