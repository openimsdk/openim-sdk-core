package page

import "github.com/openimsdk/protocol/sdkws"

type PageReq interface {
	GetPagination() *sdkws.RequestPagination
}
