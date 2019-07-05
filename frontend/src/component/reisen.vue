<template>
  <div>
    <div class="level">
      <div class="level-left">
        <h1 class="title">Meine Ausgaben</h1>
      </div>
      <div class="level-right">
        <div class="field is-grouped">
          <p class="control">
            <input type="text" class="input" v-model="name" placeholder="Name">
          </p>
          <p class="control">
            <button class="button" @click="add()">Erstellen</button>
          </p>
        </div>
      </div>
    </div>

    <table class="table is-hoverable is-fullwidth reisen">
      <thead>
        <th>ID</th>
        <th>Name</th>
        <th>Beteiligte</th>
        <th>Status</th>
      </thead>
      <tr v-for="r in reisen" :key="r.ID" v-if="reisen.length>0" @click="open(r.ID)">
        <td>{{r.ID}}</td>
        <td>{{r.name}}</td>
        <td>
          <div class="tags">
            <span v-for="p in r.mitreisende" class="tag is-rounded is-info" :key="p.ID">{{p.name}}</span>
          </div>
        </td>
        <td>
          {{ r.archiviert ? "Abgeschlossen" : "Offen" }}
        </td>
      </tr>
    </table>
    <span v-if="reisen.length == 0">Du hast noch keine Ausgaben in deiner Liste. Bitte andere Menschen dich einem Projekt hinzuzufügen oder erstelle ein eigenes mit dem Button oben rechts.</span>
  </div>
</template>

<script>
export default {
    data() {
      setTimeout(() => this.load())
        return {
            reisen: [],
            name: ""
        }
    },
    methods:{
        open(id) {
            this.$router.push("/reisen/"+ encodeURIComponent(id));
        },
        async load() {
          let r = await api.GET("/reisen");
          if (r.ok) this.reisen = r.content;
          if (r.status == 401) this.$router.push("/login?login=true")
        },
        async add() {
          let r = await api.PUT("/reisen", {name: this.name})
          if (r.ok) this.load();
        }
    }
}
</script>
