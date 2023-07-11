package scope

const (
	DhcpRead  = "dhcp:read"
	DhcpWrite = "dhcp:write"
	DhcpAdmin = "dhcp:admin"
)

var DhcpCanRead = []string{DhcpRead, DhcpWrite, DhcpAdmin}
var DhcpCanAdd = []string{DhcpWrite, DhcpAdmin}
var DhcpCanChange = []string{DhcpAdmin}
