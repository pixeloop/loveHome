package routers

import (
	"github.com/astaxie/beego"
	"loveHome/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	//地区请求
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetAreaInfo")

	//session请求
	beego.Router("/api/v1.0/session", &controllers.SessionController{}, "get:GetSessionName;delete:DelSessionName")
	//文件上传业务
	beego.Router("api/v1.0/user/avatar", &controllers.UserController{}, "post:UploadAvatar")
	//房屋index请求
	beego.Router("/api/v1.0/houses/index", &controllers.HouseIndexController{}, "get:GetHouseIndex")

	//user注册用户请求
	beego.Router("/api/v1.0/users", &controllers.UserController{}, "post:Reg")

	//user登录
	beego.Router("/api/v1.0/sessions", &controllers.UserController{}, "post:Login")
}
