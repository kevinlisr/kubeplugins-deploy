package cmds

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/kevinlisr/gokpdep/pkg/cache"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"os"
	"sort"
)

type V1Deployment []*appv1.Deployment
type v1Events []*v1.Event

var depName string
var ShowLables bool
var Ns string


func (this V1Deployment) Len() int {

	int := len(this)
	return int
}
func (this V1Deployment) Less(i, j int) bool {
	// gen ju shi jian pai xu    dao pei xu
	return this[i].CreationTimestamp.Time.After(this[j].CreationTimestamp.Time)
}

func (this V1Deployment) Swap(i, j int) {
	// gen ju shi jian pai xu    dao pei xu
	return
}

func (this v1Events) Len() int {

	int := len(this)
	return int
}
func (this v1Events) Less(i, j int) bool {
	// gen ju shi jian pai xu    dao pei xu
	return this[i].CreationTimestamp.Time.After(this[j].CreationTimestamp.Time)
}

func (this v1Events) Swap(i, j int) {
	// gen ju shi jian pai xu    dao pei xu
	return
}

func getLatestDeployEvent(uid types.UID, ns string) string {
	list, err := cache.Factory.Core().V1().Events().Lister().Events(ns).
		List(labels.Everything())
	if err != nil {
		return ""
	}
	sort.Sort(v1Events(list))
	for _, e := range list {
		if e.InvolvedObject.UID == uid {
			return e.Message
		}
	}
	return ""
}

// qu chu deploy lie biao
func listDeploys(ns string) []*appv1.Deployment {
	list, err := cache.Factory.Apps().V1().Deployments().
		Lister().Deployments(ns).List(labels.Everything())
	if err != nil {
		log.Println(err)
		return nil
	}
	sort.Sort(V1Deployment(list))
	return list
}

// yong yu ti shi   yong
func RecommendDeployments(ns string) (ret []prompt.Suggest) {
	depList := listDeploys(ns)
	//fmt.Printf("",depList)
	if depList == nil {
		return
	}
	for _, dep := range depList {
		ret = append(ret, prompt.Suggest{
			Text: dep.Name,
			Description: fmt.Sprintf("Replicas:%d/%d", dep.Status.Replicas,
				dep.Status.Replicas),
		})
	}
	return
}

// xuan ran deploy lie biao
func RenderDeploy(args []string, cmd *cobra.Command) {
	deplist := listDeploys(GetNameSpace(cmd))
	if deplist == nil {
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	commonHeaders := []string{"Name", "Replicas", "CreateTime", "NewEvents"}

	table.SetHeader(commonHeaders)

	for _, dep := range deplist {
		//fmt.Println(pod.Name)
		var depRow []string
		depRow = []string{dep.Name,
			fmt.Sprintf("%d/%d", dep.Status.ReadyReplicas, dep.Status.Replicas),
			dep.CreationTimestamp.Format("2006-01-02 15:04:05"),
			getLatestDeployEvent(dep.UID, Ns),
		}
		//podRow = []string{pod.Name,pod.Namespace,pod.Status.PodIP,string(pod.Status.Phase)}
		table.Append(depRow)
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
	return
}
