import CardLayout from "../layouts/CardLayout.js";
import NodeState from "../NodeState.js";
import ServerList from "../ServerList.js";
import CurrencyDisplay from "../CurrencyDisplay.js";
import NodeTypeSpec from "../NodeTypeSpec.js";
import UserLink from "../UserLink.js";
import TimestampBadge from "../TimestampBadge.js";

import { country_map, flag_map } from "../../util/country.js";

import { get_by_id, get_stats, update as update_node, get_mtservers_by_nodeid, get_latest_job } from "../../api/node.js";
import { get_hostingdomain_suffix } from "../../service/info.js";
import { get_nodetype } from "../../service/nodetype.js";
import { has_role } from "../../service/login.js";

import format_time from "../../util/format_time.js";

const bytes_in_gb = 1000 * 1000 * 1000;

function get_gb_rounded(bytes) {
	return parseInt(bytes / bytes_in_gb * 100) / 100;
}

export default {
	props: ["id"],
	components: {
		"card-layout": CardLayout,
		"node-state": NodeState,
		"server-list": ServerList,
		"currency-display": CurrencyDisplay,
		"node-type-spec": NodeTypeSpec,
		"user-link": UserLink,
		"timestamp-badge": TimestampBadge
	},
	data: function() {
		return {
			hostingdomain_suffix: get_hostingdomain_suffix(),
			servers: [],
			job: null,
			node: null,
			nodetype: null,
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "server", name: "Nodes", link: "/nodes"
			},{
				icon: "server", name: "Node detail", link: `/nodes/${this.id}`
			}],
			load_percent: 0,
			disk_gb_total: 0,
			disk_gb_used: 0,
			disk_percent: 0,
			memory_gb_total: 0,
			memory_gb_used: 0,
			memory_percent: 0,
			country_map,
			flag_map
		};
	},
	mounted: async function() {
		const nodeid = this.id;
		this.node = await get_by_id(nodeid);
		this.update_stats();
		this.handle = setInterval(() => this.update_stats(), 2000);
		this.nodetype = get_nodetype(this.node.node_type_id);
	},
	beforeUnmount: function() {
		clearInterval(this.handle);
	},
	methods: {
		format_time,
		has_role,
		update_stats: async function() {
			const nodeid = this.id;
			if (this.node.state == "PROVISIONING") {
				// fetch job progress
				this.job = await get_latest_job(this.id);
			}

			this.node = await get_by_id(nodeid);

			const list = await get_mtservers_by_nodeid(nodeid);
			this.servers = list.filter(s => s.state != "DECOMMISSIONED");
	
			if (this.node.state == "DECOMMISSIONED") {
				clearInterval(this.handle);
				return;
			}

			if (this.node.state == "RUNNING") {
				const stats = await get_stats(nodeid);
				this.load_percent = Math.min(stats.load_percent, 100);
				this.disk_gb_total = get_gb_rounded(stats.disk_size);
				this.disk_gb_used = get_gb_rounded(stats.disk_used);
				this.disk_percent = parseInt(this.disk_gb_used / this.disk_gb_total * 100);
				this.memory_gb_total = get_gb_rounded(stats.memory_size);
				this.memory_gb_used = get_gb_rounded(stats.memory_used);
				this.memory_percent = parseInt(this.memory_gb_used / this.memory_gb_total * 100);
			}
		},
		save: function() {
			update_node(this.node);
		}
	},
	template: /*html*/`
	<card-layout title="Node details" icon="server" :breadcrumb="breadcrumb">
		<h4>Details</h4>
		<table class="table" v-if="node && nodetype">
			<tbody>
				<tr>
					<td>ID</td>
					<td>{{node.id}}</td>
				</tr>
				<tr v-if="has_role('ADMIN')">
					<td>User</td>
					<td>
						<user-link :id="node.user_id"/>
					</td>
				</tr>
				<tr>
					<td>Specs</td>
					<td>
						<node-type-spec :nodetype="nodetype"/>
					</td>
				</tr>
				<tr>
					<td>Location</td>
					<td>
						{{country_map[node.location]}} {{flag_map[node.location]}}
					</td>
				</tr>
				<tr>
					<td>Hostname</td>
					<td>{{node.name}}.{{hostingdomain_suffix}}</td>
				</tr>
				<tr>
					<td>IP</td>
					<td>
						<b>v4:</b> {{node.ipv4}}
						<b>v6:</b> {{node.ipv6}}
					</td>
				</tr>
				<tr>
					<td>Created</td>
					<td>
						<timestamp-badge :timestamp="node.created" :show_duration="true"/>
					</td>
				</tr>
				<tr>
					<td>Next billing cycle</td>
					<td>
						<timestamp-badge :timestamp="node.valid_until" :show_duration="true"/>
					</td>
				</tr>
				<tr>
					<td>Daily cost</td>
					<td>
						<currency-display :eurocents="nodetype.daily_cost"/>
					</td>
				</tr>
				<tr>
					<td>State</td>
					<td>
						<node-state :state="node.state"/>
						<div class="alert alert-info" v-if="node.state == 'PROVISIONING'">
							<i class="fa-solid fa-info"></i>
							The node is still provisioning, this might take a minute or two
						</div>
					</td>
				</tr>
				<tr v-if="node.state == 'PROVISIONING' && job">
					<td>Provisioning progress</td>
					<td>
						<b>Status: </b> {{job.message}}
						<div class="progress">
							<div class="progress-bar" v-bind:style="{ width: job.progress_percent + '%' }" v-bind:class="{'bg-danger': job.state == 'DONE_FAILURE'}">
								{{job.progress_percent}}%
							</div>
						</div>
					</td>
				</tr>
				<tr v-if="node.state == 'RUNNING'">
					<td>Actions</td>
					<td>
						<router-link class="btn btn-sm btn-danger" :to="'/nodes/' + node.id + '/delete'">
							<i class="fa fa-trash"></i>
							Delete
						</router-link>
					</td>
				</tr>
				<tr v-if="node.state == 'RUNNING'">
					<td>Alias</td>
					<td>
						<div class="btn-group w-100">
							<input type="text" class="form-control" v-model="node.alias"/>
							<button class="btn btn-xs btn-info" v-on:click="save()">
								<i class="fa fa-floppy-disk"></i>
							</button>
						</div>
					</td>
				</tr>
				<tr v-if="node.state == 'RUNNING'">
					<td>CPU Usage</td>
					<td>
						<div class="progress">
							<div class="progress-bar" v-bind:style="{ width: load_percent + '%' }">
								{{load_percent}}%
							</div>
						</div>
					</td>
				</tr>
				<tr v-if="node.state == 'RUNNING'">
					<td>Disk Usage</td>
					<td>
						{{disk_gb_used}}/{{disk_gb_total}} GB
						<div class="progress">
							<div class="progress-bar bg-warning" v-bind:style="{ width: disk_percent + '%' }">
								{{disk_percent}}%
							</div>
						</div>
					</td>
				</tr>
				<tr v-if="node.state == 'RUNNING'">
					<td>Memory Usage</td>
					<td>
						{{memory_gb_used}}/{{memory_gb_total}} GB
						<div class="progress">
							<div class="progress-bar bg-danger" v-bind:style="{ width: memory_percent + '%' }">
								{{memory_percent}}%
							</div>
						</div>
					</td>
				</tr>
			</tbody>
		</table>
		<div v-if="node">
			<h4>Servers</h4>
			<server-list :list="servers" :show_stats="true"/>
			<router-link class="btn btn-success" :to="'/mtservers/create?node=' + node.id" v-if="node.state == 'RUNNING'">
				<i class="fa fa-plus"></i>
				Create server
			</router-link>
		</div>
	</card-layout>
	`
};
