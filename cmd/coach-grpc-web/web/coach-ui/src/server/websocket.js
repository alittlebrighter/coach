import uuid from "uuid/v1";

var socket = null,
  socketReady = false,
  requests = {};

function init() {
  if (socket !== null) {
    return;
  }

  socket = new WebSocket("ws://localhost:2015/ws/rpc")

  socket.onopen = () => {
    socketReady = true;
  };

  socket.onmessage = e => {
    var response = JSON.parse(e.data),
      unsub = () => {
        delete requests[response.id];
      };

    requests[response.id](response, unsub);
  };
}

const TAG_WILDCARD = "?";

function fetchScripts(query, cb) {
  if (!socketReady) {
    console.log("socket not ready yet");
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

function runScript(alias, cb) {
  if (!socketReady) {
    console.log("socket not ready yet");
    return;
  }

  var payload = {
    id: uuid(),
    method: "runScript",
    input: alias
  };
  requests[payload.id] = cb;

  socket.send(JSON.stringify(payload));
  console.log(JSON.stringify(payload));
  return payload.id;
}

function sendInput(id, input) {
  if (!socketReady) {
    console.log("socket not ready yet");
    return;
  }

  var payload = {
    id: id,
    method: "runScript",
    input: input
  };

  socket.send(JSON.stringify(payload));
  console.log(JSON.stringify(payload));
}

export default () => {
  init();
  return {
    TAG_WILDCARD,
    fetchScripts,
    runScript,
    sendInput
  }
};
