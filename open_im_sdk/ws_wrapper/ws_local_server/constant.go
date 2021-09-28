/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/8 15:39).
 */
package ws_local_server

import (
	"reflect"
	"sync"
)

const CMDLogin = "Login"

type RefRouter struct {
	refName  *map[string]reflect.Value
	wsRouter *WsFuncRouter
}

var UserRouteMap map[string]RefRouter

var UserRouteRwLock sync.RWMutex
