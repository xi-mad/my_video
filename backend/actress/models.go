package actress

import (
	"errors"
	"github.com/samber/lo"
	"github.com/xi-mad/my_video/actress_object"
	"github.com/xi-mad/my_video/commom"
	"strings"
	"time"
)

type Actress struct {
	ID         int       `gorm:"id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"name" json:"name"`
	Order      int       `gorm:"order" json:"order"`
	CreateTime time.Time `gorm:"create_time;default:(datetime('now','localtime'))" json:"create_time"`
}

func (a *Actress) TableName() string {
	return "actress"
}

type ListActressModel struct {
	Name string `json:"name" form:"name"`
}

func createActress(model CreateActressModel) (err error) {
	errorMsg := ""
	names := strings.Split(model.Name, ",")
	for _, name := range names {
		if name == "" {
			continue
		}
		var count int64
		err := commom.DB.Model(&Actress{}).Where("name = ?", name).Count(&count).Error
		if err != nil {
			errorMsg += err.Error() + ";\n"
		}
		if count > 0 {
			errorMsg += name + "演员已存在;\n"
			continue
		}
		if err = commom.DB.Create(&Actress{
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

func CreateActresses(actresses []string) []int {
	ids := make([]int, 0)
	exists := make([]Actress, 0)
	commom.DB.Model(&Actress{}).Where("name in ?", actresses).Find(&exists)
	existMap := lo.Associate(exists, func(f Actress) (string, int) {
		return f.Name, f.ID
	})
	notExists := make([]Actress, 0)
	for _, actress := range actresses {
		if id, ok := existMap[actress]; ok {
			ids = append(ids, id)
		} else {
			notExists = append(notExists, Actress{
				Name:  actress,
				Order: 0,
			})
		}
	}
	if len(notExists) > 0 {
		commom.DB.Create(&notExists)
		for _, actress := range notExists {
			ids = append(ids, actress.ID)
		}
	}
	return ids
}

func deleteActress(model DeleteActressModel) (err error) {
	if err = commom.DB.Delete(&Actress{}, model.ID).Error; err != nil {
		return err
	}
	if err = commom.DB.Delete(&actress_object.ActressObject{}, "actress_id in ?", model.ID).Error; err != nil {
		return err
	}
	return
}

type CreateActressModel struct {
	Name  string `json:"name"`
	Order int    `json:"order"`
}

type UpdateActressModel struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Order int    `json:"order"`
}

type DeleteActressModel struct {
	ID []int `json:"id"`
}

func AutoMigrate() {
	_ = commom.DB.AutoMigrate(&Actress{})
}
