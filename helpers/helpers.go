package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func wd() string {
	d, _ := os.Getwd()
	return strings.Replace(d, "consoleapps", "webapps", -1) + "/"
}
func readconfigfile(file *os.File) map[string]string {
	reader := bufio.NewReader(file)
	ret := make(map[string]string)
	for {
		line, _, e := reader.ReadLine()
		if e != nil {
			break
		}

		sval := strings.Split(string(line), "=")
		if len(sval) > 2 {
			ret[sval[0]] = strings.Replace(string(line), sval[0]+"=", "", -1)
		} else {
			ret[sval[0]] = sval[1]
		}

	}
	return ret
}
func ReadAppConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(wd() + "conf/app.conf")
	if err == nil {
		defer file.Close()
		ret = readconfigfile(file)

	} else {
		fmt.Println(err.Error())
	}

	return ret
}
