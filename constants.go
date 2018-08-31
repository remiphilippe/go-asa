package goasa

const (
	apiPOST   string = "POST"
	apiPUT    string = "PUT"
	apiDELETE string = "DELETE"
	apiGET    string = "GET"

	apiMaxResults                   int    = 100
	apiBasePath                     string = "api/"
	apiNetworkObjectsEndpoint       string = "objects/networkobjects"
	apiNetworkObjectGroupsEndpoint  string = "objects/networkobjectgroups"
	apiAccessEndpoint               string = "access"
	apiNetworkServicesEndpoint      string = "objects/networkservices"
	apiNetworkServiceGroupsEndpoint string = "objects/networkservicegroups"

	networkObjectKind            string = "object#NetworkObj"
	networkObjectRefKind         string = "objectRef#NetworkObj"
	networkObjectGroupKind       string = "object#NetworkObjGroup"
	networkObjectGroupRefKind    string = "objectRef#NetworkObjGroup"
	extendedACEKind              string = "object#ExtendedACE"
	networkProtocolObjectKind    string = "object#NetworkProtocolObj"
	networkProtocolObjectRefKind string = "objectRef#NetworkProtocolObj"
	tcpUDPServiceObjectKind      string = "object#TcpUdpServiceObj"
	tcpUDPServiceObjectRefKind   string = "objectRef#TcpUdpServiceObj"
	icmpServiceObjectKind        string = "object#ICMPServiceObj"
	icmpServiceObjectRefKind     string = "objectRef#ICMPServiceObj"

	networkObjectTypeIPv4     string = "IPv4Address"
	networkObjectTypeProtocol string = "NetworkProtocol"

	errorDuplicate string = "DUPLICATE"

	//DuplicateActionError Error on duplicate
	DuplicateActionError int = 0

	//DuplicateActionDoNothing Don't do anything
	DuplicateActionDoNothing int = 1

	//DuplicateActionReplace Replace
	DuplicateActionReplace int = 2
)
