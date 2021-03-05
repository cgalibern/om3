/*
Copyright © 2021 OPENSVC SAS <contact@opensvc.com>

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

// Package cmd defines the opensvc command line actions and options.
package cmd

import (
	"github.com/spf13/cobra"
	"opensvc.com/opensvc/core/entrypoints"
)

// nodeEventsCmd represents the nodeEvents command
var nodeEventsCmd = &cobra.Command{
	Use:     "events",
	Short:   "Print the node event stream",
	Aliases: []string{"eve", "even", "event"},
	Run:     nodeEventsCmdRun,
}

func init() {
	nodeCmd.AddCommand(nodeEventsCmd)
}

func nodeEventsCmdRun(cmd *cobra.Command, args []string) {
	e := entrypoints.Events{
		Format: formatFlag,
		Color:  colorFlag,
		Server: serverFlag,
	}
	e.Do()
}
