package method

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine/log"
	"stock/biz/dal/sql"
	"stock/biz/model"
	"stock/constant"
	"stock/util"
)

// DoFindMySQLUser 查找数据库用户
func DoFindMySQLUser(c *gin.Context, UserName string) (errCode util.HttpCode, user *model.User) {
	user = &model.User{}
	query := "SELECT * FROM user WHERE UserName = ?"
	err := config.MysqlConn.Raw(query, UserName).First(user).Error
	if err != nil {
		log.Errorf(c, "DoFindMySQLUser 操作mysql失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
		return
	}
	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: struct{}{},
	}

	return
}

// DoCreateMySQLUser 添加用户信息
func DoCreateMySQLUser(c *gin.Context, user *model.User) (errCode util.HttpCode, userNew *model.User) {
	if user != nil {
		log.Errorf(c, "DoGetRedisValue 关键信息丢失")
		errCode = util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
	}
	err := config.MysqlConn.Table("user").Create(user).Find(userNew).Error
	if err != nil {
		log.Errorf(c, "DoCreateMySQLUser 操作mysql失败")
		errCode = util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
		return
	}

	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: struct{}{},
	}
	return
}

//通过id获取用户信息

func GetUserById(c *gin.Context, id string) (errCode util.HttpCode, user *model.User) {
	user = &model.User{}
	query := "SELECT * FROM user WHERE user_id= ?"
	err := config.MysqlConn.Raw(query, id).First(user).Error
	if err != nil {
		errCode = util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
		return
	}
	return
}

// DoUpdataMySQLUser 修改数据库用户
func DoUpdataMySQLUser(c *gin.Context, userId string, deletionReason string) (errCode util.HttpCode) {
	query := "update user set delFlg =? and deletionReason = ? where userId = ?"
	err := config.MysqlConn.Exec(query, 1, deletionReason, userId).Error
	if err != nil {
		log.Errorf(c, "DoUpdataMySQLUser 操作mysql失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
		return
	}
	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: struct{}{},
	}

	return
}
