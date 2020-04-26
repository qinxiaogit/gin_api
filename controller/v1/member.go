package v1

import (
	"github.com/gin-gonic/gin"
	"go_gin_api/entity"
	"net/http"
)

func AddMember(c *gin.Context) {
	member := entity.Member{}
	res := entity.Result{}
	if err := c.ShouldBind(&member); err != nil {
		res.SetCode(entity.CODE_ERROR)
		res.SetMessage(err.Error())
		c.JSON(http.StatusForbidden,res)
		c.Abort()
		return
	}
	// 处理业务(下次再分享)

	data := map[string]interface{}{
		"name" : member.Name,
		//"age"  : member.Age,
	}

	res.SetCode(entity.CODE_ERROR)
	res.SetData(data)

	c.SecureJSON(http.StatusOK, res)
}
