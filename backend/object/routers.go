package object

import (
	"container/list"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xi-mad/my_video/commom"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func Register(router *gin.RouterGroup) {
	router.POST("/list", ListObject)
	router.POST("/create", CreateObject)
	router.PUT("/update", UpdateObject)
	router.DELETE("/delete", DeleteObject)
	router.POST("/scan", ScanObject)
	router.GET("/play", PlayObject)
	router.GET("/log", LogObject)

}

var findLog = list.New()

func LogObject(c *gin.Context) {
	l := findLog.Len()
	msg := ""
	for i := 0; i < l; i++ {
		front := findLog.Front()
		msg += fmt.Sprintf("%s\n", front.Value)
		findLog.Remove(front)
	}
	c.JSON(200, commom.CommonResultSuccess(msg))
}

func PlayObject(c *gin.Context) {
	var model PlayObjectModel
	if err := c.ShouldBindQuery(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	player := commom.DefaultConfig.Player.Path
	if player == "" {
		c.JSON(200, commom.CommonResultFailed(errors.New("player not set")))
		return
	}
	exec.Command(commom.DefaultConfig.Player.Path, model.Path).Start()
	c.JSON(200, commom.CommonResultSuccess(nil))
}

func ListObject(c *gin.Context) {
	var model ListObjectRequest
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	object, total, err := QueryObject(model)
	if err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}

	ids := commom.Map2ID(object)

	thumbMap, err := QueryThumbnail(ids)
	if err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	trees := QueryTree(ids)
	tags := QueryTags(ids)
	actress := QueryActress(ids)

	var result []ListObjectModel
	for _, v := range object {
		lom := ListObjectModel{
			ID:          v.ID,
			Type:        v.Type,
			Name:        v.Name,
			Description: v.Description,
			Path:        v.Path,
			Magnet:      v.Magnet,
			Ext:         v.Ext,
			Actress:     []int{},
			Tag:         []int{},
			Tree:        []int{},
			CreateTime:  v.CreateTime,
		}
		lom.Thumbnail = thumbMap[v.ID]
		if len(actress[v.ID]) > 0 {
			lom.Actress = actress[v.ID]
		}
		if len(tags[v.ID]) > 0 {
			lom.Tag = tags[v.ID]
		}
		if len(trees[v.ID]) > 0 {
			lom.Tree = trees[v.ID]
		}
		result = append(result, lom)

	}
	c.JSON(200, commom.CommonResultSuccess(struct {
		Total int64             `json:"total"`
		List  []ListObjectModel `json:"data"`
	}{
		Total: total,
		List:  result,
	}))
}

func CreateObject(c *gin.Context) {
	var model CreateObjectModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	c.JSON(200, commom.CommonResultAuto(createObject(model)))
}

func createObject(model CreateObjectModel) (object Object, err error) {
	fname, b64, err := detail(model.Path)
	if err != nil {
		return
	}
	object = Object{
		Type:        model.Type,
		Name:        fname,
		Description: model.Description,
		Path:        model.Path,
		Magnet:      model.Magnet,
	}
	err = commom.DB.Save(&object).Error
	if err != nil {
		return
	}
	SaveObjectActress(object.ID, model.Actress)
	SaveObjectTag(object.ID, model.Tag)
	SaveObjectTree(object.ID, model.Tree)
	thumb := Thumbnail{
		ObjectID:  object.ID,
		Thumbnail: b64,
	}
	err = commom.DB.Save(&thumb).Error
	return
}

func UpdateObject(c *gin.Context) {
	var model UpdateObjectModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	oldObj := Object{}
	err := commom.DB.Find(&oldObj, model.ID).Error
	if err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	if oldObj.Path != model.Path {
		_, b64, err := detail(model.Path)
		if err != nil {
			c.JSON(200, commom.CommonResultFailed(err))
			return
		}
		thumb := Thumbnail{
			ID:        oldObj.ID,
			Thumbnail: b64,
		}
		err = commom.DB.Updates(&thumb).Error
	}

	object := Object{
		ID:          model.ID,
		Type:        model.Type,
		Name:        model.Name,
		Description: model.Description,
		Path:        model.Path,
		Magnet:      model.Magnet,
	}
	err = commom.DB.Updates(&object).Error
	SaveObjectActress(object.ID, model.Actress)
	SaveObjectTag(object.ID, model.Tag)
	SaveObjectTree(object.ID, model.Tree)
	c.JSON(200, commom.CommonResultAuto(object, err))
}

func DeleteObject(c *gin.Context) {
	var model DeleteObjectModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	err := commom.DB.Delete(&Object{}, model.ID).Error
	if err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	deleteRelationByObjectID(model.ID, &ActressObject{})
	deleteRelationByObjectID(model.ID, &TagObject{})
	deleteRelationByObjectID(model.ID, &TreeObject{})

	err = commom.DB.Where("object_id in ?", model.ID).Delete(&Thumbnail{}).Error
	c.JSON(200, commom.CommonResultAuto(nil, err))
}

func ScanObject(c *gin.Context) {
	var model ScanObjectModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	ap, err := filepath.Abs(model.Path)
	if err != nil {
		c.JSON(200, commom.CommonResultFailed(err))
		return
	}
	model.Path = ap
	go scanObject(model)
	c.JSON(200, commom.CommonResultAuto(ap, err))
}

func scanObject(model ScanObjectModel) {
	var supportExt = []string{".mp4", ".avi", ".mkv", ".wmv", ".rmvb", ".flv", ".mov", ".mpg", ".mpeg", ".rm", ".asf", ".divx", ".vob", ".m4v", ".3gp", ".3g2", ".dat", ".m2ts", ".m2v", ".m4a", ".mj2", ".mjpg", ".mjpeg", ".moov", ".mpv", ".nut", ".ogg", ".ogm", ".qt", ".swf", ".ts", ".xvid"}

	errMsgFormat := "path: %s, error: %s"
	if err := filepath.Walk(model.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			findLog.PushBack(fmt.Sprintf(errMsgFormat, path, err.Error()))
			return nil
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		for _, v := range supportExt {
			if v == ext {
				findLog.PushBack(fmt.Sprintf("find file: %s", path))
				_, err := createObject(CreateObjectModel{
					Path:    path,
					Actress: model.Actress,
					Tag:     model.Tag,
					Tree:    model.Tree,
				})
				if err != nil {
					findLog.PushBack(fmt.Sprintf(errMsgFormat, path, err.Error()))
				}
				break
			}
		}
		return nil
	}); err != nil {
		findLog.PushBack(fmt.Sprintf(errMsgFormat, model.Path, err.Error()))
	}
	findLog.PushBack("扫描完成。")
}

func detail(path string) (fname, b64 string, err error) {
	if path, err = filepath.Abs(path); err != nil {
		return
	}
	var fsize int64
	if fi, err := os.Stat(path); err != nil || fi.IsDir() {
		return "", "", errors.New("path is not a file")
	} else {
		fname = fi.Name()
		fsize = fi.Size() / 1024 / 1024
	}
	if err = thumbnail(path, fsize); err != nil {
		return
	}
	if b64, err = thumbnailB64(fname); err != nil {
		return
	}
	return
}

func thumbnailB64(fname string) (b64 string, err error) {
	suffix := path.Ext(fname)
	prefix := fname[0 : len(fname)-len(suffix)]
	if thumbPath, err := filepath.Abs("./temp/" + prefix + "_s.jpg"); err != nil {
		return "", err
	} else {
		if f, err := os.Open(thumbPath); err != nil {
			return "", err
		} else {
			defer func() {
				_ = f.Close()
				_ = os.Remove(thumbPath)
			}()
			if bytes, err := io.ReadAll(f); err != nil {
				return "", err
			} else {
				b64 = base64.StdEncoding.EncodeToString(bytes)
			}
		}
	}
	return
}

func thumbnail(path string, fsize int64) (err error) {
	/*
		<= 50MB 2 * 2
		<= 100MB 3 * 3
		<= 500MB 4 * 4
		<= 1024MB 5 * 5
		else 6 * 6
	*/
	col, row := 2, 2
	if fsize > 50 {
		col, row = 3, 3
	} else if fsize > 100 {
		col, row = 4, 4
	} else if fsize > 500 {
		col, row = 5, 5
	} else if fsize > 1024 {
		col, row = 6, 6
	}

	thumbnailConf := commom.DefaultConfig.Thumbnail
	args := []string{"-P", "-w", thumbnailConf.Width, "-c", fmt.Sprintf("%d", col), "-r", fmt.Sprintf("%d", row), "-f", thumbnailConf.Font, "-O", "./temp", path}
	_, err = exec.Command(thumbnailConf.Mtn, args...).Output()
	return
}
