package main

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"
	"strings"

	"github.com/xyproto/term"
)

func main() {
	enableColors := runtime.GOOS != "windows"
	o := term.NewTextOutput(enableColors, true)

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalln(err)
	}

	for _, iface := range ifaces {
		fmt.Printf("%s%s%s ", o.DarkGray("["), o.LightBlue(strconv.Itoa(iface.Index)), o.DarkGray("]"))
		if strings.HasPrefix(iface.Name, "vbox") || strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "lo") {
			fmt.Print(o.LightGreen(iface.Name))
		} else {
			fmt.Print(o.White(iface.Name))
		}
		fmt.Printf(" %s", o.DarkGray(iface.Flags.String()))
		fmt.Printf("\t%s", o.LightYellow(iface.HardwareAddr.String()))
		fmt.Printf("\t%s", o.DarkPurple("MTU "+strconv.Itoa(iface.MTU)))
		fmt.Println()

		addrs, err := iface.Addrs()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, a := range addrs {
			adrstr := a.String()
			if strings.Contains(adrstr, "/") {
				parts := strings.Split(adrstr, "/")
				adrstr = strings.Replace(adrstr, parts[0], o.DarkRed(parts[0]), -1)
			}
			fmt.Printf("\t%s\t%s %s\n", o.LightBlue("adr"), o.White(adrstr), o.DarkGray("(")+a.Network()+o.DarkGray(")"))
		}

		maddrs, err := iface.MulticastAddrs()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, ma := range maddrs {
			fmt.Printf("\t%s\t%s %s\n", o.LightPurple("mul"), o.DarkBlue(ma.String()), o.DarkGray("(")+ma.Network()+o.DarkGray(")"))
		}

	}
}
