package rtc

import "errors"

func (x *GetSignalInvitationRecordsReq) Check() error {
	if x.Pagination == nil {
		return errors.New("pagination is empty")
	}
	return nil
}

func (x *DeleteSignalRecordsReq) Check() error {
	if len(x.SIDs) == 0 {
		return errors.New("sids is empty")
	}
	return nil
}

func (x *GetMeetingRecordsReq) Check() error {
	if x.Pagination == nil {
		return errors.New("pagination is empty")
	}
	return nil
}

func (x *DeleteMeetingRecordsReq) Check() error {
	if len(x.RoomIDs) == 0 {
		return errors.New("sids is empty")
	}
	return nil
}
