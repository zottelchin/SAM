<template>
  <div class="container">
    <header>
      <router-link to="/">
        <img src="../goldtopf.png" class="image">
      </router-link>
        <h1 class="title has-text-centered full">Supertolles Ausgabenmanagement</h1>
        <h1 class="title has-text-centered small">SAM</h1>
      <div class="gravatar">
        <div style="flex: 1;">
          <v-gravatar email="test@zottelchin.de" alt="Nobody" default-img="retro" :size="35"/>
        </div>
      </div>
    </header>
    <router-view id="container"></router-view>
    <footer class="footer">
      <div class="content has-text-centered">
        <p>
          <a href="https://github.com/zottelchin/sam">SAM</a>
          &#9702;
          Supertolles Ausgabenmanagement von
          <a
            href="https://zottelchin.de"
          >Phillipp Engelke</a>
        </p>
      </div>
    </footer>
  </div>
</template>

<script>
import Vue from "vue";
import VueRouter from "vue-router";
import Home from "./component/home.vue";
import Login from "./component/login.vue";
import Gravatar from "vue-gravatar";
import Übersicht from "./component/reisen.vue";
import Reise from "./component/reise.vue";
import Neu from "./component/beleg.vue"

Vue.component("v-gravatar", Gravatar);
Vue.use(VueRouter);

const router = new VueRouter({
  mode: "history",
  routes: [
    { path: "/", component: Home },
    { path: "/login", component: Login },
    { path: "/reisen", component: Übersicht },
    { path: "/reisen/:id", component: Reise },
    { path: "/reisen/:id/neu", component: Neu}
  ]
});

export default {
  router,
  data() {
    return {
      reise: {}
    };
  }
};
</script>

<style lang="sass">
@charset "utf-8"
$family-sans-serif: Coiny, BlinkMacSystemFont, -apple-system, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", "Helvetica", "Arial", sans-serif
$footer-padding: 2rem 1.5rem 2rem
$table-row-hover-background-color: hsl(179, 81%, 86%)

@import "../bulma-0.7.5/bulma.sass"
@import "./sam.sass"
</style>