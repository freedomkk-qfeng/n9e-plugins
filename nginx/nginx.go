package main

import (
	"strconv"
	"strings"

	"github.com/freedomkk-qfeng/n9e-plugins/common"
	"github.com/toolkits/pkg/logger"

	"github.com/didi/nightingale/src/dataobj"
)

type NginxStatus struct {
	ActiveConn     int
	ServerAccepts  uint64
	ServerHandled  uint64
	ServerRequests uint64
	ServerReading  int
	ServerWriting  int
	ServerWaiting  int
}

func getNginxStatus(body []byte) (NginxStatus, error) {
	var status NginxStatus
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("recovered in panic:%v", r)
		}
	}()

	str := strings.Split(string(body), "\n")
	s := strings.Split(str[0], " ")[2]
	ActiveConn, err := strconv.Atoi(s)
	s = strings.Split(str[2], " ")[1]
	ServerAccepts, err := strconv.ParseUint(s, 10, 64)
	s = strings.Split(str[2], " ")[2]
	ServerHandled, err := strconv.ParseUint(s, 10, 64)
	s = strings.Split(str[2], " ")[3]
	ServerRequests, err := strconv.ParseUint(s, 10, 64)
	s = strings.Split(str[3], " ")[1]
	ServerReading, err := strconv.Atoi(s)
	s = strings.Split(str[3], " ")[3]
	ServerWriting, err := strconv.Atoi(s)
	s = strings.Split(str[3], " ")[5]
	ServerWaiting, err := strconv.Atoi(s)
	status.ActiveConn = ActiveConn
	status.ServerAccepts = ServerAccepts
	status.ServerHandled = ServerHandled
	status.ServerReading = ServerReading
	status.ServerRequests = ServerRequests
	status.ServerWaiting = ServerWaiting
	status.ServerWriting = ServerWriting
	return status, err
}

func NginxMetrics(statusAddr string) (L []*dataobj.MetricValue) {
	res, err := common.HTTPGet(statusAddr)
	if err != nil {
		logger.Errorf("get nginx status page failed:%v", err)
		return
	}
	stat, err := getNginxStatus(res)
	if err != nil {
		logger.Errorf("get nginx status failed:%v", err)
		return
	}

	L = append(L, common.GaugeValue("nginx.conn", stat.ActiveConn))
	L = append(L, common.CounterValue("nginx.accepts", stat.ServerAccepts))
	L = append(L, common.CounterValue("nginx.handled", stat.ServerHandled))
	L = append(L, common.CounterValue("nginx.requests", stat.ServerRequests))
	L = append(L, common.GaugeValue("nginx.reading", stat.ServerReading))
	L = append(L, common.GaugeValue("nginx.waiting", stat.ServerWaiting))
	L = append(L, common.GaugeValue("nginx.writing", stat.ServerWriting))
	return
}
