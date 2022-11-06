package cmds

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevinlisr/gokpdep/pkg/cache"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"os"
	"sigs.k8s.io/yaml"
)

type depjson struct {
	title string
	path  string
}

type depmodel struct {
	items   []*depjson
	index   int
	cmd     *cobra.Command
	depName string
	ns      string
}

func (m depmodel) Init() tea.Cmd {
	return nil
}
func (m depmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.index > 0 {
				m.index--
			}
		case "down":
			if m.index < len(m.items)-1 {
				m.index++
			}
		case "enter":
			getdepDetailByJSON(m.depName, m.items[m.index].path, m.ns, m.cmd)
			return m, tea.Quit
		}

	}
	return m, nil
}
func (m depmodel) View() string {
	s := "welcome to K8S Visualization system!\n"
	for i, item := range m.items {
		selected := " "
		if m.index == i {
			selected = ">>"
		}
		s += fmt.Sprintf("%s %s\n", selected, item.title)
	}
	s += "\nEnter Q logout\n"
	return s
}

const (
	depEventType = "__event__"
	//depLogType   = "__log__"
)

func runtea(args []string, cmd *cobra.Command, ns string) {
	if len(args) == 0 {
		log.Println("dep name is required!")
		return
	}
	var depModel = depmodel{
		items:   []*depjson{},
		cmd:     cmd,
		depName: args[0],
		ns:      ns,
	}
	//v1.dep{}
	depModel.items = append(depModel.items,
		&depjson{title: "Meta Info", path: "metadata"},
		&depjson{title: "Labels", path: "metadata.labels"},
		&depjson{title: "Annotations", path: "metadata.annotations"},
		&depjson{title: "selector", path: "spec.selector"},
		&depjson{title: "dep model", path: "spec.template"},
		&depjson{title: "All Info", path: "@this"},
		&depjson{title: "status", path: "status"},
		&depjson{title: "*Events*", path: depEventType},
		//&depjson{title: "*Logs*", path: depLogType},
	)
	teaCmd := tea.NewProgram(depModel)
	if err := teaCmd.Start(); err != nil {
		fmt.Println("Start failed:", err)
		os.Exit(1)
	}
}

var eventHeaders = []string{"EVENT", "REASON", "RESOURCE", "MESSAGES"}

func printEvent(events []*v1.Event)  {
	table := tablewriter.NewWriter(os.Stdout)
	// set header
	table.SetHeader(eventHeaders)
	for _,e := range events {
		depRow := []string{e.Type, e.Reason,
			fmt.Sprintf("%s/%s",e.InvolvedObject.Kind,e.InvolvedObject.Name),e.Message}
		table.Append(depRow)
	}
	//setTable(table)
	// jin xing xuan ran
	table.Render()
}


func getdepDetailByJSON(depName, path, nameSpace string, cmd *cobra.Command) {
	//ns, err := cmd.Flags().GetString("namespace")
	//if err != nil {
	//	log.Println("error ns param")
	//	return
	//}
	if nameSpace == "" {
		nameSpace = "default"
	}
	fmt.Println("namespace is :",nameSpace)
	dep, err := cache.Factory.Apps().V1().Deployments().Lister().Deployments(nameSpace).Get(depName)
	fmt.Println("Already huo qu DEPLOYMENT")
	if err != nil {
		fmt.Println("Already huo qu DEPLOYMENT")
		log.Println(err)
		return
	}

	// get resource events
	if path == depEventType {
		eventList, err := cache.Factory.Core().V1().Events().Lister().List(labels.Everything())
		if err != nil {
			log.Println(err)
			return
		}
		depEvents := []*v1.Event{}
		for _, e := range eventList {
			if e.InvolvedObject.UID == dep.UID {
				depEvents = append(depEvents, e)
			}
		}
		printEvent(depEvents)
		// dao zhe jiu bu xu yao wang xia zou le ,zhi jie return
		return
	}

	// get resource logs
	//if path == depLogType {
	//	req := cache.Client.AppsV1().Deployments(nameSpace).GetLogs(dep.Name, &v1.depLogOptions{})
	//	//Pod, err := client.CoreV1().Pods(nameSpace).Get(context.Background(), pod.Name, metav1.GetOptions{})
	//	//fmt.Println("huo qu pod",Pod)
	//
	//	ret := req.Do(context.Background())
	//	fmt.Println(ret.Get())
	//	b, err := ret.Raw()
	//	if err != nil {
	//		fmt.Println("get logs error")
	//		log.Println(err)
	//		return
	//	}
	//	fmt.Println(string(b))
	//	return
	//
	//}

	jsonStr, _ := json.Marshal(dep)
	ret := gjson.Get(string(jsonStr), path)
	if !ret.Exists() {
		log.Println("No corresponding content was found " + path)
		return
	}
	if !ret.IsObject() && !ret.IsArray() {
		fmt.Println(ret.Raw)
		return
	}
	tempMap := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(ret.Raw), &tempMap)

	if err != nil {
		log.Println(err)
		return
	}

	b, _ := yaml.Marshal(tempMap)
	fmt.Println(string(b))
}
