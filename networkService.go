package goasa

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

// NetworkService An object represents the network (Note: The field level constraints listed here might not cover all the constraints on the field. Additional constraints might exist.)
type NetworkService struct {
	ReferenceObject
}

// Reference Returns a reference object
func (n *NetworkService) Reference() *ReferenceObject {
	var kind string

	switch n.Kind {
	case tcpUDPServiceObjectKind:
		kind = tcpUDPServiceObjectRefKind
	case icmpServiceObjectKind:
		kind = icmpServiceObjectRefKind
	case networkProtocolObjectKind:
		kind = networkProtocolObjectRefKind
	}
	r := ReferenceObject{
		Kind:     kind,
		ObjectID: n.ObjectID,
		Name:     n.Name,
	}

	return &r
}

func (a *ASA) getNetworkServices(limit, offset int) ([]*NetworkService, error) {
	var err error
	var retval []*NetworkService
	var l int

	if limit > apiMaxResults || limit < 1 {
		l = apiMaxResults
	} else {
		l = limit
	}

	query := make(map[string]string)
	query["limit"] = strconv.Itoa(l)
	query["offset"] = strconv.Itoa(offset)

	endpoint := apiNetworkServicesEndpoint
	data, err := a.Get(endpoint, query)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*NetworkService `json:"items"`
		Range rangeInfo         `json:"rangeInfo"`
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

		res, err := a.getNetworkServices(limit, v.Range.Offset+l)
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

// GetAllNetworkServices Get a list of all network objects
func (a *ASA) GetAllNetworkServices() ([]*NetworkService, error) {
	return a.getNetworkServices(0, 0)
}

// GetNetworkServices Get a list of all network objects
func (a *ASA) GetNetworkServices(limit int) ([]*NetworkService, error) {
	return a.getNetworkServices(limit, 0)
}

// GetNetworkServicesByID Get a network object by ID
func (a *ASA) GetNetworkServicesByID(id string) (*NetworkService, error) {
	var err error

	endpoint := fmt.Sprintf("%s/%s", apiNetworkServicesEndpoint, id)
	data, err := a.Get(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var v *NetworkService

	err = json.Unmarshal(data, &v)
	if err != nil {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	return v, nil
}

// CreateTCPService Create a TCP Service
func (a *ASA) CreateTCPService(name, port string, duplicateAction int) (*NetworkService, error) {
	var err error

	n := new(NetworkService)
	n.Kind = tcpUDPServiceObjectKind
	n.Name = name
	n.Value = port

	err = a.CreateNetworkService(n, duplicateAction)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// CreateUDPService Create a TCP Service
func (a *ASA) CreateUDPService(name, port string, duplicateAction int) (*NetworkService, error) {
	return a.CreateTCPService(name, port, duplicateAction)
}

// CreateNetworkService Create a new network object
func (a *ASA) CreateNetworkService(n *NetworkService, duplicateAction int) error {
	var err error
	_, err = a.Post(apiNetworkServicesEndpoint, n)
	if err != nil {
		asaErr := err.(ASAError)

		n.ObjectID = asaErr.Details
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
		n.ObjectID = n.Name
	}

	return nil
}

// DeleteNetworkService Delete a network object
func (a *ASA) DeleteNetworkService(n interface{}) error {
	var err error
	var objectID string

	switch v := n.(type) {
	case *ReferenceObject:
		objectID = v.ObjectID
	case *NetworkService:
		objectID = v.ObjectID
	case string:
		objectID = v
	default:
		return fmt.Errorf("unknown type")
	}

	if objectID == "" {
		return fmt.Errorf("error objectid is null")
	}

	err = a.Delete(fmt.Sprintf("%s/%s", apiNetworkServicesEndpoint, objectID))
	if err != nil {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}

// CreateNetworkServiceFromPorts Create Network objects from an array of IP
func (a *ASA) CreateNetworkServiceFromPorts(protocol string, ports []int) ([]*NetworkService, error) {
	var err error
	var retval []*NetworkService
	var val string

	protocol = strings.ToLower(protocol)

	objs, err := a.GetAllNetworkServices()
	if err != nil {
		if a.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	found := make(map[string]bool)

	for i := range ports {
		val = fmt.Sprintf("%s/%d", protocol, ports[i])
		for o := range objs {
			if objs[o].Value == val {
				retval = append(retval, objs[o])
				found[val] = true
				break
			}
		}
	}

	for i := range ports {
		val = fmt.Sprintf("%s/%d", protocol, ports[i])
		if _, ok := found[val]; !ok {
			switch protocol {
			case "udp":
				s, err := a.CreateUDPService(fmt.Sprintf("%s-%d", protocol, ports[i]), val, DuplicateActionDoNothing)
				if err != nil {
					if a.debug {
						glog.Errorf("Error: %s\n", err)
					}
					return nil, err
				}
				retval = append(retval, s)
			case "tcp":
				s, err := a.CreateTCPService(fmt.Sprintf("%s-%d", protocol, ports[i]), val, DuplicateActionDoNothing)
				if err != nil {
					if a.debug {
						glog.Errorf("Error: %s\n", err)
					}
					return nil, err
				}
				retval = append(retval, s)
			}
		}
	}

	return retval, nil
}
