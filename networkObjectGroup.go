package goasa

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

// NetworkObjectGroup Network Object Group
type NetworkObjectGroup struct {
	ReferenceObject
	Description string        `json:"description,omitempty"`
	Members     []interface{} `json:"members,omitempty"`
}

// UnmarshalJSON Converts JSON to struct
func (g *NetworkObjectGroup) UnmarshalJSON(data []byte) error {
	var err error
	type Alias NetworkObjectGroup

	aux := &struct {
		Members []interface{} `json:"members"`
		*Alias
	}{
		Alias: (*Alias)(g),
	}

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	for _, member := range aux.Members {
		m := member.(map[string]interface{})
		if _, ok := m["kind"]; ok {
			if m["kind"].(string) == networkObjectGroupRefKind {
				r := new(ReferenceObject)
				r.Kind = m["kind"].(string)
				r.SelfLink = m["refLink"].(string)
				r.ObjectID = m["objectId"].(string)
				g.Members = append(g.Members, r)
			}
		}
		//TODO: handle "kind": "IPv4Network", and others
	}

	return nil
}

// Reference Returns a reference object
func (g *NetworkObjectGroup) Reference() *ReferenceObject {
	r := ReferenceObject{
		ObjectID: g.ObjectID,
		Name:     g.Name,
		Kind:     networkObjectGroupRefKind,
	}

	return &r
}

func (a *ASA) getNetworkObjectGroups(limit, offset int) ([]*NetworkObjectGroup, error) {
	var err error
	var retval []*NetworkObjectGroup
	var l int

	if limit > apiMaxResults || limit < 1 {
		l = apiMaxResults
	} else {
		l = limit
	}

	query := make(map[string]string)
	query["limit"] = strconv.Itoa(l)
	query["offset"] = strconv.Itoa(offset)

	endpoint := apiNetworkObjectGroupsEndpoint
	data, err := a.Get(endpoint, query)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*NetworkObjectGroup `json:"items"`
		Range rangeInfo             `json:"rangeInfo"`
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

		res, err := a.getNetworkObjectGroups(limit, v.Range.Offset+l)
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

// GetAllNetworkObjectGroups Get a list of all network object groups
func (a *ASA) GetAllNetworkObjectGroups() ([]*NetworkObjectGroup, error) {
	return a.getNetworkObjectGroups(0, 0)
}

// GetNetworkObjectGroups Get a list of all network object groups
func (a *ASA) GetNetworkObjectGroups(limit int) ([]*NetworkObjectGroup, error) {
	return a.getNetworkObjectGroups(limit, 0)
}

// GetNetworkObjectGroupByID Get a network object by ID
func (a *ASA) GetNetworkObjectGroupByID(id string) (*NetworkObjectGroup, error) {
	var err error

	endpoint := fmt.Sprintf("%s/%s", apiNetworkObjectGroupsEndpoint, id)
	data, err := a.Get(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var v *NetworkObjectGroup

	err = json.Unmarshal(data, &v)
	if err != nil {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	return v, nil
}

// CreateNetworkObjectGroup Create a new network object
func (a *ASA) CreateNetworkObjectGroup(g *NetworkObjectGroup) error {
	var err error

	g.Kind = networkObjectGroupKind
	_, err = a.Post(apiNetworkObjectGroupsEndpoint, g)
	if err != nil {
		// ftdErr := err.(*FTDError)
		// //spew.Dump(ftdErr)
		// if len(ftdErr.Message) > 0 && (ftdErr.Message[0].Code == "duplicateName" || ftdErr.Message[0].Code == "newInstanceWithDuplicateId") {
		// 	if f.debug {
		// 		glog.Warningf("This is a duplicate\n")
		// 	}
		// 	if duplicateAction == DuplicateActionError {
		// 		return err
		// 	}
		// } else {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
		// }
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

// DeleteNetworkObjectGroup Delete a network object
func (a *ASA) DeleteNetworkObjectGroup(g *NetworkObjectGroup) error {
	var err error

	err = a.Delete(fmt.Sprintf("%s/%s", apiNetworkObjectGroupsEndpoint, g.ObjectID))
	if err != nil {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}
