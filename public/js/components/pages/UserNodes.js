import CardLayout from "../layouts/CardLayout.js";
import NodeLink from "../NodeLink.js";
import UserLink from "../UserLink.js";
import NodeState from "../NodeState.js";

import { get_all } from "../../api/node.js";
import { get_nodetype } from "../../service/nodetype.js";
import { has_role } from "../../service/login.js";
import format_time from "../../util/format_time.js";

export default {
	components: {
		"card-layout": CardLayout,
		"node-link": NodeLink,
		"node-state": NodeState,
		"user-link": UserLink
	},
	data: function() {
		return {
			nodes: [],
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "server", name: "Nodes", link: "/nodes"
			}]
		};
	},
	mounted: function() {
		this.update();
		this.handle = setInterval(() => this.update(), 2000);
	},
	beforeUnmount: function() {
		clearInterval(this.handle);
	},
	methods: {
		format_time,
		has_role,
		get_nodetype,
		update: function() {
			get_all().then(nodes => this.nodes = nodes);
		}
	},
	template: /*html*/`
	<card-layout title="Nodes" icon="server" :breadcrumb="breadcrumb">
		<table class="table">
			<thead>
				<th>Name</th>
				<th v-if="has_role('ADMIN')">User</th>
				<th>State</th>
				<th>Created</th>
				<th>Node-Type</th>
			</thead>
			<tbody>
				<tr v-for="node in nodes" :key="node.id">
					<td>
						<node-link :node="node"/>
					</td>
					<td v-if="has_role('ADMIN')">
						<user-link :id="node.user_id"/>
					</td>
					<td>
						<node-state :state="node.state"/>
					</td>
					<td>{{format_time(node.created)}}</td>
					<td>
						{{get_nodetype(node.node_type_id).name}}
					</td>
				</tr>
			</tbody>
		</table>
		<router-link class="btn btn-success" to="/nodes/create">
			<i class="fa fa-plus"></i>
			Create node
		</router-link>
	</card-layout>
	`
};
