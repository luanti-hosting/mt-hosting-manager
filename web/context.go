package web

import (
	"embed"
	"mt-hosting-manager/db"
	"mt-hosting-manager/tmpl"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/oauth"
	"mt-hosting-manager/web/usernode"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

//go:embed *
var Files embed.FS

type Context struct {
	tu          *tmpl.TemplateUtil
	repos       *db.Repositories
	GithubOauth *oauth.OAuthConfig
}

func (ctx *Context) Setup(r *mux.Router) {
	r.HandleFunc("/", ctx.Index)

	r.HandleFunc("/login", ctx.tu.OptionalSecure(ctx.Login))
	r.HandleFunc("/profile", ctx.tu.Secure(ctx.Profile))

	r.HandleFunc("/node_types", ctx.tu.Secure(ctx.NodeTypes, tmpl.RoleCheck(types.UserRoleAdmin))).Methods(http.MethodGet)
	r.HandleFunc("/node_types/{id}", ctx.tu.Secure(ctx.NodeTypeEdit, tmpl.RoleCheck(types.UserRoleAdmin)))

	usernode_ctx := usernode.New(ctx.tu, ctx.repos)
	usernode_ctx.Setup(r)

	r.PathPrefix("/assets").Handler(statigz.FileServer(Files, brotli.AddEncoding))

	if os.Getenv("GITHUB_CLIENTID") != "" {
		ctx.GithubOauth = &oauth.OAuthConfig{
			ClientID: os.Getenv("GITHUB_CLIENTID"),
			Secret:   os.Getenv("GITHUB_SECRET"),
		}

		r.Handle("/oauth_callback/github", &oauth.OauthHandler{
			Impl:     &oauth.GithubOauth{},
			BaseURL:  ctx.tu.BaseURL,
			Type:     types.UserTypeGithub,
			Tu:       ctx.tu,
			UserRepo: ctx.repos.UserRepo,
			Config:   ctx.GithubOauth,
		})
	}

	r.NotFoundHandler = ctx.NotFound()
}
