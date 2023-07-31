// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msg

import "errors"

func (x *GetMaxAndMinSeqReq) Check() error {
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	return nil
}

func (x *SendMsgReq) Check() error {
	if x.MsgData == nil {
		return errors.New("MsgData is empty")
	}
	if err := x.MsgData.Check(); err != nil {
		return err
	}
	return nil
}

func (x *SetSendMsgStatusReq) Check() error {
	if x.Status < 0 || x.Status > 3 {
		return errors.New("status is invalid")
	}
	return nil
}

func (x *GetSendMsgStatusReq) Check() error {
	return nil
}

//func (x *ModifyMessageReactionExtensionsReq) Check() error {
//	if x.ConversationID == "" {
//		return errs.ErrArgs.Wrap("conversationID is empty")
//	}
//	if x.SessionType < 1 || x.SessionType > 4 {
//		return errs.ErrArgs.Wrap("sessionType is invalid")
//	}
//	if x.ReactionExtensions == nil {
//		return errs.ErrArgs.Wrap("reactionExtensions is empty")
//	}
//	return nil
//}
//
//func (x *SetMessageReactionExtensionsReq) Check() error {
//	if x.ConversationID == "" {
//		return errs.ErrArgs.Wrap("conversationID is empty")
//	}
//	if x.SessionType < 1 || x.SessionType > 4 {
//		return errs.ErrArgs.Wrap("sessionType is invalid")
//	}
//	if x.ReactionExtensions == nil {
//		return errs.ErrArgs.Wrap("reactionExtensions is empty")
//	}
//	return nil
//}
//
//func (x *GetMessagesReactionExtensionsReq) Check() error {
//	if x.ConversationID == "" {
//		return errs.ErrArgs.Wrap("conversationID is empty")
//	}
//	if x.SessionType < 1 || x.SessionType > 4 {
//		return errs.ErrArgs.Wrap("sessionType is invalid")
//	}
//	if x.MessageReactionKeys == nil {
//		return errs.ErrArgs.Wrap("MessageReactionKeys is empty")
//	}
//	if x.TypeKeys == nil {
//		return errs.ErrArgs.Wrap("TypeKeys is empty")
//	}
//	return nil
//}
//
//func (x *DeleteMessagesReactionExtensionsReq) Check() error {
//	if x.ConversationID == "" {
//		return errs.ErrArgs.Wrap("conversationID is empty")
//	}
//	if x.SessionType < 1 || x.SessionType > 4 {
//		return errs.ErrArgs.Wrap("sessionType is invalid")
//	}
//	if x.ReactionExtensions == nil {
//		return errs.ErrArgs.Wrap("ReactionExtensions is empty")
//	}
//	return nil
//}

func (x *DelMsgsReq) Check() error {
	return nil
}

func (x *RevokeMsgReq) Check() error {
	if x.ConversationID == "" {
		return errors.New("conversationID is empty")
	}
	if x.Seq < 1 {
		return errors.New("seq is invalid")
	}
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	return nil
}

func (x *MarkMsgsAsReadReq) Check() error {
	if x.ConversationID == "" {
		return errors.New("conversationID is empty")
	}
	if x.Seqs == nil {
		return errors.New("seqs is empty")
	}
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	for _, seq := range x.Seqs {
		if seq == 0 {
			return errors.New("seqs has 0 value is invalid")
		}
	}
	return nil
}

func (x *MarkConversationAsReadReq) Check() error {
	if x.ConversationID == "" {
		return errors.New("conversationID is empty")
	}
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	if x.HasReadSeq < 1 {
		return errors.New("hasReadSeq is invalid")
	}
	for _, seq := range x.Seqs {
		if seq == 0 {
			return errors.New("seqs has 0 value is invalid")
		}
	}
	return nil
}

func (x *SetConversationHasReadSeqReq) Check() error {
	if x.ConversationID == "" {
		return errors.New("conversationID is empty")
	}
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	if x.HasReadSeq < 1 {
		return errors.New("hasReadSeq is invalid")
	}
	return nil
}

func (x *ClearConversationsMsgReq) Check() error {
	if x.ConversationIDs == nil {
		return errors.New("conversationIDs is empty")
	}
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	return nil
}

func (x *UserClearAllMsgReq) Check() error {
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	return nil
}

func (x *DeleteMsgsReq) Check() error {
	if x.ConversationID == "" {
		return errors.New("conversationID is empty")
	}
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	if x.Seqs == nil {
		return errors.New("seqs is empty")
	}
	return nil
}

func (x *DeleteMsgPhysicalReq) Check() error {
	if x.ConversationIDs == nil {
		return errors.New("conversationIDs is empty")
	}
	return nil
}

func (x *GetConversationMaxSeqReq) Check() error {
	if x.ConversationID == "" {
		return errors.New("conversationID is empty")
	}
	return nil
}

func (x *GetConversationsHasReadAndMaxSeqReq) Check() error {
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	return nil
}
