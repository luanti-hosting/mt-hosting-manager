import CardLayout from "../../layouts/CardLayout.js";
import { ServiceTicketBreadcrumb } from "./ServiceTickets.js";

import { get_all as get_all_nodes } from "../../../api/node.js";

export default {
    components: {
        "card-layout": CardLayout
    },
    data: function() {
        return {
            breadcrumb: [ServiceTicketBreadcrumb, {
                icon: "plus", name: "New service ticket", link: "/tickets/new"
            }],
            user_nodes: get_all_nodes(),
            user_node_id: null,
            title: ""
        };
    },
    template: /*html*/`
    <card-layout title="New service ticket" icon="plus" :breadcrumb="breadcrumb">
		<table class="table">
			<tbody>
				<tr>
					<td>Title</td>
					<td>
						<input type="text" class="form-control" v-model="title"/>
					</td>
				</tr>
				<tr>
					<td>Node</td>
					<td>
						<select v-model="user_node_id" class="form-control">
							<option v-for="node in user_nodes" :value="node.id">{{node.alias}} ({{node.name}})</option>
						</select>
					</td>
				</tr>
            </tbody>
        </table>
    </card-layout>
    `
};
