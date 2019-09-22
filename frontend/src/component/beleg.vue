<template>
  <div>
      <router-link class="button" :to="'/reisen/'+encodeURIComponent($route.params.id)">
        <span class="icon"><i class="remixicon-reply-line"></i></span>
        <span>zurück</span>
      </router-link>
      <div class="field">
    <h1 class="title has-text-centered">Neue Ausgabe anlegen</h1>
    </div>

    <div class="field is-horizontal">
      <div class="field-label is-normal">
        <label class="label">Betreff:</label>
      </div>
      <div class="field-body">
        <div class="field">
          <p class="control">
            <input class="input" type="text" placeholder="Betreff" v-model="name">
          </p>
        </div>
      </div>
    </div>

    <div class="field is-horizontal">
        <div class="field-label is-normal">
            <label class="label">Bezahlt von:</label>
        </div>
        <div class="field-body">
            <div class="field">
                <p class="control is-expanded">
                    <div class="select is-fullwidth">
                    <select v-model="by">
                        <option v-for="reisender in mitreisende" :key="reisender.mail" :value="reisender.mail">{{reisender.name}}</option>
                    </select>
                    </div>
                </p>
            </div>
        </div>
    </div>

    <div class="field is-horizontal">
        <div class="field-label is-normal">
            <label class="label">Bezahlt für:</label>
        </div>
        <div class="field-body">
            <div class="field">
                <div class="control" v-for="reisender in mitreisende" :key="reisender.mail">
                    <label class="checkbox">
                        <input type="checkbox" :value="reisender.mail" v-model="für">
                        {{reisender.name}}
                    </label>
                </div>
            </div>
        </div>
    </div>

    <div class="field is-horizontal">
      <div class="field-label is-normal">
        <label class="label">Betrag:</label>
      </div>
      <div class="field-body">
        <div class="field">
          <p class="control">
            <input class="input" type="number" placeholder="Betrag" v-model.number="value">
          </p>
        </div>
      </div>
    </div>

    <div class="field is-horizontal">
      <div class="field-label is-normal">
        <label class="label">Datum:</label>
      </div>
    <div class="field-body">
        <div class="field">
          <p class="control">
            <input class="input" type="date" v-model="time">
          </p>
        </div>
      </div>
    </div>

    <div class="field">
      <button class="button is-pulled-right" @click="save()">Speichern</button>
    </div>

  </div>
</template>

<script>
export default {
    data(){
        setTimeout(() => this.load());
        return { 
          value: 0,
          für: [],
          time: "",
          by: "",
          name: "",
          mitreisende: {} 
          }
        },
    methods: {
        async load() {
            let r = await api.GET("/reisen/"+ encodeURIComponent(this.$route.params.id));
            if (r.ok && r.content) this.mitreisende = r.content.mitreisende;
            if (r.status == 401) this.$router.push("/login?login=true")
        },
        async save() {
            //TODO: Error Handeling
            let r = await api.PUT("reisen/" + encodeURIComponent(this.$route.params.id) + "/beleg", {
                "name": this.name,
                "datum": this.time,
                "betrag": Math.trunc(this.value * 100),
                "von": {"mail": this.by},
                "an": this.für.map(x => {return {"mail": x};})
            });
            console.log(r);
            if (r.ok) this.$router.push("/reisen/" + encodeURIComponent(this.$route.params.id));
        }
    }
}
</script>


