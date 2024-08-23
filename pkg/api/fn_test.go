package api

import (
	"github.com/openimsdk/protocol/group"
	"testing"
)

func TestName(t *testing.T) {
	var _ group.GetGroupApplicationListResp
}

func TestName2(t *testing.T) {

	fn := (*group.GetGroupApplicationListResp).GetGroupRequests
	_ = fn
	t.Log()

}
