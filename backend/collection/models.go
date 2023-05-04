package collection

import (
	"github.com/xi-mad/my_video/common"
	"gorm.io/gorm"
	"time"
)

type Collection struct {
	ID         int       `gorm:"id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"name" json:"name"`
	Cover      string    `gorm:"cover" json:"cover"`
	Labels     string    `gorm:"labels" json:"labels"`
	Actress    string    `gorm:"actress" json:"actress"`
	CreateTime time.Time `gorm:"create_time;default:(datetime('now','localtime'))" json:"create_time"`
}

func (c *Collection) TableName() string {
	return "collection"
}

type ObjectCollection struct {
	ObjectID     int       `gorm:"object_id" json:"object_id"`
	CollectionID int       `gorm:"collection_id" json:"collection_id"`
	Order        int       `gorm:"order" json:"order"`
	CreateTime   time.Time `gorm:"create_time;default:(datetime('now','localtime'))" json:"create_time"`
}

func (o *ObjectCollection) TableName() string {
	return "object_collection"
}

func AutoMigrate() {
	_ = common.DB.AutoMigrate(&Collection{})
	_ = common.DB.AutoMigrate(&ObjectCollection{})
}

type ListCollectionRequest struct {
	Name     string `form:"name"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type ListCollectionModel struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cover      string    `json:"cover"`
	Labels     string    `json:"labels"`
	Actress    string    `json:"actress"`
	CreateTime time.Time `json:"create_time"`
}

func (l *ListCollectionRequest) Pageable() (limit, offset int) {
	if l.Page == 0 {
		l.Page = 1
	}
	if l.PageSize == 0 {
		l.PageSize = 10
	}
	limit = l.PageSize
	offset = (l.Page - 1) * l.PageSize
	return
}

func QueryCollection(model ListCollectionRequest) (res []Collection, total int64, err error) {
	var condition []func(db *gorm.DB) *gorm.DB
	if model.Name != "" {
		condition = append(condition, func(db *gorm.DB) *gorm.DB {
			return db.Where("name like ?", "%"+model.Name+"%")
		})
	}

	if err := common.DB.Model(&Collection{}).Scopes(common.PaginateQuery(&model)).Scopes(condition...).Find(&res).Error; err != nil {
		return nil, 0, err
	}

	if err := common.DB.Model(&Collection{}).Scopes(condition...).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return
}
