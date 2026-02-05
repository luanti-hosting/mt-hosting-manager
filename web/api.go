package web

import (
	"fmt"
	"mt-hosting-manager/api/coinbase"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/core"
	"mt-hosting-manager/db"
	"mt-hosting-manager/public"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/middleware"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
	"github.com/minetest-go/oauth"
	"github.com/sirupsen/logrus"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

type Api struct {
	repos   *db.Repositories
	cfg     *types.Config
	core    *core.Core
	wc      *wallee.WalleeClient
	cbc     *coinbase.CoinbaseClient
	running *atomic.Bool
}

func NewApi(repos *db.Repositories, cfg *types.Config) *Api {
	return &Api{
		repos:   repos,
		cfg:     cfg,
		wc:      wallee.New(cfg.WalleeUserID, cfg.WalleeSpaceID, cfg.WalleeKey),
		running: &atomic.Bool{},
		core:    core.New(repos, cfg),
		cbc:     coinbase.New(cfg.CoinbaseKey),
	}
}

func (api *Api) Setup() {
	api.running.Store(true)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.PrometheusMiddleware)

	// setup routes

	// public
	apir := r.PathPrefix("/api").Subrouter()
	apir.HandleFunc("/info", api.GetInfo).Methods(http.MethodGet)
	apir.HandleFunc("/healthcheck", api.Healthcheck).Methods(http.MethodGet)
	apir.HandleFunc("/login", api.Logout).Methods(http.MethodDelete)
	apir.HandleFunc("/login", api.GetLogin).Methods(http.MethodGet)
	apir.HandleFunc("/login", api.Login).Methods(http.MethodPost)
	apir.HandleFunc("/register", api.Register).Methods(http.MethodPost)
	apir.HandleFunc("/nodetype", api.GetNodeTypes).Methods(http.MethodGet)
	apir.HandleFunc("/nodetype/{id}", api.GetNodeType).Methods(http.MethodGet)
	apir.HandleFunc("/logstream/{id}", api.LogStream).Methods(http.MethodPost)
	apir.HandleFunc("/exchange_rate", api.GetExchangeRates)
	apir.HandleFunc("/exchange_rate/{currency}", api.GetExchangeRate)
	apir.HandleFunc("/geoip/{ip}", api.ResolveGeoIP)
	apir.HandleFunc("/captcha", api.CreateCaptcha).Methods(http.MethodGet)
	apir.HandleFunc("/webhook/zahlsch", api.ZahlschWebhook).Methods(http.MethodPost)
	r.PathPrefix("/api/captcha/").Handler(captcha.Server(300, 200))

	// user api
	user_api := apir.NewRoute().Subrouter()
	user_api.Use(SecureHandler(api.LoginCheck()))
	user_api.HandleFunc("/set_password", api.Secure(api.SetPassword)).Methods(http.MethodPost)
	user_api.HandleFunc("/audit_log", api.Secure(api.SearchAuditLog)).Methods(http.MethodPost)

	user_api.HandleFunc("/coupon/redeem/{code}", api.Secure(api.RedeemCoupon)).Methods(http.MethodPost)

	user_api.HandleFunc("/profile", api.Secure(api.GetUserProfile)).Methods(http.MethodGet)
	user_api.HandleFunc("/profile", api.Secure(api.UpdateUserProfile)).Methods(http.MethodPost)

	user_api.HandleFunc("/profile/settings", api.Secure(api.GetUserSettings)).Methods(http.MethodGet)
	user_api.HandleFunc("/profile/settings/{key}", api.Secure(api.SetUserSetting)).Methods(http.MethodPut)
	user_api.HandleFunc("/profile/settings/{key}", api.Secure(api.DeleteUserSetting)).Methods(http.MethodDelete)

	user_api.HandleFunc("/node", api.Secure(api.GetNodes)).Methods(http.MethodGet)
	user_api.HandleFunc("/node/search", api.Secure(api.SearchNodes)).Methods(http.MethodPost)
	user_api.HandleFunc("/node", api.Secure(api.CreateNode)).Methods(http.MethodPost)
	user_api.HandleFunc("/node/{id}", api.Secure(api.GetNode)).Methods(http.MethodGet)
	user_api.HandleFunc("/node/{id}/stats", api.Secure(api.GetNodeStats)).Methods(http.MethodGet)
	user_api.HandleFunc("/node/{id}/job", api.Secure(api.GetLatestNodeJob)).Methods(http.MethodGet)
	user_api.HandleFunc("/node/{id}/mtservers", api.Secure(api.GetNodeServers)).Methods(http.MethodGet)
	user_api.HandleFunc("/node/{id}", api.Secure(api.DeleteNode)).Methods(http.MethodDelete)
	user_api.HandleFunc("/node/{id}", api.Secure(api.UpdateNode)).Methods(http.MethodPost)

	user_api.HandleFunc("/mtserver", api.Secure(api.GetMTServers)).Methods(http.MethodGet)
	user_api.HandleFunc("/mtserver", api.Secure(api.CreateMTServer)).Methods(http.MethodPost)
	user_api.HandleFunc("/mtserver/search", api.Secure(api.SearchMTServers)).Methods(http.MethodPost)
	user_api.HandleFunc("/mtserver/validate", api.Secure(api.ValidateCreateMTServer)).Methods(http.MethodPost)
	user_api.HandleFunc("/mtserver/{id}", api.Secure(api.GetMTServer)).Methods(http.MethodGet)
	user_api.HandleFunc("/mtserver/{id}", api.Secure(api.UpdateMTServer)).Methods(http.MethodPost)
	user_api.HandleFunc("/mtserver/{id}/setup", api.Secure(api.SetupMTServer)).Methods(http.MethodPost)
	user_api.HandleFunc("/mtserver/{id}/stats", api.Secure(api.GetMTServerStats)).Methods(http.MethodGet)
	user_api.HandleFunc("/mtserver/{id}/job", api.Secure(api.GetLatestMTServerJob)).Methods(http.MethodGet)
	user_api.HandleFunc("/mtserver/{id}", api.Secure(api.DeleteMTServer)).Methods(http.MethodDelete)

	user_api.HandleFunc("/overview/{user_id}", api.Secure(api.GetOverviewData)).Methods(http.MethodGet)

	user_api.HandleFunc("/transaction", api.Secure(api.GetTransactions)).Methods(http.MethodGet)
	user_api.HandleFunc("/transaction/create", api.Secure(api.CreateTransaction)).Methods(http.MethodPost)
	user_api.HandleFunc("/transaction/search", api.Secure(api.SearchTransaction)).Methods(http.MethodPost)
	user_api.HandleFunc("/transaction/{id}", api.Secure(api.GetTransaction)).Methods(http.MethodGet)
	user_api.HandleFunc("/transaction/{id}/check", api.Secure(api.CheckTransaction)).Methods(http.MethodGet)

	user_api.HandleFunc("/backup", api.Secure(api.GetBackups)).Methods(http.MethodGet)
	user_api.HandleFunc("/backup", api.Secure(api.CreateBackup)).Methods(http.MethodPost)
	user_api.HandleFunc("/backup/{id}", api.Secure(api.UpdateBackup)).Methods(http.MethodPost)
	user_api.HandleFunc("/backup/{id}", api.Secure(api.RemoveBackup)).Methods(http.MethodDelete)
	user_api.HandleFunc("/backup/{id}", api.Secure(api.GetBackup)).Methods(http.MethodGet)
	user_api.HandleFunc("/backup/{id}/download", api.Secure(api.DownloadBackup)).Methods(http.MethodGet)
	user_api.HandleFunc("/backup/{id}/job", api.Secure(api.GetBackupJob)).Methods(http.MethodGet)

	user_api.HandleFunc("/service/ticket", api.Secure(api.CreateTicket)).Methods(http.MethodPost)
	user_api.HandleFunc("/service/ticket/search", api.Secure(api.SearchTickets)).Methods(http.MethodPost)
	user_api.HandleFunc("/service/ticket/{ticket_id}", api.Secure(api.UpdateTicket)).Methods(http.MethodPost)
	user_api.HandleFunc("/service/message", api.Secure(api.CreateTicketMessage)).Methods(http.MethodPost)
	user_api.HandleFunc("/service/message/by-ticket/{ticket_id}", api.Secure(api.GetTicketMessages)).Methods(http.MethodGet)

	user_api.HandleFunc("/user/{id}", api.Secure(api.GetUserByID)).Methods(http.MethodGet)

	// admin api
	admin_api := apir.NewRoute().Subrouter()
	admin_api.Use(SecureHandler(api.RoleCheck(types.UserRoleAdmin)))
	admin_api.HandleFunc("/user", api.Secure(api.GetUsers)).Methods(http.MethodGet)
	admin_api.HandleFunc("/user/search", api.Secure(api.SearchUser)).Methods(http.MethodPost)
	admin_api.HandleFunc("/user/{id}", api.Secure(api.SaveUser)).Methods(http.MethodPost)
	admin_api.HandleFunc("/nodetype", api.Secure(api.CreateNodeType)).Methods(http.MethodPost)
	admin_api.HandleFunc("/nodetype/{id}", api.Secure(api.UpdateNodeType)).Methods(http.MethodPost)
	admin_api.HandleFunc("/nodetype/{id}", api.Secure(api.DeleteNodeType)).Methods(http.MethodDelete)
	admin_api.HandleFunc("/job", api.Secure(api.GetJobs)).Methods(http.MethodGet)
	admin_api.HandleFunc("/job/{id}", api.Secure(api.DeleteJob)).Methods(http.MethodDelete)
	admin_api.HandleFunc("/job/{id}", api.Secure(api.RetryJob)).Methods(http.MethodPost)

	admin_api.HandleFunc("/exchange_rate", api.Secure(api.CreateExchangeRate)).Methods(http.MethodPost)
	admin_api.HandleFunc("/exchange_rate/{currency}", api.Secure(api.UpdateExchangeRate)).Methods(http.MethodPut)
	admin_api.HandleFunc("/exchange_rate/{currency}", api.Secure(api.DeleteExchangeRate)).Methods(http.MethodDelete)

	admin_api.HandleFunc("/coupon", api.Secure(api.CreateCoupon)).Methods(http.MethodPost)
	admin_api.HandleFunc("/coupon", api.Secure(api.GetCoupons)).Methods(http.MethodGet)
	admin_api.HandleFunc("/coupon/{id}", api.Secure(api.GetCoupon)).Methods(http.MethodGet)
	admin_api.HandleFunc("/coupon/{id}/users", api.Secure(api.GetCouponUsers)).Methods(http.MethodGet)
	admin_api.HandleFunc("/coupon/{id}", api.Secure(api.UpdateCoupon)).Methods(http.MethodPost)

	// oauth
	if api.cfg.GithubOauthConfig.ClientID != "" {
		oauth_handler := oauth.NewHandler(api.OauthCallback, &oauth.OAuthConfig{
			Provider:    oauth.ProviderTypeGithub,
			ClientID:    api.cfg.GithubOauthConfig.ClientID,
			Secret:      api.cfg.GithubOauthConfig.Secret,
			CallbackURL: fmt.Sprintf("%s/oauth_callback/github", api.cfg.BaseURL),
		})
		r.Handle("/oauth_callback/github", oauth_handler)
		api.cfg.GithubOauthConfig.LoginURL = oauth_handler.LoginURL()
	}

	if api.cfg.DiscordOauthConfig.ClientID != "" {
		oauth_handler := oauth.NewHandler(api.OauthCallback, &oauth.OAuthConfig{
			Provider:    oauth.ProviderTypeDiscord,
			ClientID:    api.cfg.DiscordOauthConfig.ClientID,
			Secret:      api.cfg.DiscordOauthConfig.Secret,
			CallbackURL: fmt.Sprintf("%s/oauth_callback/discord", api.cfg.BaseURL),
		})
		r.Handle("/oauth_callback/discord", oauth_handler)
		api.cfg.DiscordOauthConfig.LoginURL = oauth_handler.LoginURL()
	}

	if api.cfg.MesehubOauthConfig.ClientID != "" {
		oauth_handler := oauth.NewHandler(api.OauthCallback, &oauth.OAuthConfig{
			Provider:    oauth.ProviderTypeMesehub,
			ClientID:    api.cfg.MesehubOauthConfig.ClientID,
			Secret:      api.cfg.MesehubOauthConfig.Secret,
			CallbackURL: fmt.Sprintf("%s/oauth_callback/mesehub", api.cfg.BaseURL),
		})
		r.Handle("/oauth_callback/mesehub", oauth_handler)
		api.cfg.MesehubOauthConfig.LoginURL = oauth_handler.LoginURL()
	}

	if api.cfg.CDBOauthConfig.ClientID != "" {
		oauth_handler := oauth.NewHandler(api.OauthCallback, &oauth.OAuthConfig{
			Provider:    oauth.ProviderTypeCDB,
			ClientID:    api.cfg.CDBOauthConfig.ClientID,
			Secret:      api.cfg.CDBOauthConfig.Secret,
			CallbackURL: fmt.Sprintf("%s/oauth_callback/cdb", api.cfg.BaseURL),
		})
		r.Handle("/oauth_callback/cdb", oauth_handler)
		api.cfg.CDBOauthConfig.LoginURL = oauth_handler.LoginURL()
	}

	// static files
	if api.cfg.Webdev {
		logrus.WithFields(logrus.Fields{"dir": "public"}).Info("Using live mode")
		fs := http.FileServer(http.FS(os.DirFS("public")))
		r.PathPrefix("/").HandlerFunc(fs.ServeHTTP)

	} else {
		logrus.Info("Using embed mode")
		r.PathPrefix("/").Handler(statigz.FileServer(public.Webapp, brotli.AddEncoding))
	}

	http.Handle("/", r)
}

func (api *Api) Stop() {
	api.running.Store(false)
}
