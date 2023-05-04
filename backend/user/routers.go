package user

import (
	"github.com/gin-gonic/gin"
	"github.com/xi-mad/my_video/common"
)

func Register(r *gin.RouterGroup) {
	r.POST("/login", Login)
	r.POST("/logout", Logout)
	r.GET("/info", Info)
	r.GET("/menu", GetMenu)
}

func Login(c *gin.Context) {
	lr := LoginResponse{
		Token: "admin-token",
	}
	c.JSON(200, common.CommonResultSuccess(lr))
}

func Logout(c *gin.Context) {
	c.JSON(200, common.CommonResultSuccess(nil))
}

func Info(c *gin.Context) {
	ir := InfoResponse{
		Roles:  []string{"admin"},
		Name:   "Admin",
		Avatar: "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Token:  "admin-token",
	}
	c.JSON(200, common.CommonResultSuccess(ir))
}

func GetMenu(c *gin.Context) {
	var menu []Menu
	menu = append(menu, Menu{
		ID:        1,
		Pid:       0,
		Key:       "layout",
		Icon:      "AppleOutlined",
		Name:      "视频分组",
		Component: "BasicLayout",
		Path:      "/",
		Redirect:  "/video/collection",
		Children: []Menu{{
			ID:        2,
			Pid:       1,
			Key:       "element",
			Icon:      "ChromeOutlined",
			Name:      "我的视频",
			Component: "RouteView",
			Path:      "/video",
			Redirect:  "/video/collection",
			Children: []Menu{{
				ID:        9,
				Pid:       2,
				Key:       "video_collection",
				Name:      "合集管理",
				Path:      "/video/collection",
				Component: "/video/collection",
				Icon:      "",
				KeepAlive: true,
			}, {
				ID:        10,
				Pid:       2,
				Key:       "video_object",
				Name:      "对象管理",
				Path:      "/video/object",
				Component: "/video/object",
				Icon:      "",
			}, {
				ID:        11,
				Pid:       2,
				Key:       "video_tree",
				Name:      "树形分组",
				Path:      "/video/tree",
				Component: "/video/tree",
				Icon:      "",
			}, {
				ID:        12,
				Pid:       2,
				Key:       "video_actress",
				Name:      "演员管理",
				Path:      "/video/actress",
				Component: "/video/actress",
				Icon:      "",
			}, {
				ID:        13,
				Pid:       2,
				Key:       "video_tag",
				Name:      "标签管理",
				Path:      "/video/tag",
				Component: "/video/tag",
				Icon:      "",
			}},
		}},
	})
	c.JSON(200, common.CommonResultSuccess(menu))
}
