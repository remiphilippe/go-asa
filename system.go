package goasa

import (
	"github.com/golang/glog"
)

// Save Save ASA Configuration
func (a *ASA) Save() error {
	res, err := a.Post("commands/writemem", nil)
	if err != nil {
		if a.debug {
			glog.Infoln(res)
		}
		return err
	}

	return nil
}

// SystemObjectAny Any Object
func (a *ASA) SystemObjectAny() *ReferenceObject {
	return &ReferenceObject{Kind: "AnyIPAddress", Value: "any"}
}
