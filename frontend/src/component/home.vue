<template>
<div>
    
    <div class="home">
    <div class="items">
        <div class="register">
            <h1 class="title is-3 has-text-centered" style="margin-bottom: 1rem !important;">Login</h1>

            <article class="message is-danger" v-if="fehlerLinks">
                <div class="message-body">{{fehlerLinks}}</div>
            </article>
        
            <div class="field">
                <p class="control has-icons-left">
                    <input class="input" type="email" placeholder="Email" v-model="mail">
                    <span class="icon is-small is-left">
                        <i class="remixicon-at-line"></i>
                    </span>
                </p>
            </div>
        
            <div class="field" v-show="!reset">
                <p class="control has-icons-left">
                <input class="input" type="password" placeholder="Password" v-model="pw" @keyup.enter="login">
                <span class="icon is-small is-left">
                    <i class="remixicon-lock-password-line"></i>
                </span>
                </p>
            </div>
        
            <div class="field is-grouped">
                <p class="control" v-show="!reset">
                    <button class="button is-success" @click="login">Login</button>
                </p>
                <p class="control" v-show="!reset">
                    <a class="button is-small" @click="reset = true">Passwort vergessen</a>
                </p>
                <p class="control" v-show="reset">
                    <button class="button is-info">Passwortlink zusenden</button>
                </p>
                <p class="control" v-show="reset">
                    <a class="button is-small" @click="reset = false">
                         zurück zur Anmeldung
                    </a>
                </p>
            </div>
        </div>
    </div>

    <div class="items">
       <div class="register">
            <h1 class="title is-3 has-text-centered" style="margin-bottom: 1rem !important;">Registrieren</h1>
        
            <article class="message is-danger" v-if="fehlerRechts">
                <div class="message-body">{{fehlerRechts}}</div>
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
        
            <div class="field">
                <p class="control">
                    <a  class="button is-success" @click="register">Registrieren</a>
                </p>
            </div>
        </div>
        </div>
    </div>

    <div class="content">
        <h4 class="is-4 title">
            Was ist Sam?
        </h4>
        <p>
            Sam steht für Supertolles Aufgabenmanagement. Es ist inspiriert von <a href="https://ihatemoney.org/">iHateMoney</a>. 
            Features sind:
            <ul>
                <li>Accounts mit Ausgabensammlungen verknüpfen</li>
                <li>Ausgabensammlungen per Link freigeben</li>
                <li>"leere" Personen zu Aufgabensammlungen hinzufügen, sodass nur eine Person einen Account braucht</li>
                <li>Per Docker einfach selber hosten</li>
            </ul>
        </p>
    </div>

</div>
</template>

<script>
export default {
    data(){
        return {
            reset: false,
            fehlerRechts: "",
            fehlerLinks: "",
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
            if (!r.ok) this.fehlerRechts = r.status;
            if (r.ok) await this.login();
        },
        async login() {
            let r = await api.POST("/login", {
                "mail": this.mail,
                "password": this.pw,
            });
            if(r.ok) {
                localStorage.setItem("hash", r.content.hash); 
                api.options.headers.Authorization = "Bearer " + r.content.hash; 
                this.$router.push("/reisen");
            }
            else this.fehlerLinks = r.status;
        }
    }
}
</script>