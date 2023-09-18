package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cmdPool = &cobra.Command{
		Use:   "pool",
		Short: "Manage storage pools",
		Long:  ` A pool is a vol provider. Pools abstract the hardware and software specificities of the cluster infrastructure.`,
	}
	cmdPoolVolume = &cobra.Command{
		Use:   "volume",
		Short: "Manage storage pool volumes",
	}
)

func init() {
	root.AddCommand(
		cmdPool,
	)
	cmdPool.AddCommand(
		cmdPoolVolume,
		newCmdPoolLs(),
	)
	cmdPoolVolume.AddCommand(
		newCmdPoolVolumeLs(),
	)
}
