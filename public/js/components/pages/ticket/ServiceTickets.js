import CardLayout from "../../layouts/CardLayout.js";
import TicketStateBadge from "./TicketStateBadge.js";
import UserLink from "../../UserLink.js";

import { get_all_tickets, fetch_tickets } from "../../../service/service_ticket.js";
import format_time from "../../../util/format_time.js";
import { has_role } from "../../../service/login.js";

export const ServiceTicketBreadcrumb = {icon: "ticket", name: "Service tickets", link: "/tickets"};

export default {
    components: {
        "card-layout": CardLayout,
        "ticket-state-badge": TicketStateBadge,
        "user-link": UserLink
    },
    data: function() {
        return {
            breadcrumb: [ServiceTicketBreadcrumb]
        };
    },
    methods: {
        format_time,
        has_role
    },
    computed: {
        tickets: get_all_tickets
    },
    mounted: async function() {
        await fetch_tickets();
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
                    <th v-if="has_role('ADMIN')">User</th>
                    <th>Created</th>
                    <th style="width: 50%;">Title</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="ticket in tickets" :key="ticket.id">
                    <td>
                        <ticket-state-badge :state="ticket.state"/>
                    </td>
                    <td v-if="has_role('ADMIN')">
                        <user-link :id="ticket.user_id"/>
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
