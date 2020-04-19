package common

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/toolkits/pkg/file"
	"github.com/toolkits/pkg/runner"
)

type Status struct {
	Addr string `yaml:"addr"`
}

var (
	alock sync.Once
	mods  map[string]Status
)

func GetStatus(modKey string) Status {
	alock.Do(func() {
		parseConf()
	})

	mod, has := mods[modKey]
	if !has {
		fmt.Printf("module(%s) status section not found", modKey)
		os.Exit(1)
	}

	return mod
}

func parseConf() {
	conf := getConf()

	var c map[string]Status
	err := file.ReadYaml(conf, &c)
	if err != nil {
		fmt.Println("cannot parse file:", conf)
		os.Exit(1)
	}

	mods = c
}

func getConf() string {
	conf := path.Join(runner.Cwd, "etc", "status.yml")
	if file.IsExist(conf) {
		return conf
	}

	fmt.Println("configuration file status.yml not found")
	os.Exit(1)
	return ""
}
