import { search_tickets } from "../api/service_ticket.js";

const store = Vue.reactive({
    tickets: []
});

export const fetch_tickets = () => search_tickets({}).then(t => store.tickets = t);

export const get_open_tickets = () => store.tickets.filter(t => t.state == 'OPEN');

export const get_all_tickets = () => store.tickets;