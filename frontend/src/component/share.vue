<template>
   <div>

    <div v-if="toggle">

        <div class="content">
            <p>Diese Freigabe ist durch ein Passwort geschützt. Wenn Sie kein Passwort erhalten haben, wenden Sie sich bitte an den Sender der Links.</p>
        </div>
        <div class="field has-addons">
            <p class="control has-icons-left">
            <input class="input" type="password" placeholder="Password" v-model="passwort" @keyup.enter="load">
            <span class="icon is-small is-left">
                <i class="remixicon-lock-password-line"></i>
            </span>
            </p>
            <p class="control">
                <span class="button" @click="load()">Absenden</span>
            </p>
        </div>

    </div>
    
    <div class="level" v-if="!toggle">
      <div class="level-left"></div>

      <div class="level-left">
        <h1 class="title">Share: {{reise.name}}</h1>
      </div>
      
      <div class="level-right"></div>
    </div>


    <table class="table is-fullwidth" v-if="!toggle">
        <thead>
            <th>Datum</th>
            <th>Betreff</th>
            <th>Bezahlt von</th>
            <th>Bezahlt für</th>
            <th>Betrag</th>
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
          <td width="1">
            <span>{{ Number(beleg.betrag / 100).toFixed(2) }}€</span>
          </td>
        </tr>
        <tr style="cursor: default; background-color: transparent;">
            <td colspan="4"><span class="is-pulled-right">Summe:</span></td>
            <td><strong>{{summe}}</strong></td>
        </tr>
      </tbody>
    </table>

    <div v-if="!toggle">
    <h3 class="title is-5 ">Folgende &Uuml;berweisungen begleichen alle Schulden:</h3>
    <table class="table is-fullwidth">
      <tr v-for="r in Überweisen" :style="{'text-decoration':r.done ? 'line-through' : ''}">
        <td>{{reise.mitreisende.find(x => x.mail == r.from).name}}</td>
        <td style="text-decoration: none !important;"><i class="remixicon-arrow-right-line"></i></td>
        <td>{{reise.mitreisende.find(x => x.mail == r.to).name}}</td>
        <td>{{ Number(r.amount / 100).toFixed(2) }}€</td>
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
            passwort: "",
            toggle: true,
        }
    },
    methods: {
      async load() {
            this.toggle = false;
            let r = await api.GET("/share/"+ encodeURIComponent(this.$route.params.key) + "/" + encodeURIComponent(this.passwort));
            console.log(r.content);
            if (r.ok && r.content) this.reise = r.content;
            if (r.status == 404) this.$router.push("/404");
            if (r.status == 402) this.toggle = true;
            this.Überweisen = this.getTransactionsFromDeltas(this.calucalteDelta());
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
    }},
    computed:{
        summe(){
            return (Number(this.reise.belege ? this.reise.belege.reduce((pv, cv) => pv + parseFloat(cv.betrag), 0) : 0)/100).toFixed(2) + "€";
        }
    },
}
</script>


