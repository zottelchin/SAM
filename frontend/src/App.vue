<template>
<div class="body">
      <h1>
        <img src="../goldtopf.png">SAM
      </h1>
    <router-view id="container"></router-view>
    <footer>
      <div class="container">
        <div>
          <i class="remixicon-information-line"></i>
          <router-link to="/infos">Informationen</router-link>
        </div>
        <div>
          Logo von&nbsp;<a href="https://icons8.com/icon/35176/goldtopf">Icons8</a>
        </div>
        <div>
          <i>&copy;</i> <a href="http://zottelchin.de">Phillipp</a>
        </div>
        <div>
          <i class="remixicon-git-repository-line"></i>
          <a href="https://github.com/zottelchin/sam">SAM</a>
        </div>
        <div>
          <i class="remixicon-git-commit-line"></i>
          <a :href="github">Build vom: {{version.build}}</a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
import Vue from "vue";
import VueRouter from "vue-router";
import Home from "./component/home.vue";
import Übersicht from "./component/reisen.vue";
import Reise from "./component/reise.vue";
import Neu from "./component/beleg.vue";
import Info from "./component/Info.vue";
import ReiseEinstellungen from "./component/reise-einstellungen.vue";

Vue.use(VueRouter);

const router = new VueRouter({
  mode: "history",
  routes: [
    { path: "/", component: Home },
    { path: "/reisen", component: Übersicht },
    { path: "/reisen/:id", component: Reise },
    { path: "/reisen/:id/neu", component: Neu},
    { path: "/infos", component: Info},
    { path: "/reisen/:id/einstellungen", component: ReiseEinstellungen }
  ]
});

export default {
  router,
  data() {
     this.load();
    return {
      reise: {},
      version: {
        commit: "",
        version: "", 
        build: ""
      }
    };
  },
  methods: {
    async load() {
      let res = await api.GET("/version");
      this.version = res.content;
      let r = await api.GET("/config");
      window.key = r.content;
    }
  },
  computed: {
    github () {
      return "https://github.com/zottelchin/sam/commit/"+ this.version.commit;
    }
  }
};
</script>

<style lang="sass">
@charset "utf-8"
$footer-padding: 2rem 1.5rem 2rem
$table-row-hover-background-color: hsl(179, 81%, 86%)

@import "../bulma-0.7.5/bulma.sass"
@import "./sam.sass"
</style>