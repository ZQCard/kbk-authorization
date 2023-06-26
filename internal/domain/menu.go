package domain

type MenuBtn struct {
	Id          int64
	MenuId      int64
	Name        string
	Description string
	Identifier  string
	CreatedAt   string
	UpdatedAt   string
}

type Menu struct {
	Id        int64
	Name      string
	Path      string
	ParentId  int64
	ParentIds string
	Hidden    bool
	Sort      int64
	Component string
	Title     string
	Icon      string
	CreatedAt string
	UpdatedAt string
	MenuBtns  []MenuBtn
	Children  []Menu
}
