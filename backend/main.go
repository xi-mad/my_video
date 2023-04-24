package main

import (
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/gin-gonic/gin"
	"github.com/xi-mad/my_video/actress"
	conf "github.com/xi-mad/my_video/commom"
	"github.com/xi-mad/my_video/object"
	"github.com/xi-mad/my_video/plantform"
	"github.com/xi-mad/my_video/tag"
	"github.com/xi-mad/my_video/tree"
	"github.com/xi-mad/my_video/user"
	"io"
	"log"
	"os"
)

func main() {
	Mkdir()

	gin.DefaultWriter = io.MultiWriter(conf.Logfile, os.Stdout)
	log.SetOutput(io.MultiWriter(conf.Logfile, os.Stdout))

	conf.DefaultConfig = conf.NewConfig("./config/config.yaml")

	if conf.DefaultConfig.App.Mode != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

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
	log.Printf("server start at %s\n", conf.DefaultConfig.App.Port)
	log.Printf("please open http://localhost:%s or http://127.0.0.1:%s\n", conf.DefaultConfig.App.Port, conf.DefaultConfig.App.Port)
	go aNotify()
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

	webPage := systray.AddMenuItem("打开my_video", "打开my_video")
	logFolder := systray.AddMenuItem("打开日志文件夹", "打开日志文件夹")
	configFolder := systray.AddMenuItem("打开配置文件夹", "打开配置文件夹")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出程序")

	go func() {
		for {
			select {
			case <-logFolder.ClickedCh:
				dir, _ := os.Getwd()
				_ = plantform.OpenFolder(dir + "\\log")
			case <-configFolder.ClickedCh:
				dir, _ := os.Getwd()
				_ = plantform.OpenFolder(dir + "\\config")
			case <-webPage.ClickedCh:
				openInBrowser()
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}

func openInBrowser() {
	err := plantform.OpenInBrowser("http://127.0.0.1:" + conf.DefaultConfig.App.Port)
	if err != nil {
		log.Println(err)
	}
}

func aNotify() {
	err := beeep.Notify("my_video", "my_video已启动, 请打开浏览器访问 http://127.0.0.1:"+conf.DefaultConfig.App.Port, "./static/favicon.ico")
	if err != nil {
		log.Println(err)
	}
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
