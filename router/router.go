package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go_gin_api/common"
	v1 "go_gin_api/controller/v1"
	v2 "go_gin_api/controller/v2"
	"go_gin_api/validator/member"
	"net/url"
	"strconv"
)

func InitRouter(r *gin.Engine)  {
	r.GET("/sn",SignDemo)
	//v1
	GroupV1 := r.Group("/v1")
	{
		GroupV1.Any("/index",v1.Index)
		GroupV1.Any("/grpc",v1.Grpc)
		GroupV1.Any("/member/add",v1.AddMember)
	}
	GroupV2 := r.Group("/v2")
	{
		GroupV2.Any("/index",v2.Index)

	}
	//绑定验证器
	if v,ok := binding.Validator.Engine().(*validator.Validate);ok{
		v.RegisterValidation("NameValid",member.NameValid)
		v.RegisterValidation("AgeValid",member.AgeValid)
	}
}
func SignDemo(c *gin.Context)  {
	ts := strconv.FormatInt(common.GetTimeUnix(),10)
	res := map[string]interface{}{}
	params := url.Values{
		"name":[]string{"a"},
		"price":[]string{"10"},
		"ts":[]string{ts},
	}
	res["sn"] = common.CreateSign(params)
	c.JSON(200,res)
}
