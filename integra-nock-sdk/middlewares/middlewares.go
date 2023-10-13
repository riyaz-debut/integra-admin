package middlewares

import (
	"integra-nock-sdk/helpers/token"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"integra-nock-sdk/utils"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	log.Println("Inside middleware fx of admin side")
	return func(c *gin.Context) {
		receivedToken := c.Request.Header.Get("Authorization")
		if len(strings.Split(receivedToken, " ")) == 2 {
			receivedToken = strings.Split(receivedToken, " ")[1]
		}
		log.Println("token received in admin midddleware", receivedToken)
		tokenString, err := token.TokenValid(receivedToken)
		if err != nil {
			c.String(http.StatusUnauthorized, "Find unauthorized token in middleware check at admin side")
			c.Abort()
			return
		}
		log.Println("token in middleware: ", tokenString)
		user_id, org_id, err := token.ExtractTokenID(tokenString)
		log.Println("user_id in middleware: ", user_id)

		log.Println("user_id in middleware: ", org_id)
		if err != nil {
			log.Println("err in token ", err)
			c.IndentedJSON(401, utils.Response{
				Status:  401,
				Message: "Find unauthorized token in middleware check at client side",
				Data:    nil,
			})
			c.Abort()
			return
		}
		c.Set("user_id", user_id)
		c.Set("org_id", org_id)
		c.Next()
	}
}
