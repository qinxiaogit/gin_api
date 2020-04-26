package v2

import "github.com/gin-gonic/gin"

func Index(c *gin.Context)  {
	//获取get参数
	name := c.Query("name")
	price:= c.DefaultQuery("price","200")

	c.JSON(200,gin.H{
		"v2" :"index",
		"name":name,
		"price":price,
	})
}