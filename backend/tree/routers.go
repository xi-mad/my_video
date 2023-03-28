package tree

import (
	"github.com/gin-gonic/gin"
	"github.com/xi-mad/my_video/commom"
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
	err := commom.DB.Find(&trees).Error
	c.JSON(200, commom.CommonResultAuto(trees, err))
}

func CreateTree(c *gin.Context) {
	var model CreateTreeModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	Tree := Tree{
		Name:     model.Name,
		Order:    model.Order,
		ParentID: model.ParentID,
	}
	err := commom.DB.Create(&Tree).Error
	c.JSON(200, commom.CommonResultAuto(Tree, err))
}

func UpdateTree(c *gin.Context) {
	var model UpdateTreeModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	Tree := Tree{
		ID:    model.ID,
		Name:  model.Name,
		Order: model.Order,
	}
	err := commom.DB.Updates(&Tree).Error
	c.JSON(200, commom.CommonResultAuto(Tree, err))
}

func DeleteTree(c *gin.Context) {
	var model DeleteTreeModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	c.JSON(200, commom.CommonResultAuto(nil, deleteActress(model)))
}

func Options(c *gin.Context) {
	var trees []Tree
	err := commom.DB.Model(&Tree{}).Select("id, name, parent_id").Find(&trees).Error
	options := make([]commom.TreeSelectOption, 0)
	for _, tree := range trees {
		options = append(options, commom.TreeSelectOption{
			Value:    tree.ID,
			Label:    tree.Name,
			ParentID: tree.ParentID,
		})
	}
	c.JSON(200, commom.CommonResultAuto(options, err))
}
