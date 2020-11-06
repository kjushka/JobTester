package model

type Project struct {
	Id int
	Name string
}

func (p *Project) String() string {
	return p.Name
}
