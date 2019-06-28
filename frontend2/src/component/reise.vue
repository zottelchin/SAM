<template>
    <div>
        <div class="level">
      <div class="level-left">
        <h1 class="title">{{reise.name}}</h1>
      </div>
      <div class="level-right">
        <div class="field">
          <p class="control">
            <router-link class="button" :to="$route.path + '/neu'">Neue Ausgabe</router-link>
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
            <span>{{beleg.betrag}}€</span>
          </td>
          <td>
            <div class="buttons has-addons are-small">
              <span class="button is-warning">
                <i class="remixicon-pencil-line"></i>
              </span>
              <span class="button is-danger">
                <i class="remixicon-delete-bin-5-line"></i>
              </span>
            </div>
          </td>
        </tr>
        <tr>
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
    computed:{
        summe(){
            return this.reise.belege ? this.reise.belege.reduce((pv, cv) => pv + cv, 0) + "€" : 0 + "€";
        },
        async load() {
          let r = await api.GET("/reisen/"+ encodeURIComponent(this.$route.params.id));
            if (r.ok && r.content) this.reise = r.content;
        }
    },
}
</script>
