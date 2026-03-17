import { createRouter, createWebHistory } from "vue-router";
import App_Body from "./../components/App_Body.vue";
import App_Login from "./../components/App_Login.vue";
import App_Books from "./../components/App_Books.vue";
import App_Book from "./../components/App_Book.vue";
import App_Books_Admin from "./../components/App_Books_Admin.vue";
import App_Book_Edit from "./../components/App_Book_Edit.vue";
import App_Users from "./../components/App_Users.vue";
import App_User_Edit from "./../components/App_User_Edit.vue";

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
    {
        path: "/books",
        name: "Books",
        component: App_Books,
    },
    {
        path: "/books/:bookName",
        name: "Book",
        component: App_Book,
    },
    {
        path: "/admin/books",
        name: "Books_Admin",
        component: App_Books_Admin,
    },
    {
        path: "/admin/books/:bookId",
        name: "Book_Edit",
        component: App_Book_Edit,
    },
    {
        path: "/admin/users",
        name: "Users",
        component: App_Users,
    },
    {
        path: "/admin/users/:userId",
        name: "User_Edit",
        component: App_User_Edit,
    },
];

const router = createRouter({ history: createWebHistory(), routes });
export default router;
