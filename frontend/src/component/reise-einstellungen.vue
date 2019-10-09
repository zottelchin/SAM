<template>
    <div>

        <div class="level">
        <div class="level-left">
            <router-link class="button" :to="'/reisen/' + $route.params.id">
            <span class="icon"><i class="remixicon-reply-line"></i></span>
            <span>zurück</span>
        </router-link>
        </div>
        <div class="level-left">
            <h1 class="title">{{reise.name}}</h1>
        </div>
        <div class="level-right">
        </div>
        </div>

        <article class="message is-danger" v-if="error">
                <div class="message-body">{{error}}</div>
            </article>

        <label class="label">Titel</label>
        <div class="field has-addons">
        <div class="control is-expanded">
            <input class="input" type="text" v-model="reise.name">
        </div>
        <div class="control">
            <span class="button" @click="changeName">Ändern</span>
        </div>
        </div>

        

        <label class="label">Beteiligte Personen</label>
        <div class="field has-addons" v-for="u in reise.mitreisende">
            <div class="control is-expanded">
                <input class="input" disabled :value="u.name">
            </div>
            <div class="control">
                <span class="button" @click="remove(u.mail)">Entfernen</span>
            </div>
        </div>

        <label class="label">Nutzer hinzufügen <small style="opacity: 0.5;">(bestehenden Benutzer über seine E-Mail hinzufügen)</small></label>
        <div class="field has-addons">
        <div class="control is-expanded">
            <input class="input" type="email" v-model="mail" placeholder="mail@example.com">
        </div>
        <div class="control">
            <span class="button" @click="add()">Hinzufügen</span>
        </div>
        </div>

        <label  class="label">Anzeigenutzer hinzufügen <small style="opacity: 0.5;">(ist nur in diesem Projekt verwendbar)</small></label>
        <div class="field has-addons">
            <div class="control is-expanded">
                <input type="text" class="input" v-model="dummyName" placeholder="Peter">
            </div>
            <div class="control">
                <span class="button" @click="addDisplayUser">Hinzufügen</span>
            </div>
        </div>

        <label class="label">Reise teilen:</label>
        <div class="field has-addons">
            <div class="control is-expanded">
                <input type="text" class="input" v-model="shareName">
            </div>
            <div class="control">
                <span class="button" @click="addShare">Hinzufügen</span>
            </div>
        </div>

        <div class="field">
            <label class="label">Projektaktionen</label>
            <div class="buttons">
                <span class="button is-danger" @click="archive()">Archivieren</span>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    data() {
      setTimeout(() => this.load());
        return {
            reise: {},
            mail: "",
            error: "",
            dummyName: "",
            shareName: "",
        }
    },
    methods: {
      async load() {
          let r = await api.GET("/reisen/"+ encodeURIComponent(this.$route.params.id));
            if (r.ok && r.content) this.reise = r.content;
            if (r.status == 401) this.$router.push("/")
        },
        async remove(mail) {
            this.error = "";
            let r = await api.DELETE("/reisen/" + encodeURIComponent(this.$route.params.id) + "/mitreisende/" + encodeURIComponent(mail))
            if (!r.ok) this.error = "Nutzer Entfernen: " + r.status + " - " + r.content;
            if (r.ok) this.reise = r.content;
        },
        async add() {
            this.error = "";
            let r = await api.POST("/reisen/" + encodeURIComponent(this.$route.params.id) + "/mitreisende/", {"mail": this.mail})
            if (!r.ok) this.error = "Nutzer hinzufügen: " + r.status + " - " + r.content;
            if (r.ok) this.reise = r.content;
        },
        async archive() {
            this.error = "";
            let r = await api.POST("/reisen/" + encodeURIComponent(this.$route.params.id), {"archiviert": true});
            if (!r.ok) this.error = "Archivieren: " + r.status + " - " + r.content;
            if (r.ok) this.reise = r.content;
        },
        async changeName() {
            this.error = "";
            let r = await api.POST("/reisen/" + encodeURIComponent(this.$route.params.id), {"name": this.reise.name});
            if (!r.ok) this.error = "Name ändern: " + r.status + " - " + r.content;
            if (r.ok) this.reise = r.content;
        },
        async addDisplayUser() {
            this.error = "";
            let r = await api.PUT("/reisen/" + encodeURIComponent(this.$route.params.id) + "/mitreisende/", {"name": this.dummyName});
            if (!r.ok) this.error = "Nutzer hinzufügen: " + r.status + " - " + r.content;
            if (r.ok) this.reise = r.content;
        },
        async addShare() {
            this.error = "";
            let r = await api.PUT("/reisen/" + encodeURIComponent(this.$route.params.id) + "/share", {"key": this.shareName});
            if (!r.ok) this.error = "Nutzer hinzufügen: " + r.status + " - " + r.content;
        }
    }
}
</script>