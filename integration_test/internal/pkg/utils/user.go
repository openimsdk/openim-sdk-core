package utils

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"strconv"
	"strings"
)

func GetUserID(num int) string {
	return vars.UserIDPrefix + strconv.Itoa(num)
}

func GetUserNum(id string) (int, error) {
	if !strings.HasPrefix(id, vars.UserIDPrefix) {
		return -1, errs.New("invalid user id in GetUserNum").Wrap()
	}
	num, err := strconv.Atoi(strings.TrimPrefix(id, vars.UserIDPrefix))
	if err != nil {
		return -1, errs.WrapMsg(err, "invalid user id in GetUserNum")
	}
	return num, nil
}

func MustGetUserNum(id string) int {
	num, err := GetUserNum(id)
	if err != nil {
		log.ZError(context.TODO(), "err in MustGetUserNum", err)
	}
	return num
}

// IsSuperUser check if the user has all friends
func IsSuperUser(id string) bool {
	num := MustGetUserNum(id)
	return datautil.BetweenLEq(num, 0, vars.SuperUserNum)
}

// NextOffsetNums get num with an offset behind.
func NextOffsetNums(userNum, offset int) []int {
	ids := make([]int, offset)
	for i := 1; i <= offset; i++ {
		ids[i-1] = NextOffsetNum(userNum, i)
	}
	return ids
}

// NextOffsetNum get num with an offset behind.
func NextOffsetNum(num, offset int) int {
	return (num + offset) % vars.UserNum
}

// NextNum get next num.
func NextNum(num int) int {
	return NextOffsetNum(num, 1)
}

// PreOffsetNum get num with an offset forward.
func PreOffsetNum(num, offset int) int {
	return (num + offset) % vars.UserNum
}

// NextOffsetUserIDs get userIDs with an offset behind.
func NextOffsetUserIDs(userNum, offset int) []string {
	ids := make([]string, offset)
	for i := 1; i <= offset; i++ {
		ids[i-1] = GetUserID(NextOffsetNum(userNum, i))
	}
	return ids
}
