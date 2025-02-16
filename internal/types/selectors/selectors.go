package selectors

type indexPage struct {
	BaseHTML string
}

type usersView struct {
	UsersView string
}

var IndexPage = indexPage{
	BaseHTML: "base-html",
}

var UsersView = usersView{
	UsersView: "base-users-view",
}
