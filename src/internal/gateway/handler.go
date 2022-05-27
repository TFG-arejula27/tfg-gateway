package gateway

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
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
	h.Receptor.State.ChangeOcupation(1)
	//call pymemo
	err := callPymemo(int(h.Receptor.State.Threshold))
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	//calculate time
	h.Receptor.State.ChangeOcupation(-1)
	t := time.Now()
	time := t.Sub(start)
	log.Println(float64(time))
	h.Receptor.State.ChangeExecutionTime(float64(time))
	c.String(http.StatusOK, "ok")

}

func callPymemo(threshold int) error {
	bodyContent := "-t " + strconv.Itoa(threshold)
	body := strings.NewReader(bodyContent)
	req, err := http.NewRequest("POST", "http://localhost:8080/function/threshold", body)

	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Println(err)
		return err
	}

	print(string(respDump))
	defer resp.Body.Close()

	return nil

}
