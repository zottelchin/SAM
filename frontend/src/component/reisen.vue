<template>
  <div>
    <h2 class="title is-4">Meine Ausgaben</h2>
    
    <table class="table is-fullwidth reisen" v-if="reisen.length>0">
      <thead>
        <th>ID</th>
        <th>Name</th>
        <th>Beteiligte</th>
        <th>Status</th>
      </thead>
      <tr v-for="r in reisen" :key="r.ID"  @click="open(r.ID)">
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
    <small v-if="reisen.length == 0">Du hast noch keine Ausgaben in deiner Liste. Bitte andere Menschen dich einem Projekt hinzuzufügen oder erstelle ein eigenes mit dem Button oben rechts.</small>

    <h2 class="title is-4">Erstelle eine neue Ausgabensammlung</h2>
    <div class="field is-grouped">
      <p class="control">
        <input type="text" class="input" v-model="name" placeholder="Name">
      </p>
      <p class="control">
        <button class="button" @click="add()">Erstellen</button>
      </p>
    </div>

    <hr>

    <div>
      <p><strong>Du bist eingelogged als:&nbsp;</strong> {{meinName}}</p>
      <p>Zu deinen Einstellungen:&nbsp; <router-link to="/einstellungen">Einstellungen</router-link></p>
      <p><a><i class="remixicon-logout-box-line" styl></i>Abmelden</a></p>
    </div>
  
  </div>
</template>

<script>
export default {
    data() {
      setTimeout(() => this.load())
        return {
            reisen: [{"test":"53"}],
            name: "",
            meinName: ""
        }
    },
    methods:{
        open(id) {
            this.$router.push("/reisen/"+ encodeURIComponent(id));
        },
        async load() {
          let r = await api.GET("/reisen");
          if (r.ok) {
            this.reisen = r.content;
            let b = await api.GET("/me");
            if (b.ok) this.meinName = b.content.name;
          }
          if (r.status == 401) this.$router.push("/")

        },
        async add() {
          let r = await api.PUT("/reisen", {name: this.name})
          if (r.ok) this.load();
        }
    }
}
</script>
