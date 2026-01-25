<template>
    <form autocomplete="off" :method="method" :action="action" class="needs-validation" novalidate>
        <slot></slot>
    </form>
</template>

<script>
    export default {
        name: 'FormTag',
        props: ["method", "action"],
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
        }
    }
</script>
