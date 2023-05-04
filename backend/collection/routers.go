package collection

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/xi-mad/my_video/common"
)

func Register(r *gin.RouterGroup) {
	r.GET("/", ListCollection)
	r.POST("/", CreateCollection)
	r.PUT("/", UpdateCollection)
	r.DELETE("/", DeleteCollection)
	r.GET("/options", OptionCollection)
	r.POST("/associate", AssociateCollection)
	r.GET("/detail", DetailCollection)
}

func ListCollection(c *gin.Context) {
	var model ListCollectionRequest
	if err := c.ShouldBindQuery(&model); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	res, total, err := QueryCollection(model)
	if err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}

	c.JSON(200, common.CommonResultSuccess(struct {
		Total int64                 `json:"total"`
		List  []ListCollectionModel `json:"data"`
	}{
		Total: total,
		List: lo.Map(res, func(item Collection, index int) ListCollectionModel {
			return ListCollectionModel{
				ID:         item.ID,
				Name:       item.Name,
				Cover:      item.Cover,
				Labels:     item.Labels,
				Actress:    item.Actress,
				CreateTime: item.CreateTime,
			}
		}),
	}))
}

func CreateCollection(c *gin.Context) {
	type CreateCollectionRequest struct {
		Name    string `json:"name"`
		Cover   string `json:"cover"`
		Labels  string `json:"labels"`
		Actress string `json:"actress"`
	}

	var req CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	create := Collection{
		Name:    req.Name,
		Cover:   req.Cover,
		Labels:  req.Labels,
		Actress: req.Actress,
	}
	if err := common.DB.Create(&create).Error; err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	c.JSON(200, common.CommonResultSuccess(nil))
}

func UpdateCollection(c *gin.Context) {
	type UpdateCollectionRequest struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Cover   string `json:"cover"`
		Labels  string `json:"labels"`
		Actress string `json:"actress"`
	}

	var req UpdateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	update := Collection{
		ID:      req.ID,
		Name:    req.Name,
		Cover:   req.Cover,
		Actress: req.Actress,
		Labels:  req.Labels,
	}
	if err := common.DB.Model(&Collection{}).Where("id = ?", req.ID).Updates(&update).Error; err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	c.JSON(200, common.CommonResultSuccess(nil))
}

func DeleteCollection(c *gin.Context) {
	type DeleteCollectionRequest struct {
		ID []int `json:"id"`
	}

	var req DeleteCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	if err := common.DB.Where("id in ?", req.ID).Delete(&Collection{}).Error; err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	c.JSON(200, common.CommonResultSuccess(nil))
}

func OptionCollection(c *gin.Context) {
	var req ListCollectionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	req.Page = 1
	req.PageSize = 100000
	res, _, err := QueryCollection(req)
	if err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	c.JSON(200, common.CommonResultSuccess(lo.Map(res, func(item Collection, index int) common.SelectOption {
		return common.SelectOption{
			Label: item.Name,
			Value: item.ID,
		}
	})))
}

func AssociateCollection(c *gin.Context) {
	type AssociateCollectionRequest struct {
		CollectionID int   `json:"collection_id"`
		ObjectID     []int `json:"object_id"`
	}

	var req AssociateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	oc := lo.Map(req.ObjectID, func(o int, index int) ObjectCollection {
		return ObjectCollection{
			ObjectID:     o,
			CollectionID: req.CollectionID,
			Order:        0,
		}
	})
	if err := common.DB.Where("collection_id = ? and object_id in ?", req.CollectionID, req.ObjectID).Delete(&ObjectCollection{}).Error; err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	if err := common.DB.Create(oc).Error; err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	c.JSON(200, common.CommonResultSuccess(nil))
}

func DetailCollection(c *gin.Context) {
	type DetailCollectionRequest struct {
		ID int `form:"id"`
	}
	var req DetailCollectionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	oc := make([]ObjectCollection, 0)
	if err := common.DB.Model(&ObjectCollection{}).Where("collection_id = ?", req.ID).Find(&oc).Error; err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	c.JSON(200, common.CommonResultSuccess(oc))

}
