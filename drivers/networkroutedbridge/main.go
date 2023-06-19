//go:build linux

package networkroutedbridge

import (
	"fmt"
	"math"
	"net"
	"net/netip"
	"strings"

	"github.com/vishvananda/netlink"

	"github.com/opensvc/om3/core/driver"
	"github.com/opensvc/om3/core/keyop"
	"github.com/opensvc/om3/core/network"
	"github.com/opensvc/om3/util/hostname"
	"github.com/opensvc/om3/util/key"
)

type (
	T struct {
		network.T
	}
)

var (
	drvID = driver.NewID(driver.GroupNetwork, "routed_bridge")
)

func init() {
	driver.Register(drvID, NewNetworker)
}

func NewNetworker() network.Networker {
	t := New()
	var i interface{} = t
	return i.(network.Networker)
}

func New() *T {
	t := T{}
	return &t
}

func (t T) Usage() (network.StatusUsage, error) {
	usage := network.StatusUsage{}
	return usage, nil
}

// CNIConfigData returns a cni network configuration, like
//
//	{
//	   "cniVersion": "0.3.0",
//	   "name": "net1",
//	   "type": "bridge",
//	   "bridge": "obr_net1",
//	   "isGateway": true,
//	   "ipMasq": false,
//	   "ipam": {
//	       "type": "host-local",
//	       "subnet": "10.23.0.0/26",
//	       "routes": [
//	           {
//	               "dst": "0.0.0.0/0"
//	           },
//	           {
//	               "dst": "10.23.0.0/24",
//	               "gw": "10.23.0.1"
//	           }
//	       ]
//	   }
//	}
func (t T) CNIConfigData() (interface{}, error) {
	name := t.Name()
	nwStr := t.Network()
	brName := t.brName()
	brIP, err := t.bridgeIP()
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{
		"cniVersion": network.CNIVersion,
		"name":       name,
		"type":       "bridge",
		"bridge":     brName,
		"isGateway":  true,
		"ipMasq":     false,
		"ipam": map[string]interface{}{
			"type": "host-local",
			"routes": []map[string]interface{}{
				{"dst": defaultRouteDst(nwStr)},
				{"dst": nwStr, "gw": brIP.String()},
			},
			"subnet": t.subnet(),
		},
	}
	return m, nil
}

func defaultRouteDst(cidr string) string {
	if isIP6(cidr) {
		return "::/0"
	} else {
		return "0.0.0.0/0"
	}
}

func isIP6(cidr string) bool {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	return ip.To4() == nil
}

func (t T) bridgeIP() (net.IP, error) {
	subnetStr := t.subnet()
	if subnetStr == "" {
		return nil, fmt.Errorf("network#%s.subnet is required", t.Name())
	}
	ip, _, err := net.ParseCIDR(subnetStr)
	if err != nil {
		return nil, err
	}
	ip[len(ip)-1]++
	return ip, nil
}

func (t T) allocateSubnets() error {
	subnetMap, err := t.subnetMap()
	if err != nil {
		return err
	}
	for _, s := range subnetMap {
		if s != "" {
			// A subnet setup is already configured. Don't change, even if bogus.
			return nil
		}
	}

	ipsPerNode := t.GetInt("ips_per_node")
	if ipsPerNode == 0 {
		return fmt.Errorf("ips_per_node must be greater than 0")
	}
	if ipsPerNode&(ipsPerNode-1) == 1 {
		return fmt.Errorf("ips_per_node must be a power of 2")
	}

	ip, network, err := net.ParseCIDR(t.Network())
	if err != nil {
		return err
	}
	nodes, err := t.Nodes()
	if err != nil {
		return err
	}
	ones, bits := network.Mask.Size()
	maxIps := 1 << (bits - ones)
	maxIpsPerNodes := maxIps / len(nodes)
	subnetOnes := int(math.Log2(float64(maxIpsPerNodes)))
	maxIpsPerNodes = 1 << subnetOnes
	if ipsPerNode > maxIpsPerNodes {
		return fmt.Errorf("ips_per_node must be <%d (%d ips in %s divided by %d nodes)", maxIpsPerNodes, maxIps, network, len(nodes))
	}
	addr, err := netip.ParseAddr(ip.String())
	if err != nil {
		return err
	}

	kops := make([]keyop.T, 0)
	for _, nodename := range nodes {
		subnet := fmt.Sprintf("%s/%d", addr, bits-subnetOnes)
		for i := 0; i < maxIpsPerNodes; i += 1 {
			addr = addr.Next()
		}
		kops = append(kops, keyop.T{
			Key:   key.New("network#"+t.Name(), "subnet@"+nodename),
			Op:    keyop.Set,
			Value: subnet,
		})
		t.Log().Info().Msgf("assign subnet %s to node %s", subnet, nodename)
	}
	return t.Config().SetKeys(kops...)
}

func (t *T) Setup() error {
	var (
		localIP net.IP
		brIP    net.IP
		link    netlink.Link
		err     error
	)
	if err := t.allocateSubnets(); err != nil {
		return err
	}
	if brIP, err = t.bridgeIP(); err != nil {
		return err
	}
	if link, err = t.setupBridge(); err != nil {
		return fmt.Errorf("setup br: %w", err)
	}
	if err := t.setupBridgeIP(link, brIP); err != nil {
		return fmt.Errorf("setup br ip: %w", err)
	}
	if err := t.setupBridgeMAC(link, brIP); err != nil {
		return fmt.Errorf("setup mac: %w", err)
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("link up: %w", err)
	}
	if localIP, err = t.getLocalIP(); err != nil {
		return fmt.Errorf("get local ip: %w", err)
	}
	nodes, err := t.Nodes()
	if err != nil {
		return err
	}
	for idx, nodename := range nodes {
		if err := t.setupNode(nodename, idx, localIP, brIP); err != nil {
			return fmt.Errorf("setup network to node %s: %w", nodename, err)
		}
	}
	return nil
}

func (t *T) setupNode(nodename string, nodeIndex int, localIP, brIP net.IP) error {
	var route network.Route
	tunnel := t.tunnel()
	peerIP, err := t.getNodeIP(nodename)
	if err != nil {
		return fmt.Errorf("get peer ip: %w", err)
	}
	dst, err := t.NodeSubnet(nodename)
	if err != nil {
		return fmt.Errorf("get peer subnet: %w", err)
	}
	if dst == nil {
		return fmt.Errorf("no peer subnet: %w", err)
	}
	if hostname.Hostname() == nodename {
		route = network.Route{
			Nodename: nodename,
			Dst:      dst,
			Dev:      t.brName(),
		}
	} else if v, err := t.mustTunnel(tunnel, peerIP); err != nil {
		return fmt.Errorf("must tunnel: %w", err)
	} else if v {
		name := tunName(peerIP, nodeIndex)
		if err := t.setupNodeTunnelLink(nodename, name, localIP, peerIP); err != nil {
			return fmt.Errorf("setup tunnel: %w", err)
		}
		if err := t.setupNodeTunnelLinkUp(name); err != nil {
			return fmt.Errorf("setup tunnel: %w", err)
		}
		route = network.Route{
			Nodename: nodename,
			Dst:      dst,
			Dev:      name,
			Src:      brIP,
		}
	} else {
		route = network.Route{
			Nodename: nodename,
			Dst:      dst,
			Gateway:  peerIP,
		}
	}
	if err := t.setupNodeRoutes(route); err != nil {
		return fmt.Errorf("setup route: %w", err)
	}
	return nil
}

func (t T) mustTunnel(tunnel string, peerIP net.IP) (bool, error) {
	if tunnel == "never" {
		return true, nil
	}
	if tunnel == "always" {
		return false, nil
	}
	if _, ipnet, err := network.IPReachableFrom(peerIP); err != nil {
		return false, err
	} else if ipnet != nil {
		t.Log().Debug().Msgf("%s is reachable from %s. tunnel not needed", peerIP, ipnet)
		return false, nil
	}
	return true, nil
}

func (t *T) setupNodeTunnelLink(nodename, name string, localIP, peerIP net.IP) error {
	// clean up existing tunnels with same endpoints but different name
	link, err := t.getTunnelByEndpoints(localIP, peerIP)
	if err != nil {
		return fmt.Errorf("get tunnel from %s to %s: %w", localIP, peerIP, err)
	}
	if link != nil {
		if link.Attrs().Name == name {
			t.Log().Info().Msgf("tunnel to %s is already setup", nodename)
			return nil
		} else {
			t.Log().Info().Msgf("delete conflicting tunnel %s from %s to %s", link.Attrs().Name, localIP, peerIP)
			if err := netlink.LinkDel(link); err != nil {
				return err
			}
		}
	}

	// modify up existing tunnels with same name but different endpoints
	// or add a new tunnel
	link, err = netlink.LinkByName(name)
	switch {
	case err != nil:
		if _, ok := err.(netlink.LinkNotFoundError); !ok {
			return err
		}
		fallthrough
	case link == nil:
		if err := t.addTunnel(name, localIP, peerIP); err != nil {
			return fmt.Errorf("add tunnel: %w", err)
		}
	case link != nil && t.isSameTunnel(link, localIP, peerIP):
		t.Log().Info().Msgf("tunnel to %s is already configured", nodename)
		return nil
	default:
		if err := t.modTunnel(name, localIP, peerIP); err != nil {
			return fmt.Errorf("modify tunnel to %s: %w", nodename, err)
		}
	}
	return nil
}

func (t *T) setupNodeTunnelLinkUp(name string) error {
	if link, err := netlink.LinkByName(name); err != nil {
		return fmt.Errorf("link up: %w", err)
	} else if link != nil {
		t.Log().Info().Msgf("link up %s", name)
		netlink.LinkSetUp(link)
	}

	return nil
}

func (t T) getTunnelByEndpoints(localIP, peerIP net.IP) (netlink.Link, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		var local, remote net.IP
		switch tun := link.(type) {
		case *netlink.Iptun:
			local = tun.Local
			remote = tun.Remote
		case *netlink.Ip6tnl:
			local = tun.Local
			remote = tun.Remote
		}
		if local.Equal(localIP) && remote.Equal(peerIP) {
			return link, nil
		}
	}
	return nil, nil
}

func (t T) isSameTunnel(link netlink.Link, localIP, peerIP net.IP) bool {
	name := link.Attrs().Name
	var local, remote net.IP
	switch tun := link.(type) {
	case *netlink.Iptun:
		if localIP.To4() == nil {
			t.Log().Info().Msgf("link %s is not a ipip tunnel", name)
			return false
		}
		local = tun.Local
		remote = tun.Remote
	case *netlink.Ip6tnl:
		if localIP.To4() != nil {
			t.Log().Info().Msgf("link %s is not a ip6ip6 tunnel", name)
			return false
		}
		local = tun.Local
		remote = tun.Remote
	}
	if !local.Equal(localIP) {
		t.Log().Info().Msgf("tunnel %s local ip is %s, should be %s", name, local, localIP)
		return false
	}
	if !remote.Equal(peerIP) {
		t.Log().Info().Msgf("tunnel %s remote ip is %s, should be %s", name, remote, peerIP)
		return false
	}
	return true
}

func (t T) modTunnel(name string, localIP, peerIP net.IP) error {
	if localIP.To4() == nil {
		return t.modTunnel6(name, localIP, peerIP)
	} else {
		return t.modTunnel4(name, localIP, peerIP)
	}
}

func (t T) addTunnel(name string, localIP, peerIP net.IP) error {
	if localIP.To4() == nil {
		return t.addTunnel6(name, localIP, peerIP)
	} else {
		return t.addTunnel4(name, localIP, peerIP)
	}
}

func (t T) modTunnel6(name string, localIP, peerIP net.IP) error {
	link := &netlink.Ip6tnl{
		LinkAttrs: netlink.LinkAttrs{
			Name:      name,
			EncapType: "ip6tnl",
		},
		Local:  localIP,
		Remote: peerIP,
	}
	t.Log().Info().Interface("link", link).Msgf("modify ipip tun %s", name)
	if h, err := netlink.NewHandle(); err != nil {
		defer h.Delete()
		return h.LinkModify(link)
	} else {
		return err
	}
}

func (t T) modTunnel4(name string, localIP, peerIP net.IP) error {
	link := &netlink.Iptun{
		LinkAttrs: netlink.LinkAttrs{
			Name:      name,
			EncapType: "ipip",
		},
		Local:  localIP,
		Remote: peerIP,
	}
	t.Log().Info().Interface("link", link).Msgf("modify ipip tun %s", name)
	if h, err := netlink.NewHandle(); err != nil {
		defer h.Delete()
		return h.LinkModify(link)
	} else {
		return err
	}
}

func (t T) addTunnel6(name string, localIP, peerIP net.IP) error {
	link := &netlink.Ip6tnl{
		LinkAttrs: netlink.LinkAttrs{
			Name:      name,
			EncapType: "ip6tnl",
		},
		Local:  localIP,
		Remote: peerIP,
	}
	t.Log().Info().Interface("link", link).Msgf("add ipip tun %s", name)
	return netlink.LinkAdd(link)
}

func (t T) addTunnel4(name string, localIP, peerIP net.IP) error {
	link := &netlink.Iptun{
		LinkAttrs: netlink.LinkAttrs{
			Name:      name,
			EncapType: "ipip",
		},
		Local:  localIP,
		Remote: peerIP,
	}
	t.Log().Info().Interface("link", link).Msgf("add ipip tun %s", name)
	return netlink.LinkAdd(link)
}

func tunName(peerIP net.IP, nodeIndex int) string {
	if peerIP.To4() == nil {
		return fmt.Sprintf("otun%d", nodeIndex)
	} else {
		return fmt.Sprintf("tun%s", strings.ReplaceAll(peerIP.String(), ".", ""))
	}
}

func (t T) setupBridge() (netlink.Link, error) {
	la := netlink.NewLinkAttrs()
	la.Name = t.brName()
	link, err := netlink.LinkByName(la.Name)
	_, linkNotFound := err.(netlink.LinkNotFoundError)
	switch {
	case linkNotFound:
	case err != nil:
		return nil, err
	case link != nil:
		t.Log().Info().Msgf("bridge link %s already exists", la.Name)
		return link, nil
	}
	br := &netlink.Bridge{LinkAttrs: la}
	err = netlink.LinkAdd(br)
	if err != nil {
		return nil, fmt.Errorf("failed to add bridge link %s: %v", la.Name, err)
	}
	t.Log().Info().Msgf("added bridge link %s", la.Name)
	return br, nil
}

func (t T) subnet() string {
	return t.GetString("subnet")
}

func (t T) subnetMap() (map[string]string, error) {
	m := make(map[string]string)
	nodes, err := t.Nodes()
	if err != nil {
		return nil, err
	}
	for _, nodename := range nodes {
		m[nodename] = t.GetString("subnet@" + nodename)
	}
	return m, nil
}

func (t T) tunnel() string {
	return t.GetString("tunnel")
}

func (t T) brName() string {
	return "obr_" + t.Name()
}

func (t T) BackendDevName() string {
	return t.brName()
}

func (t T) setupBridgeMAC(br netlink.Link, brIP net.IP) error {
	var (
		mac net.HardwareAddr
		err error
	)
	if br == nil {
		return nil
	}
	if mac, err = network.MACFromIP4(brIP); err != nil {
		return err
	}
	if br.Attrs().HardwareAddr.String() == mac.String() {
		t.Log().Info().Msgf("bridge %s mac is already %s", br.Attrs().Name, mac)
		return nil
	}
	t.Log().Info().Msgf("bridge %s set mac to %s", br.Attrs().Name, mac)
	return netlink.LinkSetHardwareAddr(br, mac)
}

func (t T) setupBridgeIP(br netlink.Link, brIP net.IP) error {
	if br == nil {
		return nil
	}
	brName := t.brName()
	subnetStr := t.subnet()
	_, ipnet, err := net.ParseCIDR(subnetStr)
	if err != nil {
		return err
	}
	ipnet.IP = brIP
	ipnetStr := ipnet.String()

	if intf, err := net.InterfaceByName(brName); err != nil {
		return err
	} else if addrs, err := intf.Addrs(); err != nil {
		return err
	} else {
		for _, addr := range addrs {
			if addr.String() == ipnetStr {
				t.Log().Info().Msgf("bridge ip %s already added to %s", ipnet, brName)
				return nil
			}
		}
	}
	addr := &netlink.Addr{IPNet: ipnet}
	if err := netlink.AddrAdd(br, addr); err != nil {
		return err
	}
	t.Log().Info().Msgf("added ip %s to bridge %s", ipnet, brName)
	return nil
}

// getNodeIP returns the addr scoped for nodename from the network config.
// Defaults to the first resolved ip address with the network address family (ip4 or ip6).
func (t T) getNodeIP(nodename string) (net.IP, error) {
	var keyName string
	if nodename == hostname.Hostname() {
		keyName = "addr"
	} else {
		keyName = "addr@" + nodename
	}
	if addr := t.GetString(keyName); addr != "" {
		return net.ParseIP(addr), nil
	}
	af := t.getAF()
	return network.GetNodeAddr(nodename, af)
}

// getLocalIP returns the addr set in the network config.
// Defaults to the first resolved ip address with the network address family (ip4 or ip6).
func (t T) getLocalIP() (net.IP, error) {
	return t.getNodeIP(hostname.Hostname())
}

// getAF returns the network address family (ip4 or ip6).
func (t T) getAF() (af string) {
	if t.IsIP6() {
		af = "ip6"
	} else {
		af = "ip4"
	}
	return
}

func (t *T) setupNodeRoutes(route network.Route) error {
	for _, table := range t.Tables() {
		route.Table = table
		t.Log().Info().Msgf("route add %s", route)
		if err := route.Add(); err != nil {
			return fmt.Errorf("route add %s: %w", route, err)
		}
	}
	return nil
}
