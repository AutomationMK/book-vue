<template>
    <div class="container">
        <div class="row">
            <div class="col">
                <h1 class="mt-5">Login</h1>
                <hr>

                <form method="post" action="/login" class="needs-validation" novalidate>
                    <text-input label="email" type="email" name="email" required="true"></text-input>

                    <text-input label="password" type="password" name="password" required="true"></text-input>

                    <hr>
                    <input type="submit" class="btn btn-primary" value="Login">
                </form>
            </div>
        </div>
    </div>
</template>

<script>
    import TextInput from './forms/TextInput.vue'

    export default {
        name: 'App_Login',
        components: {
            TextInput,
        },
        mounted() {
            (function () {
                'use strict'

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
