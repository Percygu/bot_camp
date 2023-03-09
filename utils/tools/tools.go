package tools

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

func Prettier(val interface{}) {
	resp2, _ := json.MarshalIndent(val, "", "    ")
	fmt.Println(string(resp2))
	logrus.Info(string(resp2))
}
