/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/4/8 15:09).
 */
package utils

import (
	"encoding/json"
	"net"
	"strconv"
)

func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func StringToInt(i string) int {
	j, _ := strconv.Atoi(i)
	return j
}
func StringToInt64(i string) int64 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return j
}

// judge a string whether in the  string list
func IsContain(target string, List []string) bool {

	for _, element := range List {

		if target == element {
			return true
		}
	}
	return false

}
func InterfaceArrayToStringArray(data []interface{}) (i []string) {
	for _, param := range data {
		i = append(i, param.(string))
	}
	return i
}
func StructToJsonString(param interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

// The incoming parameter must be a pointer
func JsonStringToStruct(s string, args interface{}) error {
	err := json.Unmarshal([]byte(s), args)
	return err
}

var ServerIP = ""

func Init() {
	//fixme In the configuration file, ip takes precedence, if not, get the valid network card ip of the machine
	//if config.Config.ServerIP != "" {
	//	ServerIP = config.Config.ServerIP
	//	return
	//}
	//fixme Get the ip of the local network card
	netInterfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(netInterfaces); i++ {
		//Exclude useless network cards by judging the net.flag Up flag
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			address, _ := netInterfaces[i].Addrs()
			for _, addr := range address {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						ServerIP = ipNet.IP.String()
						return
					}
				}
			}
		}
	}
}
