package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xi-mad/my_video/actress"
	conf "github.com/xi-mad/my_video/commom"
	"github.com/xi-mad/my_video/object"
	"github.com/xi-mad/my_video/tag"
	"github.com/xi-mad/my_video/tree"
	"github.com/xi-mad/my_video/user"
	"os"
)

func main() {
	Mkdir()

	conf.DefaultConfig = conf.NewConfig("./config/config.yaml")

	conf.Load(conf.DefaultConfig.Sqlite.Path)
	AutoMigrate()

	r := gin.Default()

	r.Static("/assets", "./static/assets")
	r.StaticFile("/", "./static/index.html")

	user.Register(r.Group("/"))

	api := r.Group("/api")
	tag.Register(api.Group("/tag"))
	tree.Register(api.Group("/tree"))
	object.Register(api.Group("/object"))
	actress.Register(api.Group("/actress"))

	_ = r.Run()
}

func Mkdir() {
	_ = os.Mkdir("./data", 0755)
	_ = os.Mkdir("./temp", 0755)
	_ = os.Mkdir("./config", 0755)
}

func AutoMigrate() {
	tag.AutoMigrate()
	tree.AutoMigrate()
	object.AutoMigrate()
	actress.AutoMigrate()
}
