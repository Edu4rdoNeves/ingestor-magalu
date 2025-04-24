package login

import (
	"net/http"

	"github.com/Edu4rdoNeves/ingestor-magalu/application/usecases/login"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/gin-gonic/gin"
)

type ILoginController interface {
	Login(context *gin.Context)
}

type LoginController struct {
	usecases login.ILoginUseCase
}

func NewLoginController(usecases login.ILoginUseCase) ILoginController {
	return &LoginController{usecases}
}

func (c *LoginController) Login(context *gin.Context) {
	var input *dto.LoginInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	loginAuth, err := c.usecases.Login(input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error generating token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"Auth": loginAuth})
}
