'use strict';

// nsContent
const nsCt = { ns: 'nsCt' };
nsCt.init = function() {
  nsCt.el_navToWallListing = document.getElementById('navToWallListing');
  nsCt.el_navToWallView = document.getElementById('navToWallView');
  nsCt.el_submitDelete = document.getElementById('submitDelete');

  nsCt.fn_doSubmitDelete = function() {
    nsCt.el_submitDelete.click();
  };

  nsCt.el_navToWallListing.addEventListener('click', function() {
    window.location = nsCD.nav_wallListing;
  });

  if (nsCD == undefined || nsCD.isAddMode == undefined) {
    nsCt.el_navToWallView.disabled = true;
  } else {
    if (nsCD.isAddMode) {
      nsCt.el_navToWallView.disabled = true;
    } else {
      nsCt.el_navToWallView.addEventListener('click', function() {
        window.location = nsCD.nav_wallView;
      });
    }
  }
};
nsCt.fn_formatTimestamp = function(ts) {
  // temp, make date util stuffs
  // https://moment.github.io/luxon/#/formatting?id=table-of-tokens
  return luxon.DateTime.fromISO(ts).toFormat('dd LLLL yyyy, HH:mm:ss ZZZZ');
}

// websocket
const nsCtWs = { ns: 'nsCtWs' };
nsCtWs.CONSTANT = {
  WS_REQCODE_WALL_INFO: '/wall/info',
  WS_REQCODE_WALL_ITEM_INFO: '/wall/item/info',

  WS_REQCODE_NODE_INFO_LISTING  : '/node/info/listing',
  WS_REQCODE_STREAM_INFO_LISTING: '/stream/info/listing',
};
nsCtWs.init = function() {
  nsCtWs.dat_wsInstance = null;
  nsCtWs.dat_wsObss = {};
  nsCtWs.dat_wsSubs = {};

  nsCtWs.fn_wsSetup = function() {
    const wsId = 'wall-detail';
    const wsUri = nsCD.ws_uri;

    nsCtWs.dat_wsInstance = _nsWs.newWsInstance(wsId, wsUri);

    nsCtWs.dat_wsObss['open'] = nsCtWs.dat_wsInstance.onopenSubj.asObservable();
    nsCtWs.dat_wsObss['msg'] = nsCtWs.dat_wsInstance.onmessageSubj.asObservable().pipe(_nsWs.plainJsonParser);

    nsCtWs.dat_wsObss['msg:ev'] = new rxjs.ReplaySubject(3);
    nsCtWs.dat_wsObss['msg:rr:wi'] = new rxjs.Subject();
    nsCtWs.dat_wsObss['msg:rr:dlg'] = new rxjs.Subject();

    nsCtWs.dat_wsSubs['msg:router'] = nsCtWs.dat_wsObss['msg']
      .pipe(
        rxjs.tap((msg) => {
          switch(msg._bt) {
            case '_bev':
              nsCtWs.dat_wsObss['msg:ev'].next(msg._bp);
              break;

            case '_brr':
              switch(msg._bp._brc) {
                case nsCtWs.CONSTANT.WS_REQCODE_WALL_INFO:
                case nsCtWs.CONSTANT.WS_REQCODE_WALL_ITEM_INFO:
                  nsCtWs.dat_wsObss['msg:rr:wi'].next(msg._bp);
                  break;

                case nsCtWs.CONSTANT.WS_REQCODE_NODE_INFO_LISTING:
                case nsCtWs.CONSTANT.WS_REQCODE_STREAM_INFO_LISTING:
                  nsCtWs.dat_wsObss['msg:rr:dlg'].next(msg._bp);
                  break;

                default:
                  console.warn(`${nsCtWs.ns}`, `no route for req_code '${msg._bp._brc}'`);
                  break;
              }
              break;

            default:
              console.warn(`${nsCtWs.ns}`, `no handler for msg_type '${msg._bt}'`);
              break;
          }
        }),
      )
      .subscribe();
  };

  nsCtWs.fn_sendMessage = function(msgObj) {
    nsCtWs.dat_wsInstance.wsock.send(JSON.stringify(msgObj));
  };
};

// wall items
const nsCtWi = { ns: 'nsCtWi' };
nsCtWi.CONSTANT = {
  // wall layout code
  WLC__DEFAULT_4   : 'DEFAULT_4',
  WLC__DEFAULT_12  : 'DEFAULT_12',
  WLC__DEFAULT_16  : 'DEFAULT_16',
  WLC__DEFAULT_1B7S: 'DEFAULT_1B7S',
};
nsCtWi.init = function() {
  // TODO: use layout formula instead of making template for each layout
  nsCtWi.tmpl_viszD4 = `
<div class="row g-1 mb-1">
  <div class="col-6"><div id="viszItem1" class="wall-visz-d4-item"></div></div>
  <div class="col-6"><div id="viszItem2" class="wall-visz-d4-item"></div></div>
</div>
<div class="row g-1 mb-1">
  <div class="col-6"><div id="viszItem3" class="wall-visz-d4-item"></div></div>
  <div class="col-6"><div id="viszItem4" class="wall-visz-d4-item"></div></div>
</div>
`;
  nsCtWi.tmpl_viszD12 = `
  <div class="row g-1 mb-1">
  <div class="col-3"><div id="viszItem1" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem2" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem3" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem4" class="wall-visz-item-small"></div></div>
</div>
<div class="row g-1 mb-1">
  <div class="col-3"><div id="viszItem5" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem6" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem7" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem8" class="wall-visz-item-small"></div></div>
</div>
<div class="row g-1 mb-1">
  <div class="col-3"><div id="viszItem9" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem10" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem11" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem12" class="wall-visz-item-small"></div></div>
</div>
<div class="row g-1 mb-1">
  <div class="col-3"><div class="wall-visz-item-small invisible"></div></div>
</div>
`;
  nsCtWi.tmpl_viszD16 = `
<div class="row g-1 mb-1">
  <div class="col-3"><div id="viszItem1" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem2" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem3" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem4" class="wall-visz-item-small"></div></div>
</div>
<div class="row g-1 mb-1">
  <div class="col-3"><div id="viszItem5" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem6" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem7" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem8" class="wall-visz-item-small"></div></div>
</div>
<div class="row g-1 mb-1">
  <div class="col-3"><div id="viszItem9" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem10" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem11" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem12" class="wall-visz-item-small"></div></div>
</div>
<div class="row g-1 mb-1">
  <div class="col-3"><div id="viszItem13" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem14" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem15" class="wall-visz-item-small"></div></div>
  <div class="col-3"><div id="viszItem16" class="wall-visz-item-small"></div></div>
</div>
`;
  nsCtWi.tmpl_viszD1B7S = `
<div class="row g-1 mb-1">
  <div class="col-9"><div id="viszItem1" class="wall-visz-d1b7s-item-big"></div></div>
  <div class="col-3">
    <div class="row g-1 mb-1">
      <div class="col-12"><div id="viszItem2" class="wall-visz-d1b7s-item-small"></div></div>
    </div>
    <div class="row g-1 mb-1">
      <div class="col-12"><div id="viszItem3" class="wall-visz-d1b7s-item-small"></div></div>
    </div>
    <div class="row g-1">
      <div class="col-12"><div id="viszItem4" class="wall-visz-d1b7s-item-small"></div></div>
    </div>
  </div>
</div>
<div class="row g-1 mb-1">
  <div class="col-3"><div id="viszItem5" class="wall-visz-d1b7s-item-small"></div></div>
  <div class="col-3"><div id="viszItem6" class="wall-visz-d1b7s-item-small"></div></div>
  <div class="col-3"><div id="viszItem7" class="wall-visz-d1b7s-item-small"></div></div>
  <div class="col-3"><div id="viszItem8" class="wall-visz-d1b7s-item-small"></div></div>
</div>
`;

  nsCtWi.tmpl_wallItem = `
<div class="list-group-item p-3 wall-item">
  <div class="row">
    <!-- infos -->
    <div class="col">
      <div class="row">
        <div class="col-1">
          <div class="row">
            <div class="col-auto status-indicator status-off wi_status TODO" style="display:none;">
              <span class="status-indicator-circle"></span>
              <span class="status-indicator-circle"></span>
              <span class="status-indicator-circle"></span>
            </div>
            <div class="col card-title wi_index"></div>
          </div>
        </div>
        <div class="col-3">
          <div class="col card-title wi_nodeId"></div>
        </div>
        <div class="col-8">
          <div class="card-title wi_streamCode"></div>
        </div>
      </div>
      <div class="hr-text hr-text-center hr-text-spaceless mt-1 mb-2">Details</div>
      <div class="row mb-1">
        <div class="col-1"></div>
        <div class="col-3"><div class="text-secondary">Source type</div></div>
        <div class="col-8"><div class="text-secondary text-truncate wi_sourceType"></div></div>
      </div>
      <div class="row mb-1">
        <div class="col-1"></div>
        <div class="col-3"><div class="text-secondary">Estimated video bitrate</div></div>
        <div class="col-8"><div class="text-secondary text-truncate wi_videoBitrate"></div></div>
      </div>
    </div>
    <!-- aktions -->
    <div class="col-auto">
      <div class="row">
        <div class="col-auto">
          <div class="btn-list dropdown">
            <button class="btn btn-icon btn-danger wall-item-unlink wi_unlink" data-bs-toggle="dropdown">
            </button>
            <button class="btn btn-icon btn-primary wall-item-link wi_link"
                type="button" data-bs-toggle="modal" data-bs-target="#dlgLinkStream">
              <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-link"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M9 15l6 -6" /><path d="M11 6l.463 -.536a5 5 0 0 1 7.071 7.072l-.534 .464" /><path d="M13 18l-.397 .534a5.068 5.068 0 0 1 -7.127 0a4.972 4.972 0 0 1 0 -7.071l.524 -.463" /></svg>
            </button>

            <div class="dropdown-menu dropdown-menu-start dropdown-menu-arrow dropdown-confirmation">
              <div class="dropdown-confirmation-title">Unlink ?</div>
              <div class="dropdown-confirmation-menu">
                <button class="btn btn-outline-secondary dropdown-confirmation-button">No</button>
                <button class="btn btn-outline-danger dropdown-confirmation-button wi_unlinkYes">Yes</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
`;

  nsCtWi.tmpl_itemUnlink = `<svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-link-off"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M9 15l3 -3m2 -2l1 -1" /><path d="M11 6l.463 -.536a5 5 0 0 1 7.071 7.072l-.534 .464" /><path d="M3 3l18 18" /><path d="M13 18l-.397 .534a5.068 5.068 0 0 1 -7.127 0a4.972 4.972 0 0 1 0 -7.071l.524 -.463" /></svg>`;
  nsCtWi.tmpl_itemUnlinkPending = `<div class="spinner-border spinner-border-sm" role="status"></div>`;

  nsCtWi.el_wallVisz = document.getElementById('wallVisz');
  nsCtWi.el_wallItemListing = document.getElementById('wallItemListing');

  nsCtWi.dat_wallItemInstances = new Map();

  nsCtWi.dat_wsSubs = {};

  nsCtWi.fn_wsSetup = function() {
    nsCtWi.dat_wsSubs['ws:open'] = nsCtWs.dat_wsObss['open']
      .subscribe({
        next: (v) => {
          nsCtWi.ws_reqWallInfo();
        },
      });

    nsCtWi.dat_wsSubs['ws:msg:wi'] = nsCtWs.dat_wsObss['msg:rr:wi']
      .subscribe({
        next: (p_rep) => {
          switch(p_rep._brc) {
            case nsCtWs.CONSTANT.WS_REQCODE_WALL_INFO:
              const wallInfo = p_rep.wall_info;
              nsCtWi.fn_renderWallVisz(wallInfo.layout_code);
              nsCtWi.fn_populateWallItemInstances(wallInfo.items);
              nsCtWi.fn_renderWallItemListing();
              break;

            case nsCtWs.CONSTANT.WS_REQCODE_WALL_ITEM_INFO:
              const wallItemInfo = p_rep.wall_item_info;
              nsCtWi.fn_updateWallItemInstance(wallItemInfo);
              break;

            default:
              console.warn(`${nsCtWi.ns}`, `no handler for req_code '${p_rep._brc}'`);
              break;
          }
        },
      });
  };

  nsCtWi.fn_renderWallVisz = function(layoutCode) {
    switch(layoutCode) {
      case nsCtWi.CONSTANT.WLC__DEFAULT_4:
        nsCtWi.el_wallVisz.innerHTML = nsCtWi.tmpl_viszD4;
        break;

      case nsCtWi.CONSTANT.WLC__DEFAULT_12:
        nsCtWi.el_wallVisz.innerHTML = nsCtWi.tmpl_viszD12;
        break;

      case nsCtWi.CONSTANT.WLC__DEFAULT_16:
        nsCtWi.el_wallVisz.innerHTML = nsCtWi.tmpl_viszD16;
        break;

      case nsCtWi.CONSTANT.WLC__DEFAULT_1B7S:
        nsCtWi.el_wallVisz.innerHTML = nsCtWi.tmpl_viszD1B7S;
        break;

      default:
        console.warn(`${nsCtWi.ns}`, `cannot render unknown layout code '${layoutCode}'`);
        break;
    }
  };

  nsCtWi.fn_populateWallItemInstances = function(wallInfoItems) {
    nsCtWi.dat_wallItemInstances.clear();

    for (let i = 0; i < wallInfoItems.length; i++) {
      const wallItemInstance = nsCtWi.fn_newWallItemInstance(i, wallInfoItems[i]);

      nsCtWi.dat_wallItemInstances.set(wallItemInstance.id, wallItemInstance);
    }
  };

  nsCtWi.fn_renderWallItemListing = function() {
    nsCtWi.el_wallItemListing.innerHTML = '';

    nsCtWi.dat_wallItemInstances.forEach((inst) => {
      nsCtWi.el_wallItemListing.appendChild(inst.el_root);

      inst.fn_renderContent();
    });
  };

  nsCtWi.fn_newWallItemInstance = function(idx, wallItemInfo) {
    const inst = {};
    inst.id = wallItemInfo.id;
    inst.dat_item = wallItemInfo;

    inst.el_viszItem = document.getElementById(`viszItem${inst.dat_item.index}`);

    const tempRoot = document.createElement('div');
    tempRoot.innerHTML = nsCtWi.tmpl_wallItem;
    inst.el_root = tempRoot.children[0];

    inst.el_status = inst.el_root.getElementsByClassName('wi_status')[0];
    inst.el_index = inst.el_root.getElementsByClassName('wi_index')[0];
    inst.el_nodeId = inst.el_root.getElementsByClassName('wi_nodeId')[0];
    inst.el_streamCode = inst.el_root.getElementsByClassName('wi_streamCode')[0];
    inst.el_sourceType = inst.el_root.getElementsByClassName('wi_sourceType')[0];
    inst.el_vidBitrate = inst.el_root.getElementsByClassName('wi_videoBitrate')[0];
    inst.el_unlink = inst.el_root.getElementsByClassName('wi_unlink')[0];
    inst.el_unlinkYes = inst.el_root.getElementsByClassName('wi_unlinkYes')[0];
    inst.el_link = inst.el_root.getElementsByClassName('wi_link')[0];

    inst.dat_pendingAction = null; // null | 'unlink'

    inst.fn_renderContent = function() {
      const itemIdx = `${inst.dat_item.index}`;

      let sourceNodeText = inst.dat_item.source_node;
      if (sourceNodeText === '') sourceNodeText = '-';

      let sourceStreamText = inst.dat_item.source_stream;
      if (sourceStreamText === '') sourceStreamText = '-';

      let viszItemText = `#${itemIdx}`;
      if (inst.dat_item.source_node !== null) {
        viszItemText += `<br>${sourceNodeText}<br>${sourceStreamText}`;
      }


      let sourceTypeText = '-';
      let vidBitrateText = '-';
      if (inst.dat_item.stream_info !== null) {
        sourceTypeText = inst.dat_item.stream_info.source_type;
        vidBitrateText = `${inst.dat_item.stream_info.estimated_video_bitrate} bit/s`;
      }


      inst.el_viszItem.innerHTML = viszItemText;

      inst.el_index.innerHTML = `#${itemIdx}`;
      inst.el_nodeId.innerHTML = sourceNodeText;
      inst.el_streamCode.innerHTML = sourceStreamText;
      inst.el_sourceType.innerHTML = sourceTypeText;
      inst.el_vidBitrate.innerHTML = vidBitrateText;

      inst.fn_renderStatus();
      inst.fn_renderUnlinkButton();
    };

    inst.fn_updateStatus = function() {
      inst.fn_renderStatus();
    };
    inst.fn_renderStatus = function() {
      // TODO:
      // inst.el_status.innerHTML = 'Coming soon...';
    };

    inst.fn_renderUnlinkButton = function() {
      inst.el_unlink.innerHTML = nsCtWi.tmpl_itemUnlink;
      
      if (inst.dat_item.node_info === null) {
        inst.el_unlink.disabled = true;
      } else {
        inst.el_unlink.disabled = false;
      }
    };

    inst.fn_renderPendingAction = function() {
      if (inst.dat_pendingAction === null) {
        inst.el_unlink.disabled = false;
        inst.el_link.disabled = false;

        inst.fn_renderUnlinkButton();

      } else {
        inst.el_unlink.disabled = true;
        inst.el_link.disabled = true;

        switch (inst.dat_pendingAction) {
          case 'unlink':
            inst.el_unlink.innerHTML = nsCtWi.tmpl_itemUnlinkPending;
        }
      }
    };
    inst.fn_markPendingUnlink = function() {
      inst.dat_pendingAction = 'unlink';
      inst.fn_renderPendingAction();
    };
    inst.fn_clearPendingAction = function() {
      inst.dat_pendingAction = null;
      inst.fn_renderPendingAction();
    };

    // === instance - html event stuffs ===
    // TODO: what happens to all these event listeners, when the element is removed from DOM ?
    if (inst.el_root !== null) {
      inst.el_root.addEventListener('mouseover', function() {
        if (inst.el_viszItem !== null) inst.el_viszItem.classList.add('highlight');
      });
      inst.el_root.addEventListener('mouseout', function() {
        if (inst.el_viszItem !== null) inst.el_viszItem.classList.remove('highlight');
      });
    }

    if (inst.el_viszItem !== null) {
      inst.el_viszItem.addEventListener('mouseover', function() {
        if (inst.el_root !== null) inst.el_root.classList.add('highlight');
      });
      inst.el_viszItem.addEventListener('mouseout', function() {
        if (inst.el_root !== null) inst.el_root.classList.remove('highlight');
      });
    }

    inst.el_unlinkYes.addEventListener('click', async function() {
      // hmm... hide must go first
      // if called after mark, the hide is blocked by something (bootstrap related)
      // revisit later...
      inst.fn_hideUnlinkDropdown();
      inst.fn_markPendingUnlink();

      const respBundle = await nsCtWi.lapi_updateWallItem(inst.dat_item.id, '', '');
      if (respBundle.isOk()) {
        // TODO: visual indicator for new updated item
        nsCtWi.ws_reqWallItemInfo(inst.dat_item.id);

      } else {
        // TODO:
        console.error('TODO: show feedback somewhere ', respBundle);
      }

      inst.fn_clearPendingAction();
    });

    inst.el_link.addEventListener('click', function() {
      nsDlgLS.fn_openDialog(inst);
    });

    inst.fn_hideUnlinkDropdown = function() {
      const bsInst = bootstrap.Dropdown.getInstance(inst.el_unlink);
      if (bsInst === null) return;

      bsInst.hide();
    };

    return inst;
  };

  nsCtWi.fn_updateWallItemInstance = function(wallItemInfo) {
    nsCtWi.dat_wallItemInstances.forEach((inst) => {
      if (inst.id !== wallItemInfo.id) return;

      inst.dat_item = wallItemInfo
      inst.fn_renderContent();
    });
  };

  nsCtWi.ws_reqWallInfo = function() {
    nsCtWs.fn_sendMessage({
      _bid: '',
      _brc: nsCtWs.CONSTANT.WS_REQCODE_WALL_INFO,
    });
  };
  nsCtWi.ws_reqWallItemInfo = function(wallItemId) {
    nsCtWs.fn_sendMessage({
      _bid: '',
      _brc: nsCtWs.CONSTANT.WS_REQCODE_WALL_ITEM_INFO,

      wall_item_id: wallItemId,
    });
  };

  // TODO: replace functionality with ws ?
  nsCtWi.lapi_updateWallItem = async function(wallItemId, sourceNodeId, sourceStreamCode) {
    const reqUri = '/wall/local-api/update-wall-item';
    const reqBody = {
      wall_item_id: wallItemId,
      source_node_id: sourceNodeId,
      stream_code: sourceStreamCode,
    };

    let respBundle = null;
    await nsComm.lapi_defaultPOST(reqUri, JSON.stringify(reqBody))
            .then((_respBundle) => { respBundle = _respBundle })
            .catch((_err) => { console.error(`TODO: deadend ${_err}`); });
  
    return respBundle
  }
};

const nsDlgLS = { ns: 'nsDlgLS' };
nsDlgLS.init = function() {
  nsDlgLS.tmpl_nodeItem = `
<div class="list-group-item px-3 py-2 node-item faux-selection">
  <div class="row mb-1">
    <div class="col-auto status-indicator ni_status">
      <span class="status-indicator-circle"></span>
      <span class="status-indicator-circle"></span>
      <span class="status-indicator-circle"></span>
    </div>
    <div class="col">
      <div class="card-title mb-0 ni_title"></div>
    </div>
  </div>
  <div class="row">
    <div class="col-3"><div class="text-secondary">Last activity at</div></div>
    <div class="col-9"><div class="text-secondary text-truncate ni_last-activity"></div></div>
  </div>
  <div class="hr-text hr-text-center hr-text-spaceless mt-2 mb-2">Details</div>
  <div class="row mb-1">
    <div class="col-3"><div class="text-secondary">Stream Count</div></div>
    <div class="col-9"><div class="text-secondary text-truncate ni_stream-count"></div></div>
  </div>
</div>
`;

  nsDlgLS.tmpl_streamItem = `
<div class="list-group-item px-3 py-2 stream-item">
  <div class="row">
    <!-- col - stream thumbnail -->
    <div class="col-auto">
      <div class="rounded stream-item-thumbnail-container">
        <img class="stream-item-thumbnail" src="" alt="">
      </div>
    </div>

    <!-- col - stream details -->
    <div class="col">
      <div class="row">
        <div class="col-auto status-indicator si_status">
          <span class="status-indicator-circle"></span>
          <span class="status-indicator-circle"></span>
          <span class="status-indicator-circle"></span>
        </div>
        <div class="col">
          <div class="card-title mb-2 si_title"></div>
        </div>
      </div>
      <div class="row">
        <div class="col-3"><div class="text-secondary">Last activity at</div></div>
        <div class="col-9"><div class="text-secondary text-truncate si_last-activity"></div></div>
      </div>
      <div class="hr-text hr-text-center hr-text-spaceless mt-2 mb-2">Details</div>
      <div class="row mb-1">
        <div class="col-3"><div class="text-secondary">Source type</div></div>
        <div class="col-9"><div class="text-secondary text-truncate si_source-type"></div></div>
      </div>
      <div class="row mb-1">
        <div class="col-3"><div class="text-secondary">Estimated video bitrate</div></div>
        <div class="col-9"><div class="text-secondary text-truncate si_video-bitrate"></div></div>
      </div>
    </div>

    <!-- col - stream actions -->
    <div class="col-auto">
      <!-- stream item action - main act -->
      <div class="row mb-4">
        <div class="col px-0"></div>
        <div class="col-auto">
          <div class="btn-list">
            <button class="btn btn-icon btn-primary stream-item-link si_link">
              <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-link"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M9 15l6 -6" /><path d="M11 6l.463 -.536a5 5 0 0 1 7.071 7.072l-.534 .464" /><path d="M13 18l-.397 .534a5.068 5.068 0 0 1 -7.127 0a4.972 4.972 0 0 1 0 -7.071l.524 -.463" /></svg>
            </button>
          </div>
        </div>
      </div>

      <!-- stream item action - preview stream -->
      <div class="row">
        <div class="col">
          <button class="btn btn-icon stream-item-preview si_preview">
            <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-eye"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M10 12a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" /><path d="M21 12c-2.4 4 -5.4 6 -9 6c-3.6 0 -6.6 -2 -9 -6c2.4 -4 5.4 -6 9 -6c3.6 0 6.6 2 9 6" /></svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</div>
`;

  nsDlgLS.el_dialog = document.getElementById('dlgLinkStream');
  nsDlgLS.el_dialogTitle = document.getElementById('dlgLS_dialogTitle');
  nsDlgLS.el_nodeListingTitle = document.getElementById('dlgLS_nodeListingTitle');
  nsDlgLS.el_nodeListingSearch = document.getElementById('dlgLS_nodeListingSearch');
  nsDlgLS.el_nodeListingReload = document.getElementById('dlgLS_nodeListingReload');
  nsDlgLS.el_nodeListing = document.getElementById('dlgLS_nodeListing');
  nsDlgLS.el_streamListingTitle = document.getElementById('dlgLS_streamListingTitle');
  nsDlgLS.el_streamListingSearch = document.getElementById('dlgLS_streamListingSearch');
  nsDlgLS.el_streamListingReload = document.getElementById('dlgLS_streamListingReload');
  nsDlgLS.el_streamListing = document.getElementById('dlgLS_streamListing');
  nsDlgLS.el_streamPreviewTitle = document.getElementById('dlgLS_streamPreviewTitle');
  nsDlgLS.el_streamPreviewStop = document.getElementById('dlgLS_streamPreviewStop');
  nsDlgLS.el_streamPreview = document.getElementById('dlgLS_streamPreview');

  nsDlgLS.bs_dialog = bootstrap.Modal.getOrCreateInstance(`#${nsDlgLS.el_dialog.id}`);

  nsDlgLS.dat_currWallItemInstance = null;

  nsDlgLS.dat_nodeInstances = new Map();
  nsDlgLS.dat_selectedNodeInstance = null;
  nsDlgLS.dat_nodeSearchSubj = new rxjs.Subject();

  nsDlgLS.dat_streamInstances = new Map();
  nsDlgLS.dat_streamSearchSubj = new rxjs.Subject();

  nsDlgLS.dat_wsSubs = {};

  nsDlgLS.fn_wsSetup = function() {
    nsDlgLS.dat_wsSubs['ws:msg:dlg'] = nsCtWs.dat_wsObss['msg:rr:dlg']
      .subscribe({
        next: (p_rep) => {
          switch(p_rep._brc) {
            case nsCtWs.CONSTANT.WS_REQCODE_NODE_INFO_LISTING:
              nsDlgLS.fn_populateNodeInstances(p_rep.node_info_listing);
              nsDlgLS.fn_renderNodeListing();
              nsDlgLS.fn_nodeSearch(nsDlgLS.el_nodeListingSearch.value);
              break;

            case nsCtWs.CONSTANT.WS_REQCODE_STREAM_INFO_LISTING:
              nsDlgLS.fn_populateStreamInstances(p_rep.stream_info_listing);
              nsDlgLS.fn_renderStreamListing();
              nsDlgLS.fn_streamSearch(nsDlgLS.el_streamListingSearch.value);
              break;

            default:
              console.warn(`${nsDlgLS.ns}`, `no handler for req_code '${p_rep._brc}'`);
              break;
          }
        },
      });
  };

  // === node stuffs ===
  nsDlgLS.fn_populateNodeInstances = function(nodeInfoListing) {
    nsDlgLS.dat_nodeInstances.clear();

    for (let i = 0; i < nodeInfoListing.length; i++) {
      const nodeInstance = nsDlgLS.fn_newNodeInstance(i, nodeInfoListing[i]);

      nsDlgLS.dat_nodeInstances.set(nodeInstance.id, nodeInstance);
    }
  };

  nsDlgLS.fn_renderNodeListing = function() {
    nsDlgLS.el_nodeListingTitle.innerHTML = `Nodes (${nsDlgLS.dat_nodeInstances.size} / ${nsDlgLS.dat_nodeInstances.size})`;
    nsDlgLS.el_nodeListing.innerHTML = '';

    nsDlgLS.dat_nodeInstances.forEach((inst) => {
      nsDlgLS.el_nodeListing.appendChild(inst.el_root);

      inst.fn_renderContent();
    });
  };

  nsDlgLS.fn_newNodeInstance = function(idx, nodeInfo) {
    const inst = {};
    inst.id = nodeInfo.id;
    inst.dat_item = nodeInfo;

    const tempRoot = document.createElement('div');
    tempRoot.innerHTML = nsDlgLS.tmpl_nodeItem;
    inst.el_root = tempRoot.children[0];

    inst.el_status = inst.el_root.getElementsByClassName('ni_status')[0];
    inst.el_title = inst.el_root.getElementsByClassName('ni_title')[0];
    inst.el_lastActivityAt = inst.el_root.getElementsByClassName('ni_last-activity')[0];
    inst.el_streamCount = inst.el_root.getElementsByClassName('ni_stream-count')[0];

    inst.fn_renderContent = function() {
      switch (inst.dat_item.state) {
        case 'n_s:ready':
          inst.el_status.classList.add('status-green');
          break;

        default:
          inst.el_status.classList.add('status-off');
          break;
      }

      const fmt1 = nsCt.fn_formatTimestamp(inst.dat_item.last_activity_at);

      inst.el_title.innerHTML = inst.id;
      inst.el_lastActivityAt.innerHTML = `${fmt1}`;
      inst.el_streamCount.innerHTML = inst.dat_item.stream_count;
    };

    inst.fn_show = function() { inst.el_root.style.display = ''; };
    inst.fn_hide = function() { inst.el_root.style.display = 'none'; };

    inst.el_root.addEventListener('click', function() {
      nsDlgLS.fn_updateNodeSelection(inst);
    });

    return inst;
  };

  nsDlgLS.fn_updateNodeSelection = function(nodeInst) {
    nsDlgLS.fn_resetNodeSelection();

    nsDlgLS.dat_selectedNodeInstance = nodeInst;

    nodeInst.el_root.classList.add('active');

    nsDlgLS.ws_reqStreamInfoListing(nsDlgLS.dat_selectedNodeInstance.id);
  };

  nsDlgLS.fn_resetNodeSelection = function() {
    nsDlgLS.dat_selectedNodeInstance = null;

    nsDlgLS.dat_nodeInstances.forEach((inst) => {
      inst.el_root.classList.remove('active');
    });
  };

  nsDlgLS.fn_resetNodeListing = function() {
    nsDlgLS.el_nodeListingSearch.value = '';

    nsDlgLS.fn_resetNodeSelection();
    nsDlgLS.dat_nodeInstances.clear();

    nsDlgLS.fn_renderNodeListing();
  };

  // === stream stuffs ===
  nsDlgLS.fn_populateStreamInstances = function(streamInfoListing) {
    nsDlgLS.dat_streamInstances.clear();

    for (let i = 0; i < streamInfoListing.length; i++) {
      const streamInstance = nsDlgLS.fn_newStreamInstance(i, streamInfoListing[i]);

      nsDlgLS.dat_streamInstances.set(streamInstance.id, streamInstance);
    }
  };

  nsDlgLS.fn_renderStreamListing = function() {
    nsDlgLS.fn_renderStreamListingTitle();
    nsDlgLS.el_streamListing.innerHTML = '';

    nsDlgLS.dat_streamInstances.forEach((inst) => {
      nsDlgLS.el_streamListing.appendChild(inst.el_root);

      inst.fn_renderContent();
    });
  };
  nsDlgLS.fn_renderStreamListingTitle = function(shownCount) {
    let titleText = 'Streams';
    if (nsDlgLS.dat_selectedNodeInstance !== null) {
      titleText += ` on node '${nsDlgLS.dat_selectedNodeInstance.dat_item.id}'`;

      if (shownCount === undefined || shownCount === null) {
        titleText += ` (${nsDlgLS.dat_streamInstances.size} / ${nsDlgLS.dat_streamInstances.size})`;
      } else {
        titleText += ` (${shownCount} / ${nsDlgLS.dat_streamInstances.size})`;
      }
    }

    nsDlgLS.el_streamListingTitle.innerHTML = titleText;
  };

  nsDlgLS.fn_newStreamInstance = function(idx, streamInfo) {
    const inst = {};
    inst.id = streamInfo.id;
    inst.dat_item = streamInfo;

    const tempRoot = document.createElement('div');
    tempRoot.innerHTML = nsDlgLS.tmpl_streamItem;
    inst.el_root = tempRoot.children[0];

    inst.el_status = inst.el_root.getElementsByClassName('si_status')[0];
    inst.el_title = inst.el_root.getElementsByClassName('si_title')[0];
    inst.el_lastActivityAt = inst.el_root.getElementsByClassName('si_last-activity')[0];
    inst.el_sourceType = inst.el_root.getElementsByClassName('si_source-type')[0];
    inst.el_videoBitrate = inst.el_root.getElementsByClassName('si_video-bitrate')[0];
    inst.el_link = inst.el_root.getElementsByClassName('si_link')[0];
    inst.el_preview = inst.el_root.getElementsByClassName('si_preview')[0];

    inst.fn_renderContent = function() {
      switch (inst.dat_item.streamer_state) {
        case 'lss:start':
          // temp workaround for state change to 'running' does not send any event
          // hence cache is not updated (state is as last state, most likely 'start')
          inst.el_status.classList.add('status-green', 'status-indicator-animated');
          break;

        case 'lss:running':
          inst.el_status.classList.add('status-green');
          break;

        case 'lss:stop:normal':
          inst.el_status.classList.add('status-off');
          break;

        case 'lss:stop:unexpected':
          inst.el_status.classList.add('status-red');
          break;

        default:
          inst.el_status.classList.add('status-off');
          break;
      }

      const fmt1 = nsCt.fn_formatTimestamp(inst.dat_item.last_activity_at);

      inst.el_title.innerHTML = `${inst.dat_item.code}`;
      inst.el_lastActivityAt.innerHTML = `${fmt1}`
      inst.el_sourceType.innerHTML = inst.dat_item.source_type;
      inst.el_videoBitrate.innerHTML = `${inst.dat_item.estimated_video_bitrate} bit/s`; // TODO: formatting
    };

    inst.fn_show = function() { inst.el_root.style.display = ''; };
    inst.fn_hide = function() { inst.el_root.style.display = 'none'; };

    inst.el_link.addEventListener('click', async function() {
      if (nsDlgLS.dat_selectedNodeInstance === null) return;
      
      const respBundle = await nsDlgLS.lapi_updateWallItem(
        nsDlgLS.dat_currWallItemInstance.dat_item.id,
        nsDlgLS.dat_selectedNodeInstance.dat_item.id,
        inst.dat_item.code,
      );

      if (respBundle.isOk()) {
        // TODO: visual indicator for new updated item
        nsCtWi.ws_reqWallItemInfo(nsDlgLS.dat_currWallItemInstance.dat_item.id);
        nsDlgLS.fn_closeDialog();

      } else {
        // TODO:
        console.error('TODO: show feedback somewhere ', respBundle);
      }
    });

    inst.el_preview.addEventListener('click', function() {
      const previewUrl = inst.dat_item.preview_url;

      // TODO TODO
      if (previewUrl === null) return;

      let nodeText = '???';
      if (nsDlgLS.dat_selectedNodeInstance !== null) {
        nodeText = nsDlgLS.dat_selectedNodeInstance.id;
      }

      nsDlgLS.fn_previewStream(nodeText, inst.dat_item.code, previewUrl);
    });

    return inst;
  };

  nsDlgLS.fn_resetStreamListing = function() {
    nsDlgLS.el_streamListingSearch.value = '';

    nsDlgLS.dat_streamInstances.clear();

    nsDlgLS.fn_renderStreamListing();
  };

  // === node search stuffs ===
  nsDlgLS.fn_nodeSearch = function(searchStr) {
    nsDlgLS.fn_resetNodeSelection();

    nsDlgLS.fn_resetStreamListing();
    nsDlgLS.fn_resetStreamPreview();

    let hasSearchStr = true;
    let shownCount = 0;

    if (searchStr === undefined || searchStr === null || searchStr === '') {
      hasSearchStr = false;
    } else {
      searchStr = searchStr.toLowerCase();
    }

    nsDlgLS.dat_nodeInstances.forEach((inst) => {
      let doShow = true;

      if (hasSearchStr) {
        doShow = (
          inst.dat_item.id.toLowerCase().includes(searchStr)
        );
      }

      if (doShow) {
        inst.fn_show();
        shownCount += 1;

      } else {
        inst.fn_hide();
      }
    });

    nsDlgLS.el_nodeListingTitle.innerHTML = `Nodes (${shownCount} / ${nsDlgLS.dat_nodeInstances.size})`;
  };

  nsDlgLS.dat_nodeSearchSubj
    .pipe(
      rxjs.debounceTime(333),
    )
    .subscribe({
      next: searchStr => {
        nsDlgLS.fn_nodeSearch(searchStr);
      },
    });

  // === stream search stuffs ===
  nsDlgLS.fn_streamSearch = function(searchStr) {
    let hasSearchStr = true;
    let shownCount = 0;

    if (searchStr === undefined || searchStr === null || searchStr === '') {
      hasSearchStr = false;
    } else {
      searchStr = searchStr.toLowerCase();
    }

    nsDlgLS.dat_streamInstances.forEach((inst) => {
      let doShow = true;

      if (hasSearchStr) {
        doShow = (
          inst.dat_item.code.toLowerCase().includes(searchStr)
        );
      }

      if (doShow) {
        inst.fn_show();
        shownCount += 1;

      } else {
        inst.fn_hide();
      }
    });

    nsDlgLS.fn_renderStreamListingTitle(shownCount);
  };

  nsDlgLS.dat_streamSearchSubj
    .pipe(
      rxjs.debounceTime(333),
    )
    .subscribe({
      next: searchStr => {
        nsDlgLS.fn_streamSearch(searchStr);
      },
    });

  // === preview stuffs ===
  nsDlgLS.fn_resetStreamPreview = function() {
    const componentId = 'dlgLS_streamPreview';
    nsDlgLS.el_streamPreviewTitle.innerHTML = 'Stream Preview';

    if (_nsDLS.componentExists(componentId)) {
      _nsDLS.destroyComponent(componentId);
    }
  };
  nsDlgLS.fn_previewStream = function(nodeId, streamCode, previewUrl) {
    nsDlgLS.el_streamPreviewTitle.innerHTML = `Stream Preview<br>${nodeId} - ${streamCode}`;

    const componentId = 'dlgLS_streamPreview';
    if (_nsDLS.componentExists(componentId)) {
      _nsDLS.changeHlsSource(componentId, previewUrl);
    } else {
      _nsDLS.newComponent(componentId, nsDlgLS.el_streamPreview, previewUrl);
    }
  };

  nsDlgLS.el_streamPreviewStop.addEventListener('click', function() {
    nsDlgLS.fn_resetStreamPreview();
  });


  // === dialog stuffs ===
  nsDlgLS.fn_openDialog = function(wallItemInst) {
    // hmm... somehow somewhat this line below is always called after show/shown.bs.modal fired...
    // something to do with all bootstrap API is async
    // so... show/shown event is unused for now
    // revisit later...
    nsDlgLS.dat_currWallItemInstance = wallItemInst;

    if (nsDlgLS.dat_currWallItemInstance !== null) {
      nsDlgLS.el_dialogTitle.innerHTML = `Select Stream for slot #${nsDlgLS.dat_currWallItemInstance.dat_item.index}`;
    }

    if (nsDlgLS.dat_nodeInstances.size === 0) {
      nsDlgLS.ws_reqNodeInfoListing();
    }

    nsDlgLS.fn_resetStreamPreview();

    nsDlgLS.bs_dialog.show();
  };

  nsDlgLS.fn_closeDialog = function() {
    nsDlgLS.bs_dialog.hide();
  };

  nsDlgLS.el_dialog.addEventListener('hidden.bs.modal', function() {
    nsDlgLS.dat_currWallItemInstance = null;

    nsDlgLS.fn_resetStreamPreview();
  });

  nsDlgLS.ws_reqNodeInfoListing = function() {
    nsCtWs.fn_sendMessage({
      _bid: '',
      _brc: nsCtWs.CONSTANT.WS_REQCODE_NODE_INFO_LISTING,
    });
  };

  nsDlgLS.ws_reqStreamInfoListing = function(nodeId) {
    nsCtWs.fn_sendMessage({
      _bid: '',
      _brc: nsCtWs.CONSTANT.WS_REQCODE_STREAM_INFO_LISTING,

      node_id: nodeId,
    });
  };

  nsDlgLS.lapi_updateWallItem = async function(wallItemID, sourceNodeID, sourceStreamCode) {
    const reqURI = '/wall/local-api/update-wall-item';
    const reqBody = {
      wall_item_id: wallItemID,
      source_node_id: sourceNodeID,
      stream_code: sourceStreamCode,
    };

    let respBundle = null;
    await nsComm.lapi_defaultPOST(reqURI, JSON.stringify(reqBody))
            .then((_respBundle) => { respBundle = _respBundle })
            .catch((_err) => { console.error(`TODO: deadend ${_err}`); });
  
    return respBundle
  }

  // === event listener stuffs ===
  nsDlgLS.el_nodeListingSearch.addEventListener('input', function() {
    nsDlgLS.dat_nodeSearchSubj.next(this.value);
  });

  // TODO: loading indicator
  nsDlgLS.el_nodeListingReload.addEventListener('click', function() {
    nsDlgLS.fn_resetNodeListing();
    nsDlgLS.ws_reqNodeInfoListing();

    nsDlgLS.fn_resetStreamListing();
  });

  nsDlgLS.el_streamListingSearch.addEventListener('input', function() {
    nsDlgLS.dat_streamSearchSubj.next(this.value);
  });

  // TODO: loading indicator
  nsDlgLS.el_streamListingReload.addEventListener('click', function() {
    nsDlgLS.fn_resetStreamListing();

    if (nsDlgLS.dat_selectedNodeInstance === null) return;

    nsDlgLS.ws_reqStreamInfoListing(nsDlgLS.dat_selectedNodeInstance.id);
  });
};

document.addEventListener('DOMContentLoaded', function() {
  _nsWs.init();

  nsCt.init();
  if (nsCD && nsCD.isEditMode) {
    nsCtWs.init();
    nsCtWi.init();
    nsDlgLS.init();

    if (nsCD.ws_uri) {
      nsCtWs.fn_wsSetup();
      nsCtWi.fn_wsSetup();
      nsDlgLS.fn_wsSetup();
    }
  }

  nsMain.focusOn('code');
});
