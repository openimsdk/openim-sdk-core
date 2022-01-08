package init

import (
	"github.com/golang/protobuf/proto"
	ws "open_im_sdk/internal/controller/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"sync"
)

type MsgSync struct {
	*db.DataBase
	seqMsg      map[int32]server_api_params.MsgData
	seqMsgMutex sync.RWMutex
	*ws.Ws
	loginUserID string
}

func NewMsgSync(dataBase *db.DataBase, ws *ws.Ws, loginUserID string, ch chan common.Cmd2Value) *MsgSync {
	return &MsgSync{DataBase: dataBase, seqMsg: make(map[int32]server_api_params.MsgData, 1000), Ws: ws, loginUserID: loginUserID}
}

func (u *MsgSync) getNeedSyncSeq(svrMinSeq, svrMaxSeq int32) []int32 {
	localMinSeq := u.GetNeedSyncLocalMinSeq()
	var startSeq int32
	if localMinSeq > svrMinSeq {
		startSeq = localMinSeq
	} else {
		startSeq = svrMinSeq
	}

	seqList := make([]int32, 0)

	var maxConsequentSeq int32
	isBreakFlag := false
	normalSeq := u.getNormalChatLogSeq(startSeq)
	errorSeq := u.getErrorChatLogSeq(startSeq)
	for seq := startSeq; seq <= svrMaxSeq; seq++ {
		_, ok1 := normalSeq[seq]
		_, ok2 := errorSeq[seq]
		if ok1 || ok2 {
			if !isBreakFlag {
				maxConsequentSeq = seq
			}
			continue
		} else {
			isBreakFlag = true
			if seq != 0 {
				seqList = append(seqList, seq)
			}
		}
	}

	var firstSeq int32
	if len(seqList) > 0 {
		firstSeq = seqList[0]
	} else {
		if maxConsequentSeq > startSeq {
			firstSeq = maxConsequentSeq
		} else {
			firstSeq = startSeq
		}
	}
	if firstSeq > localMinSeq {
		u.SetNeedSyncLocalMinSeq(firstSeq)
	}

	return seqList
}

func (u *MsgSync) syncMsgFromServer(needSyncSeqList []int32) (err error) {
	notInCache := u.getNotInSeq(needSyncSeqList)
	if len(notInCache) == 0 {
		log.Info("notInCache is null, don't sync from svr")
		return nil
	}
	log.Info("notInCache ", notInCache)
	var SPLIT int = 100
	for i := 0; i < len(notInCache)/SPLIT; i++ {
		//0-99 100-199
		u.syncMsgFromServerSplit(notInCache[i*SPLIT : (i+1)*SPLIT])
		log.Info("syncMsgFromServerSplit idx: ", i*SPLIT, (i+1)*SPLIT)
	}
	u.syncMsgFromServerSplit(notInCache[SPLIT*(len(notInCache)/SPLIT):])
	log.Info("syncMsgFromServerSplit idx: ", SPLIT*(len(notInCache)/SPLIT), len(notInCache))
	return nil
}

func (u *MsgSync) getNotInSeq(needSyncSeqList []int32) (seqList []int64) {
	u.seqMsgMutex.RLock()
	defer u.seqMsgMutex.RUnlock()

	for _, v := range needSyncSeqList {
		_, ok := u.seqMsg[v]
		if !ok {
			seqList = append(seqList, int64(v))
		}
	}
	return seqList
}

type PullUserMsgResp struct {
	ErrCode       int                       `json:"errCode"`
	ErrMsg        string                    `json:"errMsg"`
	ReqIdentifier int                       `json:"reqIdentifier"`
	MsgIncr       int                       `json:"msgIncr"`
	Data          paramsPullUserMsgDataResp `json:"data"`
}
type paramsPullUserMsgDataResp struct {
	Group  []*server_api_params.GatherFormat `json:"group"`
	MaxSeq int64                             `json:"maxSeq"`
	MinSeq int64                             `json:"minSeq"`
	Single []*server_api_params.GatherFormat `json:"single"`
}

type ArrMsg struct {
	SingleData []server_api_params.MsgData
	GroupData  []server_api_params.MsgData
}

func (u *MsgSync) syncMsgFromServerSplit(needSyncSeqList []int64) (err error) {
	if len(needSyncSeqList) == 0 {
		log.Info("len(needSyncSeqList) == 0  don't pull from svr")
		return nil
	}

	var pullMsgReq server_api_params.PullMessageBySeqListReq
	pullMsgReq.SeqList = needSyncSeqList
	buff, err := proto.Marshal(&pullMsgReq)
	resp, err, operationID := u.SendReqWaitResp(buff, constant.WSPullMsgBySeqList, 30, u.loginUserID)
	if err != nil {
		log.Error(operationID, "SendReqWaitResp failed ", err.Error())
		return err
	}

	var pullMsg PullUserMsgResp
	var pullMsgResp server_api_params.PullMessageBySeqListResp
	err := proto.Unmarshal(resp.Data, &pullMsgResp)
	if err != nil {
		log.Error(operationID, "Unmarshal failed ", err.Error())
		return err
	}
	pullMsg.Data.Group = pullMsgResp.GroupUserMsg
	pullMsg.Data.Single = pullMsgResp.SingleUserMsg
	pullMsg.Data.MaxSeq = pullMsgResp.MaxSeq
	pullMsg.Data.MinSeq = pullMsgResp.MinSeq

	u.seqMsgMutex.Lock()
	isInmap := false
	arrMsg := ArrMsg{}
	//	sdkLog("pullmsg data: ", pullMsgResp.SingleUserMsg, pullMsg.Data.Single)
	for i := 0; i < len(pullMsg.Data.Single); i++ {
		for j := 0; j < len(pullMsg.Data.Single[i].List); j++ {
			log.Info(operationID, "open_im pull one msg: |", pullMsg.Data.Single[i].List[j].ClientMsgID, "|")
			log.Info(operationID, "pull all: |", pullMsg.Data.Single[i].List[j].Seq, pullMsg.Data.Single[i].List[j])

			singleMsg := server_api_params.MsgData{
				SendID:           pullMsg.Data.Single[i].List[j].SendID,
				RecvID:           pullMsg.Data.Single[i].List[j].RecvID,
				SessionType:      constant.SingleChatType,
				MsgFrom:          pullMsg.Data.Single[i].List[j].MsgFrom,
				ContentType:      pullMsg.Data.Single[i].List[j].ContentType,
				ServerMsgID:      pullMsg.Data.Single[i].List[j].ServerMsgID,
				Content:          pullMsg.Data.Single[i].List[j].Content,
				SendTime:         pullMsg.Data.Single[i].List[j].SendTime,
				Seq:              pullMsg.Data.Single[i].List[j].Seq,
				SenderNickname:   pullMsg.Data.Single[i].List[j].SenderNickname,
				SenderFaceURL:    pullMsg.Data.Single[i].List[j].SenderFaceURL,
				ClientMsgID:      pullMsg.Data.Single[i].List[j].ClientMsgID,
				SenderPlatformID: pullMsg.Data.Single[i].List[j].SenderPlatformID,
			}

			b1 := u.isExistsInErrChatLogBySeq(pullMsg.Data.Single[i].List[j].Seq)
			b2 := u.judgeMessageIfExistsBySeq(pullMsg.Data.Single[i].List[j].Seq)
			_, ok := u.seqMsg[int32(pullMsg.Data.Single[i].List[j].Seq)]
			if b1 || b2 || ok {
				log.Info(operationID, "seq in : ", pullMsg.Data.Single[i].List[j].Seq, b1, b2, ok)
			} else {
				isInmap = true
				u.seqMsg[int32(pullMsg.Data.Single[i].List[j].Seq)] = singleMsg
				log.Info(operationID, "into map, seq: ", pullMsg.Data.Single[i].List[j].Seq, pullMsg.Data.Single[i].List[j].ClientMsgID, pullMsg.Data.Single[i].List[j].ServerMsgID, pullMsg.Data.Single[i].List[j])
			}
		}
	}

	for i := 0; i < len(pullMsg.Data.Group); i++ {
		for j := 0; j < len(pullMsg.Data.Group[i].List); j++ {
			groupMsg := server_api_params.MsgData{
				SendID:           pullMsg.Data.Group[i].List[j].SendID,
				RecvID:           pullMsg.Data.Group[i].List[j].RecvID,
				SessionType:      constant.GroupChatType,
				MsgFrom:          pullMsg.Data.Group[i].List[j].MsgFrom,
				ContentType:      pullMsg.Data.Group[i].List[j].ContentType,
				ServerMsgID:      pullMsg.Data.Group[i].List[j].ServerMsgID,
				Content:          pullMsg.Data.Group[i].List[j].Content,
				SendTime:         pullMsg.Data.Group[i].List[j].SendTime,
				Seq:              pullMsg.Data.Group[i].List[j].Seq,
				SenderNickname:   pullMsg.Data.Group[i].List[j].SenderNickname,
				SenderFaceURL:    pullMsg.Data.Group[i].List[j].SenderFaceURL,
				ClientMsgID:      pullMsg.Data.Group[i].List[j].ClientMsgID,
				SenderPlatformID: pullMsg.Data.Group[i].List[j].SenderPlatformID,
			}

			b1 := u.isExistsInErrChatLogBySeq(pullMsg.Data.Group[i].List[j].Seq)
			b2 := u.judgeMessageIfExistsBySeq(pullMsg.Data.Group[i].List[j].Seq)
			_, ok := u.seqMsg[int32(pullMsg.Data.Group[i].List[j].Seq)]
			if b1 || b2 || ok {
				log.Info(operationID, "seq in : ", pullMsg.Data.Group[i].List[j].Seq, b1, b2, ok)
			} else {
				isInmap = true
				u.seqMsg[int32(pullMsg.Data.Group[i].List[j].Seq)] = groupMsg
				log.Info(operationID, "into map, seq: ", pullMsg.Data.Group[i].List[j].Seq, pullMsg.Data.Group[i].List[j].ClientMsgID, pullMsg.Data.Group[i].List[j].ServerMsgID)
				log.Info(operationID, "pull all: |", pullMsg.Data.Group[i].List[j].Seq, pullMsg.Data.Group[i].List[j])

			}
		}
	}
	u.seqMsgMutex.Unlock()

	if isInmap {
		err = common.TriggerCmdNewMsgCome(arrMsg, )
		if err != nil {
		}
	}
	return nil
}
