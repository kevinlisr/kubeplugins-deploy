package cmds

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/kevinlisr/gokpdep/pkg/cache"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/informers"
	"log"
	"os"
	"regexp"
	"strings"
)

var fields string
var PodName string
var Ns string



func setNamespace(c *cobra.Command,a []string) string {
	Ns, err = c.Flags().GetString("namespace")
	if err != nil {
		log.Fatalln(err)
	}

	if Ns == ""{

		fmt.Printf("Ns is null set Ns == %s", a[1])
		Ns = a[1]
	}
	return Ns
}

func clearConsole()  {
	MyConsoleWriter.EraseScreen()
	MyConsoleWriter.CursorGoTo(0,0)
	MyConsoleWriter.Flush()
}



func executorCmd(cmd *cobra.Command,ns *string) func(in string) {

	return func(in string) {
		in = strings.TrimSpace(in)
		blocks := strings.Split(in, " ")
		args := []string{}
		if len(blocks) > 1{
			args=blocks[1:]
		}

		switch blocks[0] {
		case "exit":
			fmt.Print("Bye!")
			os.Exit(0)
		case "list":
			//getPodDetail(args)
			RenderDeploy(args,cmd)
		case "set":
			fmt.Printf("ns is null set ns == %s", args[1])
			*ns = args[1]
		case "clear":
			clearConsole()

		//case "del":
		//	delPod(args,cmd)
		//case "ns":
		//	showNameSpace(cmd)
		}

	}
}
var suggestions = []prompt.Suggest{
	{"get","GET "},
	{"list","LIST"},
	{"exit","EXIT the interactive window"},
	{"exec","exec the container"},
}

//var suggestions = []prompt.Suggest{
//	// Command
//	{"exec", "pod shell cao zuo "},
//	{"get","get pod xiang xi xin xi"},
//	{"use", "she zhi dang qian ming ming kong jian"},
//	{"del", "shan chu mou ge pod"},
//	{"list","xian shi pod lie biao"},
//	{"clear", "qing chu ping mu"},
//	{"exit","tui chu jiao hu shi"},
//}

var podSuggestions = []prompt.Suggest{
	{"get","get pods details"},
	{"list","show pods list"},
	{"exit","exit the interactive window"},
	{"exec","exec the container"},
}



func parseCmd(w string) (string, string) {
	w = regexp.MustCompile("\\s+").ReplaceAllString(w, " ")
	l := strings.Split(w," ")
	if len(l)>= 2{
		return l[0],strings.Join(l[1:]," ")
	}
	return w,""
}

//func completer(in prompt.Document) []prompt.Suggest {
//	w := in.GetWordBeforeCursor()
//	if w == ""{
//		return []prompt.Suggest{}
//	}
//
//	cmd, opt := parseCmd(in.TextBeforeCursor())
//	//fmt.Println(cmd)
//	if cmd == "get"{
//		return prompt.FilterHasPrefix(getPodsList(ns),opt,true)
//	}
//
//	return prompt.FilterHasPrefix(suggestions,w,true)
//}

func completerCmd(ns *string) func (in prompt.Document) []prompt.Suggest {
	return func (in prompt.Document) []prompt.Suggest {
		w := in.GetWordBeforeCursor()
		if w == ""{
			return []prompt.Suggest{}
		}
		cmd, opt := parseCmd(in.TextBeforeCursor())
		//fmt.Println(cmd)

		//if inArray([]string{"get","del","exec"},cmd){
		//	return prompt.FilterHasPrefix(getPodsList(ns),opt,true)
		//}

		if cmd == "list"{
			return prompt.FilterHasPrefix(RecommendDeployments(*ns),opt,true)
		}


		return prompt.FilterHasPrefix(suggestions,w,true)
	}
}
//var ns string
var err error
var MyConsoleWriter = prompt.NewStderrWriter()
var Labels string
var Factory informers.SharedInformerFactory

var PromptCmd = &cobra.Command{
	Use: "prompt",
	Short: "prompt pods",
	Example: "kubectl pods prompt [flags]",
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {

		Ns, err = c.Flags().GetString("namespace")
		fmt.Println("promptCmd  get namespace is:", Ns)
		if err != nil {
			log.Fatalln(err)
		}

		if len(args) > 1{
			Ns = args[1]
		}

		if Ns == ""{
			Ns = "default"
		}

		cache.InitCache()
		p := prompt.New(
			executorCmd(c,&Ns),
			completerCmd(&Ns),
			prompt.OptionPrefix(">>>"),
			// she zhi  "clear" ming ling lai qing ping
			prompt.OptionWriter(MyConsoleWriter),

		)
		p.Run()
		return nil
	},
}

var cacheCmd = &cobra.Command{
	Use:    "cache",
	Short:  "deps by cache",
	Hidden: true,
	RunE: func(c *cobra.Command, args []string) error {
		//ns, err = c.Flags().GetString("namespace")
		//fmt.Println("cacheCmd  get namespace is:", ns)
		//if err != nil {
		//	log.Fatalln(err)
		//}

		//arg := fmt.Sprintf("%s", args[1])
		//if Ns == ""{Ns=args[1]}
		if len(args) != 0 {
			fmt.Println("args[1] fu zhi")
			Ns = args[1]
			fmt.Printf("++++++++++++++++%s\n", args[1])
		} else {
			if Ns == "" {
				Ns = "default"
			}
			//Ns="default"
		}

		pods, err := Factory.Core().V1().Pods().Lister().Pods(Ns).
			List(labels.Everything())
		if err != nil {
			return err
		}
		fmt.Println("cong huan cun zhong qu")
		table := tablewriter.NewWriter(os.Stdout)

		commonHeaders := []string{"Name", "Namespace", "Ip", "Phase","hello"}
		//
		//if ShowLables{
		//	commonHeaders = append(commonHeaders,"tag")
		//}

		table.SetHeader(commonHeaders)

		for _, pod := range pods {
			//fmt.Println(pod.Name)
			p, err := json.Marshal(pod)
			if err != nil {
				log.Fatalln(err)
			}
			ret := gjson.Get(string(p), "metadata.name")

			var podRow []string
			if m, err := regexp.MatchString(PodName, ret.String()); err == nil && m {
				podRow = []string{pod.Name, pod.Namespace, pod.Status.PodIP, string(pod.Status.Phase)}

			}

			table.Append(podRow)
		}
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t") // pad with tabs
		table.SetNoWhiteSpace(true)
		table.Render()
		return nil
	},
}


//func MergeFlags(cmd,  prompt, cacheCmd *cobra.Command) {
//	cmd.Flags().StringP("namespace", "n", "", "kubectl pods --namespace=\"kube-system\"")
//	//cmd.Flags().Bool("show-labels",false,"kubectl pods --show-labels")
//	cmd.Flags().BoolVar(&ShowLables, "show-labels", false, "kubectl pods --show-labels")
//	cmd.Flags().StringVar(&Labels, "labels", "",
//		"kubectl pods --labels app=ngx or kubectl pods --labels=\"app=ngx,version=v1\"")
//	cmd.Flags().StringVar(&fields, "fields", "",
//		"kubectl pods --fields=\"status.phase=Running\"")
//	cmd.Flags().StringVar(&PodName, "name", "",
//		"kubectl pods --name=\"^ng\"")
//	////
//	//listCmd.Flags().StringP("namespace", "n", "", "kubectl pods --namespace=\"kube-system\"")
//	////cmd.Flags().Bool("show-labels",false,"kubectl pods --show-labels")
//	//listCmd.Flags().BoolVar(&ShowLables, "show-labels", false, "kubectl pods --show-labels")
//	//listCmd.Flags().StringVar(&Labels, "labels", "",
//	//	"kubectl pods --labels app=ngx or kubectl pods --labels=\"app=ngx,version=v1\"")
//	//listCmd.Flags().StringVar(&fields, "fields", "",
//	//	"kubectl pods --fields=\"status.phase=Running\"")
//	//listCmd.Flags().StringVar(&PodName, "name", "",
//	//	"kubectl pods --name=\"^ng\"")
//	prompt.Flags().StringP("namespace", "n", "", "kubectl pods --namespace=\"kube-system\"")
//	//cache.Flags().StringP("namespace", "n", "", "kubectl pods --namespace=\"kube-system\"")
//	cacheCmd.Flags().StringP("namespace", "n", "", "kubectl pods --namespace=\"kube-system\"")
//
//}
//
//func RunCmd() {
//	cmd := &cobra.Command{
//		Use:          "kubectl deps [flags]",
//		Short:        "list pods",
//		Example:      "kubectl deps [flags]",
//		SilenceUsage: true,
//		//RunE:         run,
//	}
//	MergeFlags(cmd, cmds.PromptCmd, cacheCmd)
//
//	//添加参数
//	//BoolVar用来支持 是否
//	//cmd.Flags().BoolVar(&ShowLabels, "show-labels", false, "kubectl pods --show-labels")
//	//cmd.Flags().StringVar(&Labels, "labels", "", "kubectl pods --labels=\"app=nginx\"")
//
//	cmd.AddCommand( cmds.PromptCmd, cacheCmd)
//
//	err := cmd.Execute()
//	fmt.Println("stop exec  cmd")
//	if err != nil {
//		log.Fatalln(err, "exec bao cuo")
//	}
//}

func getNameSpace(c *cobra.Command) string {
	Ns, err = c.Flags().GetString("namespace")
	if Ns == "" {
		Ns = "plugiNs"
	}
	if err != nil {
		log.Println(err)

	}
	return Ns
}
