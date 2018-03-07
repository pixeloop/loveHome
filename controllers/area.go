package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"loveHome/models"
	"time"
)

type AreaController struct {
	beego.Controller
}

func (this *AreaController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *AreaController) GetAreaInfo() {
	beego.Info("===========get areainfo succ...===============")

	//返回前端的map结构体
	resp := make(map[string]interface{})
	//返回的默认值
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)

	cache_conn, err := cache.NewCache("redis", `{"key":"lovehome","conn":":6379","dbNum":"0"}`)
	if err != nil {
		beego.Info("cache redis conn err, err =", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	/*
		value := cache_conn.Get("haha")
		if value != nil {
			beego.Info("cache get value = ", value)
			fmt.Printf("value = %s\n", value)
		}
	*/

	area_info_value := cache_conn.Get("area_info")
	if area_info_value != nil {
		beego.Info("=====get area_info from cache!!!===")
		var area_info interface{}

		json.Unmarshal(area_info_value.([]byte), &area_info)
		resp["data"] = area_info
		return
	}

	o := orm.NewOrm()

	var areas []models.Area

	qs := o.QueryTable("area")
	num, err := qs.All(&areas)
	if err != nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	if num == 0 {
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
		return
	}
	//成功
	resp["data"] = areas

	//areajson字符串存入redis
	areas_info_str, _ := json.Marshal(areas)
	if err := cache_conn.Put("area_info", areas_info_str, time.Second*3600); err != nil {
		beego.Info("set area_info--> redis fail. err = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	return
}
