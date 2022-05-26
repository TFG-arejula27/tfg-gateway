package gateway

import (
	"log"
	"time"

	"github.com/arejula27/energy-cluster-manager/internal/receptor"
	"github.com/gin-gonic/gin"
)

type handler struct {
	Gin      *gin.Engine
	Receptor *receptor.Receptor
}

func NewHandler(receptor *receptor.Receptor) *handler {
	r := gin.Default()
	h := &handler{
		Gin:      r,
		Receptor: receptor,
	}

	r.POST("/pymemo", h.handlerPymemo)
	return h

}

func (h *handler) handlerPymemo(c *gin.Context) {

	start := time.Now()
	//call pymemo
	time.Sleep(time.Second * 3)
	//calculate time
	t := time.Now()
	time := t.Sub(start)
	log.Println(float64(time))
	h.Receptor.State.ChangeExecutionTime(float64(time))

}
