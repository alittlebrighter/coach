<template>
  <div>
    <h2 class="inline">Editing {{ script.alias }}</h2>

    <form class="mdl-grid" onsubmit="return false;">
        <button @click="save" class="mdl-button mdl-js-button mdl-button--raised mdl-js-ripple-effect mdl-button--colored mdl-cell mdl-cell--1-col action">
          <i class="fas fa-save"></i>
        </button>
        <button @click="run" class="mdl-button mdl-js-button mdl-button--raised mdl-js-ripple-effect mdl-button--colored mdl-cell mdl-cell--1-col action">
          <i class="fas fa-play"></i>
        </button>
        <div class="mdl-cell mdl-cell--10-col"></div>

        <label for="alias" class="mdl-cell mdl-cell--1-col">Alias</label>
        <input v-model="script.alias" class="mdl-textfield__input mdl-cell mdl-cell--4-col" type="text" id="alias">
        <div class="mdl-cell mdl-cell--7-col"></div>
      
        <label for="tags" class="mdl-cell mdl-cell--1-col">Tags</label>
        <input v-model="tagsString" class="mdl-textfield__input mdl-cell mdl-cell--4-col" type="text" id="tags" />
        <div class="mdl-cell mdl-cell--7-col"></div>
      
        <label for="shell" class="mdl-cell mdl-cell--1-col">Shell</label>
        <input v-model="script.script.shell" v-on:change="applyShell" class="mdl-textfield__input mdl-cell mdl-cell--4-col" type="text" id="shell">
        <div class="mdl-cell mdl-cell--7-col"></div>
      
        <label for="documentation" class="mdl-cell mdl-cell--1-col">Documentation</label>
        <textarea v-model="script.documentation" class="mdl-textfield__input mdl-cell mdl-cell--4-col" type="text" rows= "3" id="documentation" ></textarea>
        <div class="mdl-cell mdl-cell--7-col"></div>
      
      <label for="script-content" class="mdl-cell mdl-cell--1-col">Script</label>
      <codemirror v-model="script.script.content" :options="cmOptions" id="script-content" class="mdl-cell mdl-cell--10-col"></codemirror>
    </form>
  </div>
</template>

<script>
import { codemirror } from 'vue-codemirror';

// require styles
import 'codemirror/lib/codemirror.css';
// language js
import 'codemirror/mode/shell/shell.js';
import 'codemirror/mode/powershell/powershell.js';
import 'codemirror/mode/python/python.js';
import 'codemirror/mode/ruby/ruby.js';
import 'codemirror/mode/javascript/javascript.js';

// theme css
//import 'codemirror/theme/base16-light.css';

import store from "@/store";
import server from "@/server/websocket";
import router from "@/router";

const ws = server(),
  storeKey = "edit-script";

var empty = {
    alias: "new-script",
    tags: ["new"],
    documentation: "",
    script: {
        content: "",
        shell: navigator.userAgent.toLowerCase().indexOf("windows") > -1 ? "windowsCMD" : "bash"
    }
}

export default {
  data () {
    const cached = store.get(storeKey);

    var script = null;
    if (this.$route.params.script) {
      script = this.$route.params.script
    } else if (cached && cached !== "undefined") {
      script = JSON.parse(cached);
    } else {
      script = empty;
    }

    store.set(storeKey, JSON.stringify(script));    

    var data = {
      script: script,
      requestId: "",
      cmOptions: {
        lineNumbers: true,
        line: true,
        styleActiveLine: true
      }
    };

    var opts = this.newShell(data.script.script.shell);
    data.cmOptions.tabSize = opts.tabSize;
    data.cmOptions.mode = opts.mode;

    data.tagsString = this.tagsToString(data.script.tags);

    return data;
  },
  methods: {
    update () {
      var query = store.get("tag-query");
      this.tagQuery = query === "undefined" ? "" : query;
      this.fetchScripts();
    },
    save () {
        this.script.tags = this.stringToTags(this.tagsString);
        ws.saveScript(this.script, true, function(response, unsub) {
            console.log(JSON.stringify(response));
        });
    },
    run () {
        this.save();
        router.push('/run/' + this.script.alias);
    },
    newShell (shell) {
        var opts = {tabSize: 4, mode:"text/x-sh"}
        if (!shell) {
            return opts;
        }

        shell = shell.toLowerCase();

        if (shell.indexOf("node") > -1) {
            opts.tabSize = 2;
            opts.mode = "text/javascript";
        } else if (shell.indexOf("python") > -1) {
            opts.tabSize = 4;
            opts.mode="text/x-python";
        } else if (shell.indexOf("ruby") > -1) {
            opts.tabSize = 4;
            opts.mode="application/x-ruby";
        } else if (shell.indexOf("powershell") > -1) {
            opts.tabSize = 4;
            opts.mode="application/x-powershell";
        }

        return opts
    },
    applyShell () {
        var newOpts = this.newShell(this.script.script.shell);

        this.cmOptions.tabSize = newOpts.tabSize;
        this.cmOptions.mode = newOpts.mode;
    },
    tagsToString (tags) {
        return tags.join(", ")
    },
    stringToTags (tagsStr) {
        var tags = tagsStr.split(",");
        for (var i = 0; i < tags.length; i++) {
            tags[i] = tags[i].replace(/(^\s+|\s+$)/g, "");
        }
        return tags;
    }
  },
  components: {
      codemirror
  }
};
</script>

<style scoped>
label {
    text-align: right;
}

label:after {
    content: ":";
}
</style>
