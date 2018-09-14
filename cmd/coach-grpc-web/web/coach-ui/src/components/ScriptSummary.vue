<template>
  <div class="mdl-card mdl-shadow--2dp">
  <div class="mdl-card__title mdl-card--expand mdl-color--blue-grey-100">
    <h2 class="mdl-card__title-text">
      {{ script.alias }}
    </h2>
  </div>
  <div class="mdl-card__supporting-text mdl-color--blue-grey-50" style="width: 100%;">
    <em>{{ script.tags }}</em><br><br>
    {{ script.documentation }}
  </div>
  <div class="mdl-card__actions mdl-card--border mdl-color--blue-grey-50">
    <router-link :to="'/run/' + script.alias" class="mdl-button mdl-js-button mdl-button--icon">
        <i class="fas fa-play"></i>
    </router-link>
    <router-link :to="{ name: 'edit', params: { script: script }}" class="mdl-button mdl-js-button mdl-button--icon">
      <i class="far fa-edit"></i>
    </router-link>
  </div>
</div>
</template>

<script>
export default {
  name: "ScriptSummary",
  props: {
    script: Object
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="css">
b.label {
  width: 14em;
  text-align: right;
}
</style>
