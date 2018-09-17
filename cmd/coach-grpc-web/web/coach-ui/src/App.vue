<template>
<div id="app" class="mdl-layout mdl-js-layout mdl-layout--fixed-header">
  <header class="mdl-layout__header">
    <div class="mdl-layout__header-row">
      <!-- Title -->
      <span class="mdl-layout-title">{{ title }}</span>
    </div>
    <!-- Tabs -->
    <div class="mdl-layout__tab-bar mdl-js-ripple-effect" id="tab-bar">
      <router-link :class="{'mdl-layout__tab': true, 'is-active': $route.path == '/'}" to="/"><i class="far fa-search">
        </i>&nbsp;Scripts
      </router-link>
      <router-link v-show="$route.path.indexOf('/run') > -1" :class="{'mdl-layout__tab': true, 'is-active': $route.path.indexOf('/run') > -1 }" to="/run">
        <i class="fas fa-play"></i>&nbsp;Run
      </router-link>
      <router-link v-show="$route.path.indexOf('/edit') > -1" :class="{'mdl-layout__tab': true, 'is-active': $route.path.indexOf('/edit') > -1 }" to="/edit">
        <i class="far fa-edit"></i>&nbsp;Edit
      </router-link>
    </div>
  </header>
  <div class="mdl-layout__drawer">
    <span class="mdl-layout-title">{{ title }}</span>
  </div>
  <main class="mdl-layout__content">
    <section class="mdl-layout__tab-panel is-active">
      <div class="page-content"><router-view /></div>
    </section>
  </main>
</div>
</template>

<script>
export default {
  data () {
    return {
      title: "CoachUI"
    };
  }
}
</script>


<style lang="scss">
@import "@/assets/styles/material-design-lite/material-design-lite.scss";
@import "assets/fontawesome/css/all.min.css";

main {
  margin-left: 1em;
}

#tab-bar {
  height: auto;
}

.fab {
  position: absolute;
  right: 3em;
  bottom: 3em;
  z-index: 100;
}

.success {
    background-color: green !important;
}

.fail {
    background-color: red !important;
}
</style>
