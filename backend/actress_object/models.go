package actress_object

import "time"

type ActressObject struct {
	ActressID  int       `gorm:"actress_id" json:"actress_id"`
	ObjectID   int       `gorm:"object_id" json:"object_id"`
	CreateTime time.Time `gorm:"create_time;default:(datetime('now','localtime'))" json:"create_time"`
}

func (a *ActressObject) TableName() string {
	return "actress_object"
}
