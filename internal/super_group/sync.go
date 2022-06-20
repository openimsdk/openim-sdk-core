package super_group

import (
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
)

//func (s *SuperGroup) SyncJoinedGroupList(operationID string) {
//	groupList, err := s.getJoinedGroupListFromSvr(operationID)
//	log.Debug(operationID, "getJoinedGroupListFromSvr result ", groupList)
//	localGroupList := common.TransferToLocalGroupInfo(groupList)
//	if err == nil {
//		s.db.DeleteAllSuperGroup()
//		for _, v := range localGroupList {
//			err = s.db.InsertSuperGroup(v)
//			if err != nil {
//				log.Error(operationID, "InsertSuperGroup  failed ", err.Error(), v)
//			} else {
//				log.Debug(operationID, "InsertSuperGroup  ok  ", v)
//			}
//		}
//	}
//}

func (s *SuperGroup) SyncJoinedGroupList(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := s.getJoinedGroupListFromSvr(operationID)
	log.Info(operationID, "getJoinedGroupListFromSvr", svrList, s.loginUserID)
	if err != nil {
		log.NewError(operationID, "getJoinedGroupListFromSvr failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalGroupInfo(svrList)
	onLocal, err := s.db.GetJoinedSuperGroupList()
	if err != nil {
		log.NewError(operationID, "GetJoinedSuperGroupList failed ", err.Error())
		return
	}

	log.NewInfo(operationID, " onLocal ", onLocal, s.loginUserID)
	aInBNot, bInANot, sameA, sameB := common.CheckGroupInfoDiff(onServer, onLocal)
	log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB, s.loginUserID)
	for _, index := range aInBNot {
		log.Info(operationID, "InsertSuperGroup ", *onServer[index])
		err := s.db.InsertSuperGroup(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertGroup failed ", err.Error(), *onServer[index])
			continue
		}
	}

	for _, index := range sameA {
		log.Info(operationID, "UpdateSuperGroup ", *onServer[index])
		err := s.db.UpdateSuperGroup(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateGroup failed ", err.Error(), onServer[index])
			continue
		}
	}

	for _, index := range bInANot {
		log.Info(operationID, "DeleteSuperGroup: ", onLocal[index].GroupID, s.loginUserID)
		err := s.db.DeleteSuperGroup(onLocal[index].GroupID)
		if err != nil {
			log.NewError(operationID, "DeleteSuperGroup failed ", err.Error(), onLocal[index].GroupID)
			continue
		}
	}
}
