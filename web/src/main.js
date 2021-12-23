import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import JsonViewer from 'vue-json-viewer'
import Axios from "axios";


Vue.prototype.$http = Axios;
Vue.config.productionTip = false
Vue.use(JsonViewer)
Vue.prototype.$http.defaults.baseURL = process.env.VUE_APP_API_ENDPOINT;

new Vue({
  vuetify,
  render: h => h(App)
}).$mount('#app')
