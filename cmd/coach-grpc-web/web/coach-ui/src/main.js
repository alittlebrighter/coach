import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import "material-design-lite";

Vue.config.productionTip = false;

new Vue({
  router,
  render: h => h(App)
}).$mount("#app");
