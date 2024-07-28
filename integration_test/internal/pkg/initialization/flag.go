package initialization

import (
	"flag"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
)

func InitFlag() {
	flag.IntVar(&vars.UserNum, "u", 1000, "user num")
	flag.IntVar(&vars.SuperUserNum, "su", 50, "number of users with all friends")
	flag.IntVar(&vars.LargeGroupNum, "gl", 50, "number of big group")
	flag.IntVar(&vars.CommonGroupNum, "gs", 500, "number of small group")
	flag.IntVar(&vars.CommonGroupMemberNum, "gsm", 10, "number of small group member num")
	flag.IntVar(&vars.SingleMessageNum, "sm", 5, "number of single message each user send")
	flag.IntVar(&vars.GroupMessageNum, "gm", 1, "number of group message each user send")

	flag.BoolVar(&vars.ShouldRegister, "reg", false, "determine whether register")
	flag.BoolVar(&vars.ShouldCreateGroup, "crg", false, "determine whether create group")

	flag.BoolVar(&vars.ShouldCheckGroupNum, "ckgn", false, "determine whether check group num")

	flag.Float64Var(&vars.LoginRate, "lgr", 100, "number of login user rate")
}

// SetFlagLimit prevent parameters from exceeding the limit
func SetFlagLimit() {
	vars.UserNum = min(vars.UserNum, vars.MaxUserNum)
	vars.CommonGroupNum = min(vars.CommonGroupMemberNum, vars.UserNum)

	vars.LoginRate = min(vars.LoginRate, 100)
}
