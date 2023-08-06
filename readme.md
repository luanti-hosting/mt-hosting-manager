
# Minetest hosting orchestrator

State: **WIP**

## Roadmap

* [ ] Login options
  * [x] github (MVP)
  * [ ] discord
  * [ ] mesehub
  * [ ] gitlab
* [ ] Payment options
  * [x] Wallee (MVP)
  * [ ] Stripe
  * [ ] Coinbase crypto
* [ ] Host creation
  * [x] Hetzner (MVP)
  * [ ] Contabo
* [x] Host provisioning (MVP)
* [ ] Instance setup
* [ ] Backups?

# Dev

```sh
# setup
docker-compose up
# set all users as admin
sudo sqlite3 mt-hosting.sqlite "update user set role = 'ADMIN'"
```

# Environment variables

* `LOGLEVEL` "debug" / "info"
* `ENABLE_WORKER`
* `STAGE` "prod" / "dev"

* `CSRF_KEY`
* `JWT_KEY`
* `BASEURL`
* `WEBDEV`
* `COOKIE_DOMAIN`
* `COOKIE_PATH`
* `COOKIE_SECURE`

* `GITHUB_CLIENTID`
* `GITHUB_SECRET`

* `ADMIN_USER_MAIL` mail of the user that gets the admin role on register
* `DISABLE_SIGNUP`

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
* Media
  * "docs/minetest.svg" `CC BY-SA 3.0`