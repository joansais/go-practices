package wiki

import (
	"html/template"
)

type PageModel struct {
	Id			PageId
	Title		string
	BodyToEdit	string
	BodyAsHtml	template.HTML
}

type PageListModel []*PageModel

func (list PageListModel) Len() int           { return len(list) }
func (list PageListModel) Less(i, j int) bool { return list[i].Title < list[j].Title }
func (list PageListModel) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }

