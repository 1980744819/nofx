package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// api/proxy.go
func proxyNofxosData(c *gin.Context) {
	url := "https://nofxos.ai" + c.Request.URL.Path
	log.Printf("Proxying request to %s", url)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Proxy request failed"})
		return
	}
	defer resp.Body.Close()
	log.Printf("Received response with status %d", resp.StatusCode)

	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

// 在api/server.go中注册路由
