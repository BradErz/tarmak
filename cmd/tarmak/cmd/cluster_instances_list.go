package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/jetstack/tarmak/pkg/tarmak"
)

var clusterInstancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "list instances of the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		t := tarmak.New(cmd)
		hosts, err := t.Context().Environment().Provider().ListHosts()
		if err != nil {
			logrus.Fatal(err)
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)

		for _, host := range hosts {
			fmt.Fprintf(
				w,
				"%s\t%s\t%s\n",
				host.ID(),
				host.Hostname(),
				strings.Join(host.Roles(), ", "),
			)
		}
		w.Flush()
	},
}

func init() {
	clusterInstancesCmd.AddCommand(clusterInstancesListCmd)
}
