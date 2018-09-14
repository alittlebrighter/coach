<template>
  <div>
    <form onsubmit="return false;">
      <div class="mdl-textfield">
        <input 
          v-model="tagQuery" 
          v-on:keyup.enter="fetchScripts(tagQuery);" 
          class="mdl-textfield__input" 
          type="text" 
          id="script-search" 
          placeholder="Search by tag" />
      </div>
      <button @click="fetchScripts" class="mdl-button mdl-js-button mdl-button--icon">
        <i class="far fa-search"></i>
      </button>
    </form>

    <h4>Results:</h4>
    <em v-show="scripts && scripts.length == 0">None</em>
    <div class="mdl-grid">
      <script-summary v-for="script in scripts" :key="script.id" :script="script" class="mdl-cell mdl-cell--4-col" />
    </div>
  </div>
</template>

<script>
import store from "@/store";
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
  created () {
    this.update();
  },
  watch: {
    '$route': 'update'
  },
  methods: {
    update () {
      var query = store.get("tag-query");
      this.tagQuery = query === "undefined" ? "" : query;
      this.fetchScripts();
    },
    fetchScripts() {
      store.set("tag-query", this.tagQuery);
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
