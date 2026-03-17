<template>
    <div class="container">
        <div class="row">
            <div class="col">
                <h1 class="mt-3">User Edit</h1>
                <hr>
                <form-tag v-if="this.ready" @userEditEvent="submitHandler" name="userform" event="userEditEvent">
                    <text-input v-model="user.first_name" type="text" required="true" label="First Name"
                        :value="user.first_name" name="first-name"></text-input>
                    <text-input v-model="user.last_name" type="text" required="true" label="Last Name"
                        :value="user.last_name" name="last-name"></text-input>
                    <text-input v-model="user.email" type="email" required="true" label="Email" :value="user.email"
                        name="email"></text-input>
                    <text-input v-if="this.user.id === 0" v-model="user.password" type="password" required="true"
                        label="Password" :value="user.password" name="password"></text-input>
                    <text-input v-else v-model="user.password" type="password" label="Password"
                        help="Leave empty to keep existing password" :value="user.password"
                        name="password"></text-input>
                    <hr>
                    <div class="float-start">
                        <input type="submit" class="btn btn-primary me-2" value="Save">
                        <router-link to="/admin/users" class="btn btn-outline-secondary">Cancel</router-link>
                    </div>
                    <div class="float-end">

                        <a v-if="(this.$route.params.userId > 0) && (parseInt(String(this.$route.params.userId), 10) !== store.user.id)"
                            class="btn btn-danger" href="javascript:void(0);"
                            @click="confirmDelete(this.user.id)">Delete</a>
                    </div>
                    <div class="clearfix"></div>
                </form-tag>
                <p v-else>Loading...</p>
            </div>
        </div>
    </div>
</template>

<script>
    import Security from './security.js'
    import FormTag from './forms/FormTag.vue'
    import TextInput from './forms/TextInput.vue'
    import notie from 'notie'
    import {store} from './store.js'
    import router from './../router/index.js'

    export default {
        beforeMount() {
            Security.requireToken();

            if (parseInt(String(this.$route.params.userId), 10) > 0) {
                // editing an existing user
                fetch(`${process.env.VUE_APP_API_URL}/admin/users/get/${this.$route.params.userId}`, Security.requestOptions(""))
                    .then((response) => response.json())
                    .then((data) => {
                        if (data.error) {
                            this.$emit('error', data.message);
                        } else {
                            this.user = data;
                            // we want password to be empty for existing users
                            this.user.password = "";
                        }
                        this.ready = true;
                    })
            } else {
                this.ready = true;
            }
        },
        data() {
            return {
                user: {
                    id: 0,
                    first_name: "",
                    last_name: "",
                    email: "",
                    password: "",
                },
                store,
                ready: false,
            }
        },
        components: {
            'form-tag': FormTag,
            'text-input': TextInput,
        },
        methods: {
            submitHandler() {
                const payload = {
                    id: parseInt(String(this.$route.params.userId), 10),
                    first_name: this.user.first_name,
                    last_name: this.user.last_name,
                    email: this.user.email,
                    password: this.user.password,
                }

                fetch(`${process.env.VUE_APP_API_URL}/admin/users/save`, Security.requestOptions(payload))
                    .then((response) => response.json())
                    .then((data) => {
                        if (data.error) {
                            this.$emit('error', data.message);
                        } else {
                            this.$emit('success', 'Changes saved!');
                            router.push("/admin/users");
                        }
                    }).catch((error) => {
                        this.$emit('error', error);
                    });
            },
            confirmDelete(id) {
                notie.confirm({
                    text: "Are you sure you want to delete this user?",
                    submitText: "Delete",
                    submitCallback: () => {
                        let payload = {
                            id: id,
                        }

                        fetch(`${process.env.VUE_APP_API_URL}/admin/users/delete`, Security.requestOptions(payload))
                            .then((response) => response.json())
                            .then((data) => {
                                if (data.error) {
                                    this.$emit('error', data.message);
                                } else {
                                    this.$emit('success', 'User deleted');
                                    router.push("/admin/users")
                                }
                            })
                    }
                })
            }
        }
    }
</script>
