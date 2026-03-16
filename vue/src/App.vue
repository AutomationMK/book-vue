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
        }
    }
</script>

<style>

</style>
