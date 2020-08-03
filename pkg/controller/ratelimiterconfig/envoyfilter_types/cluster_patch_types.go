package envoyfilter_types

type SocketAddress struct {
	Address   string `json:"address" yaml:"address"`
	PortValue int32  `json:"port_value" yaml:"port_value"`
}

type Address struct {
	SocketAddress SocketAddress `json:"socket_address" yaml:"socket_address"`
}

type Endpoint struct {
	Address Address `json:"address" yaml:"address"`
}

type LbEndpoint struct {
	Endpoint Endpoint `json:"endpoint" yaml:"endpoint"`
}

type LoadAssignmentEndpoints struct {
	LbEndpoints []LbEndpoint `json:"lb_endpoints" yaml:"lb_endpoints"`
}

type LoadAssignment struct {
	ClusterName string                    `json:"cluster_name" yaml:"cluster_name"`
	Endpoints   []LoadAssignmentEndpoints `json:"endpoints" yaml:"endpoints"`
}

type Http2ProtocolOption struct{}

type ClusterPatchValues struct {
	Name                 string              `json:"name" yaml:"name"`
}
