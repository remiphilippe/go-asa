package goasa

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

// NetworkObject An object represents the network (Note: The field level constraints listed here might not cover all the constraints on the field. Additional constraints might exist.)
type NetworkObject struct {
	ReferenceObject
	Host struct {
		Kind  string `json:"kind"`
		Value string `json:"value"`
	} `json:"host"`
}

// Reference Returns a reference object
func (n *NetworkObject) Reference() *ReferenceObject {
	r := ReferenceObject{
		Kind:     networkObjectRefKind,
		ObjectID: n.ObjectID,
		Name:     n.Name,
	}

	return &r
}

func (a *ASA) getNetworkObjects(limit, offset int) ([]*NetworkObject, error) {
	var err error
	var retval []*NetworkObject
	var l int

	if limit > apiMaxResults || limit < 1 {
		l = apiMaxResults
	} else {
		l = limit
	}

	query := make(map[string]string)
	query["limit"] = strconv.Itoa(l)
	query["offset"] = strconv.Itoa(offset)

	endpoint := apiNetworkObjectsEndpoint
	data, err := a.Get(endpoint, query)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*NetworkObject `json:"items"`
		Range rangeInfo        `json:"rangeInfo"`
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	retval = v.Items

	var upperBound int
	if limit == 0 {
		upperBound = v.Range.Total
	} else {
		if limit > v.Range.Total {
			upperBound = v.Range.Total
		} else {
			upperBound = limit
		}
	}

	glog.Infof("Upperbound: %d, Offset: %d, limit: %d\n", upperBound, v.Range.Offset, l)
	if v.Range.Offset+l < upperBound {

		res, err := a.getNetworkObjects(limit, v.Range.Offset+l)
		if err != nil {
			if a.debug {
				glog.Errorf("Error: %s\n", err)
			}
			return nil, err
		}

		retval = append(retval, res...)
	}

	return retval, nil
}

// GetAllNetworkObjects Get a list of all network objects
func (a *ASA) GetAllNetworkObjects() ([]*NetworkObject, error) {
	return a.getNetworkObjects(0, 0)
}

// GetNetworkObjects Get a list of all network objects
func (a *ASA) GetNetworkObjects(limit int) ([]*NetworkObject, error) {
	return a.getNetworkObjects(limit, 0)
}

// GetNetworkObjectByID Get a network object by ID
func (a *ASA) GetNetworkObjectByID(id string) (*NetworkObject, error) {
	var err error

	endpoint := fmt.Sprintf("%s/%s", apiNetworkObjectsEndpoint, id)
	data, err := a.Get(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var v *NetworkObject

	err = json.Unmarshal(data, &v)
	if err != nil {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	return v, nil
}

// CreateNetworkObject Create a new network object
func (a *ASA) CreateNetworkObject(n *NetworkObject, duplicateAction int) error {
	var err error

	n.Kind = networkObjectKind
	_, err = a.Post(apiNetworkObjectsEndpoint, n)
	if err != nil {
		asaErr := err.(ASAError)
		//spew.Dump(asaErr)
		if asaErr.Code == errorDuplicate {
			if a.debug {
				glog.Warningf("This is a duplicate\n")
			}
			if duplicateAction == DuplicateActionError {
				return err
			}
		} else {
			if a.debug {
				glog.Errorf("Error: %s\n", err)
			}
			return err
		}
	}

	// query := fmt.Sprintf("name:%s", n.Name)
	// obj, err := f.getNetworkObjectBy(query, 0)
	// if err != nil {
	// 	if f.debug {
	// 		glog.Errorf("Error: %s\n", err)
	// 	}
	// 	return err
	// }

	// var o *NetworkObject
	// if len(obj) == 1 {
	// 	o = obj[0]
	// } else {
	// 	if f.debug {
	// 		glog.Errorf("Error: length of object is not 1\n")
	// 	}
	// 	return err
	// }

	// switch duplicateAction {
	// case DuplicateActionReplace:
	// 	o.Value = n.Value
	// 	o.SubType = n.SubType

	// 	err = f.UpdateNetworkObject(o)
	// 	if err != nil {
	// 		if f.debug {
	// 			glog.Errorf("Error: %s\n", err)
	// 		}
	// 		return err
	// 	}
	// }

	// *n = *o
	return nil
}

// DeleteNetworkObject Delete a network object
func (a *ASA) DeleteNetworkObject(n interface{}) error {
	var err error
	var objectID string

	switch v := n.(type) {
	case *ReferenceObject:
		objectID = v.ObjectID
	case *NetworkObject:
		objectID = v.ObjectID
	case string:
		objectID = v
	default:
		return fmt.Errorf("unknown type")
	}

	err = a.Delete(fmt.Sprintf("%s/%s", apiNetworkObjectsEndpoint, objectID))
	if err != nil {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}
