package networking

import (
	"context"
	"crypto/rand"
	"database/sql"
	"embed"
	"github.com/cloudflare/cloudflare-go"
	"github.com/daycat/daycatapi/config"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oschwald/maxminddb-golang"
	"math/big"
	"net"
	"net/http"
)

//go:embed GeoLite2-City.mmdb
//go:embed GeoLite2-ASN.mmdb
var f embed.FS

type IpCity struct {
	City struct {
		Names struct {
			En string `maxminddb:"en"`
		} `maxminddb:"names"`
	} `maxminddb:"city"`
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
		Names   struct {
			En string `maxminddb:"en"`
		} `maxminddb:"names"`
	} `maxminddb:"country"`
}
type IpAsn struct {
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number"`
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization"`
}
type IpRecord struct {
	Ip      string
	City    string
	Country string
	Asn     uint
	AsnOrg  string
	ISOCode string
}
type GeneralStatus struct {
	Success bool
	Error   string
}
type DomainResponse struct {
	Domain      string
	ReferenceID string
}
type RecordDBHandler struct {
	db *sql.DB
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		c, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			panic(err)
		}
		b[i] = letterBytes[c.Int64()]
	}
	return string(b)
}

func Whoami(c *gin.Context) {
	c.String(http.StatusOK, c.ClientIP()+"\n")
}

func IpInfo(c *gin.Context) {
	ip := c.Query("ip")
	city, _ := f.ReadFile("GeoLite2-City.mmdb")
	asn, _ := f.ReadFile("GeoLite2-ASN.mmdb")
	dbc, err := maxminddb.FromBytes(city)
	dba, err := maxminddb.FromBytes(asn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GeneralStatus{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	defer dbc.Close()
	defer dba.Close()
	query := net.ParseIP(ip)
	var resultc IpCity
	var resulta IpAsn
	err = dbc.Lookup(query, &resultc)
	err = dba.Lookup(query, &resulta)
	var result IpRecord
	result.Ip = ip
	result.City = resultc.City.Names.En
	result.Country = resultc.Country.Names.En
	result.Asn = resulta.AutonomousSystemNumber
	result.AsnOrg = resulta.AutonomousSystemOrganization
	result.ISOCode = resultc.Country.ISOCode
	if err != nil {
		c.JSON(http.StatusInternalServerError, GeneralStatus{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
	return
}

func AssignDomain(c *gin.Context) {
	ip := c.Query("ip")
	rtype := c.Query("type")
	if rtype != "A" && rtype != "AAAA" {
		c.JSON(http.StatusBadRequest, GeneralStatus{
			Success: false,
			Error:   "Invalid Record Type",
		})
		return
	}

	Domain := RandString(6) + config.RootDomain
	ReferenceID := RandString(16)
	var result DomainResponse
	result.Domain = Domain
	result.ReferenceID = ReferenceID
	// stores the domain and reference id in a database
	go MkRecord(Domain, ReferenceID, ip, rtype)
	c.JSON(http.StatusOK, result)

}

func MkRecord(domain string, referenceid string, ip string, rtype string) {
	api, err := cloudflare.New(config.ApiKey, config.ApiEmail)
	ctx := context.Background()
	resp, err := api.CreateDNSRecord(ctx, config.Zoneid, cloudflare.DNSRecord{
		Type:    rtype,
		Name:    domain,
		Content: ip,
		ZoneID:  config.Zoneid,
		TTL:     60,
	})
	if err != nil {
		return
	}
	var CloudflareRecord string
	CloudflareRecord = resp.Result.ID

	sqlDB, err := sql.Open("sqlite3", "./networking/domains.db")
	if err != nil {
		return
	}
	defer sqlDB.Close()
	_, errdb := sqlDB.Exec("INSERT INTO domains (Domains, ReferenceID, CloudflareRecord) VALUES (?, ?, ?)", domain, referenceid, CloudflareRecord)
	if errdb != nil {
		print("Error inserting record")
		print(errdb)
	}

}

func BoolAddr(b bool) *bool {
	boolVar := b
	return &boolVar
}

func ToggleProxy(c *gin.Context) {
	//referenceid := c.Query("referenceid")
	var proxyState bool
	if c.Query("proxy") == "true" {
		proxyState = true
	} else if c.Query("proxy") == "false" {
		proxyState = false
	} else {
		c.JSON(http.StatusBadRequest, GeneralStatus{Success: false, Error: "Invalid Proxy State: " + c.Query("proxy") + " is not a valid state (true/false)"})
		return
	}
	proxy := BoolAddr(proxyState)
	sqlDB, err := sql.Open("sqlite3", "./networking/domains.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, GeneralStatus{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	defer sqlDB.Close()
	row := sqlDB.QueryRow("SELECT CloudflareRecord FROM Domains WHERE ReferenceID = ?", c.Query("referenceid"))
	var CloudflareRecord string
	err = row.Scan(&CloudflareRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GeneralStatus{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	api, err := cloudflare.New(config.ApiKey, config.ApiEmail)
	ctx := context.Background()
	err = api.UpdateDNSRecord(ctx, config.Zoneid, CloudflareRecord, cloudflare.DNSRecord{
		Proxied: proxy,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, GeneralStatus{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GeneralStatus{
		Success: true,
		Error:   "",
	})
}
