package cmd

import (
	"fmt"
	"strings"

	"github.com/opensvc/om3/core/array"
	"github.com/opensvc/om3/core/object"
	"github.com/opensvc/om3/util/key"
	"github.com/spf13/cobra"
)

var (
	arrayName string
	cmdArray  = &cobra.Command{
		Use:   "array",
		Short: "Manage storage arrays",
		Long:  ` A array is backend storage provider for pools.`,
		RunE: func(_ *cobra.Command, args []string) error {
			return runArray(args)
		},
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
	}
)

func init() {
	root.AddCommand(
		cmdArray,
	)
	cmdArray.AddCommand(
		newCmdArrayLs(),
	)
	cmdArray.PersistentFlags().StringVar(&arrayName, "array", "", "the section name or index identifying the array")
}

func runArray(args []string) error {
	o, err := object.NewCluster(object.WithVolatile(true))
	if err != nil {
		return err
	}
	if !strings.HasPrefix(arrayName, "array#") {
		arrayName = "array#" + arrayName
	}
	if !o.Config().HasSectionString(arrayName) {
		return fmt.Errorf("no section found matching %s in the cluster config", arrayName)
	}
	arrayType, err := o.Config().GetStringStrict(key.T{arrayName, "type"})
	if err != nil {
		return err
	}
	drv := array.GetDriver(arrayType)
	if drv == nil {
		return fmt.Errorf("no array driver found matching type %s", arrayType)
	}
	drv.SetName(arrayName)
	drv.SetConfig(o.Config())
	return drv.Run(args)
}
