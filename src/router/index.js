import { createRouter, createWebHistory } from "vue-router";
import App_Body from "./../components/App_Body.vue";
import App_Login from "./../components/App_Login.vue";

const routes = [
    {
        path: "/",
        name: "Home",
        component: App_Body,
    },
    {
        path: "/login",
        name: "Login",
        component: App_Login,
    },
];

const router = createRouter({ history: createWebHistory(), routes });
export default router;
