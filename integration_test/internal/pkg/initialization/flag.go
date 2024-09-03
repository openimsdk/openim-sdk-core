package initialization

import (
	"context"
	"flag"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/internal/flagconst"
	"github.com/openimsdk/tools/log"
)

func InitFlag() {
	flag.BoolVar(&flagconst.TestMode, vars.FlagMap["TestMode"], false, "mark is test mode")

	flag.IntVar(&vars.UserNum, vars.FlagMap["UserNum"], 100, "user num")
	flag.IntVar(&vars.SuperUserNum, vars.FlagMap["SuperUserNum"], 10, "number of users with all friends")
	flag.IntVar(&vars.LargeGroupNum, vars.FlagMap["LargeGroupNum"], 5, "number of big group")
	flag.IntVar(&vars.LargeGroupMemberNum, vars.FlagMap["LargeGroupMemberNum"], 100, "number of big group member")
	flag.IntVar(&vars.CommonGroupNum, vars.FlagMap["CommonGroupNum"], 10, "number of small group")
	flag.IntVar(&vars.CommonGroupMemberNum, vars.FlagMap["CommonGroupMemberNum"], 20, "number of small group member")
	flag.IntVar(&vars.SingleMessageNum, vars.FlagMap["SingleMessageNum"], 5, "number of single message each user send")
	flag.IntVar(&vars.GroupMessageNum, vars.FlagMap["GroupMessageNum"], 1, "number of group message each user send")

	flag.BoolVar(&vars.ShouldRegister, vars.FlagMap["ShouldRegister"], false, "determine whether register")
	flag.BoolVar(&vars.ShouldImportFriends, vars.FlagMap["ShouldImportFriends"], false, "determine whether import friends")
	flag.BoolVar(&vars.ShouldCreateGroup, vars.FlagMap["ShouldCreateGroup"], false, "determine whether create group")
	flag.BoolVar(&vars.ShouldSendMsg, vars.FlagMap["ShouldSendMsg"], false, "determine whether send messages")

	flag.BoolVar(&vars.ShouldCheckGroupNum, vars.FlagMap["ShouldCheckGroupNum"], false, "determine whether check group num")
	flag.BoolVar(&vars.ShouldCheckConversationNum, vars.FlagMap["ShouldCheckConversationNum"], false, "determine whether check conversation num")
	flag.BoolVar(&vars.ShouldCheckMessageNum, vars.FlagMap["ShouldCheckMessageNum"], false, "determine whether check message num")
	flag.BoolVar(&vars.ShouldCheckUninsAndReins, vars.FlagMap["ShouldCheckUninsAndReins"], false, "determine whether check again after uninstall and reinstall")

	flag.Float64Var(&vars.LoginRate, vars.FlagMap["LoginRate"], 0, "number of login user rate")

}

// SetFlagLimit prevent parameters from exceeding the limit
func SetFlagLimit() {
	vars.UserNum = min(vars.UserNum, config.MaxUserNum)
	vars.CommonGroupMemberNum = min(vars.CommonGroupMemberNum, vars.UserNum)

	vars.LoginRate = min(vars.LoginRate, 1)

	if !isSet(vars.FlagMap["LargeGroupMemberNum"]) {
		vars.LargeGroupMemberNum = vars.UserNum
	}

	if isSet(vars.FlagMap["LoginRate"]) {
		vars.IsLogin = true
	}
}

func PrintFlag() {
	result := fmt.Sprintf(
		"TestMode-%s:%t, UserNum-%s:%d, SuperUserNum-%s:%d, LargeGroupNum-%s:%d, LargeGroupMemberNum-%s:%d, CommonGroupNum-%s:%d, CommonGroupMemberNum-%s:%d, SingleMessageNum-%s:%d, GroupMessageNum-%s:%d, ShouldRegister-%s:%t, ShouldImportFriends-%s:%t, ShouldCreateGroup-%s:%t, ShouldSendMsg-%s:%t, ShouldCheckGroupNum-%s:%t, ShouldCheckConversationNum-%s:%t, ShouldCheckMessageNum-%s:%t, ShouldCheckUninsAndReins-%s:%t, LoginRate-%s:%.2f",
		vars.FlagMap["TestMode"], flagconst.TestMode,
		vars.FlagMap["UserNum"], vars.UserNum,
		vars.FlagMap["SuperUserNum"], vars.SuperUserNum,
		vars.FlagMap["LargeGroupNum"], vars.LargeGroupNum,
		vars.FlagMap["LargeGroupMemberNum"], vars.LargeGroupMemberNum,
		vars.FlagMap["CommonGroupNum"], vars.CommonGroupNum,
		vars.FlagMap["CommonGroupMemberNum"], vars.CommonGroupMemberNum,
		vars.FlagMap["SingleMessageNum"], vars.SingleMessageNum,
		vars.FlagMap["GroupMessageNum"], vars.GroupMessageNum,
		vars.FlagMap["ShouldRegister"], vars.ShouldRegister,
		vars.FlagMap["ShouldImportFriends"], vars.ShouldImportFriends,
		vars.FlagMap["ShouldCreateGroup"], vars.ShouldCreateGroup,
		vars.FlagMap["ShouldSendMsg"], vars.ShouldSendMsg,
		vars.FlagMap["ShouldCheckGroupNum"], vars.ShouldCheckGroupNum,
		vars.FlagMap["ShouldCheckConversationNum"], vars.ShouldCheckConversationNum,
		vars.FlagMap["ShouldCheckMessageNum"], vars.ShouldCheckMessageNum,
		vars.FlagMap["ShouldCheckUninsAndReins"], vars.ShouldCheckUninsAndReins,
		vars.FlagMap["LoginRate"], vars.LoginRate,
	)
	fmt.Println(result)
	log.ZWarn(context.TODO(), "flags", nil, "flag params", result)
}

func isSet(fg string) bool {
	set := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == fg {
			set = true
		}
	})
	return set
}
