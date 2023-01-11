package external

import (
	"crypto/tls"
	"encoding/json"
	"github.com/daycat/daycatapi/networking"
	"github.com/k0kubun/pp/v3"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	DefaultTransport = &http.Transport{
		// Match app's TLS config or API will reject us with code 403 error 1020
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS12},
		ForceAttemptHTTP2: false,
		// From http.DefaultTransport
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	} // from wgcf
	DefaultHeader = http.Header{
		"User-Agent":        {"okhttp/3.12.1"},
		"Accept":            {"application/json"},
		"Cf-Client-Version": {"a-6.3-1922"},
		"Content-Type":      {"application/json"},
	}
)

func Register(pk string) account {
	t := time.Now()
	timenow := t.Format(time.RFC3339Nano)
	rD := registerData{
		InstallID: "",
		Key:       pk,
		Locale:    "en_US",
		Model:     "dayCat api instance",
		Tos:       timenow,
		Type:      "Android",
	}
	rDJson, _ := json.Marshal(rD)

	// registers account
	client := http.Client{Transport: DefaultTransport}
	req, err := http.NewRequest("POST", "https://api.cloudflareclient.com/v0a1922/reg", strings.NewReader(string(rDJson)))
	req.Header = DefaultHeader
	resp, err := client.Do(req)
	if err != nil {
		pp.Print(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		pp.Print(err)
	}
	var accountData account
	err = json.Unmarshal(body, &accountData)
	if err != nil {
		pp.Print(err)
	}
	return accountData
}
func Activate(accountData account) {
	// register device
	client := http.Client{
		Transport: DefaultTransport,
	}
	deviceName := networking.RandString(6)
	rD2 := registerDevice{
		Name: deviceName,
	}
	rD2Json, _ := json.Marshal(rD2)
	req, err := http.NewRequest("PATCH", "https://api.cloudflareclient.com/v0a1922/reg/"+accountData.ID+"/account/reg/"+accountData.ID, strings.NewReader(string(rD2Json)))
	if err != nil {
		return
	}
	req.Header = DefaultHeader
	req.Header.Set("Authorization", "Bearer "+accountData.Token)
	resp, err := client.Do(req)
	if err != nil {
		// handle err
		return
		//pp.Print(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
		//pp.Print(err)
	}
	var device registeredDevice
	json.Unmarshal(body, &device)
	req, err = http.NewRequest("GET", "https://api.cloudflareclient.com/v0a1922/reg/"+accountData.ID, nil)
	req.Header = DefaultHeader
	req.Header.Set("Authorization", "Bearer "+accountData.Token)
	resp, err = client.Do(req)
	if err != nil {
		return
		//pp.Print(err)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
		//pp.Print(err)
	}
	req, err = http.NewRequest("PATCH", "https://api.cloudflareclient.com/v0a1922/reg/"+accountData.ID+"/account/reg/"+accountData.ID, strings.NewReader("{\"active\":true}"))
	req.Header = DefaultHeader
	req.Header.Set("Authorization", "Bearer "+accountData.Token)
	resp, err = client.Do(req)
	if err != nil {
		//pp.Print(err)
		return
	}
	defer resp.Body.Close()
	//pp.Print(rDed)

}
