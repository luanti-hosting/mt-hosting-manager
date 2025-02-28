import CardLayout from "../layouts/CardLayout.js";
import CurrencyDisplay from "../CurrencyDisplay.js";

import { update } from "../../service/user.js";
import format_time from "../../util/format_time.js";
import { get_user_profile } from "../../service/user.js";
import { get_rates } from "../../service/exchange_rate.js";

export default {
	components: {
		"card-layout": CardLayout,
		"currency-display": CurrencyDisplay
	},
	data: function() {
		return {
			profile: get_user_profile(),
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "user", name: "Profile", link: "/profile"
			}]
		};
	},
	methods: {
		get_rates,
		format_time,
		save: function() {
			update(this.profile);
		}
	},
	template: /*html*/`
	<card-layout title="Profile" icon="user" :breadcrumb="breadcrumb">
		<table class="table" v-if="profile && profile.name">
			<tr>
				<td>Name</td>
				<td>{{profile.name}}</td>
			</tr>
			<tr>
				<td>Role</td>
				<td>{{profile.role}}</td>
			</tr>
			<tr>
				<td>Type</td>
				<td>{{profile.type}}</td>
			</tr>
			<tr>
				<td>Created</td>
				<td>
					{{format_time(profile.created)}}
				</td>
			</tr>
			<tr>
				<td>Balance</td>
				<td>
					<currency-display :eurocents="profile.balance" :enable_warning="true"/>
				</td>
			</tr>
			<tr>
				<td>E-Mail (for notifications only)</td>
				<td>
					<input class="form-control" v-model="profile.mail" placeholder="E-Mail address"/>
				</td>
			</tr>
			<tr>
				<td>Alternative currency</td>
				<td>
					<select class="form-control" v-model="profile.currency">
						<option v-for="r in get_rates()" :key="r.currency" :value="r.currency">{{r.display_name}}</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>Action</td>
				<td>
					<a class="btn btn-success w-100" v-on:click="save">
						<i class="fa-solid fa-floppy-disk"></i>
						Save
					</a>
				</td>
			</tr>
		</table>
	</card-layout>
    `
};
