package object

import (
	"fmt"
	"github.com/xi-mad/my_video/actress_object"
	"github.com/xi-mad/my_video/commom"
	"github.com/xi-mad/my_video/tag_object"
	"gorm.io/gorm"
	"time"
)

type Object struct {
	ID          int       `gorm:"id;primaryKey;autoIncrement" json:"id"`
	Num         string    `gorm:"num" json:"num"`
	Type        string    `gorm:"type" json:"type"`
	Name        string    `gorm:"name" json:"name"`
	Md5Value    string    `gorm:"md5_value;index:md5_index" json:"md5_value"`
	Description string    `gorm:"description" json:"description"`
	Path        string    `gorm:"path" json:"path"`
	ExistNFO    bool      `gorm:"exist_nfo" json:"exist_nfo"`
	Rating      string    `gorm:"rating" json:"rating"`
	Release     string    `gorm:"release" json:"release"`
	Label       string    `gorm:"label" json:"label"`
	Magnet      string    `gorm:"magnet" json:"magnet"`
	Ext         string    `gorm:"ext" json:"ext"`
	ViewCount   int       `gorm:"view_count;default:0" json:"view_count"`
	CreateTime  time.Time `gorm:"create_time;default:(datetime('now','localtime'))" json:"create_time"`
}

func (o *Object) TableName() string {
	return "object"
}

type PlayObjectModel struct {
	Path string `json:"path" form:"path"`
}

type ListObjectRequest struct {
	Path     string `json:"path"`
	Actress  []int  `json:"actress"`
	Tag      []int  `json:"tag"`
	NFO      bool   `json:"nfo"`
	Tree     []int  `json:"tree"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

func (l *ListObjectRequest) Pageable() (limit, offset int) {
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

func QueryObject(model ListObjectRequest) (object []Object, total int64, err error) {
	var condition []func(db *gorm.DB) *gorm.DB
	if model.Path != "" {
		condition = append(condition, func(db *gorm.DB) *gorm.DB {
			return db.Where("path like ?", fmt.Sprintf("%%%s%%", model.Path))
		})
	}
	if len(model.Actress) > 0 {
		for i := range model.Actress {
			v := model.Actress[i]
			condition = append(condition, func(db *gorm.DB) *gorm.DB {
				return db.Where("id in (select object_id from actress_object where actress_id = ?)", v)
			})
		}
	}
	if len(model.Tag) != 0 {
		for i := range model.Tag {
			v := model.Tag[i]
			condition = append(condition, func(db *gorm.DB) *gorm.DB {
				return db.Where("id in (select object_id from tag_object where tag_id = ?)", v)
			})
		}
	}
	if len(model.Tree) != 0 {
		for i := range model.Tree {
			v := model.Tree[i]
			condition = append(condition, func(db *gorm.DB) *gorm.DB {
				return db.Where("id in (select object_id from tree_object where tree_id = ?)", v)
			})
		}
	}
	if model.NFO {
		condition = append(condition, func(db *gorm.DB) *gorm.DB {
			return db.Where("exist_nfo = ?", model.NFO)
		})
	}

	if err := commom.DB.Model(&Object{}).Scopes(commom.PaginateQuery(&model)).Scopes(condition...).Find(&object).Error; err != nil {
		return nil, 0, err
	}

	if err := commom.DB.Model(&Object{}).Scopes(condition...).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return
}

type ListObjectModel struct {
	ID          int       `json:"id"`
	Num         string    `json:"num"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Md5Value    string    `json:"md5_value"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	Magnet      string    `json:"magnet"`
	ExistNFO    bool      `json:"exist_nfo"`
	Rating      string    `json:"rating"`
	Release     string    `json:"release"`
	Label       string    `json:"label"`
	Ext         string    `json:"ext"`
	ViewCount   int       `json:"view_count"`
	Actress     []int     `json:"actress"`
	Tag         []int     `json:"tag"`
	Tree        []int     `json:"tree"`
	Thumbnail   string    `json:"thumbnail"`
	Fanart      string    `json:"fanart"`
	Poster      string    `json:"poster"`
	Thumb       string    `json:"thumb"`
	CreateTime  time.Time `json:"create_time"`
}

type CreateObjectModel struct {
	Type        string `json:"type"`
	Num         string `json:"num"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	ExistNFO    bool   `json:"exist_nfo"`
	Rating      string `json:"rating"`
	Release     string `json:"release"`
	Label       string `json:"label"`
	Magnet      string `json:"magnet"`
	Actress     []int  `json:"actress"`
	Tag         []int  `json:"tag"`
	Tree        []int  `json:"tree"`
}

type UpdateObjectModel struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	ExistNFO    bool   `json:"exist_nfo"`
	Rating      string `json:"rating"`
	Magnet      string `json:"magnet"`
	Actress     []int  `json:"actress"`
	Tag         []int  `json:"tag"`
	Tree        []int  `json:"tree"`
}

type DeleteObjectModel struct {
	ID []int `json:"id"`
}

type ScanObjectModel struct {
	Path    string `json:"path"`
	Actress []int  `json:"actress"`
	Tag     []int  `json:"tag"`
	Tree    []int  `json:"tree"`
}

type Thumbnail struct {
	ID         int       `gorm:"id;primaryKey;autoIncrement" json:"id"`
	ObjectID   int       `gorm:"object_id" json:"object_id"`
	Thumbnail  string    `gorm:"thumbnail" json:"thumbnail"`
	Fanart     string    `gorm:"fanart" json:"fanart"`
	Poster     string    `gorm:"poster" json:"poster"`
	Thumb      string    `gorm:"thumb" json:"thumb"`
	CreateTime time.Time `gorm:"create_time;default:(datetime('now','localtime'))" json:"create_time"`
}

func (t *Thumbnail) TableName() string {
	return "thumbnail"
}

func QueryThumbnail(ids []int) (res map[int]Thumbnail, err error) {
	var thumb []Thumbnail
	if err := commom.DB.Model(&Thumbnail{}).Find(&thumb, "object_id in ?", ids).Error; err != nil {
		return nil, err
	}
	res = make(map[int]Thumbnail)
	for _, v := range thumb {
		res[v.ObjectID] = v
	}
	return
}

type TreeObject struct {
	TreeID     int       `gorm:"tree_id" json:"tree_id"`
	ObjectID   int       `gorm:"object_id" json:"object_id"`
	Order      int       `gorm:"order" json:"order"`
	CreateTime time.Time `gorm:"create_time;default:(datetime('now','localtime'))" json:"create_time"`
}

func (t *TreeObject) TableName() string {
	return "tree_object"
}

func AutoMigrate() {
	_ = commom.DB.AutoMigrate(&Object{})
	_ = commom.DB.AutoMigrate(&Thumbnail{})
	_ = commom.DB.AutoMigrate(&tag_object.TagObject{})
	_ = commom.DB.AutoMigrate(&actress_object.ActressObject{})
	_ = commom.DB.AutoMigrate(&TreeObject{})
}

func SaveObjectTag(objectID int, tag []int) {
	commom.DB.Where("object_id = ?", objectID).Delete(&tag_object.TagObject{})
	if len(tag) == 0 {
		return
	}
	var to []tag_object.TagObject
	for _, v := range tag {
		to = append(to, tag_object.TagObject{TagID: v, ObjectID: objectID})
	}
	_ = commom.DB.Create(&to)
}

func SaveObjectActress(objectID int, actress []int) {
	commom.DB.Where("object_id = ?", objectID).Delete(&actress_object.ActressObject{})
	if len(actress) == 0 {
		return
	}
	var ao []actress_object.ActressObject
	for _, v := range actress {
		ao = append(ao, actress_object.ActressObject{ActressID: v, ObjectID: objectID})
	}
	_ = commom.DB.Create(&ao)
}

func SaveObjectTree(objectID int, tree []int) {
	commom.DB.Where("object_id = ?", objectID).Delete(&TreeObject{})
	if len(tree) == 0 {
		return
	}
	var to []TreeObject
	for _, v := range tree {
		to = append(to, TreeObject{TreeID: v, ObjectID: objectID})
	}
	_ = commom.DB.Create(&to)
}

func deleteRelationByObjectID(objectId []int, t interface{}) {
	if len(objectId) > 0 {
		commom.DB.Delete(t, "object_id in ?", objectId)
	}
}

func QueryTags(objectID []int) map[int][]int {
	if len(objectID) == 0 {
		return map[int][]int{}
	}
	var tag []tag_object.TagObject
	if err := commom.DB.Model(&tag_object.TagObject{}).Where("object_id in (?)", objectID).Find(&tag).Error; err != nil {
		return map[int][]int{}
	}
	res := make(map[int][]int)
	for _, v := range tag {
		res[v.ObjectID] = append(res[v.ObjectID], v.TagID)
	}
	return res
}

func QueryActress(objectID []int) map[int][]int {
	if len(objectID) == 0 {
		return map[int][]int{}
	}
	var actress []actress_object.ActressObject
	if err := commom.DB.Model(&actress_object.ActressObject{}).Where("object_id in (?)", objectID).Find(&actress).Error; err != nil {
		return map[int][]int{}
	}
	res := make(map[int][]int)
	for _, v := range actress {
		res[v.ObjectID] = append(res[v.ObjectID], v.ActressID)
	}
	return res
}

func QueryTree(objectID []int) map[int][]int {
	if len(objectID) == 0 {
		return map[int][]int{}
	}
	var tree []TreeObject
	if err := commom.DB.Model(&TreeObject{}).Where("object_id in (?)", objectID).Find(&tree).Error; err != nil {
		return map[int][]int{}
	}
	res := make(map[int][]int)
	for _, v := range tree {
		res[v.ObjectID] = append(res[v.ObjectID], v.TreeID)
	}
	return res
}

func PathExist(path string) (bool, error) {
	var count int64
	if err := commom.DB.Model(&Object{}).Where("path = ?", path).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
