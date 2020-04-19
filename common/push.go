package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/didi/nightingale/src/dataobj"
	"github.com/toolkits/pkg/logger"
)

func Push(metricItems []*dataobj.MetricValue) error {
	var err error
	var items []*dataobj.MetricValue
	for _, item := range metricItems {
		if GetConfig().Specify != "" {
			item.Endpoint = GetConfig().Specify
		}
		item.Step = GetConfig().Step
		item.Timestamp = time.Now().Unix()

		logger.Debugf("push item: %v", item)
		items = append(items, item)
	}
	bs, err := json.Marshal(items)
	if err != nil {
		msg := fmt.Errorf("json marshal failed:%v", err)
		return msg
	}
	data := bytes.NewBuffer(bs)
	return N9ePush(GetConfig().PushAPI, data)
}
