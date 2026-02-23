import CardLayout from "../../layouts/CardLayout.js";
import TicketStateBadge from "./TicketStateBadge.js";

import format_time from "../../../util/format_time.js";
import { search_tickets, get_messages, create_message, update_ticket } from "../../../api/service_ticket.js";
import { fetch_tickets } from "../../../service/service_ticket.js";
import { get_cached_user } from "../../../service/user_cache.js";

export const ServiceTicketBreadcrumb = {icon: "ticket", name: "Service tickets", link: "/tickets"};

export default {
    props: ["id"],
    components: {
        "card-layout": CardLayout,
        "ticket-state-badge": TicketStateBadge
    },
    data: function() {
        return {
            breadcrumb: [ServiceTicketBreadcrumb, {
                icon: "ticket",
                name: `Ticket ${this.id}`,
                link: `/ticket/${this.id}`
            }],
            ticket: null,
            response_text: "",
            messages: []
        };
    },
    methods: {
        format_time,
        send_message: async function() {
            const msg = await create_message({
                ticket_id: this.id,
                message: this.response_text
            });
            this.response_text = "";
            msg.user = await get_cached_user(msg.user_id);
            this.messages.push(msg);
        },
        mark_resolved: async function(){
            this.ticket.state = "RESOLVED";
            await update_ticket(this.ticket);
            await fetch_tickets();
        },
        close: async function() {
            this.ticket.state = "CLOSED";
            await update_ticket(this.ticket);
            await fetch_tickets();
        },
        reopen: async function() {
            this.ticket.state = "OPEN";
            await update_ticket(this.ticket);
            await fetch_tickets();
        }
    },
    mounted: async function() {
        const tickets = await search_tickets({ ticket_id: this.id });
        if (tickets) {
            this.ticket = tickets[0];
        }

        const messages = await get_messages(this.id);
        for (let i=0; i<messages.length; i++) {
            const msg = messages[i];
            msg.user = await get_cached_user(msg.user_id);
        }

        this.messages = messages;
    },
    template: /*html*/`
    <card-layout title="Service ticket" icon="ticket" :breadcrumb="breadcrumb">
        <table class="table table-condensed" v-if="ticket">
            <tbody>
                <tr>
                    <th>State</th>
                    <td>
                        <ticket-state-badge :state="ticket.state"/>
                    </td>
                </tr>
                <tr>
                    <th>Created</th>
                    <td>{{ format_time(ticket.created) }}</td>
                </tr>
                <tr>
                    <th>Title</th>
                    <td>{{ ticket.title }}</td>
                </tr>
            </tbody>
        </table>
        <i class="fa fa-spinner fa-spin" v-else></i>

        <table class="table table-condensed">
            <thead>
                <tr>
                    <th class="col-md-5">Received</th>
                    <th>Message</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="msg in messages" :key="msg.id">
                    <td>
                        {{ format_time(msg.timestamp) }}
                        <span class="badge bg-secondary">{{ msg.user.name }}</span>
                    </td>
                    <td>
                        <pre>{{ msg.message }}</pre>
                    </td>
                </tr>
            </tbody>
        </table>
        <div class="input-group" v-if="ticket && ticket.state == 'OPEN'">
            <textarea placeholder="Response text" class="form-control" v-model="response_text"/>
            <button class="btn btn-outline-secondary" v-on:click="send_message" :disabled="!response_text">
                <i class="fa fa-envelope"></i>
            </button>
        </div>
        <br>
        <div class="btn-group w-100" v-if="ticket && ticket.state == 'OPEN'">
            <button class="btn btn-success" v-on:click="mark_resolved">
                Mark as resolved
            </button>
            <button class="btn btn-secondary" v-on:click="close">
                Close ticket
            </button>
        </div>
        <div class="btn-group w-100" v-if="ticket && ticket.state != 'OPEN'">
            <button class="btn btn-secondary" v-on:click="reopen">
                Reopen
            </button>
        </div>
    </card-layout>
    `
};
