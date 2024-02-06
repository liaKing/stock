package service

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine/log"
	"net/http"
	"stock/biz/model"
	"stock/constant"
	"stock/method"
	"stock/util"
	"strconv"
	"time"
)

func UserLogin(c *gin.Context) {
	user := &model.K2SLoginUser{}
	err := c.ShouldBind(&user)
	if err != nil {
		//log.Errorf(c, "UserLogin ShouldBind解析出错 err%d", err)
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
		//log.Errorf(c, "DoUserLogin 关键信息丢失")
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

func UserRegister(c *gin.Context) {
	user := &model.K2SRegisterUser{}
	err := c.ShouldBind(&user)
	if err != nil {
		log.Errorf(c, "UserRegister ShouldBind解析出错 err%d", err)
		c.JSON(http.StatusOK, util.HttpCode{
			Code: constant.ERRSHOULDBIND,
			Data: struct{}{},
		})
		return
	}

	errCode := DoUserRegister(c, user)
	if errCode.Code != constant.ERRSUCCER {
		c.JSON(http.StatusOK, errCode)
		return
	}

	c.JSON(http.StatusOK, errCode)
}

func DoUserRegister(c *gin.Context, user *model.K2SRegisterUser) (errCode util.HttpCode) {
	if user.UserName == "" || user.PassWord == "" || user.RealName == "" {
		log.Errorf(c, "DoUserRegister 关键信息丢失")
		errCode = util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
		return errCode
	}

	errCode, mUser := method.DoFindMySQLUser(c, user.UserName)
	if errCode.Code != constant.ERRSUCCER {
		return errCode
	}

	if mUser.UserName != "" {
		log.Errorf(c, "DoUserRegister 账号重复")
		errCode = util.HttpCode{
			Code: constant.ERRACCREPEAT,
			Data: struct{}{},
		}
		return errCode
	}

	user.PassWord = util.HashPassword(user.PassWord)

	snowflake, _ := util.NewSnowflake(1)
	uuid := snowflake.Generate()

	userNew := &model.User{
		UserId:      strconv.FormatInt(uuid, 10),
		Ctime:       time.Now().Unix(),
		UserName:    user.UserName,
		PassWord:    user.PassWord,
		RealName:    user.RealName,
		Name:        user.Name,
		WeChat:      user.WeChat,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
		Referrer:    user.Referrer,
	}

	errCode, userNew = method.DoCreateMySQLUser(c, userNew)
	if errCode.Code != constant.ERRSUCCER {
		return errCode
	}

	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: userNew,
	}
	return errCode
}
