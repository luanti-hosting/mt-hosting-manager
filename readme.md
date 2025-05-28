
# Luanti hosting orchestrator

![](https://github.com/luanti-hosting/mt-hosting-manager/workflows/build/badge.svg)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/luanti-hosting/mt-hosting-manager)
[![Go Report Card](https://goreportcard.com/badge/github.com/luanti-hosting/mt-hosting-manager)](https://goreportcard.com/report/github.com/luanti-hosting/mt-hosting-manager)
[![Coverage Status](https://coveralls.io/repos/github/luanti-hosting/mt-hosting-manager/badge.svg)](https://coveralls.io/github/luanti-hosting/mt-hosting-manager)


State: **WIP**

# Dev

```sh
# start postgres
docker compose up -d postgres
# ui assets
docker compose up hosting_webapp
# main app
docker compose up hosting
```

# Environment variables

* `LOGLEVEL` "debug" / "info"
* `ENABLE_WORKER`
* `STAGE` "prod" / "dev"

* `JWT_KEY`
* `BASEURL`
* `WEBDEV`
* `COOKIE_DOMAIN`
* `COOKIE_PATH`
* `COOKIE_SECURE`

* `GITHUB_CLIENTID`
* `GITHUB_SECRET`

* `SIGNUP_WHITELIST`

* `WALLEE_USERID`
* `WALLEE_SPACEID`
* `WALLEE_KEY`

* `NTFY_URL`
* `NTFY_TOPIC`
* `NTFY_USERNAME`
* `NTFY_PASSWORD`

* `SSH_KEY`

# License

* Code: `MIT`

## Other assets

* `public/assets/default_mese_crystal.png` CC BY-SA 3.0 https://github.com/minetest/minetest_game
