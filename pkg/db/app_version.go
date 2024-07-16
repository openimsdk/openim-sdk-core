// Copyright © 2024 OpenIM SDK. All rights reserved.
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

//go:build !js
// +build !js

package db

import (
	"context"

	"gorm.io/gorm"

	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
)

func (d *DataBase) GetAppSDKVersion(ctx context.Context) (*model_struct.LocalAppSDKVersion, error) {
	var appVersion model_struct.LocalAppSDKVersion
	err := d.conn.WithContext(ctx).Take(&appVersion).Error
	if err == gorm.ErrRecordNotFound {
		err = errs.ErrRecordNotFound
	}
	return &appVersion, errs.Wrap(err)
}

func (d *DataBase) SetAppSDKVersion(ctx context.Context, appVersion *model_struct.LocalAppSDKVersion) error {
	var exist model_struct.LocalAppSDKVersion
	err := d.conn.WithContext(ctx).First(&exist).Error
	if err == gorm.ErrRecordNotFound {
		if createErr := d.conn.WithContext(ctx).Create(appVersion).Error; createErr != nil {
			return errs.Wrap(createErr)
		}
		return nil
	} else if err != nil {
		return errs.Wrap(err)
	}

	if updateErr := d.conn.WithContext(ctx).Model(&exist).Updates(appVersion).Error; updateErr != nil {
		return errs.Wrap(updateErr)
	}

	return nil
}
