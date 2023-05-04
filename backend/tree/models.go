package tree

import (
	"errors"
	"github.com/xi-mad/my_video/common"
	"github.com/xi-mad/my_video/object"
	"time"
)

type Tree struct {
	ID         int       `gorm:"id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"name" json:"name"`
	ParentID   int       `gorm:"parent_id" json:"parent_id"`
	Order      int       `gorm:"order" json:"order"`
	CreateTime time.Time `gorm:"create_time;default:(datetime('now','localtime'))" json:"create_time"`
	Children   []Tree    `gorm:"-" json:"children"`
}

func (t *Tree) TableName() string {
	return "tree"
}

func listTree(parentID int) (trees []Tree, err error) {
	err = common.DB.Model(&Tree{}).Find(&trees, "parent_id = ?", parentID).Error
	return
}

func deleteActress(model DeleteTreeModel) (err error) {
	trees, err := listTree(model.ID)
	if err != nil {
		return err
	}
	if len(trees) > 0 {
		return errors.New("该节点下有子节点，无法删除")
	}
	if err = common.DB.Delete(&Tree{}, model.ID).Error; err != nil {
		return err
	}
	if err = common.DB.Delete(&object.TreeObject{}, "tree_id = ?", model.ID).Error; err != nil {
		return err
	}
	return
}

type CreateTreeModel struct {
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
	Order    int    `json:"order"`
}

type UpdateTreeModel struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
	Order    int    `json:"order"`
}

type DeleteTreeModel struct {
	ID int `json:"id"`
}

func AutoMigrate() {
	_ = common.DB.AutoMigrate(&Tree{})
}
