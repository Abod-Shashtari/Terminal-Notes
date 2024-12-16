package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var(
	docStyle = lipgloss.NewStyle().Margin(1,2)
	choiceStyle = lipgloss.NewStyle().Padding(0,1).Background(lipgloss.Color("#2D2C68"))

	//First window=false ; Second window=true//
	windowx=false
	write=false
	textAreaFlag=false
	inputFlag=false
	setTextFlag=true
	createFileFlag=true
	//insert mode=false ; edit mode=true//
	textAreaMode=false

	path="C:\\Users\\abods\\OneDrive\\Documents\\Development\\Go\\Terminal-Notes\\text files\\"
	listFile="C:\\Users\\abods\\OneDrive\\Documents\\Development\\Go\\Terminal-Notes\\text files\\list\\list.txt"
)
//item sturcture
type item struct{
	title ,desc string
}
func (i item) Title() string{
	return i.title
}
func (i item) Description() string{
	return i.desc
}
func (i item) FilterValue() string{
	return i.title
}

//create list item
func addItem(m *model) item{
	return item{title: m.input.Value(),desc: ""}
}
func initList(m *model){
	content,err:=os.ReadFile(listFile)
	if err!=nil {
		print("ERROR")	
	}
	items:=strings.Split(string(content),",")
	//Reverse
	rev:=len(items)
	revItems:=make([]string,len(items))
	for i := 0; i < len(items); i++ {
		rev--
		revItems[i]=items[rev]
	}
	//create and insert items
	for i:=0;i<len(items)-1;i++{
		m.list.InsertItem(0,item{title: items[i],desc: ""})
	}
}

//Enable/DisableList
func edList(m *model,enable bool) {
	if enable==true {
		m.list.KeyMap.CursorUp.SetEnabled(true)
		m.list.KeyMap.CursorDown.SetEnabled(true)
		m.list.KeyMap.NextPage.SetEnabled(true)
		m.list.KeyMap.PrevPage.SetEnabled(true)
		m.list.KeyMap.Filter.SetEnabled(true)
	}else{
		m.list.KeyMap.CursorUp.SetEnabled(false)
		m.list.KeyMap.CursorDown.SetEnabled(false)
		m.list.KeyMap.NextPage.SetEnabled(false)
		m.list.KeyMap.PrevPage.SetEnabled(false)
		m.list.KeyMap.Filter.SetEnabled(false)
		
	}
}
//edit mode
func editMode(m* model,x bool){
		m.text.KeyMap.WordBackward.SetEnabled(x)
		m.text.KeyMap.WordForward.SetEnabled(x)
		m.text.KeyMap.DeleteWordBackward.SetEnabled(x)
}

//New keys (FOR HELP)
type listKeyMap struct {
	createFile    key.Binding
	quit   key.Binding
}
func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		createFile: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "Create Note"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c","Quit"),
		),
	}
}

//bubbleTea model
type model struct{
	//list
	list list.Model
	choice string
	input textinput.Model
	//text area
	text textarea.Model
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model,tea.Cmd){
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg:=msg.(type){
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m,tea.Quit
		case "t":
			if createFileFlag {
				createFileFlag=false
				
				if(m.list.FilterState()!=list.Filtering){
					if(windowx==false){
					m.input.SetValue("")
					m.input.Focus()
					write=true
					}

				}
			}
		case "enter":
			if write {
				write=false
				var itemVar item =addItem(&m)
				createFileFlag=true
				if(m.input.Value()!=""){
					m.list.InsertItem(0,itemVar)
					//save item in file
					listF,err:=os.OpenFile(listFile,os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
					if err!=nil {
						print("ERROR")
					}
					itemStringWithComma:=fmt.Sprintf("%s,",m.input.Value())
					listF.WriteString(itemStringWithComma)
					//create file to write in
					fileName:=fmt.Sprintf("%s%s.txt",path,m.input.Value())
					f, err:=os.Create(fileName)
					if err !=nil {
						print("ERROR")
					}
					defer f.Close()
				}
				inputFlag=false
			}else{
			i,ok:=m.list.SelectedItem().(item)
			if ok {
				m.choice=i.title
				windowx=true
				fileNameRead:=fmt.Sprintf("%s%s.txt",path,m.choice)
				content,err:= ioutil.ReadFile(fileNameRead)
				if err!=nil {
					print("ERROR")
				}
				if setTextFlag==true {	
					m.text.SetValue(string(content))
					setTextFlag=false
				}
			}
		}
		case "ctrl+s":
			data := []byte(m.text.Value())
			fileNameWrite:=fmt.Sprintf("%s%s.txt",path,m.choice)
			err:=os.WriteFile(fileNameWrite,data,0644)
			if err!=nil {
				print("ERROR")
			}
		case "esc":
			if(windowx==true){
			if(textAreaMode==false){
				textAreaMode=true
				editMode(&m,true)
				m.text.KeyMap.LineNext.SetKeys("j")
				m.text.KeyMap.LinePrevious.SetKeys("k")
			}else{
				textAreaMode=false
				editMode(&m,false)
				m.text.KeyMap.LineNext.SetKeys("down")
				m.text.KeyMap.LinePrevious.SetKeys("up")
			}
		}
		case "ctrl+q":
			setTextFlag=true
			windowx=false
			textAreaFlag=false
			m.choice=""
		}
	case tea.WindowSizeMsg:
		h,v:=docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h,msg.Height-v)
		m.text.SetWidth(msg.Width-5)
		m.text.SetHeight(msg.Height-7)
	}

	//fix pre-typing (text input)
	if inputFlag==true {
		m.input,cmd=m.input.Update(msg)
	}
	//fix pre-typing (text area)
	if(textAreaFlag==true){
		m.text,cmd=m.text.Update(msg)
		cmds=append(cmds, cmd)
	}
	//make list fixed while typing in text area
	if windowx==true ||write==true {
		edList(&m,false)
	}else{
		edList(&m,true)
	}

	m.list,cmd=m.list.Update(msg)
	return m, cmd
}
func (m model) View() string{
    var s string;
    mode:=""
    if textAreaMode==false {
        mode="INSER MODE"
    }else{
        mode="EDIT MODE"
    }
    if windowx==true {
        s=fmt.Sprintf("%s\n\n%s\n%s",docStyle.Render(choiceStyle.Render(m.choice)),docStyle.Render(m.text.View()),mode)
        m.list.KeyMap.CursorDown.SetEnabled(false)
        textAreaFlag=true
    }else{

        s=fmt.Sprintf("%s",docStyle.Render(m.list.View()))
        if write {
            s=fmt.Sprintf("%s\n%s",docStyle.Render(m.list.View()),m.input.View())
            inputFlag=true
        }
    }
    return s 
}

func main(){ 
	
	items:=[]list.Item{}
	listKeys:=newListKeyMap()

	//text area
	ti:=textarea.New();
	ti.Focus()
	ti.CharLimit=10000;
	ti.Placeholder="TYPE..."
	ti.KeyMap.WordBackward.SetKeys("b")
	ti.KeyMap.WordForward.SetKeys("w")
	ti.KeyMap.DeleteWordBackward.SetKeys("d")
	if(textAreaMode==true){
		ti.KeyMap.LineNext.SetKeys("j")
		ti.KeyMap.LinePrevious.SetKeys("k")
	}else{
		ti.KeyMap.LineNext.SetKeys("down")
		ti.KeyMap.LinePrevious.SetKeys("up")
	}

	//text input
	tin:=textinput.New()
	tin.Placeholder="TEXT..."
	
	m := model{list: list.New(items,list.NewDefaultDelegate(),0,0),text: ti,input: tin}
	//make edit mode off in initial
	editMode(&m,false)
	m.list.Title="Text-Doc List"
	m.list.DisableQuitKeybindings()
	m.list.AdditionalShortHelpKeys=func() []key.Binding{
		return []key.Binding{
			listKeys.createFile,
			listKeys.quit,
		}
	}
	initList(&m)
	//Run the program//
	tea.NewProgram(m,tea.WithAltScreen()).Start()
}
