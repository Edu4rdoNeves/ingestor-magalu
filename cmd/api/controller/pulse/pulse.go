package pulse

import (
	"net/http"
	"strconv"

	"github.com/Edu4rdoNeves/ingestor-magalu/application/usecases/pulse"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/gin-gonic/gin"
)

type IPulseController interface {
	GetPulses(context *gin.Context)
	GetPulseByID(context *gin.Context)
	PopulateQueueWithPulses(ctx *gin.Context)
}

type PulseController struct {
	usecases pulse.IPulseUseCase
}

func NewPulseController(usecases pulse.IPulseUseCase) IPulseController {
	return &PulseController{usecases}
}

func (c *PulseController) GetPulses(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'page' parameter"})
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'limit' parameter"})
		return
	}

	pulses, err := c.usecases.GetPulses(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pulses)
}

func (c *PulseController) GetPulseByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	pulseID, err := strconv.Atoi(idParam)
	if err != nil || pulseID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	pulse, err := c.usecases.GetPulseByID(pulseID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Pulse not found: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pulse)
}

func (c *PulseController) PopulateQueueWithPulses(ctx *gin.Context) {
	var params *dto.PopulateQueueParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters", "details": err.Error()})
		return
	}

	if err := c.usecases.PopulateQueueWithPulses(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Simulation started successfully",
	})
}
