package ox

import (
	"github.com/spf13/cobra"
)

var (
	cmdNetwork = &cobra.Command{
		Use:     "network",
		Short:   "Manage backend networks",
		Aliases: []string{"net"},
		Long:    `A backend network provides ip addresses to svc objects via ip.cni resources. These addresses are automatically allocated, accessible from all cluster nodes, and resolved by the cluster dns.`,
	}
	cmdNetworkIP = &cobra.Command{
		Use:   "ip",
		Short: "Manage ip on backend networks",
	}
)

func init() {
	root.AddCommand(
		cmdNetwork,
	)
	cmdNetwork.AddCommand(
		cmdNetworkIP,
		newCmdNetworkLs(),
		newCmdNetworkSetup(),
	)
	cmdNetworkIP.AddCommand(
		newCmdNetworkIPLs(),
	)
}
