// api
package api

import (
	"net/http"

	"strings"

	"lottery-server/constants"
	"lottery-server/model"
	"lottery-server/response"
	"lottery-server/util"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func Register(c *gin.Context) {

	var response response.RegisterReponse

	userName := c.PostForm("userName")
	userPwd := c.PostForm("userPwd")

	if strings.Compare(userName, "") == 0 || strings.Compare(userPwd, "") == 0 {

		response.Status = constants.ApiCodeFail
		response.Message = "empty"

		c.JSON(http.StatusOK, response)

		return
	}

	db, err := util.OpenLotteryDb()
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	rows, err := db.Query("SELECT * FROM user WHERE UserName=?", userName)
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	if rows.Next() {
		response.Status = constants.ApiCodeFail
		response.Message = "exist"
		c.JSON(http.StatusOK, response)

		db.Close()
		return
	}

	userToken := util.UUID()
	stmt, err := db.Prepare("INSERT user SET UserName=?,UserId=?,UserPwd=?,UserToken=?")
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	userId := util.UUID()

	_, err = stmt.Exec(userName, userId, userPwd, userToken)
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	db.Close()

	response.Status = constants.ApiCodeSuccess
	response.Message = "success"
	response.Data = model.User{userId, userName, userPwd, userToken}

	c.JSON(http.StatusOK, response)

}

func Login(c *gin.Context) {
	var response response.RegisterReponse

	userName := c.Query("userName")
	userPwd := c.Query("userPwd")

	if strings.Compare(userName, "") == 0 || strings.Compare(userPwd, "") == 0 {

		response.Status = constants.ApiCodeFail
		response.Message = "empty"

		c.JSON(http.StatusOK, response)

		return
	}

	db, err := util.OpenLotteryDb()
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	rows, err := db.Query("SELECT * FROM user WHERE UserName=?", userName)
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	if rows.Next() == false {
		response.Status = constants.ApiCodeFail
		response.Message = "user not exist"
		c.JSON(http.StatusOK, response)

		db.Close()
		return
	}

	var user model.User
	err = rows.Scan(&user.UserId, &user.UserName, &user.UserPwd)
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	if userPwd != user.UserPwd {
		response.Status = constants.ApiCodeFail
		response.Message = "user pwd wrong"
		c.JSON(http.StatusOK, response)

		db.Close()
		return
	}

	userToken := util.UUID()

	stmt, err := db.Prepare("update user set UserToken=? where UserId=?")
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}
	_, err = stmt.Exec(userToken, user.UserId)
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	db.Close()

	user.UserToken = userToken
	response.Status = constants.ApiCodeSuccess
	response.Message = "success"
	response.Data = user

	c.JSON(http.StatusOK, response)
}
