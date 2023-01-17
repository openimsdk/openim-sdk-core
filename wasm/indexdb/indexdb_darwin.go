//+build !js
package indexdb

type NilIndexDB struct {
}

func NewIndexDB(loginUserID string) *NilIndexDB {
	return &NilIndexDB{}
}
