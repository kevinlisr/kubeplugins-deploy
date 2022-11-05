package cmds

import (
	"fmt"
	"github.com/kevinlisr/gokpdep/pkg/cache"
	"github.com/spf13/cobra"
	"log"
)

func MergeFlags(cmds ...*cobra.Command) {
	for _,cmd := range cmds{
		cache.CfgFlags.AddFlags(cmd.Flags())
	}
}

func RunCmd() {
	cmd := &cobra.Command{
		Use:          "kubectl ingress prompt",
		Short:        "list ingress",
		Example:      "kubectl ingress prompt",
		SilenceUsage: true,
	}
	cache.InitCache()

	MergeFlags(cmd, PromptCmd)

	// jia ru zi ming ling
	cmd.AddCommand(PromptCmd)
	err := cmd.Execute()
	fmt.Println("stop exec  cmd")
	if err != nil {
		log.Fatalln(err, "exec bao cuo")
	}
}