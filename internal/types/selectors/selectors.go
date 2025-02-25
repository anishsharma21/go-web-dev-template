package selectors

type indexPage struct {
	IndexHtml string
	IndexBody string
}

var IndexPage = indexPage{
	IndexHtml: "index-html",
	IndexBody: "index-body",
}

type loginView struct {
	LoginView string
}

var LoginView = loginView{
	LoginView: "base-login-view",
}

type usersView struct {
	UsersView string
}

var UsersView = usersView{
	UsersView: "base-users-view",
}
