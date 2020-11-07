package model

import (
	"fmt"
	"time"
)

type Branch struct {
	Idbranch  int
	Name      string
	Idcompany int
	Themes    []*Theme
	Username  string
}

func (b *Branch) String() string {
	str := b.Name + "\n"
	for i, _ := range b.Themes {
		temp := b.Themes[i].String()
		str += temp + "\n"
		fmt.Println(temp)
	}
	return str
}

type Theme struct {
	Idtheme  int
	Name     string
	Idbranch int
	Index    int
	Tasks    []*Task
}

func (t *Theme) String() string {
	str := t.Name + "\n"
	for i, _ := range t.Tasks {
		temp := t.Tasks[i].String()
		str += temp + "\n"
		fmt.Println(temp)
	}
	return str
}

type Task struct {
	Idtask  int
	Name    string
	Text    string
	Idtheme int
	Answer  *Answer
}

func (t *Task) String() string {
	str := t.Name + "\n"
	if t.Answer != nil {
		str += t.Answer.String()
	}
	return str
}

type Answer struct {
	Idanswer int
	File     string
	Idsender int
	Idtask   int
	Status   int
	Date     time.Time
}

func (a *Answer) String() string {
	str := a.File + "\n"
	fmt.Print(str)
	return str
}
