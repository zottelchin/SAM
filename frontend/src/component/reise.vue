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
        <div class="field is-grouped">
          <p class="control">
            <router-link class="button" :to="$route.path + '/einstellungen'"><span class="icon"><i class="remixicon-settings-4-line"></i></span><span>Einstellungen</span></router-link>
          </p>
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
        <!-- <th>Tags</th> -->
        <th>Bezahlt von</th>
        <th>Bezahlt für</th>
        <th>Betrag</th>
        <th></th>
      </thead>
      <tbody>
        <tr v-for="beleg in reise.belege" :key="beleg.id">
          <td>{{beleg.datum}}</td>
          <td>{{beleg.name /*.replace(/#([^\s]+)/g, "")*/}}</td>
         <!--  <td>
            <div class="tags">
              <span class="tag is-rounded is-warning" v-for="tag in beleg.name.match(/#([^\s]+)/g)">{{tag.substring(1)}}</span>
            </div>
          </td> -->
          <td>
            <span class="tag is-primary is-rounded">{{beleg.von.name}}</span>
          </td>
          <td>
            <div class="tags">
              <span v-for="p in beleg.an" class="tag is-rounded is-info">{{p.name}}</span>
            </div>
          </td>
          <td width="1">
            <span>{{ Number(beleg.betrag / 100).toFixed(2) }}€</span>
          </td>
          <td width="1">
            <div class="buttons has-addons are-small is-pulled-right" style="flex-wrap: nowrap;">
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
            <td colspan="4"><span class="is-pulled-right">Summe:</span></td>
            <td><strong>{{summe}}</strong></td>
            <td><span class="button is-primary is-small has-text-weight-bold"  @click="toggle = !toggle"><i :class="{'remixicon-arrow-down-s-line': !toggle, 'remixicon-arrow-up-s-line' : toggle}"></i>Fall lösen</span></td>
        </tr>
      </tbody>
    </table>

    <div v-show="toggle">
    <h3 class="title is-5 ">Folgende &Uuml;berweisungen begleichen alle Schulden:</h3>
    <table class="table is-fullwidth">
      <tr v-for="r in Überweisen" :style="{'text-decoration':r.done ? 'line-through' : ''}">
        <td>{{reise.mitreisende.find(x => x.mail == r.from).name}}</td>
        <td style="text-decoration: none !important;"><i class="remixicon-arrow-right-line"></i></td>
        <td>{{reise.mitreisende.find(x => x.mail == r.to).name}}</td>
        <td>{{r.amount/100}}€</td>
        <td>Hier stehen Zahlungsinformationen</td>
        <td width="1"><span class="button is-success" @click="!r.done && begleichen(r)" :disabled="r.done"><i class="remixicon-check-line"></i></span></td>
      </tr>
    </table>
    <p v-if="Überweisen.length == 0">Alle Schulden sind bereits beglichen.</p>
    </div>
</div>
</template>

<script>
export default {
    data() {
      setTimeout(() => this.load(true));
        return {
            reise: {},
            Überweisen: [],
            toggle: false,
        }
    },
    methods: {
      async remove(id) {
        if(confirm("Löschen bestätigen")) {
          let r = await api.DELETE("/reisen/"+encodeURIComponent(this.$route.params.id)+"/beleg/"+encodeURIComponent(id));
          if (r.ok) this.load(true);
        }
      },
      async load(first) {
          let r = await api.GET("/reisen/"+ encodeURIComponent(this.$route.params.id));
            if (r.ok && r.content) this.reise = r.content;
            if (r.status == 401) this.$router.push("/");
            if (first) this.Überweisen = this.getTransactionsFromDeltas(this.calucalteDelta());
        },
      edit(id) {
        this.$router.push("/reisen/" + encodeURIComponent(this.$route.params.id) + "/edit/" + encodeURIComponent(id));
      },
      calucalteDelta() {
        //Wer hat wie viel Geld ausgegeben
        let PersonHaben = {};
        for (let j in this.reise.belege) {
          PersonHaben[this.reise.belege[j].von.mail] ? PersonHaben[this.reise.belege[j].von.mail] += this.reise.belege[j].betrag : PersonHaben[this.reise.belege[j].von.mail] = this.reise.belege[j].betrag ;
        }

        // Wie viel wurde für jeden ausgegeben
        let PersonSoll = {};
        for (let j in this.reise.belege) {
          let betrag = this.reise.belege[j].betrag / this.reise.belege[j].an.length;
          for (let i in this.reise.belege[j].an) {
            PersonSoll[this.reise.belege[j].an[i].mail] ? PersonSoll[this.reise.belege[j].an[i].mail] += betrag : PersonSoll[this.reise.belege[j].an[i].mail] = betrag;
          }
        }

        //Delta berechnen
        let delta = [];
        for (let j in this.reise.mitreisende) {
          delta.push({ "id": this.reise.mitreisende[j].mail, "delta": (PersonHaben[this.reise.mitreisende[j].mail] ? PersonHaben[this.reise.mitreisende[j].mail] : 0)  - (PersonSoll[this.reise.mitreisende[j].mail] ? PersonSoll[this.reise.mitreisende[j].mail] : 0)});
        }
        return delta;
      },
      getTransactionsFromDeltas(delta) {
        const debtors = delta.filter(x => x.delta < 0).sort((a, b) => a.delta - b.delta);
        const creditors = delta.filter(x => x.delta > 0);
        const result = [];

        for (let d in debtors) {
            for (let c in creditors.sort((a, b) => b.delta - a.delta)) {
                const amount = Math.min(-debtors[d].delta, creditors[c].delta);
                if (amount <= 0) continue;
                creditors[c].delta -= amount;
                debtors[d].delta += amount;
                result.push({ from: debtors[d].id, to: creditors[c].id, amount: Math.ceil(amount), done: false });
            }
        }

        return result.sort((a, b) => a.from > b.from ? 1 : a.from == b.from ? 0 : -1);
    },
    async begleichen(Überweisung) {
      //TODO: Error Handeling
            let r = await api.PUT("reisen/" + encodeURIComponent(this.$route.params.id) + "/beleg", {
                "name": "Rückzahlung",
                "datum": new Date().toISOString().slice(0,10),
                "betrag": Überweisung.amount,
                "von": {"mail": Überweisung.from},
                "an": [{"mail": Überweisung.to}]
            });
            await this.load(false);
            this.Überweisen.filter(x => x == Überweisung)[0].done = true;
    }
    },
    computed:{
        summe(){
            return (Number(this.reise.belege ? this.reise.belege.reduce((pv, cv) => pv + parseFloat(cv.betrag), 0) : 0)/100).toFixed(2) + "€";
        }
    },
}
</script>
