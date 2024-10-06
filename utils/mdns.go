// mDNS server can be used to allow for metadeta server to be discovered by devices on the network
package utils

import (
	"errors"
	"net"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/mdns"
)

const (
	service = "_sftp-ssh._tcp"
	domain  = "goback-server.local."
)

func StartmDNSServer(ip []net.IP, port int) (*mdns.Server, error) {
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	info := []string{"goback"}

	serv, err := mdns.NewMDNSService(
		host,
		service,
		domain,
		"",
		port,
		ip,
		info,
	)

	if err != nil {
		return nil, err
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: serv})
	if err != nil {
		return nil, err
	}
	return server, err
}

func LookupmDNSService() (net.IP, net.IP, int, error) {
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		query := mdns.QueryParam{
			Service: service,
			Domain:  domain,
			Timeout: time.Second,
			Entries: entriesCh,
		}

		_ = mdns.Query(&query)
		close(query.Entries)
		wg.Done()
	}()

	defer wg.Wait()
	entry := <-entriesCh
	if entry == nil {
		return net.IP{}, net.IP{}, 0, errors.New("Could not find mDNS entry")
	}
	return entry.AddrV4, entry.AddrV6, entry.Port, nil
}

func GetLocalIP() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return net.IP{}, err
	}

	for _, inter := range interfaces {
		if inter.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if inter.Flags&net.FlagLoopback != 0 {
			continue // Loopback Interface
		}

		addresses, err := inter.Addrs()
		if err != nil {
			return net.IP{}, err
		}
		for _, addr := range addresses {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}
			return ip, nil
		}
	}
	return net.IP{}, errors.New("Are you connected to the internet?")
}
