package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"stock/constant"
	"stock/method"
	"stock/util"
)

// AuthCheck 检查是否登录
func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := util.GetSession(c)
		userId, _ := data["userId"]
		userName, _ := data["userName"]
		accountString, _ := userName.(string)
		// 从 Redis 中获取存储的验证码
		key := constant.REDIS_KEY_SESSION + accountString
		errCode, session := method.DoGetRedisValue(c, key)
		if errCode.Code != constant.ERRSUCCER {
			c.JSON(http.StatusOK, errCode)
			return
		}
		val := session.(string)
		uidString := userId.(string)

		if uidString != val {
			c.Redirect(http.StatusFound, "/login")
			c.Abort() //如果用户没有登录，中间件直接返回，不再向后继续
		}

		c.Set("userName", accountString)
		c.Set("userId", uidString)
		c.Next()
	}
}

// AdminCheck 检查是否是管理员
func AdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := util.GetSession(c)
		userId, _ := data["userId"]
		uidString := userId.(string)
		if uidString != constant.AdministratorUserId {
			c.Redirect(http.StatusFound, "/login")
			c.Abort() //如果用户没有登录，中间件直接返回，不再向后继续
		}

		c.Next()
	}
}

func AuthJWTCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		//测试模式不需要验签,需要自己在请求的头部伪造一个Debug数据
		if c.GetHeader("debug") != "" {
			c.Next()
			return
		}
		auth := c.GetHeader("Authorization")
		fmt.Printf("auth:%+v\n", auth)
		data, err := util.Token.VerifyToken(auth)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.HttpCode{
				Code: constant.ERRNotLogin,
				Data: "验签失败！",
			})
		}
		fmt.Printf("data:%+v\n", data)
		if data.ID == "" || data.Name == "" {
			//如果用户没有登录，中间件直接返回，不再向后继续
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.HttpCode{
				Code: constant.ERRNotLogin,
				Data: "用户信息获取错误",
			})
			return
		}

		//将用户信息塞到Context中
		c.Set("userName", data.Name)
		c.Set("userId", data.ID)

		c.Next()
	}
}

func AuthJWTAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		//测试模式不需要验签,需要自己在请求的头部伪造一个Debug数据
		if c.GetHeader("debug") != "" {
			c.Next()
			return
		}
		auth := c.GetHeader("Authorization")
		fmt.Printf("auth:%+v\n", auth)
		data, err := util.Token.VerifyToken(auth)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.HttpCode{
				Code: constant.ERRNotLogin,
				Data: "验签失败！",
			})
		}
		fmt.Printf("data:%+v\n", data)
		if data.ID != "001" || data.Name == "" {
			//如果用户没有登录，中间件直接返回，不再向后继续
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.HttpCode{
				Code: constant.ERRNotLogin,
				Data: "用户信息获取错误",
			})
			return
		}

		//将用户信息塞到Context中
		c.Set("userName", data.Name)
		c.Set("userId", data.ID)

		c.Next()
	}
}
