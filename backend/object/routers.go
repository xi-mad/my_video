package object

import (
	"container/list"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/xi-mad/my_video/actress"
	"github.com/xi-mad/my_video/actress_object"
	"github.com/xi-mad/my_video/common"
	"github.com/xi-mad/my_video/media"
	"github.com/xi-mad/my_video/tag"
	"github.com/xi-mad/my_video/tag_object"
	"github.com/xi-mad/my_video/util"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Register(router *gin.RouterGroup) {
	router.POST("/list", ListObject)
	router.POST("/create", CreateObject)
	router.PUT("/update", UpdateObject)
	router.DELETE("/delete", DeleteObject)
	router.POST("/scan", ScanObject)
	router.GET("/play", PlayObject)
	router.GET("/log", LogObject)
	router.GET("/video", VideoObject)
	router.GET("/viewinc", ViewObjectInc)
	router.GET("/random", RandomObject)
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
	c.JSON(200, common.CommonResultSuccess(msg))
}

func ViewObjectInc(c *gin.Context) {
	if id, exist := c.GetQuery("id"); !exist {
		c.JSON(200, common.CommonResultFailed(errors.New("id is empty")))
		return
	} else {
		if err := common.DB.Model(&Object{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
			c.JSON(200, common.CommonResultFailed(err))
			return
		}
	}
	c.JSON(200, common.CommonResultSuccess(nil))
	return
}
func VideoObject(c *gin.Context) {
	var model PlayObjectModel
	if err := c.ShouldBindQuery(&model); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	c.File(model.Path)
	return
}

func PlayObject(c *gin.Context) {
	var model PlayObjectModel
	if err := c.ShouldBindQuery(&model); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	player := common.DefaultConfig.Player.Path
	if player == "" {
		c.JSON(200, common.CommonResultFailed(errors.New("player not set")))
		return
	}
	err := util.ExecCmd(common.DefaultConfig.Player.Path, model.Path)
	c.JSON(200, common.CommonResultAuto(nil, err))
}

func ListObject(c *gin.Context) {
	var model ListObjectRequest
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	object, total, err := QueryObject(model)
	if err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}

	ids := common.Map2ID(object)

	thumbMap, err := QueryThumbnail(ids)
	if err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	trees := QueryTree(ids)
	tags := QueryTags(ids)
	actressMap := QueryActress(ids)

	var result []ListObjectModel
	for _, v := range object {
		lom := ListObjectModel{
			ID:          v.ID,
			Num:         v.Num,
			Type:        v.Type,
			Name:        v.Name,
			Description: v.Description,
			Md5Value:    v.Md5Value,
			Path:        v.Path,
			Magnet:      v.Magnet,
			ExistNFO:    v.ExistNFO,
			Release:     v.Release,
			Label:       v.Label,
			Ext:         v.Ext,
			ViewCount:   v.ViewCount,
			Rating:      v.Rating,
			Actress:     []int{},
			Tag:         []int{},
			Tree:        []int{},
			CreateTime:  v.CreateTime,
		}
		lom.Thumbnail = thumbMap[v.ID].Thumbnail
		lom.Fanart = thumbMap[v.ID].Fanart
		lom.Thumb = thumbMap[v.ID].Thumb
		lom.Poster = thumbMap[v.ID].Poster
		if len(actressMap[v.ID]) > 0 {
			lom.Actress = actressMap[v.ID]
		}
		if len(tags[v.ID]) > 0 {
			lom.Tag = tags[v.ID]
		}
		if len(trees[v.ID]) > 0 {
			lom.Tree = trees[v.ID]
		}
		result = append(result, lom)
	}
	c.JSON(200, common.CommonResultSuccess(struct {
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
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	c.JSON(200, common.CommonResultAuto(createObject(model)))
}

func createObject(model CreateObjectModel) (object Object, err error) {
	if exist, err := PathExist(model.Path); err != nil {
		return object, err
	} else if exist {
		return object, errors.New("file already exist, skipped")
	}
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
		Num:         model.Num,
		ExistNFO:    model.ExistNFO,
		Release:     model.Release,
		Label:       model.Label,
	}
	err = common.DB.Save(&object).Error
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
	err = common.DB.Save(&thumb).Error
	return
}

func UpdateObject(c *gin.Context) {
	var model UpdateObjectModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	oldObj := Object{}
	err := common.DB.Find(&oldObj, model.ID).Error
	if err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	object := Object{
		ID:          model.ID,
		Type:        model.Type,
		Name:        model.Name,
		Description: model.Description,
		Path:        model.Path,
		Rating:      model.Rating,
		Magnet:      model.Magnet,
	}
	if oldObj.Path != model.Path {
		_, b64, err := detail(model.Path)
		if err != nil {
			c.JSON(200, common.CommonResultFailed(err))
			return
		}
		thumb := Thumbnail{
			ID:        oldObj.ID,
			Thumbnail: b64,
		}
		err = common.DB.Updates(&thumb).Error
	}

	err = common.DB.Updates(&object).Error
	SaveObjectActress(object.ID, model.Actress)
	SaveObjectTag(object.ID, model.Tag)
	SaveObjectTree(object.ID, model.Tree)
	c.JSON(200, common.CommonResultAuto(object, err))
}

func DeleteObject(c *gin.Context) {
	var model DeleteObjectModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	err := common.DB.Delete(&Object{}, model.ID).Error
	if err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	deleteRelationByObjectID(model.ID, &actress_object.ActressObject{})
	deleteRelationByObjectID(model.ID, &tag_object.TagObject{})
	deleteRelationByObjectID(model.ID, &TreeObject{})

	err = common.DB.Where("object_id in ?", model.ID).Delete(&Thumbnail{}).Error
	c.JSON(200, common.CommonResultAuto(nil, err))
}

func ScanObject(c *gin.Context) {
	var model ScanObjectModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	ap, err := filepath.Abs(model.Path)
	if err != nil {
		c.JSON(200, common.CommonResultFailed(err))
		return
	}
	model.Path = ap
	go scanObject(model)
	c.JSON(200, common.CommonResultAuto(ap, err))
}

func scanObject(model ScanObjectModel) {
	var supportExt = []string{".mp4", ".avi", ".mkv", ".wmv", ".rmvb", ".flv", ".mov", ".mpg", ".mpeg", ".rm", ".asf", ".divx", ".vob", ".m4v", ".3gp", ".3g2", ".dat", ".m2ts", ".m2v", ".m4a", ".mj2", ".mjpg", ".mjpeg", ".moov", ".mpv", ".nut", ".ogg", ".ogm", ".qt", ".swf", ".ts", ".xvid"}

	errMsgFormat := "path: %s, error: %s"
	if err := filepath.Walk(model.Path, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			findLog.PushBack(fmt.Sprintf(errMsgFormat, fpath, err.Error()))
			return nil
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(fpath)
		for _, v := range supportExt {
			if v == ext {
				findLog.PushBack(fmt.Sprintf("find file: %s", fpath))
				filePrefix := fpath[0 : len(fpath)-len(path.Ext(fpath))]
				nfo := filePrefix + ".nfo"
				if _, err := os.Stat(nfo); err == nil {
					findLog.PushBack(fmt.Sprintf("find nfo: %s", nfo))
					b, err := os.ReadFile(nfo)
					if err != nil {
						findLog.PushBack(fmt.Sprintf("open nfo: %s, faild: %s", nfo, err))
						continue
					}
					m := media.NFO{}
					if err := xml.Unmarshal(b, &m); err != nil {
						findLog.PushBack(fmt.Sprintf("parse nfo: %s, faild: %s", nfo, err))
						continue
					}
					tags := lo.Map(m.Tag, func(x media.Inner, index int) string {
						return x.Inner
					})
					actors := lo.Map(m.Actor, func(x media.Actor, index int) string {
						return x.Name
					})
					tagsId := tag.CreateTags(tags)
					actorsId := actress.CreateActresses(actors)
					rating, err := strconv.ParseFloat(m.Rating, 64)
					if err != nil {
						rating = 0
					}

					obj, err := createObject(CreateObjectModel{
						Path:        fpath,
						Actress:     actorsId,
						Tag:         tagsId,
						ExistNFO:    true,
						Description: strings.ReplaceAll(m.Plot.Inner, " ", ""),
						Name:        strings.ReplaceAll(m.Title.Inner, " ", ""),
						Num:         m.Number,
						Release:     m.Release,
						Label:       m.Label,
						Rating:      RateScore(rating),
					})
					if err != nil {
						findLog.PushBack(fmt.Sprintf(errMsgFormat, fpath, err.Error()))
						continue
					}
					dir := filepath.Dir(fpath)
					fanart := dir + "\\fanart.jpg"
					poster := dir + "\\poster.jpg"
					thumb := dir + "\\thumb.jpg"
					fanartB64, err := util.Image2Base64(fanart)
					if err == nil {
						common.DB.Model(&Thumbnail{}).Where("object_id = ?", obj.ID).
							Update("fanart", fanartB64)
					}
					posterB64, err := util.Image2Base64(poster)
					if err == nil {
						common.DB.Model(&Thumbnail{}).Where("object_id = ?", obj.ID).
							Update("poster", posterB64)
					}
					thumbB64, err := util.Image2Base64(thumb)
					if err == nil {
						common.DB.Model(&Thumbnail{}).Where("object_id = ?", obj.ID).
							Update("thumb", thumbB64)
					}
				} else {
					_, err := createObject(CreateObjectModel{
						Path:    fpath,
						Actress: model.Actress,
						Tag:     model.Tag,
						Tree:    model.Tree,
					})
					if err != nil {
						findLog.PushBack(fmt.Sprintf(errMsgFormat, fpath, err.Error()))
					}
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

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	if err = thumbnail(path, fsize, suffix); err != nil {
		return
	}
	if b64, err = thumbnailB64(fname, suffix); err != nil {
		return
	}
	return
}

func thumbnailB64(fname, suffix string) (b64 string, err error) {
	prefix := fname[0 : len(fname)-len(path.Ext(fname))]
	if thumbPath, err := filepath.Abs("./temp/" + prefix + fmt.Sprintf("_%s.jpg", suffix)); err != nil {
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

func thumbnail(path string, fsize int64, suffix string) (err error) {
	thumbnailConf := common.DefaultConfig.Thumbnail
	width, col, row := thumbnailConf.Width, thumbnailConf.Col, thumbnailConf.Row
	for _, optional := range thumbnailConf.Optional {
		if fsize <= optional.FSizeLess {
			width, col, row = optional.Width, optional.Col, optional.Row
			break
		}
	}
	if width <= 0 {
		width = 2048
	}
	if col <= 0 {
		col = 4
	}
	if row <= 0 {
		row = 4
	}
	args := []string{"-P",
		"-w", fmt.Sprintf("%d", width),
		"-c", fmt.Sprintf("%d", col),
		"-r", fmt.Sprintf("%d", row),
		"-f", thumbnailConf.Font,
		"-o", fmt.Sprintf("_%s.jpg", suffix),
		"-O", "./temp",
		path}
	err = util.ExecCmd(thumbnailConf.Mtn, args...)
	if err != nil {
		log.Printf("mtn error: %s \n", err)
	}
	return
}

func RandomObject(c *gin.Context) {
	var obj Object
	if err := common.DB.Model(&Object{}).Order("random()").Limit(1).First(&obj).Error; err != nil {
		c.JSON(http.StatusOK, common.CommonResultFailed(err))
		return
	}
	if err := common.DB.Model(&Object{}).Where("id = ?", obj.ID).Update("view_count", obj.ViewCount+1).Error; err != nil {
		c.JSON(http.StatusOK, common.CommonResultFailed(err))
		return
	}
	c.JSON(http.StatusOK, common.CommonResultSuccess(obj))
}
