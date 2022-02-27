package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/xyproto/textoutput"
)

const versionString = "interfaces 1.2.1"

var (
	noHighlightPrefixes = []string{"vbox", "docker", "lo"}
)

func pad(s string, n int) string {
	var padding string
	for i := 0; i < (n - len(s)); i++ {
		padding += " "
	}
	return s + padding
}

func main() {
	usage := `interfaces

Usage:
  interfaces
  interfaces -l | --long
  interfaces -h | --help
  interfaces -v | --version

Options:
  -h --help     This help screen
  -v --version  Version information
  -l --long     Longer output`

	enableColors := runtime.GOOS != "windows"
	o := textoutput.NewTextOutput(enableColors, true)

	// Parse arguments
	arguments, err := docopt.Parse(usage, nil, true, versionString, false)
	if err != nil {
		log.Fatalln(err)
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalln(err)
	}

	for _, iface := range ifaces {
		var w bytes.Buffer
		fmt.Fprintf(&w, "%s%s%s ", o.DarkGray("["), o.LightBlue(strconv.Itoa(iface.Index)), o.DarkGray("]"))

		highlight := true
		for _, noh := range noHighlightPrefixes {
			if strings.HasPrefix(iface.Name, noh) {
				highlight = false
				break
			}
		}

		paddedName := pad(iface.Name, 8)

		if highlight {
			fmt.Fprintf(&w, o.DarkRed(paddedName))
		} else {
			fmt.Fprintf(&w, o.LightGreen(paddedName))
		}

		hwAddr := iface.HardwareAddr.String()
		if hwAddr == "" {
			fmt.Fprintf(&w, "\t\t%s", o.LightYellow(pad("-", 17)))
		} else {
			fmt.Fprintf(&w, "\t\t%s", o.LightYellow(hwAddr))
		}

		fmt.Fprintf(&w, "  %s", o.DarkPurple(pad("MTU "+strconv.Itoa(iface.MTU), 9)))
		fmt.Fprint(&w, "  ")
		flags := strings.Split(iface.Flags.String(), "|")
		if len(flags) > 0 && flags[0] != "up" {
			fmt.Fprintf(&w, o.DarkGray("↓    | "))
		}
		for i, flag := range flags {
			if i > 0 {
				fmt.Fprintf(&w, o.DarkGray(" | "))
			}
			if flag == "up" {
				fmt.Fprintf(&w, o.LightGreen("↑ ") + o.DarkGreen(flag))
			} else if flag == "loopback" {
				fmt.Fprintf(&w, o.DarkBlue(flag))
			} else {
				fmt.Fprintf(&w, o.DarkCyan(flag))
			}
		}

		fmt.Println(w.String())

		if !arguments["--long"].(bool) {
			// Skip the interface details
			continue
		}

		w = bytes.Buffer{}

		addrs, err := iface.Addrs()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, a := range addrs {
			adrstr := a.String()
			if strings.Contains(adrstr, "/") {
				parts := strings.Split(adrstr, "/")
				adrstr = strings.Replace(adrstr, parts[0], o.White(parts[0]), -1)
			}
			fmt.Fprintf(&w, "  %s\t%s\t%s\n", o.LightBlue("adr"), o.White(pad(adrstr, 32+11)), o.DarkGray("(")+a.Network()+o.DarkGray(")"))
		}

		mAddrs, err := iface.MulticastAddrs()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, ma := range mAddrs {
			fmt.Fprintf(&w, "  %s\t%s\t%s\n", o.LightPurple("mul"), o.DarkBlue(pad(ma.String(), 32)), o.DarkGray("(")+ma.Network()+o.DarkGray(")"))
		}

		fmt.Println(w.String())

	}
}
