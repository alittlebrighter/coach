<template>
  <div>
    <h2>Running {{ alias }}</h2>
    <div class="mdl-grid">
      <form class="mdl-cell mdl-cell--4-col">
        <div class="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
        <input v-model="args" class="mdl-textfield__input" type="text" id="args">
        <label class="mdl-textfield__label" for="args">command line arguments...</label>
        </div>
      </form>
      <div class="mdl-cell mdl-cell--8-col"></div>
    <button @click="run" class="mdl-button mdl-js-button mdl-button--raised mdl-js-ripple-effect mdl-button--colored mdl-cell mdl-cell--1-col action">
      <i v-show="lines.length == 0" class="fas fa-play"></i>
      <i v-show="lines.length > 0" class="fas fa-redo-alt"></i>
    </button>
    <button @click="edit" class="mdl-button mdl-js-button mdl-button--raised mdl-js-ripple-effect mdl-button--colored mdl-cell mdl-cell--1-col action">
      <i class="fas fa-edit"></i>
    </button>
    </div>

    <div v-show="lines.length > 0" class="run mdl-grid">
    <!--
    <ul class="mdl-list mdl-cell mdl-cell--10-col mdl-color--blue-grey-100">
      <li v-for="(line, key, i) in lines" :key="i" v-show="line.content !== EOF" class="mdl-list__item">
        <span :class="{'mdl-list__item-primary-content': true, red: line.error}">
          <span v-html="line.content"></span>
        </span>
      </li>
    </ul>
    -->
    <div class="mdl-color--blue-grey-800 mdl-color-text--light-green-200 mdl-cell mdl-cell--10-col console">
      <div v-for="(line, key, i) in lines" :key="i" v-show="line.content !== EOF" v-html="line.content" :class="{red: line.error}"></div>
      <form v-show="lines.length > 0 && (!stdoutEOF || !stderrEOF)" onsubmit="return false;" class="mdl-cell mdl-cell--4-col">
      > <div id="stdin-form" class="mdl-textfield mdl-js-textfield">
        <input v-model="stdin" v-on:keyup.enter="sendInput(stdin)" class="mdl-textfield__input" type="text" id="stdin-box">
      </div>
      <button @click="sendInput(stdin)" class="mdl-button mdl-js-button mdl-button--icon">
        <i class="fas fa-sign-in-alt"></i>
      </button>
    </form>
    </div>
    
    </div>

    <button @click="sendInput('coach::cancelRun')" v-show="isRunning" class="mdl-button mdl-js-button mdl-button--fab mdl-js-ripple-effect mdl-button--colored fab fail">
      <i class="fas fa-times"></i>
    </button>
  </div>
</template>

<script>
import server from "@/server/websocket";
import router from "@/router";

const ws = server();

export default {
  data() {
    var alias = this.$route.params.alias || "No script selected.";

    return {
      EOF: "!!!EOF!!!",
      alias: alias,
      args: "",
      stdin: "",
      isRunning: false,
      lines: [],
      requestId: "",
      stdoutEOF: true,
      stderrEOF: true
    };
  },
  methods: {
    run() {
      this.isRunning = true;
      this.stdoutEOF = false;
      this.stderrEOF = false;
      this.lines = [];
      this.stdin = "";
      this.requestId = ws.runScript(
        this.alias + " " + this.args,
        this.parseResponse
      );
      this.lines.push({ 
        content: "--- Running: '" + this.alias +
          (this.args.length > 0 ? "' with args '" + this.args : "") + "' ---",
        error: false 
      });
    },
    parseResponse(response, unsub) {
      if (response["output"]) {
        this.lines.push({
          content: response.output.replace(/\\n/g, "<br>"),
          error: false
        });
        this.stdoutEOF = response.output === this.EOF;
      }

      if (response["error"]) {
        this.lines.push({
          content: response.error.replace(/\\n/g, "<br>"),
          error: true
        });
        this.stderrEOF = response.error === this.EOF;
      }

      if (this.stdoutEOF && this.stderrEOF) {
        this.lines.push({ content: "--- Stopped ---", error: false });
        this.isRunning = false;
        unsub();
      }
    },
    sendInput(input) {
      ws.sendInput(this.requestId, input || this.stdin);
      this.stdin = "";
    },
    edit() {
      router.push("/edit/" + this.alias);
    }
  }
};
</script>

<style scoped>
.red {
  color: #ff4c4c;
}

form {
  vertical-align: bottom;
}

li {
  margin-top: 0.1em;
  margin-bottom: 0.1em;
}

#start-button {
  margin-left: 1.5em;
}

.console {
  padding: 1em;
  font-size: 1.5em;
  font-family: monospace !important;
}

.console > div {
  margin: 0.5em;
}

#stdin-box {
  border-color: rgb(178, 255, 89);
}
</style>

