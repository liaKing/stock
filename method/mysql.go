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
