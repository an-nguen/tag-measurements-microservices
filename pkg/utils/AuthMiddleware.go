package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"tag-measurements-microservices/pkg/models"
)

type JWTAuthMiddleware struct {
	Secret string
	UserDB *gorm.DB
}

func (c JWTAuthMiddleware) New() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr, err := ParseToken(ctx.GetHeader("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}

		var claims jwt.StandardClaims
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				ctx.Abort()
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}
			return []byte(c.Secret), nil
		})
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		} else if token.Valid == false {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		} else if claims.ExpiresAt < time.Now().Unix() {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (c JWTAuthMiddleware) NewWithRole(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, done := c.getUserByToken(ctx)
		if done {
			return
		}
		if len(user.Roles) == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no roles"})
			ctx.Abort()
			return
		}

		userRoles := make(map[string]*models.Role)
		for _, role := range user.Roles {
			userRoles[role.Name] = role
		}
		if _, ok := userRoles[role]; !ok {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (c JWTAuthMiddleware) NewWithPrivileges(privilege string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr, err := ParseToken(ctx.GetHeader("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}

		var claims jwt.StandardClaims
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				ctx.Abort()
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}
			return []byte(c.Secret), nil
		})
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		} else if token.Valid == false {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		} else if claims.ExpiresAt < time.Now().Unix() {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
			ctx.Abort()
			return
		}

		id, err := strconv.Atoi(claims.Audience)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		var user models.User
		if err := c.UserDB.
			Where("id = ?", id).
			Preload("Roles").
			Preload("Roles.Privileges").
			First(&user).
			Error; err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "can't find user"})
			ctx.Abort()
			return
		}

		if len(user.Roles) == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no roles"})
			ctx.Abort()
			return
		}

		userPrivileges := make(map[string]*models.Privilege)
		for _, role := range user.Roles {
			for _, privilege := range role.Privileges {
				userPrivileges[privilege.Name] = privilege
			}
		}
		if _, ok := userPrivileges[privilege]; !ok {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (c JWTAuthMiddleware) getUserByToken(ctx *gin.Context) (models.User, bool) {
	tokenStr, err := ParseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
		return models.User{}, true
	}

	var claims jwt.StandardClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			ctx.Abort()
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(c.Secret), nil
	})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		ctx.Abort()
		return models.User{}, true
	} else if token.Valid == false {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		ctx.Abort()
		return models.User{}, true
	} else if claims.ExpiresAt < time.Now().Unix() {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
		ctx.Abort()
		return models.User{}, true
	}

	id, err := strconv.Atoi(claims.Audience)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		ctx.Abort()
		return models.User{}, true
	}

	var user models.User
	if err := c.UserDB.
		Where("id = ?", id).
		Preload("Roles").
		Preload("Roles.Privilege").
		First(&user).
		Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "can't find user"})
		ctx.Abort()
		return models.User{}, true
	}
	return user, false
}
