import CardLayout from "../../layouts/CardLayout.js";
import { ServiceTicketBreadcrumb } from "./ServiceTickets.js";

import format_size from "../../../util/format_size.js";
import format_time from "../../../util/format_time.js";

import { get_all as get_all_nodes } from "../../../api/node.js";
import { get_all as get_all_servers } from "../../../api/mtserver.js";
import { get_all as get_all_backups } from "../../../api/backup.js";

export default {
    components: {
        "card-layout": CardLayout
    },
    data: function() {
        return {
            breadcrumb: [ServiceTicketBreadcrumb, {
                icon: "plus", name: "New service ticket", link: "/tickets/new"
            }],
            user_nodes: [],
            user_node_id: null,
            minetest_servers: [],
            minetest_server_id: null,
            backups: [],
            backup_id: null,
            title: "",
            message: ""
        };
    },
    mounted: async function() {
        this.user_nodes = (await get_all_nodes()).filter(node => node.state == "RUNNING");
        this.minetest_servers = (await get_all_servers()).filter(server => server.state == "RUNNING");
        this.backups = await get_all_backups();
    },
    methods: {
        format_size,
        format_time
    },
    template: /*html*/`
    <card-layout title="New service ticket" icon="plus" :breadcrumb="breadcrumb">
		<table class="table">
			<tbody>
				<tr>
					<td>Title</td>
					<td>
						<input type="text" class="form-control" placeholder="Short summary of your issue" v-model="title"/>
					</td>
				</tr>
				<tr>
					<td>Node (optional)</td>
					<td>
						<select v-model="user_node_id" class="form-control">
							<option v-for="node in user_nodes" :value="node.id">{{node.alias}} ({{node.name}})</option>
						</select>
					</td>
				</tr>
				<tr>
					<td>Server (optional)</td>
					<td>
						<select v-model="minetest_server_id" class="form-control">
							<option v-for="server in minetest_servers" :value="server.id">{{server.dns_name}}:{{server.port}}</option>
						</select>
					</td>
				</tr>
				<tr>
					<td>Backup (optional)</td>
					<td>
						<select v-model="backup_id" class="form-control">
							<option v-for="backup in backups" :value="backup.id">
                                Id: {{backup.id}} ({{backup.state}}, {{format_time(backup.created)}} / {{format_size(backup.size)}})
                            </option>
						</select>
					</td>
				</tr>
                <tr>
                    <td>Message</td>
                    <td>
                        <textarea class="form-control" style="height: 250px;" placeholder="Describe your issue here"></textarea>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">
                        <button class="btn btn-primary w-100">
                            <i class="fa fa-plus"></i>
                            Create ticket
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
    </card-layout>
    `
};
