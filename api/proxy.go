package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// api/proxy.go
func proxyNofxosData(c *gin.Context) {
	url := "https://nofxos.ai" + c.Request.URL.Path
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Proxy request failed"})
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

// 在api/server.go中注册路由
