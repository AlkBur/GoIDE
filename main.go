package main

import (
	"flag"
	"github.com/AlkBur/GoIDE/log"
	"runtime"
	"github.com/AlkBur/GoIDE/conf"
)

func init() {
	confPath := flag.String("conf", "config.json", "path of config.json")
	confIP := flag.String("ip", "127.0.0.1", "this will overwrite IP if specified")
	confPort := flag.Int("port", 8080, "this will overwrite Port if specified")
	confLogLevel := flag.String("log_level", "debug", "this will overwrite LogLevel if specified")
	flag.Parse()

	log.SetLevel(log.LevelInfo)

	conf.Load(*confPath, *confIP, *confPort, *confLogLevel)
	conf.CheckEnv()

	log.Debug("host ("+runtime.Version()+", "+runtime.GOOS+"_"+runtime.GOARCH+")")
}

func main() {

}
