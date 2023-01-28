package config

import (
	"encoding/json"
	"fmt"
	"strings"
)

func Assert(condition bool, v interface{}, args ...string) {
	k := strings.Join(args, ".")
	if !condition {
		panic(fmt.Errorf("invalid config value [%v] for key [%v] ", v, k))
	}
}

func JSON(args ...string) string {
	s := strings.Join(args, ".")
	data, err := json.Marshal(Config.Sub(s).AllSettings())
	if err != nil {
		panic(err)
	}
	return string(data)
}
