package main

import (
	"flag"
	"fmt"
	"github.com/AlkBur/GoIDE/conf"
	"github.com/AlkBur/GoIDE/handlers"
	"github.com/AlkBur/GoIDE/i18n"
	"github.com/AlkBur/GoIDE/log"
	"github.com/AlkBur/GoIDE/session"
	"github.com/AlkBur/GoIDE/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

func init() {
	confPath := flag.String("conf", "config.json", "path of config.json")
	confIP := flag.String("ip", "", "this will overwrite IP if specified")
	confPort := flag.Int("port", 8080, "this will overwrite Port if specified")
	confLogLevel := flag.String("log_level", "info", "this will overwrite LogLevel if specified")
	flag.Parse()

	log.SetLevel(log.LevelInfo)

	i18n.Load()
	conf.Load(*confPath, *confIP, *confPort, *confLogLevel)
	conf.CheckEnv()
	session.StartUserMonitor()

	log.Debug("host (" + runtime.Version() + ", " + runtime.GOOS + "_" + runtime.GOARCH + ")")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	handleSignal()

	// IDE
	if log.GetLevel() == log.LevelDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	///////////////Router///////////////
	router := http.NewServeMux()
	router.HandleFunc(conf.IDE.Context+"/", handlerWrapper(handlers.Index, true))

	// static resources
	fs := http.FileServer(http.Dir(filepath.Join(conf.IDE.WD, "static")))
	router.Handle(conf.IDE.Context+"/static/", http.StripPrefix(conf.IDE.Context+"/static/", fs))
	router.HandleFunc(conf.IDE.Context+"/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, conf.IDE.WD+"/static/favicon.ico")
	})

	// user
	router.HandleFunc(conf.IDE.Context+"/login", handlerWrapper(handlers.LoginHandler, false))

	///////////////Start server///////////////
	url := fmt.Sprintf("%s:%d%s", conf.IDE.IP, conf.IDE.Port, conf.IDE.Context)
	log.Info("IDE is running [%s]", url)

	err := http.ListenAndServe(url, router)
	if err != nil {
		log.Error(err)
	}

	//[[range .css]]
	//<script type="text/javascript" src="[[.conf.Context]]/static/css/pp">
	//	[[end]]
}

func handleSignal() {
	go func() {
		c := make(chan os.Signal)

		signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		s := <-c
		log.Debug("Got signal [%s]", s)

		session.SaveActiveUsers()
		log.Debug("Saved all online user, exit")

		os.Exit(0)
	}()
}

func handlerWrapper(f http.HandlerFunc, gz bool) http.HandlerFunc {
	handler := panicRecover(f)
	if gz {
		handler = util.GzipHandler(handler)
	}
	handler = loging(handler)
	//handler = i18nLoad(handler)

	return handler
}

func panicRecover(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer util.Recover()
		handler(w, r)
	}
}

func loging(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Debug("[%s, %s, %s]", r.Method, r.RequestURI, time.Since(start))
		}()
		handler(w, r)
	}
}
