package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	authmodel "github.com/Popov-Dmitriy-Ivanovich/Diplom_auth/auth_model"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_auth/models"
)

func AuthMiddleware(requiredAr ... models.AccessRight) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		claims := &authmodel.Claims{}
		token = strings.TrimPrefix(token, "Bearer ")
		token = strings.TrimPrefix(token, "bearer ")

		jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil || !jwtToken.Valid {
			c.JSON(401, err.Error())
			c.Abort()
			return
		}

		db := models.GetDb()
		user := models.User{}

		if err := db.First(&user, claims.UserID).Error; err != nil {
			c.AbortWithStatus(404)
			return
		}

		for _, ar := range requiredAr{
			if !user.HasAccessRight(ar) {
				c.AbortWithStatus(401)
				c.Abort()
				return
			}
		}
		c.Set("UserID", claims.UserID)
		c.Set("UserAR", user.AccessRights)
		
		c.Next()
	}
}