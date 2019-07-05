<template>
    <div>
        <div class="register">
            <h1 class="title is-3 has-text-centered" style="margin-bottom: 1rem !important;">Login</h1>

            <article class="message is-danger" v-if="fehler">
                <div class="message-body">{{fehler}}</div>
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
                <input class="input" type="password" placeholder="Password" v-model="pw">
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
                <p class="control" v-show="!reset">
                    <router-link to="/register" class="button is-info is-small">Registrieren</router-link>
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
</template>

<script>
export default {
  data() {
    return {
      reset: false,
      register: false,
      fehler: "",
      mail: "",
      pw: "",
    };
  },
  methods: {
      async login() {
          let r = await api.POST("/login", {
              "mail": this.mail,
              "password": this.password,
          });
          if(r.ok) localStorage.setItem("hash", r.content.hash);
          else this.fehler = r.status;
      }
  }
};
</script>
