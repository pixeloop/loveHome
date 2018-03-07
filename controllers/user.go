package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"loveHome/models"
	"path"
)

type UserController struct {
	beego.Controller
}

//返回结构变成json返回前端
func (this *UserController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

//用户注册 [post]
func (this *UserController) Reg() {
	beego.Info("========get user succ....")

	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	//获得客户端注册请求的json数据
	var regRequestData = make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &regRequestData)

	beego.Info("modile= ", regRequestData["mobile"])
	beego.Info("password= ", regRequestData["password"])
	beego.Info("sms_code=", regRequestData["sms_code"])

	if regRequestData["mobile"] == "" || regRequestData["password"] == "" || regRequestData["sms_code"] == "" {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}

	//数据存入数据库
	user := models.User{}
	user.Mobile = regRequestData["mobile"].(string)
	user.Password_hash = regRequestData["password"].(string)
	user.Name = regRequestData["mobile"].(string)

	o := orm.NewOrm()

	id, err := o.Insert(&user)
	if err != nil {
		beego.Info("Insert userinfo err ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	beego.Info("reg succ...  user id = ", id)

	//将用户的信息存到session中
	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", id)
	this.SetSession("mobile", user.Mobile)

	return
}

//用户登录
func (this *UserController) Login() {
	beego.Info("======== user Login  succ....======")

	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	//获得客户端注册请求的json数据
	var loginRequestData = make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &loginRequestData)

	beego.Info("modile= ", loginRequestData["mobile"])
	beego.Info("password= ", loginRequestData["password"])
	beego.Info("sms_code=", loginRequestData["sms_code"])

	if loginRequestData["mobile"] == "" || loginRequestData["password"] == "" || loginRequestData["sms_code"] == "" {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}

	//查询数据库
	var user models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	if err := qs.Filter("mobile", loginRequestData["mobile"]).One(&user); err != nil {
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
		return
	}
	if user.Password_hash != loginRequestData["password"].(string) {
		resp["errno"] = models.RECODE_PWDERR
		resp["errmsg"] = models.RecodeText(models.RECODE_PWDERR)
		return
	}

	beego.Info("====login succ ==username = ", user.Name)

	//将用户的信息存到session中
	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", user.Id)
	this.SetSession("mobile", user.Mobile)

	return
}

//处理上传头像的业务
func (this *UserController) UploadAvatar() {

	//返回给前端的map结构体
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	//得到文件二进制数据
	file, header, err := this.GetFile("avatar")
	if err != nil {
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}

	fileBuffer := make([]byte, header.Size)
	if _, err := file.Read(fileBuffer); err != nil {
		resp["errno"] = models.RECODE_IOERR
		resp["errmsg"] = models.RecodeText(models.RECODE_IOERR)
		return
	}

	suffix := path.Ext(header.Filename) // home.jpg.rmvb--->  .rmvb

	//将文件的二进制数据上传到fastdfs中

	groupName, fileId, err := models.FDFSUploadByBuffer(fileBuffer, suffix[1:]) //"rmvb"
	if err != nil {
		resp["errno"] = models.RECODE_IOERR
		resp["errmsg"] = models.RecodeText(models.RECODE_IOERR)
		beego.Info("upload file to fastdfs error err = ", err)
		return
	}

	beego.Info("fdfs upload succ groupname = ", groupName, "  fileid = ", fileId)

	//fileid ---> user 表里avatar_ur字段中
	//可以从seession中获得user.Id
	user_id := this.GetSession("user_id")
	user := models.User{Id: user_id.(int), Avatar_url: fileId}

	//数据库的操作，
	o := orm.NewOrm()
	if _, err := o.Update(&user, "avatar_url"); err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	//将fileid拼接成一个完整的url路径
	avatar_url := "http://192.168.182.130:8080/" + fileId

	//安装协议做出json返回给前端

	url_map := make(map[string]interface{})
	url_map["avatar_url"] = avatar_url
	resp["data"] = url_map

	return

}
