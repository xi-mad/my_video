package tag

import (
	"errors"
	"github.com/xi-mad/my_video/commom"
	"github.com/xi-mad/my_video/object"
	"strings"
	"time"
)

type Tag struct {
	ID         int       `gorm:"id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"name" json:"name"`
	Order      int       `gorm:"order" json:"order"`
	CreateTime time.Time `gorm:"create_time;default:(datetime('now','localtime'))" json:"create_time"`
}

func (t *Tag) TableName() string {
	return "tag"
}

func createTag(model CreateTagModel) (err error) {
	errorMsg := ""
	names := strings.Split(model.Name, ",")
	for _, name := range names {
		if name == "" {
			continue
		}
		var count int64
		err := commom.DB.Model(&Tag{}).Where("name = ?", name).Count(&count).Error
		if err != nil {
			errorMsg += err.Error() + ";\n"
		}
		if count > 0 {
			errorMsg += name + "标签已存在;\n"
			continue
		}
		if err = commom.DB.Create(&Tag{
			Name:  name,
			Order: model.Order,
		}).Error; err != nil {
			errorMsg += err.Error() + ";\n"
		}
	}
	if errorMsg != "" {
		err = errors.New("以下内容已存在，其余内容已添加;\n" + errorMsg)
	}
	return
}

func deleteTag(model DeleteTagModel) (err error) {
	if err = commom.DB.Delete(&Tag{}, model.ID).Error; err != nil {
		return err
	}
	if err = commom.DB.Delete(&object.TagObject{}, "tag_id in ?", model.ID).Error; err != nil {
		return err
	}
	return
}

type ListTagModel struct {
	Name string `json:"name" form:"name"`
}

type CreateTagModel struct {
	Name  string `json:"name"`
	Order int    `json:"order"`
}

type UpdateTagModel struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Order int    `json:"order"`
}

type DeleteTagModel struct {
	ID []int `json:"id"`
}

func AutoMigrate() {
	_ = commom.DB.AutoMigrate(&Tag{})
}
