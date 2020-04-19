package main

import (
	"flag"
	"fmt"

	"os"

	tlogger "github.com/didi/nightingale/src/toolkits/logger"
	"github.com/freedomkk-qfeng/n9e-plugins/common"

	"github.com/toolkits/pkg/file"
	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/runner"
)

const (
	VERSION = "0.1.0"
)

var (
	vers *bool
	help *bool
	conf *string
)

func init() {
	vers = flag.Bool("v", false, "display the version.")
	help = flag.Bool("h", false, "print this help.")
	conf = flag.String("f", "", "specify configuration file.")
	flag.Parse()

	if *vers {
		fmt.Println("version:", VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	aconf()
	pconf()
	start()
	cfg := common.GetConfig()

	tlogger.Init(cfg.Logger)

	metrics := NginxMetrics(common.GetStatus("nginx").Addr)
	if metrics == nil {
		os.Exit(1)
	}
	if err := common.Push(metrics); err != nil {
		logger.Error(err)
	}
	logger.Close()
}

// auto detect configuration file
func aconf() {
	if *conf != "" && file.IsExist(*conf) {
		return
	}

	*conf = "etc/config.yml"
	if file.IsExist(*conf) {
		return
	}

	fmt.Println("no configuration file for nginx collector")
	os.Exit(1)
}

// parse configuration file
func pconf() {
	if err := common.Parse(*conf); err != nil {
		fmt.Println("cannot parse configuration file:", err)
		os.Exit(1)
	}
}

func start() {
	runner.Init()
	fmt.Println("collector start, use configuration file:", *conf)
	fmt.Println("runner.cwd:", runner.Cwd)
}
