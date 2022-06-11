package gateway

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"github.com/arejula27/energy-cluster-manager/internal/configuration"
	"github.com/arejula27/energy-cluster-manager/internal/receptor"
	"github.com/gin-gonic/gin"
)

type handler struct {
	Gin      *gin.Engine
	Receptor *receptor.Receptor
	config   configuration.Configuration
}

func NewHandler(receptor *receptor.Receptor, conf configuration.Configuration) *handler {
	r := gin.Default()
	h := &handler{
		Gin:      r,
		Receptor: receptor,
		config:   conf,
	}

	r.POST("/pymemo", h.handlerPymemo)

	return h

}

func (h *handler) handlerPymemo(c *gin.Context) {

	start := time.Now()
	h.Receptor.Manager.Eval()
	//can execute localy
	if !h.Receptor.Manager.Forward {
		log.Println("request queued")
		waitUntil := make(chan bool, 1)
		h.Receptor.Manager.AddExecution(&waitUntil)
		<-waitUntil
		log.Println("request executed")
		h.Receptor.Manager.ChangeOcupation(1)
		defer h.Receptor.Manager.ChangeOcupation(-1)
		err := h.callPymemo(int(h.Receptor.Manager.Threshold))
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		//calculate time

		t := time.Now()
		time := t.Sub(start)
		log.Println(float64(time))
		h.Receptor.Manager.ChangeExecutionTime(float64(time))

	} else {
		//forward

		log.Println("request forwarded")
		h.forwardPymemo()
	}
	//call pymemo

	c.String(http.StatusOK, "ok")

}

func (h *handler) forwardPymemo() error {
	address := h.config.ForwardAdress
	req, err := http.NewRequest("POST", address, nil)
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func (h *handler) callPymemo(threshold int) error {
	time.Sleep(time.Second * 15)
	return nil

	address := "http://localhost:8080/function/threshold"

	bodyContent := "-t " + strconv.Itoa(threshold)
	body := strings.NewReader(bodyContent)
	req, err := http.NewRequest("POST", address, body)

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
