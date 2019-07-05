<template>
    <div>
        <h1 class="title is-3 has-text-centered" style="margin-bottom: 1rem !important;">Login</h1>
        <article class="message is-danger" v-if="fehler">
            <div class="message-body" >{{fehler}}</div>
        </article>
      
        <div class="field">
            <p class="control has-icons-left">
                <input class="input" type="email" placeholder="Email">
                <span class="icon is-small is-left">
                    <i class="remixicon-at-line"></i>
                </span>
            </p>
        </div>
      
        <div class="field">
            <p class="control has-icons-left">
            <input class="input" type="password" placeholder="Password">
            <span class="icon is-small is-left">
                <i class="remixicon-lock-password-line"></i>
            </span>
            </p>
        </div>
      
        <div class="field is-grouped">
            <p class="control">
                <button class="button is-success">Login</button>
            </p>
            <p class="control">
                <a class="button is-small" @click="reset = true">Passwort vergessen</a>
            </p>
            <p class="control">
                <a class="button is-info is-small" @click="register = true">Registrieren</a>
            </p>
        </div>

    <div class="login" v-if="reset">
      <a @click="reset = false">
        <i class="remixicon-arrow-go-back-line"></i> Zurück
      </a>
      <h1
        class="title is-3 has-text-centered"
        style="margin-bottom: 1rem !important;"
      >Passwort zurück setzten</h1>
      <article class="message is-danger">
  <div class="message-body" v-if="fehler">
    {{fehler}}
  </div>
</article>
      <div class="field">
        <p class="control has-icons-left">
          <input class="input" type="email" placeholder="Email">
          <span class="icon is-small is-left">
            <i class="remixicon-at-line"></i>
          </span>
        </p>
      </div>
      <div class="field">
        <p class="control">
          <button class="button is-info">Passwortlink zusenden</button>
        </p>
      </div>
    </div>
</template>

<script>
export default {
  data() {
    return {
      reset: false,
      register: false,
      fehler: ""
    };
  }
};
</script>
