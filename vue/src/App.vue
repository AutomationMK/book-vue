<template>
    <App_Header />
    <div>
        <router-view />
    </div>
    <App_Footer />
</template>

<script>
    import App_Header from "./components/App_Header.vue"
    import App_Footer from "./components/App_Footer.vue"
    import {store} from './components/store.js';

    const getCookie = (name) => {
        return document.cookie.split("; ").reduce((r, v) => {
            const parts = v.split("=");
            return parts[0] === name ? decodeURIComponent(parts[1]) : r;
        }, "");
    }

    export default {
        name: 'App',
        components: {
            App_Header,
            App_Footer,
        },
        data() {
            return {
                store
            }
        },
        beforeMount() {
            // check for a cookie
            let data = getCookie("_site_data");

            if (data !== "") {
                let cookieData = JSON.parse(data);

                // update store
                store.token = cookieData.token.token;
                store.user = {
                    id: cookieData.user.id,
                    first_name: cookieData.user.first_name,
                    last_name: cookieData.user.last_name,
                    email: cookieData.user.email,
                }
            }
        }
    }
</script>

<style>

</style>
