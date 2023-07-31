package testv3new

import (
	"flag"
	"strings"
)

type sliceValue []string

func (s *sliceValue) String() string {
	return ""
}

func (s *sliceValue) Set(val string) error {
	*s = sliceValue(strings.Split(val, ","))
	return nil
}

type PressureTestAttribute struct {
	messageNumber int
	sendUserIDs   sliceValue
	recvUserIDs   sliceValue
	groupIDs      sliceValue
	timeInterval  int64
}

func (r *PressureTestAttribute) InitWithFlag() {
	flag.IntVar(&r.messageNumber, "m", 0, "messageNumber for single sender")
	flag.Var(&r.sendUserIDs, "s", "sender id list")
	flag.Var(&r.recvUserIDs, "r", "recv id list")
	flag.Var(&r.groupIDs, "g", "groupID for testing")
	flag.Int64Var(&r.timeInterval, "t", 0, "timeInterVal during sending message")
}

func ParseFlag() {
	flag.Parse()
}
