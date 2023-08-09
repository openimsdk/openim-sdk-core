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
	sendNums      int
	recvNums      int
	groupNums     int
	timeInterval  int64
}

func (r *PressureTestAttribute) InitWithFlag() {
	flag.IntVar(&r.messageNumber, "m", 0, "messageNumber for single sender")
	flag.IntVar(&r.sendNums, "s", 0, "the number of senders for testing")
	flag.IntVar(&r.recvNums, "r", 0, "the number of receivers for testing")
	flag.IntVar(&r.groupNums, "g", 0, "group for testing")
	flag.Int64Var(&r.timeInterval, "t", 0, "timeInterVal during sending message")
}

func ParseFlag() {
	flag.Parse()
}
