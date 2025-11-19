'use strict';

// requires rxjs

// TODO: state handling + auto reconnect ?

const _nsWs = {};
_nsWs.init = function() {
  _nsWs.id = '_nsWs';
  _nsWs.operational = false;
  _nsWs.instances = {}; // map[wsInstance.id] = wsInstance // TODO: use map

  _nsWs.requirementCheck = function() {
    if (window["WebSocket"]) {
      _nsWs.operational = true;
    } else {
      _nsWs.operational = false;
      console.warn(_nsWs.id, 'Your browser does not support WebSockets.');
    }

    // TODO: rxjs check

    return _nsWs.operational;
  };

  _nsWs.newWsInstance = function(id, uri) {
    const wsProtocol = window.location.protocol.includes('https') ? 'wss://' : 'ws://';
    const wsock = new WebSocket(wsProtocol + window.location.host + uri);

    const wsInstance = {
      id: id,
      uri: uri,
      wsock: wsock,

      // use standard subject for generic behavior
      onopenSubj: new rxjs.Subject(),
      oncloseSubj: new rxjs.Subject(),
      onmessageSubj: new rxjs.Subject(),
      onerrorSubj: new rxjs.Subject(),
    };

    wsInstance.wsock.onopen = function(ev) {
      wsInstance.onopenSubj.next(ev);
    };
    wsInstance.wsock.onclose = function(ev) {
      wsInstance.oncloseSubj.next(ev);

      // TODO: renew connection
    };
    wsInstance.wsock.onmessage = function(ev) {
      wsInstance.onmessageSubj.next(ev.data);
    };
    wsInstance.wsock.onerror = function(ev) {
      wsInstance.onerrorSubj.next(ev);

      console.error(_nsWs.id, `ws err: `, ev);
    };

    _nsWs.instances[id] = wsInstance;

    return wsInstance;
  };

  _nsWs.removeWsInstance = function(id) {
    const wsInstance = _nsWs.instances[id];

    if (wsInstance == undefined || wsInstance == null) {
    } else {
      wsInstance.wsock.close();
      wsInstance.onopenSubj.complete();
      wsInstance.oncloseSubj.complete();
      wsInstance.onmessageSubj.complete();
      wsInstance.onerrorSubj.complete();
      _nsWs.instances[id] = null;
    }
  };

  _nsWs.reconnect = function(id) {
    const wsInstance = _nsWs.instances[id];

    if (wsInstance == undefined || wsInstance == null) {
      console.warn(_nsWs.id, `Unable to reconnect on non existing WsInstance ${id}. use newWsInstance first.`);
      return null;
    }

    if (wsInstance.wsock == undefined || wsInstance.wsock == null) {
    } else {
      wsInstance.wsock.close();
    }

    const wsProtocol = window.location.protocol.includes('https') ? 'wss://' : 'ws://';
    const wsock = new WebSocket(wsProtocol + window.location.host + wsInstance.uri);
    wsInstance.wsock = wsock;

    wsInstance.wsock.onopen = function(ev) {
      wsInstance.onopenSubj.next(ev);
    };
    wsInstance.wsock.onclose = function(ev) {
      wsInstance.oncloseSubj.next(ev);
    };
    wsInstance.wsock.onmessage = function(ev) {
      wsInstance.onmessageSubj.next(ev.data);
    };
    wsInstance.wsock.onerror = function(ev) {
      console.error(_nsWs.id, `ws err: `, ev);
      wsInstance.onerrorSubj.next(ev);
    };
  };

  _nsWs.plainJsonParser = new rxjs.map((v) => {
    try {
      return JSON.parse(v);
    } catch(err) {
      console.warn(_nsWs.id, `Invalid JSON discarded...`, v);
      return null;
    }
  });

  _nsWs.requirementCheck();
};
