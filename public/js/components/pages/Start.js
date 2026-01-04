import CardLayout from "../layouts/CardLayout.js";
import Breadcrumb from "../Breadcrumb.js";

import { is_logged_in } from "../../service/login.js";

export default {
	components: {
		"card-layout": CardLayout,
		"bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [{icon: "home", name: "Home", link: "/"}]
		};
	},
	computed: {
		is_logged_in
	},
	template: /*html*/`
	<card-layout title="Home" icon="home" :breadcrumb="breadcrumb">
		<div class="text-center">
			<h4>Luanti hosting</h4>
			<img src="assets/luanti-hosting-600px.png" class="img img-rounded"/>
			<hr/>
			<router-link class="btn btn-secondary" to="/pricing">
				<i class="fa fa-money-bill"></i> Pricing
			</router-link>
			&nbsp;
			<a class="btn btn-secondary" href="mailto:hosting@luanti.ch">
				<i class="fa fa-envelope"></i> Contact
			</a>
			&nbsp;
			<a class="btn btn-secondary" href="https://forum.luanti.org/viewtopic.php?t=31291">
				<i class="fa fa-link"></i> Luanti forum topic
			</a>
			&nbsp;
			<a class="btn btn-secondary" href="https://discord.gg/Xj62fUbQkn" target="new">
				<i class="fa-brands fa-discord"></i>
				Join the discord server
			</a>
			&nbsp;
			<a class="btn btn-outline-secondary" href="https://github.com/luanti-hosting/mt-hosting-manager" target="new">
				<i class="fa-brands fa-github"></i> Source
			</a>
			&nbsp;
			<a class="btn btn-outline-secondary" href="https://github.com/minetest-go/mtui" target="new">
				<i class="fa-brands fa-github"></i> Powered by mtui
			</a>
			&nbsp;
			<router-link class="btn btn btn-outline-success" :to="'/mtservers/create'" v-if="is_logged_in">
				<i class="fa fa-plus"></i>
				Create server
			</router-link>
		</div>
	</card-layout>
	`
};
