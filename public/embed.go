package public

import (
	"embed"
)

//go:embed assets/*
//go:embed js/* index.html
//go:embed node_modules/bootstrap/dist/css/bootstrap.min.css
//go:embed node_modules/vue/dist/vue.global.prod.js
//go:embed node_modules/vue/dist/vue.global.js
//go:embed node_modules/vue-router/dist/vue-router.global.prod.js
//go:embed node_modules/vue-router/dist/vue-router.global.js
//go:embed node_modules/@fortawesome/fontawesome-free/css/all.min.css
//go:embed node_modules/@fortawesome/fontawesome-free/webfonts/*
//go:embed node_modules/@vuepic/vue-datepicker/dist/vue-datepicker.iife.js
//go:embed node_modules/@vuepic/vue-datepicker/dist/main.css
var Webapp embed.FS
