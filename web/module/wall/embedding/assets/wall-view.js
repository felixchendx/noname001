'use strict';

// nsContent
const nsCt = {};
nsCt.init = function() {
  nsCt.CONSTANT = {
    WS_REQCODE_WALL_INFO: '/wall/info',
    WS_REQCODE_WALL_ITEM_INFO: '/wall/item/info',
  };

  nsCt.dat_panelInstances = [];

  nsCt.dat_wsInstance = null;
  nsCt.dat_wsObss = {};
  nsCt.dat_wsSubs = {};

  nsCt.fn_goToListing = function() {
    window.location = window.location.protocol + '//' + window.location.host + '/wall/wall/listing';
  };
  nsCt.fn_goToEdit = function() {
    window.location.pathname = '/wall/wall/detail';
  };

  nsCt.fn_populatePanelStuffs = function() {
    const panelElements = document.getElementsByClassName('wall-panel');

    for (let i = 0; i < panelElements.length; i++) {
      const el_panel = panelElements.item(i);
      
      const panelInstance = {
        el_panel: el_panel,
        el_status: el_panel.getElementsByClassName('wall-panel-status')[0],
        el_title: el_panel.getElementsByClassName('wall-panel-title')[0],
        el_btnReload: el_panel.getElementsByClassName('wall-panel-reload')[0],
        // el_btnInfo: el_panel.getElementsByClassName('wall-panel-info')[0],

        dat_index: i,
        dat_displayIndex: i + 1,
        dat_item: {
          source_node: '',
          stream_code: '',

          stream_state: '',
        },
      };

      panelInstance.fn_updateSource = function(sourceNode, streamCode) {
        panelInstance.dat_item.source_node = sourceNode;
        panelInstance.dat_item.stream_code = streamCode;

        let title = `#${panelInstance.dat_displayIndex}`;
        if (sourceNode === '' && streamCode === '') {
        } else {
          title += ` ${sourceNode} - ${streamCode}`;
        }

        panelInstance.el_title.innerHTML = title;
      };

      panelInstance.fn_updateState = function(streamState) {
        panelInstance.dat_item.stream_state = streamState;

        panelInstance.el_status.classList.remove(
          'status-indicator-animated',
          'status-off',
          'status-green', 'status-yellow', 'status-red',
        );

        switch (streamState) {
          case 'ls:new':
            panelInstance.el_status.classList.add('status-green', 'status-indicator-animated');
            break;

          case 'ls:inactive':
            panelInstance.el_status.classList.add('status-off');
            break;

          case 'ls:init:begin':
            panelInstance.el_status.classList.add('status-green', 'status-indicator-animated');
            break;

          case 'ls:init:fail':
            panelInstance.el_status.classList.add('status-red');
            break;

          case 'ls:init:ok':
            panelInstance.el_status.classList.add('status-green');
            break;

          case 'ls:reload:begin':
            panelInstance.el_status.classList.add('status-green', 'status-indicator-animated');
            break;

          case 'ls:reload:fail':
            panelInstance.el_status.classList.add('status-red');
            break;

          case 'ls:reload:ok':
            panelInstance.el_status.classList.add('status-green');
            break;

          case 'ls:destroy':
            panelInstance.el_status.classList.add('status-off');
            break;

          case 'ls:bg:fail':
            panelInstance.el_status.classList.add('status-red');
            break;

          default:
            panelInstance.el_status.classList.add('status-off');
            break;
        }
      };

      panelInstance.el_btnReload.addEventListener('click', () => { nsCt.fn_reloadStream(i); });
      // panelInstance.el_btnInfo.addEventListener('click', () => { nsCt.fn_showStreamInfo(i); });

      nsCt.dat_panelInstances[i] = panelInstance;
    }
  };

  nsCt.fn_reloadStream = function(idx) {
    _nsDLS.reloadStream(idx);
  };
  nsCt.fn_showStreamInfo = function(idx) {
    console.log(`showStreamInfo ${idx}`);
  };

  nsCt.fn_setupWs = function() {
    const wsId = 'wall-view';
    const wsUri = '/wall/wall/view/ws' + window.location.search;
    nsCt.dat_wsInstance = _nsWs.newWsInstance(wsId, wsUri);

    nsCt.dat_wsObss['open'] = nsCt.dat_wsInstance.onopenSubj.asObservable();
    nsCt.dat_wsObss['close'] = nsCt.dat_wsInstance.oncloseSubj.asObservable();
    nsCt.dat_wsObss['msg'] = nsCt.dat_wsInstance.onmessageSubj.asObservable().pipe(_nsWs.plainJsonParser);
    // nsCt.dat_wsObss['error'] = nsCt.dat_wsInstance.onerrorSubj.asObservable();

    nsCt.dat_wsObss['msg:ev'] = new rxjs.ReplaySubject(5);
    nsCt.dat_wsObss['msg:rr'] = new rxjs.Subject();

    nsCt.dat_wsSubs['open'] = nsCt.dat_wsObss['open']
      .subscribe({
        next: (v) => {
          const req = {
            _bid: '123',
            _brc: nsCt.CONSTANT.WS_REQCODE_WALL_INFO,
          };
          nsCt.dat_wsInstance.wsock.send(JSON.stringify(req));
        },
      });

    nsCt.dat_wsSubs['close'] = nsCt.dat_wsObss['close']
      .subscribe({
        next: (v) => {
          // TODO: parse close reason
          setTimeout(() => {
            _nsWs.reconnect('wall-view');
          }, 1000);
        },
      });

    nsCt.dat_wsSubs['msg:router'] = nsCt.dat_wsObss['msg']
      .pipe(
        rxjs.tap((msg) => {
          switch (msg._bt) {
            case '_bev': nsCt.dat_wsObss['msg:ev'].next(msg._bp); break;
            case '_brr': nsCt.dat_wsObss['msg:rr'].next(msg._bp); break;
            default:
              console.warn('nsCt', `no handler for msg_type '${msg._bt}'`);
              break;
          }
        })
      )
      .subscribe();

    nsCt.dat_wsSubs['msg:ev'] = nsCt.dat_wsObss['msg:ev']
      .pipe(
        rxjs.concatMap((p_ev) => rxjs.of(p_ev).pipe(rxjs.delay(333))),
      )
      .subscribe({
        next: (p_ev) => {
          let act = "";

          // TODO: action on state change instead of current state
          switch (p_ev.ev_code) {
            case 'ls:deactivated': act = 'do:stop'; break;
            case 'ls:init:begin': act = 'do:nothing'; break;
            case 'ls:init:fail': act = 'do:stop'; break;
            case 'ls:init:ok': act = 'do:reload'; break;
            case 'ls:disconnected': act = 'do:stop'; break;
            case 'ls:reload:begin': act = 'do:nothing'; break;
            case 'ls:reload:fail': act = 'do:stop'; break;
            case 'ls:reload:ok': act = 'do:reload'; break;
            case 'ls:destroyed': act = 'do:stop'; break;
            case 'ls:bg:fail': act = 'do:stop'; break;
            default:
              console.warn('nsCt', `no handler for ev_code '${p_ev.ev_code}'`);
              break;
          }

          switch (act) {
            case 'do:nothing':
              nsCt.dat_panelInstances[p_ev.item_index].fn_updateState(p_ev.stream_state);
              break;

            case 'do:reload':
              const req = {
                _bid: '234',
                _brc: nsCt.CONSTANT.WS_REQCODE_WALL_ITEM_INFO,
                item_index: p_ev.item_index,
              };
              nsCt.dat_wsInstance.wsock.send(JSON.stringify(req));
              break;

            case 'do:stop':
              nsCt.dat_panelInstances[p_ev.item_index].fn_updateState(p_ev.stream_state);
              // nsCt.fn_stopStream(p_ev.item_index);
              nsCt.fn_reloadStream(p_ev.item_index);
              break;

            default:
              break;
          }
        },
      });

    nsCt.dat_wsSubs['msg:rr'] = nsCt.dat_wsObss['msg:rr']
      .subscribe({
        next: (p_rep) => {
          switch (p_rep._brc) {
            case nsCt.CONSTANT.WS_REQCODE_WALL_INFO:
              const wallInfo = p_rep.wall_info;

              for (let i = 0; i < wallInfo.items.length; i++) {
                let itemInfo = wallInfo.items[i];

                nsCt.dat_panelInstances[i].fn_updateSource(itemInfo.source_node, itemInfo.stream_code);
                nsCt.dat_panelInstances[i].fn_updateState(itemInfo.stream_state);
              }
              break;
            case nsCt.CONSTANT.WS_REQCODE_WALL_ITEM_INFO:
              const itemInfo = p_rep.wall_item_info;

              nsCt.dat_panelInstances[itemInfo.item_index].fn_updateSource(itemInfo.source_node, itemInfo.stream_code);
              nsCt.dat_panelInstances[itemInfo.item_index].fn_updateState(itemInfo.stream_state);

              nsCt.fn_reloadStream(itemInfo.item_index);
              break;
            default:
              console.warn('nsCt', `no handler for req_code '${p_rep._brc}'`);
              break;
          }
        },
      });
  };
};

window.addEventListener("DOMContentLoaded", function() {
  _nsWs.init();

  // TODO: check event DOMContentLoaded, when does this fire ? after all assets done downloading ?
  // delay abit
  // so that browser's loading indicator indicates page done loading (page assets)
  // before the infinite hls traffic
  setTimeout(() => {
    _nsDLS.renderAll();
  }, 500);

  nsCt.init();
  nsCt.fn_populatePanelStuffs();
  nsCt.fn_setupWs();
});
