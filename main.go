package main

import (
	"embed"
	"io/fs"
	"net/http"
	"sw/core"
	"sw/global"
	"sw/router"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist
var f embed.FS

func main() {
	core.InitViper()
	core.InitOrm()
	core.InitRedis()
	core.InitOpcCache()
	core.InitOpc()
	core.InitWs()
	core.InitSw()
	r := router.InitRouter()
	st, _ := fs.Sub(f, "web/dist")
	r.StaticFS("/static", http.FS(st))
	r.NoRoute(func(c *gin.Context) {
		data, err := f.ReadFile("web/dist/index.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})
	r.Run(":" + global.Config.Server.Port)
}
