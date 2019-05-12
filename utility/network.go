package utility

import (
	"fmt"
	"net"
)

func GetIPADDRESS(hostname string) {
	ip, _ := net.LookupIP(hostname)
	fmt.Println(ip)
}
