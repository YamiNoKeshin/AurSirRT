package config

import "flag"


var Debug = false
func init(){
	flag.Var(&Zconnections, "zconnection", "e.g. 192.168.0.1:5555, if no port specified p2p will be enabled on the iface")
	flag.BoolVar(&Debug,"debug", false, "")

}

