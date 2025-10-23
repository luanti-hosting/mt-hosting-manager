
export default {
    props: ["items"],
    methods: {
        get_icon_class: function(item) {
            const cl = {};
            if (item.icon){
                cl.fa = true;
                cl["fa-" + item.icon] = true;
            }
            return cl;
        }
    },
    template: /*html*/`
    <nav>
        <ol class="breadcrumb bg-primary-subtle rounded p-2">
            <li class="breadcrumb-item" v-for="item in items">
                <router-link :to="item.link">
                    <i v-bind:class="get_icon_class(item)" v-if="item.icon"></i>
                    {{item.name}}
                </router-link>
            </li>
        </ol>
    </nav>

    `
};

export const REGISTER = { name: "Register", icon: "user-plus", link: "/register" };
export const HOME = { name: "Home", icon: "home", link: "/" };
