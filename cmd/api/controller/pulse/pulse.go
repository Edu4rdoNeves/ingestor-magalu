package pulse

import (
	"net/http"
	"strconv"

	"github.com/Edu4rdoNeves/ingestor-magalu/application/usecases/pulse"
	"github.com/gin-gonic/gin"
)

type IPulseController interface {
	GetPulses(context *gin.Context)
}

type PulseController struct {
	usecases pulse.IPulseUseCase
}

func NewPulseController(usecases pulse.IPulseUseCase) IPulseController {
	return &PulseController{usecases}
}

func (c *PulseController) GetPulses(context *gin.Context) {
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(context.DefaultQuery("limit", "10"))

	pulseResponse, err := c.usecases.GetPulses(page, limit)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"Error: ": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, pulseResponse)
}

func (c *PulseController) GetPulseByID(context *gin.Context) {

}
