package client

import (
	"database/sql"

	"github.com/dollarkillerx/common/pkg/conf"
	_ "github.com/taosdata/driver-go/v3/taosSql"
)

func TaosdataClient(conf conf.TDengineConfiguration) (*sql.DB, error) {
	var taosUri = "root:taosdata@tcp(localhost:6030)/"
	taos, err := sql.Open("taosSql", taosUri)
	if err != nil {
		return nil, err
	}

	err = taos.Ping()
	return taos, err
}
