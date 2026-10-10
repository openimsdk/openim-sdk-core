//go:build js

package network

import "open_im_sdk/sdk_struct"

// SetHTTPConfig is a no-op in browsers; browser networking is controlled by the host.
func SetHTTPConfig(_ sdk_struct.IMConfig) error { return nil }
