export default {
    props: ["state"],
    template: /*html*/`
        <span class="badge bg-primary" v-if="state == 'OPEN'">Open</span>
        <span class="badge bg-success" v-if="state == 'RESOLVED'">Resolved</span>
        <span class="badge bg-secondary" v-if="state == 'CLOSED'">Closed</span>
    `
};
