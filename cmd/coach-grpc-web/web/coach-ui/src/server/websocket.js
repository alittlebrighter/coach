import uuid from "uuid/v1";

var socket = null,
  socketReady = false,
  requests = {},
  todo = [];

function init() {
  if (socket !== null) {
    return;
  }

  socket = new WebSocket("ws://" + window.location.hostname + (location.port ? ':'+location.port: '') + "/ws/rpc")

  socket.onopen = () => {
    socketReady = true;

    for (var i = 0; i < todo.length; i++) {
      todo[i]();
    }
  };

  socket.onmessage = e => {
    var response = JSON.parse(e.data);
    
    if (!requests[response.id]) {
      return;
    }

    var unsub = () => {
        delete requests[response.id];
      };

    requests[response.id](response, unsub);
  };
}

const TAG_WILDCARD = "?";
const EOF = "!!!EOF!!!";

function fetchScripts(query, cb) {
  if (!socketReady) {
    todo.push(() => { fetchScripts(query, cb); });
    return;
  }

  var payload = {
    id: uuid(),
    method: "getScripts",
    input: query || TAG_WILDCARD
  };
  requests[payload.id] = cb;

  socket.send(JSON.stringify(payload));
}

function saveScript(script, overwrite, cb) {
  if (!socketReady) {
    todo.push(() => { saveScript(script, overwrite, cb); });
    return;
  }

  var payload = {
    id: uuid(),
    method: "saveScript",
    input: script || {}
  };
  requests[payload.id] = cb;

  socket.send(JSON.stringify(payload));
}

function runScript(alias, cb) {
  if (!socketReady) {
    todo.push(() => { runScript(alias, cb); });
    return;
  }

  var payload = {
    id: uuid(),
    method: "runScript",
    input: alias
  };
  requests[payload.id] = cb;

  socket.send(JSON.stringify(payload));
  
  return payload.id;
}

function sendInput(id, input) {
  if (!socketReady) {
    todo.push(() => { sendInput(id, input); });
    return;
  }

  var payload = {
    id: id,
    method: "runScript",
    input: input
  };

  socket.send(JSON.stringify(payload));
}

export default () => {
  init();
  return {
    TAG_WILDCARD,
    EOF,
    fetchScripts,
    saveScript,
    runScript,
    sendInput
  }
};
