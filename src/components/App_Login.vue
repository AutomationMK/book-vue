<template>
    <div class="container">
        <div class="row">
            <div class="col">
                <h1 class="mt-5">Login</h1>
                <hr>

                <form-tag @myevent="submitHandler">
                    <text-input v-model="email" label="Email" type="email" name="email" required="true"></text-input>
                    <text-input v-model="password" label="Password" type="password" name="password"
                        required="true"></text-input>

                    <hr>
                    <input type="submit" class="btn btn-primary" value="Login">
                </form-tag>
            </div>
        </div>
    </div>
</template>

<script>
    import FormTag from './forms/FormTag.vue';
    import TextInput from './forms/TextInput.vue';

    export default {
        name: 'App_Login',
        components: {
            FormTag,
            TextInput,
        },
        data() {
            return {
                email: "",
                password: "",
            }
        },
        methods: {
            submitHandler() {
                console.log("submitHandler called - success!");
            }
        },
        mounted() {
            (function () {
                'use strict';

                var forms = document.querySelectorAll('.needs-validation');
                Array.prototype.slice.call(forms)
                    .forEach(function (form) {
                        form.addEventListener('submit', function (event) {
                            if (!form.checkValidity()) {
                                event.preventDefault();
                                event.stopPropagation();
                            }

                            form.classList.add('was-validated');
                        }, false)
                    })
            })();
        },
    }
</script>
