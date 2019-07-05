<template>
    <div>
        <div class="register">
            <h1 class="title is-3 has-text-centered" style="margin-bottom: 1rem !important;">Registrieren</h1>
        
            <article class="message is-danger" v-if="fehler">
                <div class="message-body">{{fehler}}</div>
            </article>
            
            <div class="field">
                <p class="control has-icons-left">
                    <input class="input" type="text" placeholder="Name" v-model="name">
                    <span class="icon is-small is-left">
                        <i class="remixicon-user-line"></i>
                    </span>
                </p>
            </div>
        
            <div class="field">
                <p class="control has-icons-left">
                    <input class="input" type="email" placeholder="Email" v-model="mail">
                    <span class="icon is-small is-left">
                        <i class="remixicon-at-line"></i>
                    </span>
                </p>
            </div>
            
            <div class="field">
                <p class="control has-icons-left">
                    <input class="input" type="password" placeholder="Passwort" v-model="pw">
                    <span class="icon is-small is-left">
                        <i class="remixicon-lock-password-line"></i>
                    </span>
                </p>
            </div>

            <div class="field" v-if="keyReq">
                <p class="control has-icons-left">
                    <input type="text" placeholder="Registrierungsschlüssel" class="input">
                    <span class="icon is-left">
                        <i class="remixicon-key-line"></i>
                    </span>
                </p>
            </div>
        
            <div class="field">
                <label class="checkbox">
                    <input type="checkbox" v-model="agb">
                    Mit der Registrierung stimmen Sie zu, dass alle eingegeben Daten auf dem Server des Betreibers gespeichert werden.
                </label>
            </div>
        
            <div class="field is-grouped">
                <p class="control">
                    <a  class="button is-success" @click="register">Registrieren</a>
                </p>
                <p class="control">
                    <router-link to="/login" class="button is-small">
                         zurück zur Anmeldung
                    </router-link>
                </p>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    data(){
        return {
            fehler: "",
            agb: false,
            pw: "",
            name: "",
            mail: "",
            key: "",
            keyReq: false
        }
    },
    methods: {
        async register() {
            let r = await api.POST("/register", {
                "mail": this.mail,
                "name": this.name,
                "password": this.pw,
                "key": this.key
            })
            if (!r.ok) this.fehler = r.content;
            if (r.ok) this.$router.push("/login?register=true");
        }
    }
}
</script>
