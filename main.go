package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"
	"strings"

	"github.com/xyproto/term"
)

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
	enableColors := runtime.GOOS != "windows"
	o := term.NewTextOutput(enableColors, true)

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
			fmt.Fprintf(&w, o.White(paddedName))
		} else {
			fmt.Fprintf(&w, o.LightGreen(paddedName))
		}

		hwAddr := iface.HardwareAddr.String()
		if hwAddr == "" {
			fmt.Fprintf(&w, "\t\t%s", o.DarkRed(pad("-", 17)))
		} else {
			fmt.Fprintf(&w, "\t\t%s", o.DarkRed(hwAddr))
		}

		fmt.Fprintf(&w, "  %s", o.DarkPurple("MTU "+strconv.Itoa(iface.MTU)))
		fmt.Fprintf(&w, "  %s", o.DarkGray(iface.Flags.String()))

		fmt.Println(w.String())
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
				adrstr = strings.Replace(adrstr, parts[0], o.LightYellow(parts[0]), -1)
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
