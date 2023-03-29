package organization

import (
	"github.com/jinzhu/copier"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"sync"
)

type Organization struct {
	listener    open_im_sdk_callback.OnOrganizationListener
	loginUserID string
	db          db_interface.DataBase
	p           *ws.PostApi
	loginTime   int64

	//	memberSyncMutex sync.RWMutex
}

func (o *Organization) LoginTime() int64 {
	return o.loginTime
}

func (o *Organization) SetLoginTime(loginTime int64) {
	o.loginTime = loginTime
}

func NewOrganization(loginUserID string, db db_interface.DataBase, p *ws.PostApi) *Organization {
	return &Organization{loginUserID: loginUserID, db: db, p: p}
}

func (o *Organization) DoNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value, operationID string) {
	if o.listener == nil {
		return
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)

	if msg.SendTime < o.loginTime || o.loginTime == 0 {
		log.Warn(operationID, "ignore notification ", msg.ClientMsgID, msg.ServerMsgID, msg.Seq, msg.ContentType)
		return
	}

	go func() {
		switch msg.ContentType {
		case constant.OrganizationChangedNotification:
			o.organizationChangedNotification(msg, operationID)
		default:
			log.Error(operationID, "ContentType tip failed ", msg.ContentType)
		}
	}()
}

func (o *Organization) getSubDepartment(callback open_im_sdk_callback.Base, departmentID string, offset, count int, operationID string) sdk_params_callback.GetSubDepartmentCallback {
	subDepartmentList, err := o.db.GetSubDepartmentList(departmentID, offset, count)
	common.CheckDBErrCallback(callback, err, operationID)
	return subDepartmentList
}

func (o *Organization) getDepartmentMember(callback open_im_sdk_callback.Base, departmentID string, offset, count int, operationID string) sdk_params_callback.GetDepartmentMemberCallback {
	departmentMemberLis, err := o.db.GetDepartmentMemberListByDepartmentID(departmentID, offset, count)
	common.CheckDBErrCallback(callback, err, operationID)
	return departmentMemberLis
}

func (o *Organization) getUserInDepartment(callback open_im_sdk_callback.Base, userID string, operationID string) sdk_params_callback.GetUserInDepartmentCallback {
	departmentMemberList, err := o.db.GetDepartmentMemberListByUserID(userID)
	common.CheckDBErrCallback(callback, err, operationID)
	var userInDepartment []*sdk_params_callback.UserInDepartment
	for _, v := range departmentMemberList {
		department, err := o.db.GetDepartmentInfo(v.DepartmentID)
		if err != nil {
			continue
		}
		node := sdk_params_callback.UserInDepartment{MemberInfo: &model_struct.LocalDepartmentMember{}}
		node.DepartmentInfo = department
		copier.Copy(node.MemberInfo, v)
		userInDepartment = append(userInDepartment, &node)
	}
	return userInDepartment
}

func (o *Organization) getDepartmentMemberAndSubDepartment(callback open_im_sdk_callback.Base, departmentID string, operationID string) sdk_params_callback.GetDepartmentMemberAndSubDepartmentCallback {
	subDepartmentList, err := o.db.GetSubDepartmentList(departmentID)
	common.CheckDBErrCallback(callback, err, operationID)
	departmentMemberList, err := o.db.GetDepartmentMemberListByDepartmentID(departmentID)
	common.CheckDBErrCallback(callback, err, operationID)
	parentDepartmentList, err := o.db.GetParentDepartmentList(departmentID)
	common.CheckDBErrCallback(callback, err, operationID)
	var parentDepartmentCallbackList []sdk_params_callback.ParentDepartmentCallback
	for _, v := range parentDepartmentList {
		parentDepartmentCallbackList = append(parentDepartmentCallbackList, sdk_params_callback.ParentDepartmentCallback{
			Name:         v.Name,
			DepartmentID: v.DepartmentID,
		})
	}
	return sdk_params_callback.GetDepartmentMemberAndSubDepartmentCallback{DepartmentList: subDepartmentList, DepartmentMemberList: departmentMemberList, ParentDepartmentList: parentDepartmentCallbackList}
}

func (o *Organization) getParentDepartmentList(callback open_im_sdk_callback.Base, departmentID string, operationID string) sdk_params_callback.GetParentDepartmentListCallback {
	parentDepartmentList, err := o.db.GetParentDepartmentList(departmentID)
	common.CheckDBErrCallback(callback, err, operationID)
	return parentDepartmentList
}

func (o *Organization) getDepartmentInfo(callback open_im_sdk_callback.Base, departmentID string, operationID string) sdk_params_callback.GetDepartmentInfoCallback {
	departmentInfo, err := o.db.GetDepartmentInfo(departmentID)
	common.CheckDBErrCallback(callback, err, operationID)
	return departmentInfo
}

func (o *Organization) searchOrganization(callback open_im_sdk_callback.Base, searchParam sdk_params_callback.SearchOrganizationParams, offset, count int, operationID string) sdk_params_callback.SearchOrganizationCallback {
	departmentMemberList, err := o.db.SearchDepartmentMember(searchParam.KeyWord,
		searchParam.IsSearchUserName, searchParam.IsSearchEmail, searchParam.IsSearchMobile,
		searchParam.IsSearchPosition, searchParam.IsSearchTelephone, searchParam.IsSearchUserEnglishName, searchParam.IsSearchUserID,
		offset, count)
	common.CheckDBErrCallback(callback, err, operationID)
	departmentList, err := o.db.SearchDepartment(searchParam.KeyWord, offset, count)
	common.CheckDBErrCallback(callback, err, operationID)
	result := sdk_params_callback.SearchOrganizationCallback{
		DepartmentList: departmentList,
	}
	for _, member := range departmentMemberList {
		parentDepartmentList, err := o.db.GetParentDepartmentList(member.DepartmentID)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "GetParentDepartmentList failed", err.Error())
		}
		var path []*sdk_params_callback.ParentDepartmentCallback
		for _, department := range parentDepartmentList {
			path = append(path, &sdk_params_callback.ParentDepartmentCallback{
				Name:         department.Name,
				DepartmentID: department.DepartmentID,
			})
		}
		result.DepartmentMemberList = append(result.DepartmentMemberList, &struct {
			*model_struct.SearchDepartmentMemberResult
			ParentDepartmentList []*sdk_params_callback.ParentDepartmentCallback `json:"parentDepartmentList"`
		}{SearchDepartmentMemberResult: member, ParentDepartmentList: path})
	}
	return result
}

func (o *Organization) getSubDepartmentFromSvr(departmentID string, operationID string) ([]*api.Department, error) {
	var apiReq api.GetSubDepartmentReq
	apiReq.OperationID = operationID
	apiReq.DepartmentID = departmentID
	var realData []*api.Department
	err := o.p.PostReturn(constant.GetSubDepartmentRouter, apiReq, &realData)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return realData, nil
}

func (o *Organization) getAllDepartmentFromSvr(operationID string) ([]*api.Department, error) {
	return o.getSubDepartmentFromSvr("-1", operationID)
}

func (o *Organization) getAllDepartmentMemberFromSvr(operationID string) ([]*api.UserDepartmentMember, error) {
	return o.getDepartmentMemberFromSvr("-1", operationID)
}

func (o *Organization) getDepartmentMemberFromSvr(departmentID string, operationID string) ([]*api.UserDepartmentMember, error) {
	var apiReq api.GetDepartmentMemberReq
	apiReq.OperationID = operationID
	apiReq.DepartmentID = departmentID
	var realData []*api.UserDepartmentMember
	err := o.p.PostReturn(constant.GetDepartmentMemberRouter, apiReq, &realData)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return realData, nil
}

func (o *Organization) SyncDepartment(operationID string) {
	//	return
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := o.getAllDepartmentFromSvr(operationID)
	if err != nil {
		log.NewError(operationID, "getAllDepartmentFromSvr failed ", err.Error())
		return
	}
	log.Info(operationID, "getAllDepartmentFromSvr ", svrList)
	onServer := common.TransferToLocalDepartment(svrList)
	onLocal, err := o.db.GetAllDepartmentList()
	if err != nil {
		log.NewError(operationID, "GetAllDepartmentList failed ", err.Error())
		return
	}
	flag := 0
	log.NewInfo(operationID, "svrList onServer onLocal", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckDepartmentDiff(onServer, onLocal)
	log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := o.db.InsertDepartment(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertDepartment failed ", err.Error(), *onServer[index])
			continue
		}
		log.Info(operationID, "InsertDepartment", onServer[index])
		flag = 1
	}
	for _, index := range sameA {
		err := o.db.UpdateDepartment(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateDepartment failed ", err.Error())
			continue
		}
		log.Info(operationID, "UpdateDepartment ", onServer[index])
		flag = 1
	}
	for _, index := range bInANot {
		err := o.db.DeleteDepartment(onLocal[index].DepartmentID)
		if err != nil {
			log.NewError(operationID, "DeleteDepartment failed ", err.Error())
			continue
		}
		log.Info(operationID, "DeleteDepartment", onLocal[index].DepartmentID)
		flag = 1
	}
	if flag != 0 {
		if o.listener == nil {
			log.Error(operationID, "o.listener == nil ")
			return
		}
		o.listener.OnOrganizationUpdated()
	}
}

func (o *Organization) SyncDepartmentMemberByDepartmentID(operationID string, departmentID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", departmentID)
	svrList, err := o.getDepartmentMemberFromSvr(departmentID, operationID)
	if err != nil {
		log.NewError(operationID, "getDepartmentMemberFromSvr failed ", err.Error())
		return
	}
	log.Info(operationID, "getDepartmentMemberFromSvr result ", svrList, departmentID)
	onServer := common.TransferToLocalDepartmentMember(svrList)
	onLocal, err := o.db.GetDepartmentMemberListByDepartmentID(departmentID, 0, 1000000)
	if err != nil {
		log.NewError(operationID, "GetDepartmentMemberListByDepartmentID failed ", err.Error())
		return
	}
	flag := 0
	log.NewInfo(operationID, "svrList onServer onLocal", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckDepartmentMemberDiff(onServer, onLocal)
	log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := o.db.InsertDepartmentMember(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertDepartmentMember failed ", err.Error(), *onServer[index])
			continue
		}
		log.Info(operationID, "InsertDepartmentMember", onServer[index])
		flag = 1
	}
	for _, index := range sameA {
		err := o.db.UpdateDepartmentMember(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateDepartmentMember failed ", err.Error())
			continue
		}
		log.Info(operationID, "UpdateDepartmentMember ", onServer[index])
		flag = 1
	}
	for _, index := range bInANot {
		err := o.db.DeleteDepartmentMember(onLocal[index].DepartmentID, onLocal[index].UserID)
		if err != nil {
			log.NewError(operationID, "DeleteDepartmentMember failed ", err.Error())
			continue
		}
		log.Info(operationID, "DeleteDepartmentMember", onLocal[index].DepartmentID, onLocal[index].UserID)
		flag = 1
	}
	if flag != 0 {
		if o.listener == nil {
			log.Error(operationID, "o.listener == nil ")
			return
		}
		o.listener.OnOrganizationUpdated()
	}
}

func (o *Organization) organizationChangedNotification(msg *api.MsgData, operationID string) {
	o.SyncDepartment(operationID)
	o.SyncAllDepartmentMember(operationID)
	//	o.SyncDepartmentMember(operationID)

}
func (o *Organization) SyncAllDepartmentMember(operationID string) {
	//	return
	//	o.memberSyncMutex.Lock()
	//defer o.memberSyncMutex.Unlock()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := o.getAllDepartmentMemberFromSvr(operationID)
	if err != nil {
		log.NewError(operationID, "getAllDepartmentMemberFromSvr failed ", err.Error())
		return
	}
	log.Info(operationID, "getDepartmentMemberFromSvr result ", svrList)
	onServer := common.TransferToLocalDepartmentMember(svrList)
	onLocal, err := o.db.GetAllDepartmentMemberList()
	if err != nil {
		log.NewError(operationID, "GetAllDepartmentMemberList failed ", err.Error())
		return
	}
	flag := 0
	log.NewInfo(operationID, "svrList onServer onLocal", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckDepartmentMemberDiff(onServer, onLocal)
	log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB)

	var insertGroupMemberList []*model_struct.LocalDepartmentMember

	for _, index := range aInBNot {
		insertGroupMemberList = append(insertGroupMemberList, onServer[index])
		flag = 1
	}

	if len(insertGroupMemberList) > 0 {
		split := 1000
		idx := 0
		remain := len(insertGroupMemberList) % split
		log.Info(operationID, "BatchInsertGroupMember all: ", len(insertGroupMemberList))
		for idx = 0; idx < len(insertGroupMemberList)/split; idx++ {
			sub := insertGroupMemberList[idx*split : (idx+1)*split]
			err = o.db.BatchInsertDepartmentMember(sub)
			log.Info(operationID, "BatchInsertDepartmentMember len: ", len(sub))
			if err != nil {
				log.Error(operationID, "BatchInsertDepartmentMember failed ", err.Error(), len(sub))
				for again := 0; again < len(sub); again++ {
					if err = o.db.InsertDepartmentMember(sub[again]); err != nil {
						log.Error(operationID, "InsertDepartmentMember failed ", err.Error(), sub[again])
					}
				}
			}
		}
		if remain > 0 {
			sub := insertGroupMemberList[idx*split:]
			log.Info(operationID, "BatchInsertDepartmentMember len: ", len(sub))
			err = o.db.BatchInsertDepartmentMember(sub)
			if err != nil {
				log.Error(operationID, "BatchInsertDepartmentMember failed ", err.Error(), len(sub))
				for again := 0; again < len(sub); again++ {
					if err = o.db.InsertDepartmentMember(sub[again]); err != nil {
						log.Error(operationID, "InsertDepartmentMember failed ", err.Error(), sub[again])
					}
				}
			}
		}
	}

	for _, index := range sameA {
		err := o.db.UpdateDepartmentMember(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateDepartmentMember failed ", err.Error())
			continue
		}
		log.Info(operationID, "UpdateDepartmentMember ", onServer[index])
		flag = 1
	}
	for _, index := range bInANot {
		err := o.db.DeleteDepartmentMember(onLocal[index].DepartmentID, onLocal[index].UserID)
		if err != nil {
			log.NewError(operationID, "DeleteDepartmentMember failed ", err.Error(), onLocal[index].DepartmentID, onLocal[index].UserID)
			continue
		}
		log.Info(operationID, "DeleteDepartmentMember", onLocal[index].DepartmentID, onLocal[index].UserID)
		flag = 1
	}
	if flag != 0 {
		if o.listener == nil {
			log.Error(operationID, "o.listener == nil ")
			return
		}
		o.listener.OnOrganizationUpdated()
	}
}

func (o *Organization) SyncDepartmentMember(operationID string) {
	departmentList, err := o.db.GetAllDepartmentList()
	if err != nil {
		log.Error(operationID, "GetAllDepartmentList failed ", err.Error())
	}
	if len(departmentList) == 0 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(departmentList))
	log.Info(operationID, "SyncDepartmentMemberByDepartmentID begin", len(departmentList))
	for _, v := range departmentList {
		go func(departmentID, operationID string) {
			o.SyncDepartmentMemberByDepartmentID(operationID, departmentID)
			wg.Done()
		}(v.DepartmentID, operationID)
	}
	wg.Wait()
	log.Info(operationID, "SyncDepartmentMemberByDepartmentID end")
}

func (o *Organization) SyncOrganization(operationID string) {
	o.SyncDepartment(operationID)
	o.SyncAllDepartmentMember(operationID)
}
