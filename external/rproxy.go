package external

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func whitelisted(url *url.URL) bool {
	if url.Scheme == "http" || url.Scheme == "https" {
		if url.Host == "github.com" || url.Host == "raw.githubusercontent.com" || url.Host == "raw.githubusercontents.com" || url.Host == "objects.githubusercontent.com" || url.Host == "codeload.github.com" {
			return true
		}
	}
	return false
}
func reWriteRedirect(resp *http.Response) (err error) {
	if resp.StatusCode == 301 || resp.StatusCode == 302 {
		print(resp.Header.Get("Location"))
		//add outside url to redirected url
		resp.Header.Set("Location", "https://api.daycat.space/rproxy/"+resp.Header.Get("Location"))
	}
	return nil
}

func Rproxy(c *gin.Context) {
	proxyurl := strings.TrimLeft(c.Param("proxyurl"), "/")
	if proxyurl == "" {
		c.String(400, "No URL provided")
		return
	}
	remote, err := url.Parse(proxyurl)
	if err != nil {
		c.String(400, "Invalid URL provided")
		return
	}
	// checks if url is whitelisted
	if !whitelisted(remote) {
		c.String(400, "Bad URL")
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path
		req.Host = remote.Host
	}

	proxy.ModifyResponse = reWriteRedirect

	proxy.ServeHTTP(c.Writer, c.Request)
}
