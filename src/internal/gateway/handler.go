package gateway

import "github.com/gin-gonic/gin"

func NewHandler() *gin.Engine {

	r := gin.Default()
	r.POST("/pymemo", handlerPymemo)
	return r

}

func handlerPymemo(c *gin.Context) {
	//call pymemo

	//calculate time
}
