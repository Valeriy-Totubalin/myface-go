package auth

import (
	"time"

	"github.com/Valeriy-Totubalin/myface-go/internal/delivery/http/request"
	"github.com/Valeriy-Totubalin/myface-go/internal/domain"
	"github.com/Valeriy-Totubalin/myface-go/internal/repository/mysql_db/user_repository"
	"github.com/Valeriy-Totubalin/myface-go/pkg/token_manager"
	"github.com/gin-gonic/gin"
)

func SignUp(data request.SignUp) error {
	data.Password = generateHash(data.Password)
	user := domain.User{
		Login:    data.Login,
		Password: data.Password,
		Sex:      data.Sex,
	}

	err := user_repository.SignUp(user)
	if nil != err {
		return err
	}

	return nil
}

func SignIn(c *gin.Context, data request.SignIn) error {
	user, err := user_repository.GetByLogin(data.Login)
	if nil != err {
		return err
	}

	err = checkPassword(data.Password, user.Password)
	if nil != err {
		return err
	}

	tokens, err := createSession(user.Id)
	if nil != err {
		return err
	}

	c.Set("access_token", tokens.AccessToken)
	c.Set("refresh_token", tokens.RefreshToken)

	return nil
}

func createSession(userId int) (token_manager.Tokens, error) {
	var res token_manager.Tokens
	tokenManager, err := token_manager.NewManager("test")
	if nil != err {
		return res, err
	}

	res.AccessToken, err = tokenManager.NewJWT(string(userId), 15*time.Minute)
	if nil != err {
		return res, err
	}

	res.RefreshToken, err = tokenManager.NewRefreshToken()
	if nil != err {
		return res, err
	}

	return res, nil
}
