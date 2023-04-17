package main

import (
	"github.com/getlantern/systray"
	"github.com/gin-gonic/gin"
	"github.com/xi-mad/my_video/actress"
	conf "github.com/xi-mad/my_video/commom"
	"github.com/xi-mad/my_video/object"
	"github.com/xi-mad/my_video/tag"
	"github.com/xi-mad/my_video/tree"
	"github.com/xi-mad/my_video/user"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	Mkdir()

	conf.DefaultConfig = conf.NewConfig("./config/config.yaml")
	if !conf.DefaultConfig.App.ServerMode {
		go func() {
			systray.Run(onReady, onExit)
		}()
	}

	conf.Load(conf.DefaultConfig.Sqlite.Path)
	AutoMigrate()

	r := gin.Default()

	r.Static("/assets", "./static/assets")
	r.StaticFile("/", "./static/index.html")

	api := r.Group("/api")
	user.Register(api.Group("/user"))
	tag.Register(api.Group("/tag"))
	tree.Register(api.Group("/tree"))
	object.Register(api.Group("/object"))
	actress.Register(api.Group("/actress"))

	_ = r.Run(":" + conf.DefaultConfig.App.Port)
}

func onReady() {

	if f, err := os.Open("./static/favicon.ico"); err == nil {
		if bytes, err := io.ReadAll(f); err == nil {
			systray.SetIcon(bytes)
		}
	}
	systray.SetTitle("my_video")
	systray.SetTooltip("服务已最小化右下角, 右键点击打开菜单！")

	logFolder := systray.AddMenuItem("打开日志文件夹", "打开日志文件夹")
	configFolder := systray.AddMenuItem("打开配置文件夹", "打开配置文件夹")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出程序")

	go func() {
		for {
			select {
			case <-logFolder.ClickedCh:
				dir, _ := os.Getwd()
				_, err := exec.Command("explorer.exe", dir+"\\log").Output()
				if err != nil {
					log.Println(err)
				}

			case <-configFolder.ClickedCh:
				dir, _ := os.Getwd()
				_, err := exec.Command("explorer.exe", dir+"\\config").Output()
				if err != nil {
					log.Println(err)
				}
			case <-mQuit.ClickedCh:
				os.Exit(0)
			}
		}
	}()
}

func onExit() {
	// clean up here
}

func Mkdir() {
	_ = os.Mkdir("./data", 0755)
	_ = os.Mkdir("./temp", 0755)
	_ = os.Mkdir("./config", 0755)
	_ = os.Mkdir("./log", 0755)
}

func AutoMigrate() {
	tag.AutoMigrate()
	tree.AutoMigrate()
	object.AutoMigrate()
	actress.AutoMigrate()
}
