package user

import (
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_auth/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserData struct {
	Login        string             `validate:"required"`
	Password     string             `validate:"required"`
	AccessRights models.AccessRight `validate:"required"`
}

// Get
// @Summary      Create user
// @Description  Создает пользователя. Доступен только админу
// @Tags         User
// @Param        data    body     UserData  true  "Данные пользователя для создания"
// @Produce      json
// @Success      200  {object}   map[string]models.User
// @Failure      404  {object}   string
// @Router       /user [post]
func (u *User) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		body := UserData{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithError(422, err)
			return
		}
		pHash, err := bcrypt.GenerateFromPassword([]byte(body.Password),bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		db := models.GetDb()
		user := models.User{}

		user.Login = body.Login
		user.Password = pHash
		user.AccessRights = body.AccessRights
		if err := db.Create(&user).Error; err != nil {
			c.AbortWithError(500, err)
		}

		c.JSON(200, gin.H{"user": user})
	}
}

type UpdateUserData struct {
	ID           uint
	Login        *string
	Password     *string
	AccessRights *models.AccessRight
}

// Get
// @Summary      Update user
// @Description  Обновляет пользователя. Доступен только админу
// @Tags         User
// @Param        data    body     UpdateUserData  true  "Данные пользователя для обновления"
// @Produce      json
// @Success      200  {object}   map[string]models.User
// @Failure      404  {object}   string
// @Router       /user [put]
func (u *User) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		body := UpdateUserData{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithError(422, err)
			return
		}

		db := models.GetDb()
		user := models.User{}

		if err := db.First(&user, body.ID).Error; err != nil {
			c.AbortWithError(404, err)
			return
		}

		if body.Login != nil {
			user.Login = *body.Login
		}
		if body.Password != nil {
			pHash, err := bcrypt.GenerateFromPassword([]byte(*body.Password),bcrypt.DefaultCost)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
			user.Password = pHash
		}
		if body.AccessRights != nil {
			user.AccessRights = *body.AccessRights
		}
		
		if err := db.Save(&user).Error; err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.JSON(200, gin.H{"user": user})
	}
}

// Get
// @Summary      Delete user
// @Description  Удаляет пользователя. Доступен только админу
// @Tags         User
// @Param        id    path     int  true  "id User"
// @Produce      json
// @Success      200  {object}   map[string]models.User
// @Failure      404  {object}   string
// @Router       /user/{id} [delete]
func (u *User) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := models.User{}
		db := models.GetDb()

		if err := db.First(&user, c.Param("id")).Error; err != nil {
			c.AbortWithError(404, err)
			return
		}

		if err := db.Delete(&user).Error; err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.JSON(200, gin.H{"user": user})
	}
}
