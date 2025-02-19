package user

import (
	"os"
	"time"

	authmodel "github.com/Popov-Dmitriy-Ivanovich/Diplom_auth/auth_model"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_auth/middleware"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_auth/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {}

func (u *User) WriteRoutes (rg *gin.RouterGroup) {
	userGroup := rg.Group("/user")
	userGroup.POST("/login",u.Login())
	userGroup.Use(middleware.AuthMiddleware(models.AR_CREATE_USER))
	userGroup.GET("/",u.Get())
	userGroup.GET("/:id", u.GetId())
	userGroup.POST("/", u.Create())
	userGroup.PUT("/", u.Update())
	userGroup.DELETE("/:id", u.Delete())
}
// Get
// @Summary      Get list of user ids
// @Description  Возращает список id всех доступных user
// @Tags         User
// @Produce      json
// @Success      200  {object}   map[string][]uint
// @Failure      500  {object}   string
// @Failure      404  {object}   string
// @Router       /user [get]
func (a *User) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := models.GetDb()
		ids := []uint{}
		if err := db.Model(models.User{}).Pluck("id",&ids).Error; err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.JSON(200, gin.H{"ids": ids})
	}
}

// Get
// @Summary      Get concrete Action
// @Description  Возращает Action соответсвующую указанному ID
// @Tags         Actions
// @Param        id    path     int  true  "id Action"
// @Produce      json
// @Success      200  {object}   map[string]models.User
// @Failure      404  {object}   string
// @Router       /user/{id} [get]
func (a *User) GetId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		db := models.GetDb()
		user := models.User{}
		if err := db.First(&user,id).Error; err != nil {
			c.AbortWithError(404,err)
			return
		}
		c.JSON(200, gin.H{"user": user})
	}
}

type LoginBody struct {
	Login string
	Password string
}

// Get
// @Summary      LOGIN
// @Description  login
// @Tags         Actions
// @Param        data    body     LoginBody  true  "id Action"
// @Produce      json
// @Success      200  {object}   map[string]models.User
// @Failure      404  {object}   string
// @Router       /user/login [post]
func (a *User) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := models.User{}
		db := models.GetDb()
		body := LoginBody{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithError(422, err)
			return
		}

		if err := db.First(&user, map[string]any{"login":body.Login}).Error; err != nil {
			c.AbortWithError(404, err)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
			c.AbortWithError(422, err)
			return
		}
		jwtKey := os.Getenv("JWT_KEY")

		expTimeAccess := time.Now().Add(5 * time.Hour)

		claimsAccess := &authmodel.Claims{
			UserID: user.ID,
			Login: user.Login,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expTimeAccess),
			},
		}

		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccess)
		accessString, err := accessToken.SignedString([]byte(jwtKey))
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.JSON(200, gin.H{"token": accessString, "user": user})
	}
}