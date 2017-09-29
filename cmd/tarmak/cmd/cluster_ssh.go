package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jetstack/tarmak/pkg/tarmak"
)

var clusterSshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "ssh into instance",
	Run: func(cmd *cobra.Command, args []string) {
		t := tarmak.New(cmd)
		t.SSHPassThrough(args)
	},
}

func init() {
	clusterCmd.AddCommand(clusterSshCmd)
}
