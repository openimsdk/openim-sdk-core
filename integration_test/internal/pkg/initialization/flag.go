package initialization

import (
	"flag"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/internal/flagconst"
)

func InitFlag() {
	flag.BoolVar(&flagconst.TestMode, "test", false, "mark is test mode")

	flag.IntVar(&vars.UserNum, "u", 100, "user num")
	flag.IntVar(&vars.SuperUserNum, "su", 10, "number of users with all friends")
	flag.IntVar(&vars.LargeGroupNum, "gl", 5, "number of big group")
	flag.IntVar(&vars.CommonGroupNum, "gs", 10, "number of small group")
	flag.IntVar(&vars.CommonGroupMemberNum, "gsm", 20, "number of small group member num")
	flag.IntVar(&vars.SingleMessageNum, "sm", 5, "number of single message each user send")
	flag.IntVar(&vars.GroupMessageNum, "gm", 1, "number of group message each user send")

	flag.BoolVar(&vars.ShouldRegister, "reg", false, "determine whether register")
	flag.BoolVar(&vars.ShouldImportFriends, "imf", false, "determine whether import friends")
	flag.BoolVar(&vars.ShouldCreateGroup, "crg", false, "determine whether create group")
	flag.BoolVar(&vars.ShouldSendMsg, "sem", false, "determine whether send messages")

	flag.BoolVar(&vars.ShouldCheckGroupNum, "ckgn", false, "determine whether check group num")
	flag.BoolVar(&vars.ShouldCheckConversationNum, "ckcon", false, "determine whether check conversation num")
	flag.BoolVar(&vars.ShouldCheckMessageNum, "ckmsn", false, "determine whether check message num")
	flag.BoolVar(&vars.ShouldCheckUninsAndReins, "ckuni", false, "determine whether check again after uninstall and reinstall")

	flag.Float64Var(&vars.LoginRate, "lgr", 0, "number of login user rate")
}

// SetFlagLimit prevent parameters from exceeding the limit
func SetFlagLimit() {
	vars.UserNum = min(vars.UserNum, config.MaxUserNum)
	vars.CommonGroupMemberNum = min(vars.CommonGroupMemberNum, vars.UserNum)

	vars.LoginRate = min(vars.LoginRate, 1)
}
