package core

import (
	"mt-hosting-manager/api/coinbase"
	"mt-hosting-manager/api/hetzner_cloud"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
)

type Core struct {
	repos  *db.Repositories
	cfg    *types.Config
	wc     *wallee.WalleeClient
	hdns   *hetzner_dns.HetznerDNSClient
	hcloud *hetzner_cloud.HetznerCloudClient
	cbc    *coinbase.CoinbaseClient
	GeoIP  *GeoipResolver
}

func New(repos *db.Repositories, cfg *types.Config) *Core {
	return &Core{
		repos:  repos,
		cfg:    cfg,
		wc:     wallee.New(cfg.WalleeUserID, cfg.WalleeSpaceID, cfg.WalleeKey),
		hdns:   hetzner_dns.New(cfg.HetznerApiKey, cfg.HetznerApiZoneID),
		hcloud: hetzner_cloud.New(cfg.HetznerCloudKey),
		cbc:    coinbase.New(cfg.CoinbaseKey),
		GeoIP:  NewGeoipResolver(),
	}
}
