package main

import (
	"github.com/daycat/daycatapi/config"
	"github.com/daycat/daycatapi/external"
	"github.com/daycat/daycatapi/networking"
	"github.com/gin-gonic/gin"
)

func main() {
	config.GetConfig("config.yaml")
	r := gin.Default()
	r.GET("/whoami", networking.Whoami)
	r.GET("/ipinfo", networking.IpInfo)
	r.GET("/assign", networking.AssignDomain)
	r.GET("/toggleProxy", networking.ToggleProxy)
	r.GET("/rproxy/*proxyurl", external.Rproxy)
	r.GET("/warp", external.GetWarpConfig)
	r.Run()

}
