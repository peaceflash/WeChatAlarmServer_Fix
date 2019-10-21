// WeChatAlarm project main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/peaceflash/WeChatAlarmServer_Fix/common"
	"github.com/peaceflash/WeChatAlarmServer_Fix/config"
	"github.com/peaceflash/WeChatAlarmServer_Fix/control"


	l4g "github.com/alecthomas/log4go"
	"github.com/kardianos/service"
)

var (
	version string = "1.1.1"
)

var logger service.Logger

//  Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

//----------------------------hander ---------------
func MessageHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		l4g.Info("request:%v, %v, %v", r.RemoteAddr, r.RequestURI, string(body))
		//解析传进来的body

		resp := control.MessageNotify(string(body))

		respbody, err := json.Marshal(resp)
		if err != nil {
			l4g.Error("json response failed :%v", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		//w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(respbody))

	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func safehandle(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				l4g.Warn("WARN: panic in %v. - %v", fn, e)
				l4g.Info(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	//sync
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}

func (p *program) run() error {
	logger.Infof("I'm running %v.", service.Platform())
	go Svc()
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case tm := <-ticker.C:
			logger.Infof("Still running at %v...", tm)
		case <-p.exit:
			ticker.Stop()
			return nil
		}
	}
}

func main() {
	svcFlag := flag.String("service", "", "Control the system server")
	bVersion := flag.Bool("version", false, "Software version")
	bhelp := flag.Bool("help", false, "Software version")

	flag.Parse()

	//查看版本号
	if *bVersion {
		fmt.Printf("Current version is :\n")
		fmt.Printf("%s\n", version)
		return
	}

	if *bhelp {
		fmt.Printf("-start     : start software\n")
		fmt.Printf("-stop      : stop software\n")
		fmt.Printf("-install   : install software\n")
		fmt.Printf("-uninstall : uninstall software\n")
		fmt.Printf("-help      : help use software\n")
		fmt.Printf("-version   :  software version\n")

		return
	}

	svcConfig := &service.Config{
		Name:        "WeChatAlarmServer",
		DisplayName: "WeChatAlarmServer",
		Description: "This is a Go alarm service with wechat",
	}

	prg := &program{}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		} else {
			log.Printf("actions: %q Success \n", service.ControlAction)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

func Svc() {
	defer func() {
		err := recover()
		if err != nil {
			l4g.Error("Unrecovered Error:")
			l4g.Error("  The following error was not properly recovered, please report this ASAP!")
			l4g.Error("  %#v\n", err)
			l4g.Error("Stack Trace:")
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, true)]
			l4g.Error("%s\n", buf)
			os.Exit(1)
		}
	}()
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	dir, seq := common.GetPath()
	l4g.LoadConfiguration(dir + seq + "config/log.xml")

	l4g.Info("The version is :%s", version)

	//config init
	config.Init()

	//control  Init
	control.Init()

	//new mux
	mux := http.NewServeMux()

	//注册路由
	mux.HandleFunc("/", safehandle(MessageHandle))

	//启动 pprox
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	serveraddr := "0.0.0.0:" + config.GetServerPort()

	l4g.Info("listen server addr :%v", serveraddr)
	err := http.ListenAndServe(serveraddr, mux)
	if err != nil {
		l4g.Error("ListenAndServer: %v Fail:%v", serveraddr, err.Error())
	}
}
