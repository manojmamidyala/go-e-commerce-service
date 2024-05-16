package http

import (
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
