package wiki

type PageId string

type Page struct {
	Id		PageId
	Title	string
	Body	string
}
