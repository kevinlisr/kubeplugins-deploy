package utils

import (
	"github.com/spf13/cobra"
	"log"
)

func GetNameSpace(cmd *cobra.Command) string {
	ns, err := cmd.Flags().GetString("namespace")
	if err != nil {
		log.Println(err)
	}

	return ns
}
