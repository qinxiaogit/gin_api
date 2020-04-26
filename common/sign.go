package common

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"sort"
)

func CreateSign(params url.Values)string{
	var key []string
	var str = ""
	for k:= range params {
		if k != "sn"{
			key = append(key,k)
		}
	}
	sort.Strings(key)
	for i:=0;i<len(key) ;i++  {
		if i == 0{
			str = fmt.Sprintf("%v=%v",key[i],params.Get(key[i]))
		}else{
			str = str+fmt.Sprintf("&%v=%v",key[i],params.Get(key[i]))
		}
	}
	//自定义签名算法
	sign := MD5(MD5(str)+MD5(APP_NAME+APP_SECRET))
	return sign
}

func MD5(str string) string{
	md5 := crypto.MD5.New()
	md5.Write([]byte(str))
	return hex.EncodeToString(md5.Sum(nil))
}

func VerifySign(c *gin.Context){
	//var method = c.Request.Method
	//var ts int64
	//var sn string
	//var req url.Values
	//
	//if method == "GET"{
	//	req = c.Request.URL.Query()
	//	sn = c.Query("sn")
	//	ts,_ = strconv.ParseInt(c.Query("ts"),10,64)
	//}else if method == "POST"{
	//	req = c.Request.PostForm
	//	sn = c.PostForm("sn")
	//	ts ,_ = strconv.ParseInt(c.PostForm("ts"),10,64)
	//}else{
	//	c.JSON(405,gin.H{
	//		"message":"method is no support",
	//	})
	//}
	//exp ,_ := strconv.ParseInt("1000",10,64)
	//

}