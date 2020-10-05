package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"Thermo-WH/pkg/repository"
	"Thermo-WH/pkg/utils"
)

type AuthController struct {
	UserRepo repository.UserRepository
	Secret   string
}

type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c AuthController) AuthUser(ctx *gin.Context) {
	var userReq UserReq
	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	token, user, err := c.UserRepo.AuthUser(userReq.Username, userReq.Password, c.Secret)
	if err != nil {
		utils.LogError("AuthUser", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""

	ctx.Header("authorization", fmt.Sprintf("Bearer %s", token))
	ctx.JSON(http.StatusOK, user)
}

func (c AuthController) VerifyToken(ctx *gin.Context) {
	// get token from session id
	tokenStr, err := utils.ParseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var claims jwt.StandardClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(c.Secret), nil
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if token.Valid == false {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	} else if claims.ExpiresAt < time.Now().Unix() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token is expired"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
	})
}

func (c AuthController) RefreshToken(ctx *gin.Context) {
	// get token from session id
	tokenStr, err := utils.ParseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var claims jwt.StandardClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(c.Secret), nil
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if token.Valid == false {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	} else if claims.ExpiresAt < time.Now().Unix() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token is expired"})
		return
	}

	expTime := time.Now().Add(15 * time.Minute)
	claims.ExpiresAt = expTime.Unix()
	claims.IssuedAt = time.Now().Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshedToken, err := newToken.SignedString(c.Secret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", refreshedToken))
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
	})
}
