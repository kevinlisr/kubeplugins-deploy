package cmds

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/kevinlisr/gokpdep/pkg/cache"
	"github.com/kevinlisr/gokpdep/pkg/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"os"
	"regexp"
	"sort"
)

type V1Deployment []*appv1.Deployment
var depName string
var ShowLables bool

func (this V1Deployment) Len() int {

	int := len(this)
	return int
}
func (this V1Deployment) Less(i,j int) bool{
	// gen ju shi jian pai xu    dao pei xu
	return this[i].CreationTimestamp.Time.After(this[j].CreationTimestamp.Time)
}

func (this V1Deployment) Swap(i,j int)  {
	// gen ju shi jian pai xu    dao pei xu
	return
}

// qu chu deploy lie biao
func listDeploys(ns string) []*appv1.Deployment {
	list,err := cache.Factory.Apps().V1().Deployments().
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
	if depList == nil{
		return
	}
	for _,dep := range depList{
		ret = append(ret,prompt.Suggest{
			Text: dep.Name,
			Description: fmt.Sprintf("Replicas:%d/%d",dep.Status.Replicas,
				dep.Status.Replicas),
		})
	}
	return
}

// xuan ran deploy lie biao
func RenderDeploy(args []string, cmd *cobra.Command) {
	ingressList := listDeploys(utils.GetNameSpace(cmd))
	if ingressList == nil{
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	commonHeaders := []string{"Name", "Namespace", "Ip","Phase","hello"}

	table.SetHeader(commonHeaders)

	for _,ingress :=range ingressList{
		//fmt.Println(pod.Name)
		p, err := json.Marshal(ingress)
		if err != nil {
			log.Fatalln(err)
		}
		ret := gjson.Get(string(p), "metadata.name")

		var podRow  []string
		if m,err := regexp.MatchString(depName,ret.String());err == nil && m {
			podRow = []string{ingress.Name}

		}

		//podRow = []string{pod.Name,pod.Namespace,pod.Status.PodIP,string(pod.Status.Phase)}

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
	return
}