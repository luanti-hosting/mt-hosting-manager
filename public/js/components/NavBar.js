import { logout, get_claims, is_logged_in, has_role } from '../service/login.js';
import { get_stage } from '../service/info.js';
import CurrencyDisplay from './CurrencyDisplay.js';
import { get_balance } from '../service/user.js';

export default {
	components: {
		"currency-display": CurrencyDisplay
	},
	data: function() {
		return {
			stage: get_stage()
		};
	},
	methods: {
		logout: function(){
			logout().then(() => this.$router.push("/login"));
		},
		has_role: has_role
	},
	computed: {
		is_logged_in: is_logged_in,
		claims: get_claims,
		balance: get_balance
	},
	template: /*html*/`
		<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
			<div class="container-fluid">
				<router-link to="/" class="navbar-brand">Minetest hosting</router-link>
				<ul class="navbar-nav me-auto mb-2 mb-lg-0">
					<li class="nav-item">
						<router-link to="/" class="nav-link">
							<i class="fa fa-home"></i> Home
						</router-link>
					</li>
					<li class="nav-item">
						<router-link to="/login" class="nav-link" v-if="!is_logged_in">
							<i class="fa fa-user"></i> Login
						</router-link>
					</li>
					<li class="nav-item">
						<router-link to="/help" class="nav-link" v-if="is_logged_in">
							<i class="fa fa-question"></i> Help
						</router-link>
					</li>
					<li class="nav-item" v-if="is_logged_in">
						<router-link to="/finance" class="nav-link">
							<i class="fa-solid fa-money-bill"></i> Finance
							<currency-display :eurocents="balance" :enable_warning="true"/>
						</router-link>
					</li>
					<li class="nav-item" v-if="has_role('ADMIN')">
						<router-link to="/node_types" class="nav-link">
							<i class="fa fa-server"></i> Node-Types
						</router-link>
					</li>
					<li class="nav-item" v-if="has_role('ADMIN')">
						<router-link to="/sendmail" class="nav-link">
							<i class="fa fa-envelope"></i> Send mail
						</router-link>
					</li>
					<li class="nav-item" v-if="has_role('ADMIN')">
						<router-link to="/jobs" class="nav-link">
							<i class="fa fa-play"></i> Jobs
						</router-link>
					</li>
					<li class="nav-item" v-if="is_logged_in">
						<router-link to="/nodes" class="nav-link">
							<i class="fa fa-server"></i> Nodes
						</router-link>
					</li>
					<li class="nav-item" v-if="is_logged_in">
						<router-link to="/mtservers" class="nav-link">
							<i class="fa fa-list"></i> Servers
						</router-link>
					</li>
					<li class="nav-item" v-if="is_logged_in">
						<router-link to="/audit-logs" class="nav-link">
							<i class="fa fa-rectangle-list"></i> Audit-logs
						</router-link>
					</li>
				</ul>
				<div class="btn btn-warning" v-if="stage != 'prod'">
					<i class="fa-solid fa-triangle-exclamation"></i>
					Stage: {{stage}}
				</div>
				<div class="d-flex">
					<div class="btn-group">
						<button class="btn btn-outline-secondary" v-if="claims">
							<router-link to="/profile">
								<i class="fas fa-user"></i>
								<span>
									Logged in as <b>{{claims.mail}}</b>
								</span>
							</router-link>
						</button>
						<button class="btn btn-secondary" v-on:click="logout" v-if="is_logged_in">
							<i class="fa-solid fa-right-from-bracket"></i>
							Logout
						</button>
					</div>
				</div>
			</div>
		</nav>
	`
};
