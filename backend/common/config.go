package common

import (
	"log"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App       App       `yaml:"app"`
	Sqlite    Sqlite    `yaml:"sqlite"`
	Thumbnail Thumbnail `yaml:"thumbnail"`
	Player    Player    `yaml:"player"`
}
type App struct {
	Mode             string `yaml:"mode"`
	Port             string `yaml:"port"`
	ServerMode       bool   `yaml:"server-mode"`
	ObjectHeight     string `yaml:"object-height"`
	CollectionHeight string `yaml:"collection-height"`
}

type Sqlite struct {
	Path string `yaml:"path"`
}

type Thumbnail struct {
	Mtn      string     `yaml:"mtn"`
	Width    int        `yaml:"width"`
	Row      int        `yaml:"row"`
	Col      int        `yaml:"col"`
	Optional []Optional `yaml:"optional"`
	Font     string     `yaml:"font"`
}

type Optional struct {
	FSizeLess int64 `yaml:"fsizeless"`
	Width     int   `yaml:"width"`
	Row       int   `yaml:"row"`
	Col       int   `yaml:"col"`
}
type OptionalSlice []Optional

func (o OptionalSlice) Len() int           { return len(o) }
func (o OptionalSlice) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o OptionalSlice) Less(i, j int) bool { return o[i].FSizeLess < o[j].FSizeLess }

type Player struct {
	Path string `yaml:"path"`
}

var DefaultConfig *Config

func NewConfig(path string) *Config {
	config := &Config{
		Sqlite: Sqlite{
			Path: "./data/my_video.db",
		},
		Thumbnail: Thumbnail{
			Width: 2048,
		},
	}
	if content, err := os.ReadFile(path); err != nil {
		log.Fatal(err)
	} else {
		err = yaml.Unmarshal(content, config)
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}
	if config.Thumbnail.Mtn == "" {
		log.Println("mtn not found, use default")
		config.Thumbnail.Mtn = "mtn"
	}
	if config.Player.Path == "" {
		log.Println("播放器路径未设置")
	}
	if config.App.Port == "" {
		log.Println("端口未设置, 使用默认端口8080")
		config.App.Port = "8080"
	}
	if config.App.ServerMode {
		log.Println("服务端模式")
	}
	sort.Sort(OptionalSlice(config.Thumbnail.Optional))
	log.Println(config)
	return config
}
