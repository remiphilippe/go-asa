package goasa

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

// NetworkServiceGroup Network Object Group
type NetworkServiceGroup struct {
	ReferenceObject
	Description string             `json:"description,omitempty"`
	Members     []*ReferenceObject `json:"members,omitempty"`
}

// UnmarshalJSON Converts JSON to struct
// func (g *NetworkServiceGroup) UnmarshalJSON(data []byte) error {
// 	var err error
// 	type Alias NetworkServiceGroup

// 	aux := &struct {
// 		Members []interface{} `json:"members"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(g),
// 	}

// 	if err = json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}

// 	for _, member := range aux.Members {
// 		m := member.(map[string]interface{})
// 		if _, ok := m["kind"]; ok {
// 			if m["kind"].(string) == networkObjectGroupRefKind {
// 				r := new(ReferenceObject)
// 				r.Kind = m["kind"].(string)
// 				r.SelfLink = m["refLink"].(string)
// 				r.ObjectID = m["objectId"].(string)
// 				g.Members = append(g.Members, r)
// 			}
// 		}
// 		//TODO: handle "kind": "IPv4Network", and others
// 	}

// 	return nil
// }

// Reference Returns a reference object
func (g *NetworkServiceGroup) Reference() *ReferenceObject {
	var kind string

	switch g.Kind {
	case networkServiceGroupObjectKind:
		kind = networkServiceGroupObjectRefKind
	}

	r := ReferenceObject{
		ObjectID: g.ObjectID,
		Name:     g.Name,
		Kind:     kind,
	}

	return &r
}

func (a *ASA) getNetworkServiceGroups(limit, offset int) ([]*NetworkServiceGroup, error) {
	var err error
	var retval []*NetworkServiceGroup
	var l int

	if limit > apiMaxResults || limit < 1 {
		l = apiMaxResults
	} else {
		l = limit
	}

	query := make(map[string]string)
	query["limit"] = strconv.Itoa(l)
	query["offset"] = strconv.Itoa(offset)

	endpoint := apiNetworkServiceGroupsEndpoint
	data, err := a.Get(endpoint, query)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*NetworkServiceGroup `json:"items"`
		Range rangeInfo              `json:"rangeInfo"`
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

		res, err := a.getNetworkServiceGroups(limit, v.Range.Offset+l)
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

// GetAllNetworkServiceGroups Get a list of all network object groups
func (a *ASA) GetAllNetworkServiceGroups() ([]*NetworkServiceGroup, error) {
	return a.getNetworkServiceGroups(0, 0)
}

// GetNetworkServiceGroups Get a list of all network object groups
func (a *ASA) GetNetworkServiceGroups(limit int) ([]*NetworkServiceGroup, error) {
	return a.getNetworkServiceGroups(limit, 0)
}

// GetNetworkServiceGroupByID Get a network object by ID
func (a *ASA) GetNetworkServiceGroupByID(id string) (*NetworkServiceGroup, error) {
	var err error

	endpoint := fmt.Sprintf("%s/%s", apiNetworkObjectGroupsEndpoint, id)
	data, err := a.Get(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var v *NetworkServiceGroup

	err = json.Unmarshal(data, &v)
	if err != nil {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	return v, nil
}

// CreateNetworkServiceGroup Create a new network object
func (a *ASA) CreateNetworkServiceGroup(g *NetworkServiceGroup, duplicateAction int) error {
	var err error

	//g.Kind = networkObjectKind
	_, err = a.Post(apiNetworkServiceGroupsEndpoint, g)
	if err != nil {
		asaErr := err.(ASAError)
		//spew.Dump(asaErr)
		g.ObjectID = asaErr.Details
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
	} else {
		g.ObjectID = g.Name
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

// DeleteNetworkServiceGroup Delete a network object
func (a *ASA) DeleteNetworkServiceGroup(g interface{}) error {
	var err error
	var objectID string

	switch v := g.(type) {
	case *ReferenceObject:
		objectID = v.ObjectID
	case *NetworkServiceGroup:
		objectID = v.ObjectID
	case string:
		objectID = v
	default:
		return fmt.Errorf("unknown type")
	}

	if objectID == "" {
		return fmt.Errorf("error objectid is null")
	}

	err = a.Delete(fmt.Sprintf("%s/%s", apiNetworkServiceGroupsEndpoint, objectID))
	if err != nil {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}

// CreateNetworkServiceGroupFromMap Create an object group from an array of ip address. Network objects = ip.
func (a *ASA) CreateNetworkServiceGroupFromMap(name string, ports map[string][]int, duplicateAction int) (*NetworkServiceGroup, error) {
	var err error
	var ns []*NetworkService

	for k, v := range ports {
		s, err := a.CreateNetworkServiceFromPorts(strings.ToLower(k), v)
		if err != nil {
			if a.debug {
				glog.Errorf("Error: %s\n", err)
			}
			return nil, err
		}
		ns = append(ns, s...)
	}

	g := new(NetworkServiceGroup)
	g.Name = name
	g.Kind = "object#NetworkServiceGroup"

	for i := range ns {
		g.Members = append(g.Members, ns[i].Reference())
	}

	err = a.CreateNetworkServiceGroup(g, duplicateAction)
	if err != nil {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	return g, nil
}
