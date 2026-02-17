package types

import (
	"os"
	"strings"
)

type OAuthConfig struct {
	ClientID string
	Secret   string
	LoginURL string
}

type Config struct {
	BaseURL             string
	HostingDomainSuffix string
	DNSRecordSuffix     string
	ReservedPrefixes    []string
	Stage               string
	InitialBalance      string
	Webdev              bool
	EnableWorker        bool
	JWTKey              string

	LogStreamKey string
	LogStreamDir string

	CookieName   string
	CookiePath   string
	CookieSecure bool

	HetznerCloudKey string
	HetznerZoneName string
	HetznerSSHKeys  []string

	CoinbaseKey     string
	CoinbaseEnabled bool

	WalleeUserID  string
	WalleeSpaceID string
	WalleeKey     string
	WalleeEnabled bool

	ZahlschPageID     string
	ZahlschUser       string
	ZahlschWebhookKey string
	ZahlschEnabled    bool

	GithubOauthConfig  *OAuthConfig
	DiscordOauthConfig *OAuthConfig
	MesehubOauthConfig *OAuthConfig
	CDBOauthConfig     *OAuthConfig

	MaxBalance int //max balance in eurocents

	MailHost     string
	MailAddress  string
	MailPassword string

	// backup storage
	S3Endpoint  string
	S3KeyID     string
	S3AccessKey string
	S3Bucket    string
}

func NewConfig() *Config {
	return &Config{
		BaseURL: os.Getenv("BASEURL"),
		// entire suffix for display purposes
		HostingDomainSuffix: os.Getenv("HOSTING_DOMAIN_SUFFIX"),
		// suffix for the record itself (valid inside dns zone)
		DNSRecordSuffix:  os.Getenv("DNS_RECORD_SUFFIX"),
		ReservedPrefixes: strings.Split(os.Getenv("RESERVED_PREFIXES"), ","),
		Stage:            os.Getenv("STAGE"),
		InitialBalance:   os.Getenv("INITIAL_BALANCE"),
		Webdev:           os.Getenv("WEBDEV") == "true",
		EnableWorker:     os.Getenv("ENABLE_WORKER") == "true",
		JWTKey:           os.Getenv("JWT_KEY"),
		LogStreamKey:     os.Getenv("LOG_STREAM_KEY"),
		LogStreamDir:     os.Getenv("LOG_STREAM_DIR"),
		CookieName:       "mt-hosting-manager",
		CookiePath:       os.Getenv("COOKIE_PATH"),
		CookieSecure:     os.Getenv("COOKIE_SECURE") == "true",
		// hetzner
		HetznerCloudKey: os.Getenv("HETZNER_CLOUD_KEY"),
		HetznerZoneName: os.Getenv("HETZNER_ZONE_NAME"),
		HetznerSSHKeys:  strings.Split(os.Getenv("HETZNER_SSH_KEYS"), ","),
		// coinbase
		CoinbaseKey:     os.Getenv("COINBASE_KEY"),
		CoinbaseEnabled: os.Getenv("COINBASE_ENABLED") == "true",
		// wallee
		WalleeUserID:  os.Getenv("WALLEE_USERID"),
		WalleeSpaceID: os.Getenv("WALLEE_SPACEID"),
		WalleeKey:     os.Getenv("WALLEE_KEY"),
		WalleeEnabled: os.Getenv("WALLEE_ENABLED") == "true",
		// zahls.ch
		ZahlschPageID:     os.Getenv("ZAHLSCH_PAGEID"),
		ZahlschUser:       os.Getenv("ZAHLSCH_USER"),
		ZahlschWebhookKey: os.Getenv("ZAHLSCH_WEBHOOK_KEY"),
		ZahlschEnabled:    os.Getenv("ZAHLSCH_ENABLED") == "true",
		// oauth
		GithubOauthConfig: &OAuthConfig{
			ClientID: os.Getenv("GITHUB_CLIENTID"),
			Secret:   os.Getenv("GITHUB_SECRET"),
		},
		DiscordOauthConfig: &OAuthConfig{
			ClientID: os.Getenv("DISCORD_CLIENTID"),
			Secret:   os.Getenv("DISCORD_SECRET"),
		},
		MesehubOauthConfig: &OAuthConfig{
			ClientID: os.Getenv("MESEHUB_CLIENTID"),
			Secret:   os.Getenv("MESEHUB_SECRET"),
		},
		CDBOauthConfig: &OAuthConfig{
			ClientID: os.Getenv("CDB_CLIENTID"),
			Secret:   os.Getenv("CDB_SECRET"),
		},
		MaxBalance: 100 * 100,

		MailHost:     os.Getenv("MAIL_HOST"),
		MailAddress:  os.Getenv("MAIL_ADDRESS"),
		MailPassword: os.Getenv("MAIL_PASSWORD"),

		S3Endpoint:  os.Getenv("S3_ENDPOINT"),
		S3KeyID:     os.Getenv("S3_KEYID"),
		S3AccessKey: os.Getenv("S3_ACCESSKEY"),
		S3Bucket:    os.Getenv("S3_BUCKET"),
	}
}
