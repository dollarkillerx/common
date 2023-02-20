package client

import (
	"time"

	"github.com/dollarkillerx/common/pkg/conf"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func ETCDClient(conf conf.ETCDConfiguration) (*clientv3.Client, error) {
	return clientv3.New(ETCDOption(conf))
}

func ETCDOption(conf conf.ETCDConfiguration) clientv3.Config {
	return clientv3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: time.Duration(conf.DialTimeout) * time.Second,
		Username:    conf.Username,
		Password:    conf.Password,
	}
}
