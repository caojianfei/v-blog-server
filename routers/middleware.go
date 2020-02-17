package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"v-blog/helpers"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenArr := strings.Split(tokenStr, " ")
		if len(tokenArr) != 2 || tokenArr[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authTokenStr := tokenArr[1]
		token, err := helpers.ParseUserJWT(authTokenStr)
		if err != err {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*helpers.UserClaims); ok && token.Valid {
			email := claims.Email
			userId := claims.Id
			nickname := claims.Audience
			c.Set("user", map[string]string{
				"email": email,
				"id": userId,
				"nickname": nickname,

			})
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next();
	}
}


func TestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("TestMiddleware")
	}
}