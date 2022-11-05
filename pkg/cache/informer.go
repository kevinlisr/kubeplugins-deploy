package cache

import (
	"flag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
)



var ShowLables bool



var Cache bool
var Namespace string

// chu shi hua ke hu duan
var Client = NewK8sConfig().InitClient()
var restConfig = NewK8sConfig().K8sRestConfig()
var CfgFlags *genericclioptions.ConfigFlags

var err error
var Factory informers.SharedInformerFactory


type K8sConfig struct {
}

func NewK8sConfig() *K8sConfig {
	return &K8sConfig{}
}

func (*K8sConfig) K8sRestConfig() *rest.Config {
	// 使用当前上下文环境
	var cliKubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		cliKubeconfig = flag.String("cliKubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		cliKubeconfig = flag.String("cliKubeconfig", "", "absolute path to the kubeconfig file")

	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *cliKubeconfig)
	if err != nil {
		panic(err.Error())
	}

	return config
}

func (this *K8sConfig) InitClient() *kubernetes.Clientset {
	CfgFlags = genericclioptions.NewConfigFlags(true)
	config, err := CfgFlags.ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		log.Fatalln(err)
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}
	return c
}

type PodHandler struct {
}

func (this *PodHandler) OnAdd(obj interface{}) {}

func (this *PodHandler) OnUpdate(oldObj, newObj interface{}) {}

func (this *PodHandler) OnDelete(obj interface{}) {}



func InitCache() {
	Factory = informers.NewSharedInformerFactory(Client, 0)
	Factory.Core().V1().Pods().Informer().AddEventHandler(&PodHandler{})
	Factory.Core().V1().Events().Informer().AddEventHandler(&PodHandler{}) //wei le tou lan
	Factory.Apps().V1().Deployments().Informer().AddEventHandler(&PodHandler{})
	ch := make(chan struct{})
	Factory.Start(ch)
	Factory.WaitForCacheSync(ch)
}


//
//func run(c *cobra.Command, args []string) error {
//	//Client := NewK8sConfig().InitClient()
//	Ns, err = c.Flags().GetString("namespace")
//	fmt.Println(Ns)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	if Ns == "" {
//		Ns = "default"
//	}
//
//	list, err := Client.CoreV1().Pods(Ns).List(context.Background(),
//		v1.ListOptions{
//			LabelSelector: Labels,
//			FieldSelector: fields,
//		})
//	if err != nil {
//		log.Fatalln(err, "huo qu pod list")
//	}
//
//	//for _, p := range list.Items{
//	//	podsJson,_ := json.Marshal(p)
//	//}
//
//	//podsJson,_ := json.Marshal(list)
//
//	//err = WriteFile("pods.json", []byte(podsJson), 0666)
//	//if err != nil {
//	//	log.Fatalln(err)
//	//}
//
//	table := tablewriter.NewWriter(os.Stdout)
//	commonHeaders := []string{"Name", "Namespace", "Ip", "Phase"}
//
//	if ShowLables {
//		commonHeaders = append(commonHeaders, "tag")
//	}
//
//	table.SetHeader(commonHeaders)
//
//	for _, pod := range list.Items {
//		//fmt.Println(pod.Name)
//		p, err := json.Marshal(pod)
//		if err != nil {
//			log.Fatalln(err)
//		}
//		ret := gjson.Get(string(p), "metadata.name")
//
//		var podRow []string
//		if m, err := regexp.MatchString(PodName, ret.String()); err == nil && m {
//			podRow = []string{pod.Name, pod.Namespace, pod.Status.PodIP, string(pod.Status.Phase)}
//
//		}
//
//		//podRow = []string{pod.Name,pod.Namespace,pod.Status.PodIP,string(pod.Status.Phase)}
//		//if ShowLables {
//		//	podRow = append(podRow, Map2String(pod.Labels))
//		//}
//		table.Append(podRow)
//	}
//	table.SetAutoWrapText(false)
//	table.SetAutoFormatHeaders(true)
//	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
//	table.SetAlignment(tablewriter.ALIGN_LEFT)
//	table.SetCenterSeparator("")
//	table.SetColumnSeparator("")
//	table.SetRowSeparator("")
//	table.SetHeaderLine(false)
//	table.SetBorder(false)
//	table.SetTablePadding("\t") // pad with tabs
//	table.SetNoWhiteSpace(true)
//	table.Render()
//	return nil
//}

// 通用的文件打开函数(综合和 Create 和 Open的作用)
// OpenFile第二个参数 flag 有如下可选项
//    O_RDONLY  文件以只读模式打开
//    O_WRONLY  文件以只写模式打开
//    O_RDWR   文件以读写模式打开
//    O_APPEND 追加写入
//    O_CREATE 文件不存在时创建
//    O_EXCL   和 O_CREATE 配合使用,创建的文件必须不存在
//    O_SYNC   开启同步 I/O
//    O_TRUNC  打开时截断常规可写文件
//func WriteFile(filename string, data []byte, perm os.FileMode) error {
//	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
//	if err != nil {
//		return err
//	}
//	n, err := f.Write(data)
//	if err == nil && n < len(data) {
//		err = io.ErrShortWrite
//	}
//	if err1 := f.Close(); err == nil {
//		err = err1
//	}
//	return err
//}

