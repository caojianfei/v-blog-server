package helpers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
	"strconv"
	"time"
	"v-blog/models"
)

const (
	n      = 1 << 15
	r      = 16
	p      = 4
	keyLen = 32
)

func EncryptPassword(password string) (string, error) {

	salt := []byte("caojfcaoj")
	dk, err := scrypt.Key([]byte(password), salt, n, r, p, keyLen)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(dk), nil
}

type UserClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWTByUser(user *models.User) (string, error) {
	userClaims := &UserClaims{
		user.Email,
		jwt.StandardClaims{
			Audience:  user.Nickname,
			ExpiresAt: time.Now().Unix() + 15 * 24 * 60 * 60,
			Id:        strconv.FormatUint(uint64(user.ID),10),
			Issuer:    "v-blog",
			Subject:   "v-blog",
		},
	}

	signingKey := []byte("caojianfei")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenStr, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ParseUserJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("caojianfei"), nil
	})

	return token, err
}

type AuthorizedUser struct {
	Id int
	Nickname string
	Email string
}

func GetAuthorizedUserFromContext(c *gin.Context) (AuthorizedUser, error) {
	authorizedUser := AuthorizedUser{}
	if user, ok :=  c.Get("user"); ok {
		if loginedUser, ok := user.(map[string]string); ok {
			if idStr, ok := loginedUser["id"]; ok {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					fmt.Println("err ", err)
				} else {
					authorizedUser.Id = id
				}
			}
			if nickname, ok := loginedUser["nickname"]; ok {
				authorizedUser.Nickname = nickname;
			}
			if email, ok := loginedUser["email"]; ok {
				authorizedUser.Email = email
			}
			if authorizedUser.Id != 0 && authorizedUser.Nickname != "" && authorizedUser.Email != "" {
				return authorizedUser, nil
			}
		}
	}

	return authorizedUser, errors.New("没有获取到登录用户信息")
}

