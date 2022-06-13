package gateway

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/arejula27/energy-cluster-manager/internal/configuration"
	"github.com/arejula27/energy-cluster-manager/internal/manager"
	"github.com/gin-gonic/gin"
)

type handler struct {
	sync.Mutex
	Gin       *gin.Engine
	manager   *manager.Manager
	config    configuration.Configuration
	ocupation int
	replies   chan *chan bool
}

func NewHandler(manager *manager.Manager, conf configuration.Configuration) *handler {
	r := gin.Default()
	h := &handler{
		Gin:     r,
		manager: manager,
		config:  conf,
		replies: make(chan *chan bool, 100),
	}

	r.POST("/pymemo", h.handlerPymemo)

	return h

}

func (h *handler) handlerPymemo(c *gin.Context) {

	start := time.Now()
	defer func() {
		end := time.Now()
		time := end.Sub(start)
		h.manager.ChangeExecutionTime(time)
	}()

	//comprobar ocupación
	availability := h.checkAvailability()

	//si no hay espacio
	//delegar si hay nivel superior
	if h.config.ForwardAdress != "" && !availability {
		//TODO forward
		log.Println("Request fordwarded")
		h.forwardPymemo()
		c.String(http.StatusOK, "ok")
		return

	}

	//si se tiene que ejecutar en local
	//encolo petición
	waitUntil := make(chan bool, 4)
	h.replies <- &waitUntil

	//Si hay espacio liberar primero de cola
	if availability {
		h.popRequest()

	}
	//espero a que sea mi turno
	<-waitUntil
	h.Lock()
	h.ocupation++
	h.Unlock()

	//Free resouces when ended
	defer func() {
		h.Lock()
		h.ocupation--
		h.Unlock()
		if h.checkAvailability() {
			h.popRequest()
		}

	}()
	//ejecuto
	log.Println("Request  start execution")

	err := h.callPymemo(int(h.manager.GetThreshold()))

	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.String(http.StatusOK, "ok")
	log.Println("request executed")
	//calculate time

}

func (h *handler) checkAvailability() bool {
	h.Lock()
	availabilty := h.ocupation < h.manager.GetMaxOcupation()
	h.Unlock()
	return availabilty
}

func (h *handler) popRequest() {
	//si hay esperando
	if len(h.replies) > 0 {
		r := <-h.replies
		*r <- true
	}

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
	//respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Println(err)
		return err
	}

	//print(string(respDump))
	defer resp.Body.Close()

	return nil

}
