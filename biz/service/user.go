package service

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine/log"
	"net/http"
	"stock/biz/model"
	"stock/constant"
	"stock/method"
	"stock/util"
)

func UserLogin(c *gin.Context) {
	user := &model.K2SLoginUser{}
	err := c.ShouldBind(&user)
	if err != nil {
		log.Errorf(c, "UserLogin ShouldBind解析出错 err%d", err)
		c.JSON(http.StatusOK, util.HttpCode{
			Code: constant.ERRSHOULDBIND,
			Data: struct{}{},
		})
		return
	}

	errCode := DoUserLogin(c, user, c.ClientIP())
	if errCode.Code != constant.ERRSUCCER {
		c.JSON(http.StatusOK, errCode)
		return
	}

	c.JSON(http.StatusOK, errCode)
}

func DoUserLogin(c *gin.Context, user *model.K2SLoginUser, ip string) (errCode util.HttpCode) {
	if user.UserName == "" || user.PassWord == "" {
		log.Errorf(c, "DoUserLogin 关键信息丢失")
		errCode = util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
		return
	}

	errCode, mUser := method.DoFindMySQLUser(c, user.UserName)
	if errCode.Code != constant.ERRSUCCER {
		return
	}
	hashPassword := util.HashPassword(user.PassWord)

	if hashPassword != mUser.PassWord {
		log.Errorf(c, "DoUserLogin 密码不正确")
		errCode = util.HttpCode{
			Code: constant.ERRPSWNOTCORRECT,
			Data: struct{}{},
		}
		return
	}

	key := constant.REDIS_KEY_SESSION + user.UserName

	val := mUser.UserId + "_" + ip
	errCode = method.DoSetRedisValue(key, val, 5*60)
	if errCode.Code != constant.ERRSUCCER {
		return
	}

	err := util.SetSession(c, mUser.UserName, mUser.UserId, c.ClientIP())
	if err != nil {
		log.Errorf(c, "DoUserLogin SetSession session生成失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.ERRCREATEESSION,
			Data: struct{}{},
		}
		c.JSON(http.StatusOK, errCode)
		return
	}

	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: struct{}{},
	}

	return
}
