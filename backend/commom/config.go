package commom

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App       App       `yaml:"app"`
	Sqlite    Sqlite    `yaml:"sqlite"`
	Thumbnail Thumbnail `yaml:"thumbnail"`
	Player    Player    `yaml:"player"`
}
type App struct {
	Mode string `yaml:"mode"`
}

type Sqlite struct {
	Path string `yaml:"path"`
}

type Thumbnail struct {
	Mtn   string `yaml:"mtn"`
	Width string `yaml:"width"`
	Font  string `yaml:"font"`
}

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
			Width: "2048",
		},
	}
	if content, err := os.ReadFile(path); err != nil {
		log.Fatal(err)
	} else {
		err = yaml.Unmarshal(content, config)
		if err != nil {
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
	log.Println(config)
	return config
}
