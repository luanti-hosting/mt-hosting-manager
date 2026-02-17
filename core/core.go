package core

import (
	"mt-hosting-manager/api/coinbase"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type Core struct {
	repos *db.Repositories
	cfg   *types.Config
	wc    *wallee.WalleeClient
	hc    *hcloud.Client
	cbc   *coinbase.CoinbaseClient
	GeoIP *GeoipResolver
}

func New(repos *db.Repositories, cfg *types.Config) *Core {
	return &Core{
		repos: repos,
		cfg:   cfg,
		wc:    wallee.New(cfg.WalleeUserID, cfg.WalleeSpaceID, cfg.WalleeKey),
		cbc:   coinbase.New(cfg.CoinbaseKey),
		hc:    hcloud.NewClient(hcloud.WithToken(cfg.HetznerCloudKey)),
		GeoIP: NewGeoipResolver(),
	}
}
