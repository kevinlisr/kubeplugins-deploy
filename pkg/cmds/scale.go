package cmds

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevinlisr/gokpdep/pkg/cache"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"regexp"
	"strconv"
)

func ScaleDeploy(args []string, cmd *cobra.Command) {
	if len(args) == 0{
		fmt.Println("deploy name is required!")
		return
	}

	p := tea.NewProgram(initialScaleModel(args[0],Ns))
	if err := p.Start(); err != nil {
		log.Fatalln(err)
	}
}


type (
	errMsg error
)

type scaleModel struct {
	textInput textinput.Model
	err       error
	depname string
	ns string
}

func checkScale(v string) bool {
	if regexp.MustCompile("^([0-9]|1[0-9]|20)$").MatchString(v){
		return true
	}
	return false
}

func initialScaleModel(depname,ns string) scaleModel {
	ti := textinput.New()
	ti.Placeholder = "between 0-20"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return scaleModel{
		textInput: ti,
		err:       nil,
		depname: depname,
		ns: ns,
	}
}

func (m scaleModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m scaleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case  tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			numStr := m.textInput.Value()
			if checkScale(numStr){
				num, _ := strconv.Atoi(numStr)
				scale, err := cache.Client.AppsV1().Deployments(m.ns).GetScale(context.Background(),m.depname,metav1.GetOptions{})
				if err != nil {
					fmt.Println(err)
					return m, tea.Quit
				}
				scale.Spec.Replicas = int32(num)
				_, err = cache.Client.AppsV1().Deployments(m.ns).UpdateScale(context.Background(), m.depname, scale, metav1.UpdateOptions{})
				if err != nil {
					fmt.Println(err)
					return m, tea.Quit
				}
				fmt.Println("will modify replica is: ",num)
				fmt.Println("replicas scale successed!")

			}else {
				fmt.Println("scale value must between 0-20")
			}
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m scaleModel) View() string {
	return fmt.Sprintf(
		"please enter replaces Num(0-20): \n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}