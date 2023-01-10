# daycatapi
The backend powering https://api.daycat.space

# Endpoints

## /whoami
Returns your IP address

## /ipinfo
Returns info about queried IP address
### params
ip: IP to query (string)
### return
```json
{
"Ip": "114.5.14.2",
"City": "Jakarta",
"Country": "Indonesia",
"Asn": 4761,
"AsnOrg": "INDOSAT Internet Network Provider",
"ISOCode": "ID"
}
```

## /assign
Assigns a domain to point to an IP address
## params:
ip: IP address you want the domain to point to (string)
type: type of record (string, A or AAAA only)
### return
```json
{
"Domain": "M19bGc.dcapi.top",
"ReferenceID": "M19bGcxN9gbxXXXX"
}
```

## /toggleProxy
Changes proxy status of a domain (cloudflare CDN)

### params
referenceid: reference ID of the domain (string)
proxy: whether to enable or disable proxy (bool)

### return
```json
{
"Sucess": true,
"Error": ""
}
```
## /warp
Generates cloudflare warp credentials

### params
stack: ipv4 / ipv6 （optional, string）
format: xray / wireguard / wg-quick (string)
### return
depends on format param


# Open source projects used:
- [wgcf](https://github.com/ViRb3/wgcf)
- [wireguard-go](https://git.zx2c4.com/wireguard-go)
- [gin](https://github.com/gin-gonic/gin)
- [yaml](gopkg.in/yaml.v3)
- [viper](github.com/spf13/viper)
- [go-sqlite3](github.com/mattn/go-sqlite3)
- [maxminddb-golang](github.com/oschwald/maxminddb-golang)
- [cloudflare-go](github.com/cloudflare/cloudflare-go)