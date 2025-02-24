package selectors

type indexPage struct {
	BaseHTML string
}

var IndexPage = indexPage{
	BaseHTML: "base-html",
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
