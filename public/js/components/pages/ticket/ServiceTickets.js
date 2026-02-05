import CardLayout from "../../layouts/CardLayout.js";
import TicketStateBadge from "./TicketStateBadge.js";

import { search_tickets } from "../../../api/service_ticket.js";
import format_time from "../../../util/format_time.js";

export const ServiceTicketBreadcrumb = {icon: "ticket", name: "Service tickets", link: "/tickets"};

export default {
    components: {
        "card-layout": CardLayout,
        "ticket-state-badge": TicketStateBadge
    },
    data: function() {
        return {
            breadcrumb: [ServiceTicketBreadcrumb],
            tickets: []
        };
    },
    methods: {
        format_time
    },
    mounted: async function() {
        this.tickets = await search_tickets({});
    },
    template: /*html*/`
    <card-layout title="Service tickets" icon="ticket" :breadcrumb="breadcrumb">
        <div class="alert alert-primary">
            <i class="fa-solid fa-circle-info"></i>
            A Service Ticket is for technical or management issues around the hosting platform,
            if you have modding or client-issues please try to consult the <a href="https://forum.luanti.org/" target="new">Luanti-Forums</a> first.
        </div>
        <router-link class="btn btn-sm btn-outline-success" :to="'/tickets/new'">
            <i class="fa fa-plus"></i>
            Create ticket
        </router-link>
        <table class="table table-condensed">
            <thead>
                <tr>
                    <th>State</th>
                    <th>Created</th>
                    <th style="width: 50%;">Title</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="ticket in tickets" :key="ticket.id">
                    <td>
                        <ticket-state-badge :state="ticket.state"/>
                    </td>
                    <td>
                        {{ format_time(ticket.created) }}
                    </td>
                    <td>
                        <router-link :to="'/ticket/' + ticket.id">
                            {{ ticket.title }}
                        </router-link>
                    </td>
                </tr>
            </tbody>
        </table>
    </card-layout>
    `
};
