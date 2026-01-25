import { createRouter, createWebHistory } from "vue-router";
import App_Body from "./../components/App_Body.vue";

const routes = [
    {
        path: "/",
        name: "Home",
        component: App_Body,
    },
];

const router = createRouter({ history: createWebHistory(), routes });
export default router;
