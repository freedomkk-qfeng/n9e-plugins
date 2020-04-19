package common

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/didi/nightingale/src/toolkits/logger"
	"github.com/spf13/viper"
	"github.com/toolkits/pkg/file"
)

type ConfYaml struct {
	Logger  logger.LoggerSection `yaml:"logger"`
	Specify string               `yaml:"specify"`
	Step    int64                `yaml:"step"`
	PushAPI string               `yaml:"pushApi"`
}

var (
	Config   *ConfYaml
	lock     = new(sync.RWMutex)
	Endpoint string
	Cwd      string
)

// GetConfig configuration file
func GetConfig() *ConfYaml {
	lock.RLock()
	defer lock.RUnlock()
	return Config
}

func Parse(conf string) error {
	bs, err := file.ReadBytes(conf)
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", conf, err)
	}

	lock.Lock()
	defer lock.Unlock()

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(bs))
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", conf, err)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("Unmarshal %v", err)
	}

	return nil
}
