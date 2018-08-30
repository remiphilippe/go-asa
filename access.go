package goasa

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

// AccessRule Acces Rule object
type AccessRule struct {
	ReferenceObject
	SourceService      *ReferenceObject `json:"sourceService"`
	SrcSecurity        *ReferenceObject `json:"srcSecurity,omitempty"`
	RuleLogging        interface{}      `json:"ruleLogging,omitempty"`
	IsAccessRule       bool             `json:"isAccessRule,omitempty"`
	TimeRange          *ReferenceObject `json:"timeRange,omitempty"`
	DestinationAddress *ReferenceObject `json:"destinationAddress,omitempty"`
	Active             bool             `json:"active,omitempty"`
	DestinationService *ReferenceObject `json:"destinationService,omitempty"`
	DstSecurity        *ReferenceObject `json:"dstSecurity,omitempty"`
	User               *ReferenceObject `json:"user,omitempty"`
	Permit             bool             `json:"permit"`
	Remarks            []string         `json:"remarks,omitempty"`
	Position           int              `json:"position,omitempty"`
	SourceAddress      *ReferenceObject `json:"sourceAddress,omitempty"`
}

// Reference Returns a reference object
func (a *AccessRule) Reference() *ReferenceObject {
	r := ReferenceObject{
		Kind:     extendedACEKind,
		ObjectID: a.ObjectID,
		Name:     a.Name,
	}

	return &r
}

func (a *ASA) getAccessObjects(limit, offset int, intf string) ([]*AccessRule, error) {
	var err error
	var retval []*AccessRule
	var l int

	if limit > apiMaxResults || limit < 1 {
		l = apiMaxResults
	} else {
		l = limit
	}

	query := make(map[string]string)
	query["limit"] = strconv.Itoa(l)
	query["offset"] = strconv.Itoa(offset)

	endpoint := fmt.Sprintf("%s/%s/rules", apiAccessEndpoint, intf)
	//spew.Dump(endpoint)
	data, err := a.Get(endpoint, query)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*AccessRule `json:"items"`
		Range rangeInfo     `json:"rangeInfo"`
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

		res, err := a.getAccessObjects(limit, v.Range.Offset+l, intf)
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

// GetAllGlobalAccessObjects Get a list of all network objects
func (a *ASA) GetAllGlobalAccessObjects() ([]*AccessRule, error) {
	return a.getAccessObjects(0, 0, "global")
}

// GetGlobalAccessObjects Get a list of all network objects
func (a *ASA) GetGlobalAccessObjects(limit int) ([]*AccessRule, error) {
	return a.getAccessObjects(limit, 0, "global")
}

// GetAccessObjectByID Get a network object by ID
func (a *ASA) GetAccessObjectByID(intf, id string) (*AccessRule, error) {
	var err error

	endpoint := fmt.Sprintf("%s/%s/rules/%s", apiAccessEndpoint, intf, id)
	data, err := a.Get(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var v *AccessRule

	err = json.Unmarshal(data, &v)
	if err != nil {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	return v, nil
}

// CreateAccessRule Create a new network object
func (a *ASA) CreateAccessRule(intf string, r *AccessRule, duplicateAction int) error {
	var err error

	r.Kind = extendedACEKind
	endpoint := fmt.Sprintf("%s/%s/rules", apiAccessEndpoint, intf)
	_, err = a.Post(endpoint, r)
	if err != nil {
		asaErr := err.(ASAError)
		//spew.Dump(asaErr)
		if asaErr.Code == errorDuplicate {
			r.ObjectID = asaErr.Details
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

// DeleteAccessRule Delete a network object
func (a *ASA) DeleteAccessRule(intf string, n interface{}) error {
	var err error
	var objectID string

	switch v := n.(type) {
	case *ReferenceObject:
		objectID = v.ObjectID
	case *AccessRule:
		objectID = v.ObjectID
	case string:
		objectID = v
	default:
		return fmt.Errorf("unknown type")
	}

	if objectID == "" {
		return fmt.Errorf("error objectid is null")
	}

	endpoint := fmt.Sprintf("%s/%s/rules/%s", apiAccessEndpoint, intf, objectID)
	err = a.Delete(endpoint)
	if err != nil {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}
