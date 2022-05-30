package super_group

import "open_im_sdk/pkg/common"

func (s *SuperGroup) SyncJoinedGroupList(operationID string) {
	groupList, err := s.getJoinedGroupListFromSvr(operationID)
	localGroupList := common.TransferToLocalGroupInfo(groupList)
	if err != nil {
		s.db.DeleteAllSuperGroup()
		for _, v := range localGroupList {
			s.db.InsertSuperGroup(v)
		}
	}
}
