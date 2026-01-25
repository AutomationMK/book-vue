<template>
    <App_Header />
    <App_Body />
    <App_Footer />
</template>

<script>
    import App_Header from "./components/App_Header.vue"
    import App_Body from "./components/App_Body.vue"
    import App_Footer from "./components/App_Footer.vue"

    export default {
        name: 'App',
        components: {
            App_Header,
            App_Body,
            App_Footer,
        }
    }
</script>

<style>

</style>
