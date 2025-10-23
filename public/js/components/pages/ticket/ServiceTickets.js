import CardLayout from "../../layouts/CardLayout.js";

export const ServiceTicketBreadcrumb = {icon: "ticket", name: "Service tickets", link: "/tickets"};

export default {
    components: {
        "card-layout": CardLayout
    },
    data: function() {
        return {
            breadcrumb: [ServiceTicketBreadcrumb]
        };
    },
    template: /*html*/`
    <card-layout title="Service tickets" icon="ticket" :breadcrumb="breadcrumb">
        <router-link class="btn btn-sm btn-outline-success" :to="'/tickets/new'">
            <i class="fa fa-plus"></i>
            Create ticket
        </router-link>
    </card-layout>
    `
};
