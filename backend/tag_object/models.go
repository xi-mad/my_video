package tag_object

import "time"

type TagObject struct {
	TagID      int       `gorm:"tag_id" json:"tag_id"`
	ObjectID   int       `gorm:"object_id" json:"object_id"`
	CreateTime time.Time `gorm:"create_time;default:(datetime('now','localtime'))" json:"create_time"`
}

func (t *TagObject) TableName() string {
	return "tag_object"
}
