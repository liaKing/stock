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
		log.Errorf(c, "UserLogin ShouldBind解析出错 err%d", err)
		c.JSON(http.StatusOK, util.HttpCode{
			Code: constant.ERRSHOULDBIND,
			Data: struct{}{},
		})
		return
	}

	errCode := DoUserLogin(c, user)
	if errCode.Code != constant.ERRSUCCER {
		c.JSON(http.StatusOK, errCode)
		return
	}

	c.JSON(http.StatusOK, errCode)
}

func DoUserLogin(c *gin.Context, user *model.K2SLoginUser) (errCode util.HttpCode) {
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

	if mUser.DelFlg != 0 {
		log.Errorf(c, "DoUserLogin 用户已经注销")
		errCode = util.HttpCode{
			Code: constant.ERRSIGNOUT,
			Data: struct{}{},
		}
		return
	}

	if hashPassword != mUser.PassWord {
		log.Errorf(c, "DoUserLogin 密码不正确")
		errCode = util.HttpCode{
			Code: constant.ERRPSWNOTCORRECT,
			Data: struct{}{},
		}
		return
	}

	//key := constant.REDIS_KEY_SESSION + user.UserName
	//
	//val := mUser.UserId
	//errCode = method.DoSetRedisValue(c, key, val, 5*60)
	//if errCode.Code != constant.ERRSUCCER {
	//	return
	//}

	//err := util.SetSession(c, mUser.UserName, mUser.UserId)
	//session := util.GetSession(c)
	//if err != nil {
	//	errCode = util.HttpCode{
	//		Code: constant.ERRCREATEESSION,
	//		Data: struct{}{},
	//	}
	//	c.JSON(http.StatusOK, errCode)
	//	return
	//}
	a, r, err1 := util.Token.GetToken(mUser.UserId, mUser.UserName)
	if err1 != nil {
		c.JSON(http.StatusOK, util.HttpCode{
			Code: constant.ERRUserInfoErr,
			Data: struct{}{},
		})
		return
	}
	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: model.Token{
			AccessToken:  a,
			RefreshToken: r,
		},
	}

	c.JSON(http.StatusOK, errCode)

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
		UserId:         strconv.FormatInt(uuid, 10),
		Ctime:          time.Now().Unix(),
		DelFlg:         0,
		DeletionReason: "",
		UserName:       user.UserName,
		PassWord:       user.PassWord,
		RealName:       user.RealName,
		Name:           user.Name,
		WeChat:         user.WeChat,
		PhoneNumber:    user.PhoneNumber,
		Address:        user.Address,
		Luck:           0,
		Referrer:       user.Referrer,
	}

	errCode = method.DoCreateMySQLUser(c, userNew)
	if errCode.Code != constant.ERRSUCCER {
		return errCode
	}

	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: struct{}{},
	}
	return errCode
}

func UserDel(c *gin.Context) {
	user := &model.K2SDelUser{}
	err := c.ShouldBind(&user)
	if err != nil {
		log.Errorf(c, "UserDel ShouldBind解析出错 err%d", err)
		c.JSON(http.StatusOK, util.HttpCode{
			Code: constant.ERRSHOULDBIND,
			Data: struct{}{},
		})
		return
	}

	errCode := DoUserDel(c, user)
	if errCode.Code != constant.ERRSUCCER {
		c.JSON(http.StatusOK, errCode)
		return
	}

	c.JSON(http.StatusOK, errCode)
}

func DoUserDel(c *gin.Context, user *model.K2SDelUser) (errCode util.HttpCode) {
	if user.UserId == "" {
		log.Errorf(c, "DoUserDel 关键信息丢失")
		errCode = util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
		return errCode
	}

	errCode, mUser := method.GetUserById(c, user.UserId)
	if errCode.Code != constant.ERRSUCCER {
		return errCode
	}

	if mUser.UserName == "" {
		log.Errorf(c, "DoUserDel 用户不存在")
		errCode = util.HttpCode{
			Code: constant.ERRISNOTEXIT,
			Data: struct{}{},
		}
		return errCode
	}

	errCode = method.DoUpdataMySQLUser(c, user.UserId, user.DeletionReason)
	if errCode.Code != constant.ERRSUCCER {
		return errCode
	}

	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: struct{}{},
	}
	return errCode
}

func UserGet(c *gin.Context) {
	user := &model.K2SGetUser{}
	err := c.ShouldBind(&user)
	if err != nil {
		log.Errorf(c, "UserDel ShouldBind解析出错 err%d", err)
		c.JSON(http.StatusOK, util.HttpCode{
			Code: constant.ERRSHOULDBIND,
			Data: struct{}{},
		})
		return
	}

	errCode := DoUserGet(c, user)
	if errCode.Code != constant.ERRSUCCER {
		c.JSON(http.StatusOK, errCode)
		return
	}

	c.JSON(http.StatusOK, errCode)
}

func DoUserGet(c *gin.Context, user *model.K2SGetUser) (errCode util.HttpCode) {
	if user.UserId == "" {
		log.Errorf(c, "DoUserDel 关键信息丢失")
		errCode = util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
		return errCode
	}
	errCode, mUser := method.GetUserById(c, user.UserId)
	if errCode.Code != constant.ERRSUCCER {
		return errCode
	}
	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: mUser,
	}
	return
}

func UserDoLuck(c *gin.Context) {
	user := &model.K2SDoLuckUser{}
	err := c.ShouldBind(&user)
	if err != nil {
		log.Errorf(c, "UserDel ShouldBind解析出错 err%d", err)
		c.JSON(http.StatusOK, util.HttpCode{
			Code: constant.ERRSHOULDBIND,
			Data: struct{}{},
		})
		return
	}

	errCode := DoUserLuck(c, user)
	if errCode.Code != constant.ERRSUCCER {
		c.JSON(http.StatusOK, errCode)
		return
	}

	c.JSON(http.StatusOK, errCode)
}

func DoUserLuck(c *gin.Context, user *model.K2SDoLuckUser) (errCode util.HttpCode) {
	if user.UserId == "" || user.Luck == 0 {
		log.Errorf(c, "DoUserLuck 关键信息丢失")
		errCode = util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
		return errCode
	}
	errCode, mUser := method.GetUserById(c, user.UserId)
	if errCode.Code != constant.ERRSUCCER {
		return errCode
	}

	mUser.Luck += user.Luck

	errCode = method.DoUpdataMySQLUserLuck(c, user.UserId, mUser.Luck)
	if errCode.Code != constant.ERRSUCCER {
		return errCode
	}

	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: mUser,
	}
	return
}
