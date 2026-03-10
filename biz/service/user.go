package service

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"stock/biz/model"
	"stock/constant"
	"stock/method"
	"stock/util"
)

func UserLogin(c *gin.Context) {
	user := &model.K2SLoginUser{}
	err := c.ShouldBind(&user)
	if err != nil {
		log.Printf("[ERROR] UserLogin ShouldBind解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	errCode := DoUserLogin(c, user)
	if !errCode.IsSuccess() {
		c.JSON(http.StatusOK, errCode.EnsureMessage())
		return
	}
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

func DoUserLogin(c *gin.Context, user *model.K2SLoginUser) (errCode util.HttpCode) {
	if user.UserName == "" || user.PassWord == "" {
		log.Printf("[ERROR] DoUserLogin 关键信息丢失")
		return util.Fail(constant.ERRDATALOSE)
	}

	errCode, mUser := method.DoFindMySQLUser(c, user.UserName)
	if !errCode.IsSuccess() {
		return errCode
	}
	hashPassword := util.HashPassword(user.PassWord)

	if mUser.DelFlg != 0 {
		log.Printf("[ERROR] DoUserLogin 用户已经注销")
		return util.Fail(constant.ERRSIGNOUT)
	}

	if hashPassword != mUser.PassWord {
		log.Printf("[ERROR] DoUserLogin 密码不正确")
		return util.Fail(constant.ERRPSWNOTCORRECT)
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
		return util.Fail(constant.ERRUserInfoErr)
	}
	return util.Success(model.Token{
		AccessToken:  a,
		RefreshToken: r,
	})
}

func UserRegister(c *gin.Context) {
	user := &model.K2SRegisterUser{}
	err := c.ShouldBind(&user)
	if err != nil {
		log.Printf("[ERROR] UserRegister ShouldBind解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	errCode := DoUserRegister(c, user)
	if !errCode.IsSuccess() {
		c.JSON(http.StatusOK, errCode.EnsureMessage())
		return
	}
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

func DoUserRegister(c *gin.Context, user *model.K2SRegisterUser) (errCode util.HttpCode) {
	if user.UserName == "" || user.PassWord == "" || user.RealName == "" {
		log.Printf("[ERROR] DoUserRegister 关键信息丢失")
		return util.Fail(constant.ERRDATALOSE)
	}

	errCode, mUser := method.DoFindMySQLUser(c, user.UserName)
	if !errCode.IsSuccess() {
		return errCode
	}

	if mUser.UserName != "" {
		log.Printf("[ERROR] DoUserRegister 账号重复")
		return util.Fail(constant.ERRACCREPEAT)
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
		Luck:           user.Luck,
		Referrer:       user.Referrer,
	}

	errCode = method.DoCreateMySQLUser(c, userNew)
	if !errCode.IsSuccess() {
		return errCode
	}
	return util.Success(nil)
}

func UserDel(c *gin.Context) {
	user := &model.K2SDelUser{}
	err := c.ShouldBind(&user)
	if err != nil {
		log.Printf("[ERROR] UserDel ShouldBind解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	errCode := DoUserDel(c, user)
	if !errCode.IsSuccess() {
		c.JSON(http.StatusOK, errCode.EnsureMessage())
		return
	}
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

func DoUserDel(c *gin.Context, user *model.K2SDelUser) (errCode util.HttpCode) {
	if user.UserId == "" {
		log.Printf("[ERROR] DoUserDel 关键信息丢失")
		return util.Fail(constant.ERRDATALOSE)
	}

	errCode, mUser := method.GetUserById(c, user.UserId)
	if !errCode.IsSuccess() {
		return errCode
	}

	if mUser.UserName == "" {
		log.Printf("[ERROR] DoUserDel 用户不存在")
		return util.Fail(constant.ERRISNOTEXIT)
	}

	errCode = method.DoUpdataMySQLUser(c, user.UserId, user.DeletionReason)
	if !errCode.IsSuccess() {
		return errCode
	}
	return util.Success(nil)
}

func UserGet(c *gin.Context) {
	user := &model.K2SGetUser{}
	err := c.ShouldBindQuery(&user)
	if err != nil {
		log.Printf("[ERROR] UserGet ShouldBind解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	errCode := DoUserGet(c, user)
	if !errCode.IsSuccess() {
		c.JSON(http.StatusOK, errCode.EnsureMessage())
		return
	}
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

func DoUserGet(c *gin.Context, user *model.K2SGetUser) (errCode util.HttpCode) {
	if user.UserId == "" {
		log.Printf("[ERROR] DoUserGet 关键信息丢失")
		return util.Fail(constant.ERRDATALOSE)
	}
	errCode, mUser := method.GetUserById(c, user.UserId)
	if !errCode.IsSuccess() {
		return errCode
	}
	return util.Success(mUser.ToUserInfo())
}

// UserDoLuck 添加幸运值
func UserDoLuck(c *gin.Context) {
	user := &model.K2SDoLuckUser{}
	err := c.ShouldBind(&user)
	if err != nil {
		log.Printf("[ERROR] UserDoLuck ShouldBind解析出错 err%v", err)
		c.JSON(http.StatusOK, util.Fail(constant.ERRSHOULDBIND))
		return
	}

	errCode := DoUserLuck(c, user)
	if !errCode.IsSuccess() {
		c.JSON(http.StatusOK, errCode.EnsureMessage())
		return
	}
	c.JSON(http.StatusOK, errCode.EnsureMessage())
}

// DoUserLuck 增减幸运值（正数增加，负数减少）
func DoUserLuck(c *gin.Context, user *model.K2SDoLuckUser) (errCode util.HttpCode) {
	if user.UserId == "" || user.Luck == 0 {
		log.Printf("[ERROR] DoUserLuck 关键信息丢失")
		return util.Fail(constant.ERRDATALOSE)
	}
	errCode, mUser := method.GetUserById(c, user.UserId)
	if !errCode.IsSuccess() {
		return errCode
	}

	// 支持增减：正数加、负数减，结果不能小于 0
	newLuck := int64(mUser.Luck) + user.Luck
	if newLuck < 0 {
		log.Printf("[ERROR] DoUserLuck 幸运值不足，当前:%d 操作:%d", mUser.Luck, user.Luck)
		return util.Fail(constant.LackOfLuck)
	}
	mUser.Luck = uint64(newLuck)

	errCode = method.DoUpdataMySQLUserLuck(c, user.UserId, mUser.Luck)
	if !errCode.IsSuccess() {
		return errCode
	}
	return util.Success(mUser.ToUserInfo())
}
