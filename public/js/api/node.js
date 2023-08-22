import { protected_fetch } from "./util.js";

export const get_all = () => protected_fetch(`api/node`);

export const create = n => protected_fetch(`api/node`, {
    method: "POST",
    body: JSON.stringify(n)
});

export const update = n => protected_fetch(`api/node/${n.id}`, {
    method: "POST",
    body: JSON.stringify(n)
});

export const remove = n => protected_fetch(`api/node/${n.id}`, {
    method: "DELETE"
});