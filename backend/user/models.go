package user

type LoginModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type InfoResponse struct {
	Roles  []string `json:"roles"`
	Name   string   `json:"name"`
	Avatar string   `json:"avatar"`
	Token  string   `json:"token"`
}

type Menu struct {
	ID        int    `json:"id"`
	Pid       int    `json:"pid"`
	Key       string `json:"key"`
	Icon      string `json:"icon"`
	Name      string `json:"name"`
	Component string `json:"component"`
	Path      string `json:"path"`
	Redirect  string `json:"redirect"`
	KeepAlive bool   `json:"keepAlive"`
	Children  []Menu `json:"children"`
}
