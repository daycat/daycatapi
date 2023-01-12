package external

import "C"
import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/curve25519"
	"time"
)

type account struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Model   string `json:"model"`
	Name    string `json:"name"`
	Key     string `json:"key"`
	Account struct {
		ID                       string    `json:"id"`
		AccountType              string    `json:"account_type"`
		Created                  time.Time `json:"created"`
		Updated                  time.Time `json:"updated"`
		PremiumData              int       `json:"premium_data"`
		Quota                    int       `json:"quota"`
		Usage                    int       `json:"usage"`
		WarpPlus                 bool      `json:"warp_plus"`
		ReferralCount            int       `json:"referral_count"`
		ReferralRenewalCountdown int       `json:"referral_renewal_countdown"`
		Role                     string    `json:"role"`
		License                  string    `json:"license"`
	} `json:"account"`
	Config struct {
		ClientID string `json:"client_id"`
		Peers    []struct {
			PublicKey string `json:"public_key"`
			Endpoint  struct {
				V4   string `json:"v4"`
				V6   string `json:"v6"`
				Host string `json:"host"`
			} `json:"endpoint"`
		} `json:"peers"`
		Interface struct {
			Addresses struct {
				V4 string `json:"v4"`
				V6 string `json:"v6"`
			} `json:"addresses"`
		} `json:"interface"`
		Services struct {
			HTTPProxy string `json:"http_proxy"`
		} `json:"services"`
	} `json:"config"`
	Token           string    `json:"token"`
	WarpEnabled     bool      `json:"warp_enabled"`
	WaitlistEnabled bool      `json:"waitlist_enabled"`
	Created         time.Time `json:"created"`
	Updated         time.Time `json:"updated"`
	Tos             time.Time `json:"tos"`
	Place           int       `json:"place"`
	Locale          string    `json:"locale"`
	Enabled         bool      `json:"enabled"`
	InstallID       string    `json:"install_id"`
	FcmToken        string    `json:"fcm_token"`
}
type registerData struct {
	FcmToken  string `json:"fcm_token"`
	InstallID string `json:"install_id"`
	Key       string `json:"key"`
	Locale    string `json:"locale"`
	Model     string `json:"model"`
	Tos       string `json:"tos"`
	Type      string `json:"type"`
}
type registeredDevice []struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Model     string `json:"model"`
	Name      string `json:"name"`
	Created   string `json:"created"`
	Activated string `json:"activated"`
	Active    bool   `json:"active"`
	Role      string `json:"role"`
}
type registerDevice struct {
	Name string `json:"name"`
}
type xrayConfig struct {
	Protocol string `json:"protocol"`
	Settings struct {
		SecretKey string `json:"secretKey"`
		Peers     []struct {
			PublicKey string `json:"publicKey"`
			Endpoint  string `json:"endpoint"`
		} `json:"peers"`
		Address []string `json:"address"`
	} `json:"settings"`
	Tag string `json:"tag"`
}

func generatePrivateKey() (string, [32]byte) {
	// written according to https://cr.yp.to/ecdh.html
	k := [32]byte{}
	// generate random 32 byte key
	_, _ = rand.Read(k[:])
	// set the most significant bit of the first byte to 0
	k[0] &= 248
	// set the least significant bit of the last byte to 0
	k[31] &= 127
	// set the second least significant bit of the last byte to 1
	k[31] |= 64
	// return the base64 encoded key
	return base64.StdEncoding.EncodeToString(k[:]), k
}
func generatePublicKey(k [32]byte) string {
	// written according to https://cr.yp.to/ecdh.html
	var p [32]byte
	// generate public key
	curve25519.ScalarBaseMult(&p, &k)
	// return the base64 encoded key
	return base64.StdEncoding.EncodeToString(p[:])
}
func generateKeyPair() (string, string) {
	privateKey, k := generatePrivateKey()
	publicKey := generatePublicKey(k)
	return privateKey, publicKey
}
func generateConfig(cfgtype, PrivateKey string, account account) string {
	var wgcfg string
	print(cfgtype)
	if cfgtype == "wireguard" {
		print(1)
		// generate wireguard config
		wgcfg = "[Interface]\nPrivateKey = " + PrivateKey + "\nAddress = " + account.Config.Interface.Addresses.V4 + "/32\nAddress = " + account.Config.Interface.Addresses.V6 + "/128\nDNS = 1.1.1.1\nMTU = 1280\n[Peer]\nPublicKey = bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=\nAllowedIPs = 0.0.0.0/0\nAllowedIPs = ::/0\nEndpoint = engage.cloudflareclient.com:2408"

	} else if cfgtype == "xray" {
		print(2)
		var xraycfg xrayConfig
		xraycfg.Protocol = "wireguard"
		xraycfg.Tag = "wireguard-1"
		xraycfg.Settings.SecretKey = PrivateKey
		xraycfg.Settings.Address = append(xraycfg.Settings.Address, account.Config.Interface.Addresses.V4+"/32")
		xraycfg.Settings.Address = append(xraycfg.Settings.Address, account.Config.Interface.Addresses.V6+"/128")

		xraycfg.Settings.Peers = append(xraycfg.Settings.Peers, struct {
			PublicKey string `json:"publicKey"`
			Endpoint  string `json:"endpoint"`
		}{PublicKey: "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=", Endpoint: "engage.cloudflareclient.com:2408"})
		var b []byte
		b, err := json.Marshal(xraycfg)
		if err != nil {
			panic(0)
		}
		wgcfg = string(b)
	} else {
		print(3)
		var xraycfg xrayConfig
		xraycfg.Protocol = "wireguard"
		xraycfg.Tag = "wireguard-1"
		xraycfg.Settings.SecretKey = PrivateKey
		xraycfg.Settings.Address = append(xraycfg.Settings.Address, account.Config.Interface.Addresses.V4+"/32")
		xraycfg.Settings.Address = append(xraycfg.Settings.Address, account.Config.Interface.Addresses.V6+"/128")

		xraycfg.Settings.Peers = append(xraycfg.Settings.Peers, struct {
			PublicKey string `json:"publicKey"`
			Endpoint  string `json:"endpoint"`
		}{PublicKey: "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=", Endpoint: "engage.cloudflareclient.com:2408"})
		var b []byte
		b, err := json.Marshal(xraycfg)
		if err != nil {
			panic(0)
		}
		wgcfg = string(b)

	}

	return wgcfg
}

func GetWarpConfig(c *gin.Context) {
	privateKey, publicKey := generateKeyPair()
	wgcfgformat := c.Query("format")
	accountData := Register(publicKey)
	if accountData.Config.Interface.Addresses.V4 == "" {
		c.String(500, "Warp endpoint overloaded")
		return
	}
	wgconfig := generateConfig(wgcfgformat, privateKey, accountData)
	c.String(200, wgconfig)
	go Activate(accountData)
}
