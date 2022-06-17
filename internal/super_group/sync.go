package super_group

import (
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
)

func (s *SuperGroup) SyncJoinedGroupList(operationID string) {
	groupList, err := s.getJoinedGroupListFromSvr(operationID)
	log.Debug(operationID, "getJoinedGroupListFromSvr result ", groupList)
	localGroupList := common.TransferToLocalGroupInfo(groupList)
	if err == nil {
		s.db.DeleteAllSuperGroup()
		for _, v := range localGroupList {
			err = s.db.InsertSuperGroup(v)
			if err != nil {
				log.Error(operationID, "InsertSuperGroup  failed ", err.Error(), v)
			} else {
				log.Debug(operationID, "InsertSuperGroup  ok  ", v)
			}
		}
	}
}
