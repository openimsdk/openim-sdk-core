package organization

import (
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/db"
)

type Organization struct {
	listener    open_im_sdk_callback.OnOrganizationListener
	loginUserID string
	db          *db.DataBase
	p           *ws.PostApi
}
