import { protected_fetch } from "./protected_fetch.js";

export const get_messages = ticket_id => protected_fetch(`api/service/message/by-ticket/${ticket_id}`);

export const search_tickets = data => protected_fetch(`api/service/ticket/search`, {
    method: "POST",
    body: JSON.stringify(data)
});

export const create_ticket = data => protected_fetch(`api/service/ticket`, {
    method: "POST",
    body: JSON.stringify(data)
});

export const create_message = data => protected_fetch(`api/service/message`, {
    method: "POST",
    body: JSON.stringify(data)
});