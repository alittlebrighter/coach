<template>
  <div>
    <form onsubmit="return false;">
      <div class="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
        <input v-model="tagQuery" class="mdl-textfield__input" type="text" id="script-search">
        <label class="mdl-textfield__label" for="script-search">Search by Tag</label>
      </div>
      <button @click="fetchScripts" class="mdl-button mdl-js-button mdl-button--icon">
        <i class="far fa-search"></i>
      </button>
    </form>

    <h4>Results:</h4>
    <em v-show="scripts.length == 0">None</em>
    <div class="mdl-grid">
      <script-summary v-for="script in scripts" :key="script.id" :script="script" class="mdl-cell mdl-cell--4-col" />
    </div>
  </div>
</template>

<script>
import server from "@/server/websocket";
import ScriptSummary from "@/components/ScriptSummary.vue";

const ws = server();

export default {
  name: "home",
  data () {
    return {
      tagQuery: "",
      scripts: []
    };
  },
  methods: {
    fetchScripts () {
      console.log("fetching scripts");
      ws.fetchScripts(this.tagQuery, this.parseResponse);
    },
    parseResponse (response, unsub) {
      this.scripts = response.output.scripts;
      unsub();
    }
  },
  components: {
    ScriptSummary
  }
};
</script>
