<template>
    <div>
        <div class="level">
      <div class="level-left">
        <router-link class="button" to="/reisen">
        <span class="icon"><i class="remixicon-reply-line"></i></span>
        <span>zurück</span>
      </router-link>
      </div>
      <div class="level-left">
        <h1 class="title">{{reise.name}}</h1>
      </div>
      <div class="level-right">
        <div class="field">
          <p class="control">
            <router-link class="button" :to="$route.path + '/neu'"><span class="icon"><i class="remixicon-add-line"></i></span><span>Neue Ausgabe</span></router-link>
          </p>
        </div>
      </div>
    </div>
        <table class="table is-fullwidth">
      <thead>
        <th>Datum</th>
        <th>Betreff</th>
        <th>Bezahlt von</th>
        <th>Bezahlt für</th>
        <th>Betrag</th>
        <th></th>
      </thead>
      <tbody>
        <tr v-for="beleg in reise.belege" :key="beleg.id">
          <td>{{beleg.datum}}</td>
          <td>{{beleg.name}}</td>
          <td>
            <span class="tag is-primary is-rounded">{{beleg.von.name}}</span>
          </td>
          <td>
            <div class="tags">
              <span v-for="p in beleg.an" class="tag is-rounded is-info">{{p.name}}</span>
            </div>
          </td>
          <td>
            <span>{{ Number(beleg.betrag).toFixed(2) }}€</span>
          </td>
          <td>
            <div class="buttons has-addons are-small">
              <span class="button is-warning" @click="edit(beleg.id)">
                <i class="remixicon-pencil-line"></i>
              </span>
              <span class="button is-danger" @click="remove(beleg.id)">
                <i class="remixicon-delete-bin-5-line"></i>
              </span>
            </div>
          </td>
        </tr>
        <tr style="cursor: default; background-color: transparent;">
            <td></td>
            <td></td>
            <td></td>
            <td></td>
            <td class="summe"><strong>{{summe}}</strong></td>
        </tr>
      </tbody>
    </table>
</div>
</template>

<script>
export default {
    data() {
      setTimeout(() => this.load());
        return {
            reise: {}
        }
    },
    methods: {
      async remove(id) {
        if(confirm("Löschen bestätigen")) {
          let r = await api.DELETE("/reisen/"+encodeURIComponent(this.$route.params.id)+"/beleg/"+encodeURIComponent(id));
          if (r.ok) this.load();
        }
      },
      async load() {
          let r = await api.GET("/reisen/"+ encodeURIComponent(this.$route.params.id));
            if (r.ok && r.content) this.reise = r.content;
            if (r.status == 401) this.$router.push("/")
        },
      edit() {
        alert("Bearbeiten ist aktuell noch nicht implementiert.")
      }
    },
    computed:{
        summe(){
            return Number(this.reise.belege ? this.reise.belege.reduce((pv, cv) => pv + parseFloat(cv.betrag), 0) : 0).toFixed(2) + "€";
        }
    },
}
</script>
