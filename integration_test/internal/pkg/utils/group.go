package utils

import "github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"

func BuildGroupName(ownerID, type_ string) string {
	return vars.GroupNamePrefix + type_ + "_" + ownerID
}
