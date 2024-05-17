package http

import (
	"errors"
	"mami/e-commerce/commons/logger"
	"mami/e-commerce/pkg/utils"
	responsedto "mami/e-commerce/responseDto"
	"mami/e-commerce/user/dto"
	"mami/e-commerce/user/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.IUserService
}

func NewUserhandler(service service.IUserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get boy user.Login ", err)
		responsedto.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	user, accessToken, refreshToken, err := h.service.Login(c, &req)
	if err != nil {
		logger.Error("Failed to login ", err)
		responsedto.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.LoginRes
	utils.Copy(&user, &res.User)
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	responsedto.JSON(c, http.StatusOK, res)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		responsedto.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	user, err := h.service.Register(c, &req)
	if err != nil {
		logger.Error(err.Error())
		responsedto.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.RegisterRes
	utils.Copy(&user, &res.User)
	responsedto.JSON(c, http.StatusOK, res)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		responsedto.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}

	user, err := h.service.GetUserByID(c, userID)
	if err != nil {
		logger.Error(err.Error())
		responsedto.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.User
	utils.Copy(&user, &res)
	responsedto.JSON(c, http.StatusOK, res)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		responsedto.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}

	accessToken, err := h.service.RefreshToken(c, userID)
	if err != nil {
		logger.Error("Failed to refresh token", err)
		responsedto.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	res := dto.RefreshTokenRes{
		AccessToken: accessToken,
	}
	responsedto.JSON(c, http.StatusOK, res)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		responsedto.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	userID := c.GetString("userId")
	err := h.service.ChangePassword(c, userID, &req)
	if err != nil {
		logger.Error(err.Error())
		responsedto.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	responsedto.JSON(c, http.StatusOK, nil)
}
