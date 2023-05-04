package tree

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/xi-mad/my_video/common"
)

func Register(router *gin.RouterGroup) {
	router.GET("/list", ListTree)
	router.POST("/create", CreateTree)
	router.PUT("/update", UpdateTree)
	router.DELETE("/delete", DeleteTree)
	router.GET("/options", Options)
}

func ListTree(c *gin.Context) {
	var trees []Tree
	err := common.DB.Find(&trees).Error
	c.JSON(200, common.CommonResultAuto(trees, err))
}

func CreateTree(c *gin.Context) {
	var model CreateTreeModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	Tree := Tree{
		Name:     model.Name,
		Order:    model.Order,
		ParentID: model.ParentID,
	}
	err := common.DB.Create(&Tree).Error
	c.JSON(200, common.CommonResultAuto(Tree, err))
}

func UpdateTree(c *gin.Context) {
	var model UpdateTreeModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	Tree := Tree{
		ID:    model.ID,
		Name:  model.Name,
		Order: model.Order,
	}
	err := common.DB.Updates(&Tree).Error
	c.JSON(200, common.CommonResultAuto(Tree, err))
}

func DeleteTree(c *gin.Context) {
	var model DeleteTreeModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	c.JSON(200, common.CommonResultAuto(nil, deleteActress(model)))
}

func Options(c *gin.Context) {
	var trees []Tree
	if err := common.DB.Model(&Tree{}).Select("id, name, parent_id").Find(&trees).Error; err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	c.JSON(200, common.CommonResultSuccess(lo.Map(trees, func(tree Tree, index int) common.TreeSelectOption {
		return common.TreeSelectOption{
			Value:    tree.ID,
			Label:    tree.Name,
			ParentID: tree.ParentID,
		}
	})))
}
