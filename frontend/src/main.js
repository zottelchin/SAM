import "babel-polyfill";
import Vue from 'vue';
import App from './App.vue';

import RestAPI from './rest-api';

window.api = new RestAPI("/api", { headers: { "Authorization": "Bearer abc" } });

new Vue({
    el: '#app',
    render: h => h(App)
});