package db

import (
	"errors"
	"gorm.io/gorm"
	"open_im_sdk/pkg/utils"
	"time"
)

type LocalWorkMomentsNotification struct {
	JsonDetail string `gorm:"column:json_detail"`
	CreateTime int64  `gorm:"create_time"`
}

func (LocalWorkMomentsNotification) TableName() string {
	return "local_work_moments_notification"
}

type LocalWorkMomentsNotificationUnreadCount struct {
	UnreadCount int `gorm:"unread_count" json:"unreadCount"`
}

func (LocalWorkMomentsNotificationUnreadCount) TableName() string {
	return "local_work_moments_notification_unread_count"
}

func (d *DataBase) InsertWorkMomentsNotification(jsonDetail string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	workMomentsNotification := LocalWorkMomentsNotification{
		JsonDetail: jsonDetail,
		CreateTime: time.Now().Unix(),
	}
	return utils.Wrap(d.conn.Create(workMomentsNotification).Error, "")
}

func (d *DataBase) GetWorkMomentsNotification(offset, count int) (WorkMomentsNotifications []*LocalWorkMomentsNotification, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	WorkMomentsNotifications = []*LocalWorkMomentsNotification{}
	err = utils.Wrap(d.conn.Table("local_work_moments_notification").Order("create_time DESC").Offset(offset).Limit(count).Find(&WorkMomentsNotifications).Error, "")
	return WorkMomentsNotifications, err
}

func (d *DataBase) GetWorkMomentsNotificationLimit(pageNumber, showNumber int) (WorkMomentsNotifications []*LocalWorkMomentsNotification, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	WorkMomentsNotifications = []*LocalWorkMomentsNotification{}
	err = utils.Wrap(d.conn.Table("local_work_moments_notification").Select("json_detail").Find(WorkMomentsNotifications).Error, "")
	return WorkMomentsNotifications, err
}

func (d *DataBase) InitWorkMomentsNotificationUnreadCount() error {
	var n int64
	err := utils.Wrap(d.conn.Model(&LocalWorkMomentsNotificationUnreadCount{}).Count(&n).Error, "")
	if err == nil {
		if n == 0 {
			c := LocalWorkMomentsNotificationUnreadCount{UnreadCount: 0}
			return utils.Wrap(d.conn.Model(&LocalWorkMomentsNotificationUnreadCount{}).Create(c).Error, "IncrConversationUnreadCount failed")
		}
	}
	return err
}

func (d *DataBase) IncrWorkMomentsNotificationUnreadCount() error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	c := LocalWorkMomentsNotificationUnreadCount{}
	t := d.conn.Model(&c).Where("1=1").Update("unread_count", gorm.Expr("unread_count+?", 1))
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "IncrConversationUnreadCount failed")
}

func (d *DataBase) MarkAllWorkMomentsNotificationAsRead() (err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Model(&LocalWorkMomentsNotificationUnreadCount{}).Where("1=1").Updates(map[string]interface{}{"unread_count": 0}).Error, "")
}

func (d *DataBase) GetWorkMomentsUnReadCount() (workMomentsNotificationUnReadCount LocalWorkMomentsNotificationUnreadCount, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	workMomentsNotificationUnReadCount = LocalWorkMomentsNotificationUnreadCount{}
	err = utils.Wrap(d.conn.Model(&LocalWorkMomentsNotificationUnreadCount{}).First(&workMomentsNotificationUnReadCount).Error, "")
	return workMomentsNotificationUnReadCount, err
}

func (d *DataBase) ClearWorkMomentsNotification() (err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Table("local_work_moments_notification").Where("1=1").Delete(&LocalWorkMomentsNotification{}).Error, "")
}
