package cmd

import (
	"errors"
	// "github.com/charmbracelet/log"
	"net"

	"github.com/bsach64/goback/server"
	"github.com/bsach64/goback/utils"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var mainOptions string
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts a server",
	Long:  "starts a server",
	Run: func(cmd *cobra.Command, args []string) {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Choose an option:").
					Options(
						huh.NewOption("Listen", "Listen"),
						huh.NewOption("Log", "Log"),
						huh.NewOption("Reconstruct", "Reconstruct"),
						huh.NewOption("Exit", "Exit"),
					).
					Value(&mainOptions),
			),
		)
		err := form.Run()

		if err != nil {
			log.Error("failed", "err", err)
			// log.Fatal(err)
		}

		switch mainOptions {
		case "Listen":
			ip, err := getLocalIP()
			if err != nil {
				// log.Fatal(err)
				log.Error("failed", "err", err)
			}
			log.Info("Starting mDNS server!")
			go func(ip net.IP) {
				server, err := utils.StartmDNSServer([]net.IP{ip}, 2022)
				if err != nil {
					// log.Error("mDNS server failed:", "err", err)
					log.Error("mDNS server failed:", "err", err)
				}
				defer func() {
					err = server.Shutdown()
					if err != nil {
						// log.Error("mDNS Server exited:", "err", err)
						log.Error("mDNS server failed:", "err", err)
					}
				}()
				select {}
			}(ip)
			s := server.New(ip.String(), "private/id_rsa", 2022)
			err = server.Listen(s)
			if err != nil {
				log.Info("Could not listen on server")
			}
		case "Log":
			//TODO
		case "Exit":
			//TODO
		}
	},
}

func getLocalIP() (net.IP, error) {
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

func init() {
	rootCmd.AddCommand(serverCmd)
}
