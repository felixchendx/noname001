'use strict';

// nsContent
const nsCt = {};
nsCt.init = function() {
  nsCt.el_btnBackToListing = document.getElementById('btnBackToListing');
  nsCt.el_btnSubmitDelete = document.getElementById('btnSubmitDelete');

  nsCt.el_streamItems = document.getElementById('streamItems');

  nsCt.fn_doSubmitDelete = function() {
    nsCt.el_btnSubmitDelete.click();
  };

  nsCt.el_btnBackToListing.addEventListener('click', function() {
    window.location = window.location.protocol + '//' + window.location.host + '/stream/stream-group/listing';
  });
};

// websocket
const nsCtWs = {ns: 'nsCtWs'};
nsCtWs.init = function() {
  nsCtWs.CONSTANT = {
    WS_REQCODE_STREAM_ITEM_LISTING: '/stream-item/listing',

    WS_REQCODE_DEVICE_SNAPSHOT_LISTING: '/device/snapshot/listing',
  };

  nsCtWs.dat_wsInstance = null;
  nsCtWs.dat_wsObss = {};
  nsCtWs.dat_wsSubs = {};

  nsCtWs.fn_setupWs = function() {
    const wsId = 'stream-group-detail';
    const wsUri = '/stream/stream-group/detail-02/ws' + window.location.search;
    nsCtWs.dat_wsInstance = _nsWs.newWsInstance(wsId, wsUri);

    nsCtWs.dat_wsObss['open'] = nsCtWs.dat_wsInstance.onopenSubj.asObservable();
    nsCtWs.dat_wsObss['msg'] = nsCtWs.dat_wsInstance.onmessageSubj.asObservable().pipe(_nsWs.plainJsonParser);

    nsCtWs.dat_wsObss['msg:ev'] = new rxjs.ReplaySubject(3);
    nsCtWs.dat_wsObss['msg:rr:sil'] = new rxjs.Subject();
    nsCtWs.dat_wsObss['msg:rr:asi'] = new rxjs.Subject();

    nsCtWs.dat_wsSubs['msg:router'] = nsCtWs.dat_wsObss['msg']
      .pipe(
        rxjs.tap((msg) => {
          switch(msg._bt) {
            case '_bev': nsCtWs.dat_wsObss['msg:ev'].next(msg._bp); break;
            case '_brr':
              switch(msg._bp._brc) {
                case nsCtWs.CONSTANT.WS_REQCODE_STREAM_ITEM_LISTING:
                  nsCtWs.dat_wsObss['msg:rr:sil'].next(msg._bp);
                  break;

                case nsCtWs.CONSTANT.WS_REQCODE_DEVICE_SNAPSHOT_LISTING:
                  nsCtWs.dat_wsObss['msg:rr:asi'].next(msg._bp);
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

// stream item listing
const nsCtSil = {ns: 'nsCtSil'};
nsCtSil.init = function() {
  nsCtSil.tmpl_streamItem = `<div class="list-group-item stream-item p-3">
  <div class="row">
    <!-- col - stream item thumbnail -->
    <div class="col-auto">
      <div class="rounded item-thumbnail-container">
        <img class="item-thumbnail" src="" alt="">
      </div>
    </div>

    <!-- col - stream item details -->
    <div class="col">
      <div class="row">
        <div class="col">
          <div class="card-title mb-2 item-code"></div>
          <div class="card-subtitle mb-2 item-name"></div>
        </div>
      </div>
      <div class="hr-text hr-text-center hr-text-spaceless mt-1 mb-2">Details</div>
      <div class="row">
        <div class="col-3">
          <div class="text-secondary">Note</div>
        </div>
        <div class="col-9">
          <div class="text-secondary item-note"></div>
        </div>
      </div>
      <div class="row">
        <div class="col-3">
          <div class="text-secondary">Status</div>
        </div>
        <div class="col-9">
          <div class="text-secondary text-truncate item-status"></div>
        </div>
      </div>
      <div class="row">
        <div class="col-3">
          <div class="text-secondary">Source Type</div>
        </div>
        <div class="col-9">
          <div class="text-secondary item-source-type"></div>
        </div>
      </div>
      <div class="row">
        <div class="col-3">
          <div class="text-secondary">Source</div>
        </div>
        <div class="col-9">
          <div class="text-secondary item-source"></div>
        </div>
      </div>
    </div>

    <!-- col - stream item actions -->
    <div class="col-auto">
      <!-- stream item action - main act -->
      <div class="row mb-6">
        <div class="col px-0"></div>
        <div class="col-auto">
          <div class="btn-list dropdown">
            <button class="btn btn-icon btn-danger item-delete" data-bs-toggle="dropdown">
              <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-trash"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M4 7l16 0" /><path d="M10 11l0 6" /><path d="M14 11l0 6" /><path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12" /><path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3" /></svg>
            </button>
            <button class="btn btn-icon btn-primary item-edit">
              <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-pencil"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M4 20h4l10.5 -10.5a2.828 2.828 0 1 0 -4 -4l-10.5 10.5v4" /><path d="M13.5 6.5l4 4" /></svg>
            </button>

            <div class="dropdown-menu dropdown-menu-start dropdown-menu-arrow dropdown-confirmation">
              <div class="dropdown-confirmation-title">Delete ?</div>
              <div class="dropdown-confirmation-menu">
                <button class="btn btn-outline-secondary dropdown-confirmation-button item-delete-no">
                  No
                </button>
                <button class="btn btn-outline-danger dropdown-confirmation-button item-delete-yes">
                  Yes
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- stream item action - state toggle -->
      <div class="row mb-2">
        <div class="col">
          <button class="btn item-toggle-state">
          </button>
        </div>
      </div>

      <!-- stream item action - view stream -->
      <div class="row">
        <div class="col">
          <button class="btn item-preview w-100">
            <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-eye"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M10 12a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" /><path d="M21 12c-2.4 4 -5.4 6 -9 6c-3.6 0 -6.6 -2 -9 -6c2.4 -4 5.4 -6 9 -6c3.6 0 6.6 2 9 6" /></svg>
            View
          </button>
        </div>
      </div>
    </div>
  </div>
<div>`;

  nsCtSil.tmpl_streamInactive = `
<svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-cast-off"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M3 19h.01" /><path d="M7 19a4 4 0 0 0 -4 -4" /><path d="M11 19a8 8 0 0 0 -8 -8" /><path d="M15 19h3a3 3 0 0 0 .875 -.13m2 -2a3 3 0 0 0 .128 -.868v-8a3 3 0 0 0 -3 -3h-9m-3.865 .136a3 3 0 0 0 -1.935 1.864" /><path d="M3 3l18 18" /></svg>
OFF
`;
  nsCtSil.tmpl_streamActive = `
<svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-cast"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M3 19l.01 0" /><path d="M7 19a4 4 0 0 0 -4 -4" /><path d="M11 19a8 8 0 0 0 -8 -8" /><path d="M15 19h3a3 3 0 0 0 3 -3v-8a3 3 0 0 0 -3 -3h-12a3 3 0 0 0 -2.8 2" /></svg>
ON
`;
  nsCtSil.tmpl_streamPending = `<div class="spinner-border spinner-border-sm" role="status"></div>`;

  nsCtSil.tmpl_itemDelete = `<svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-trash"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M4 7l16 0" /><path d="M10 11l0 6" /><path d="M14 11l0 6" /><path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12" /><path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3" /></svg>`;
  nsCtSil.tmpl_itemDeletePending = `<div class="spinner-border spinner-border-sm" role="status"></div>`;

  nsCtSil.el_silNavTabbing = document.getElementById('silNavTabbing');

  nsCtSil.el_silMessage = document.getElementById('silMessage');
  nsCtSil.el_silTitle = document.getElementById('silTitle');
  nsCtSil.el_silSearch = document.getElementById('silSearch');
  nsCtSil.el_silReload = document.getElementById('silReload');
  nsCtSil.el_streamItemListing = document.getElementById('streamItemListing');
  nsCtSil.el_silPreviewTitle = document.getElementById('silPreviewTitle');
  nsCtSil.el_silPreviewStop = document.getElementById('silPreviewStop');
  nsCtSil.el_silPreviewDls = document.getElementById('silPreviewDls');

  nsCtSil.dat_streamItemInstances = new Map();
  nsCtSil.dat_silSearchSubj = new rxjs.Subject();

  nsCtSil.dat_wsSubs = {};

  nsCtSil.fn_setupWs = function() {
    nsCtSil.dat_wsSubs['ws:open'] = nsCtWs.dat_wsObss['open']
      .subscribe({
        next: (v) => {
          nsCtSil.wsReq_streamItemListing();
        },
      });

    nsCtSil.dat_wsSubs['ws:msg:rr'] = nsCtWs.dat_wsObss['msg:rr:sil']
      .subscribe({
        next: (p_rep) => {
          switch(p_rep._brc) {
            case nsCtWs.CONSTANT.WS_REQCODE_STREAM_ITEM_LISTING:
              nsCtSil.fn_populateStreamItemInstances(p_rep.items);
              nsCtSil.fn_renderStreamItemListing();
              nsCtSil.fn_searchStreamItemListing(nsCtSil.el_silSearch.value);
              break;

            default:
              console.warn(`${nsCtSil.ns}`, `no handler for req_code '${p_rep._brc}'`);
              break;
          }
        },
      });
  };

  nsCtSil.fn_populateStreamItemInstances = function(streamItems) {
    nsCtSil.dat_streamItemInstances.clear();

    for (let i = 0; i < streamItems.length; i++) {
      const streamItemInstance = nsCtSil.fn_newStreamItemInstance(i, streamItems[i]);

      nsCtSil.dat_streamItemInstances.set(streamItemInstance.id, streamItemInstance);
    }
  };

  nsCtSil.fn_renderStreamItemListing = function() {
    nsCtSil.el_silTitle.innerHTML = `Stream items (${nsCtSil.dat_streamItemInstances.size} items)`;
    nsCtSil.el_streamItemListing.innerHTML = '';

    nsCtSil.dat_streamItemInstances.forEach((inst) => {
      nsCtSil.el_streamItemListing.appendChild(inst.el_root);

      inst.fn_renderContent();
      inst.fn_renderPendingAction();
    });

    if (nsCtSil.dat_streamItemInstances.size == 0) {
      nsCtSil.fn_displayInfoMessage(
        'This stream group does not have any item yet.',
        "Add stream item from tabbing 'Add / Edit Stream Item'",
      );
    }
  };

  nsCtSil.fn_searchStreamItemListing = function(searchStr) {
    let hasSearchStr = true;
    let shownCount = 0;

    if (searchStr === undefined || searchStr === null || searchStr === '') {
      hasSearchStr = false;
    } else {
      searchStr = searchStr.toLowerCase();
    }

    nsCtSil.dat_streamItemInstances.forEach((inst) => {
      let doShow = true;

      if (hasSearchStr) {
        doShow = (
          inst.dat_item.code.toLowerCase().includes(searchStr)
          || inst.dat_item.name.toLowerCase().includes(searchStr)
        );
      }

      if (doShow) {
        inst.fn_show();
        shownCount += 1;
      } else {
        inst.fn_hide();
      }
    });

    nsCtSil.el_silTitle.innerHTML = `Stream items (${shownCount} / ${nsCtSil.dat_streamItemInstances.size} items)`;
  };
  nsCtSil.fn_resetSearch = function() {
    nsCtSil.el_silSearch.value = '';
  };

  nsCtSil.fn_newStreamItemInstance = function(idx, streamItem) {
    const inst = {};
    inst.id = streamItem.id;
    inst.dat_item = streamItem;

    const tempRoot = document.createElement('div');
    tempRoot.innerHTML = nsCtSil.tmpl_streamItem;
    inst.el_root = tempRoot.children[0];

    inst.el_thumbnail = inst.el_root.getElementsByClassName('item-thumbnail')[0];
    inst.el_code = inst.el_root.getElementsByClassName('item-code')[0];
    inst.el_name = inst.el_root.getElementsByClassName('item-name')[0];
    inst.el_note = inst.el_root.getElementsByClassName('item-note')[0];
    inst.el_status = inst.el_root.getElementsByClassName('item-status')[0];
    inst.el_sourceType = inst.el_root.getElementsByClassName('item-source-type')[0];
    inst.el_source = inst.el_root.getElementsByClassName('item-source')[0];

    inst.el_delete = inst.el_root.getElementsByClassName('item-delete')[0];
    inst.el_deleteNo = inst.el_root.getElementsByClassName('item-delete-no')[0];
    inst.el_deleteYes = inst.el_root.getElementsByClassName('item-delete-yes')[0];
    inst.el_edit = inst.el_root.getElementsByClassName('item-edit')[0];
    inst.el_toggleState = inst.el_root.getElementsByClassName('item-toggle-state')[0];
    inst.el_preview = inst.el_root.getElementsByClassName('item-preview')[0];

    // flag for instance wide block
    inst.dat_pendingAction = ''; // 'delete' | 'toggle-state'

    inst.fn_renderContent = function() {
      let sourceText = '';
      switch (inst.dat_item.source_type) {
        case 'mod_device':
          sourceText = `${inst.dat_item.device_code}, Ch ${inst.dat_item.device_channel_id} - ${inst.dat_item.device_stream_type}`;
          break;

        case 'external':
          sourceText = inst.dat_item.external_url;
          break;

        case 'file':
          sourceText = inst.dat_item.filepath;
          break;

        default:
          sourceText = 'unknown';
          break;
      }

      inst.el_thumbnail.alt = `thumbnail - ${inst.dat_item.code}`;
      inst.el_code.innerHTML = inst.dat_item.code;
      inst.el_name.innerHTML = inst.dat_item.name;
      inst.el_note.innerHTML = inst.dat_item.note;
      
      inst.el_sourceType.innerHTML = inst.dat_item.source_type;
      inst.el_source.innerHTML = sourceText;

      inst.el_status.innerHTML = 'coming soon...'; // TODO
      `<span class="badge bg-success" style="margin-bottom: 2px;"></span>
        <span>Ready</span>
        <div class="mt-1">
          <code style="display: block;">ASDASD</code>
        </div>`;

      inst.fn_renderStreamState();
      inst.fn_renderPreviewButton();
    };
    inst.fn_renderStreamState = function() {
      inst.el_toggleState.classList.remove('item-stream-active', 'item-stream-inactive');

      if (inst.dat_item.state === 'active') {
        inst.el_toggleState.innerHTML = nsCtSil.tmpl_streamActive;
        inst.el_toggleState.classList.add('item-stream-active');

      } else {
        inst.el_toggleState.innerHTML = nsCtSil.tmpl_streamInactive;
        inst.el_toggleState.classList.add('item-stream-inactive');
      }
    };
    inst.fn_renderPreviewButton = function() {
      if (inst.dat_item.state === 'active') {
        inst.el_preview.disabled = false;
      } else {
        inst.el_preview.disabled = true;
      }
    };

    inst.fn_renderPendingAction = function() {
      if (inst.dat_pendingAction === '') {
        inst.el_delete.disabled = false;
        inst.el_deleteNo.disabled = false;
        inst.el_deleteYes.disabled = false;
        inst.el_edit.disabled = false;
        inst.el_toggleState.disabled = false;

        inst.fn_renderPreviewButton();

      } else {
        inst.el_delete.disabled = true;
        inst.el_deleteNo.disabled = true;
        inst.el_deleteYes.disabled = true;
        inst.el_edit.disabled = true;
        inst.el_toggleState.disabled = true;

        inst.el_preview.disabled = true;

        switch (inst.dat_pendingAction) {
          case 'delete':
            inst.el_delete.innerHTML = nsCtSil.tmpl_itemDeletePending;
            break;

          case 'toggle-state':
            inst.el_toggleState.innerHTML = nsCtSil.tmpl_streamPending;
            break;
        }
      }
    };

    inst.fn_markPendingDelete = function() {
      inst.dat_pendingAction = 'delete';
      inst.fn_renderPendingAction();
    };
    inst.fn_markPendingToggleState = function() {
      inst.dat_pendingAction = 'toggle-state';
      inst.fn_renderPendingAction();
    };
    inst.fn_clearPendingAction = function() {
      inst.dat_pendingAction = '';
      inst.fn_renderPendingAction();
    };

    inst.fn_show = function() { inst.el_root.style.display = ''; };
    inst.fn_hide = function() { inst.el_root.style.display = 'none'; };

    // TODO: check bootstrap for way to close the dropdown instead of disabling the yes no button
    inst.el_deleteYes.addEventListener('click', async function() {
      inst.fn_markPendingDelete();

      const respBundle = await nsCtSil.lapi_removeStreamItem(inst.dat_item.id);

      if (respBundle.isOk()) {
        // TODO: refresh only this instance + visual flair
        nsCtSil.wsReq_streamItemListing();

        nsCtSil.fn_displayOkMessage('Delete ok', `${respBundle.defResp.message}`);
      } else {
        nsCtSil.fn_displayErrorMessage('Delete error', `${respBundle.defResp.message}`);
      }

      inst.fn_clearPendingAction();
    });
    inst.el_edit.addEventListener('click', function() {
      nsCtAesi.fn_editStreamItem(inst.dat_item);
    });
    inst.el_toggleState.addEventListener('click', async function() {
      inst.fn_markPendingToggleState();

      // TODO: refresh only this instance + visual flair
      // TEMP: quick and dirty (should be object shallow copy to maintain data integrity)
      const prevState = inst.dat_item.state;
      inst.dat_item.state = (inst.dat_item.state === 'active') ? 'inactive' : 'active';

      const respBundle = await nsCtSil.lapi_editStreamItem(inst.dat_item);

      if (respBundle.isOk()) {
        nsCtSil.fn_displayOkMessage('State toggle ok', `Stream item '${inst.dat_item.code}' is now '${inst.dat_item.state}'`);
      } else {
        inst.dat_item.state = prevState;
        nsCtSil.fn_displayErrorMessage('State toggle error', `Unable to change state for stream item '${inst.dat_item.code}'`);
      }

      inst.fn_clearPendingAction();
      inst.fn_renderStreamState();
      inst.fn_renderPreviewButton();
    });
    inst.el_preview.addEventListener('click', function() {
      nsCtSil.fn_previewStream(inst.dat_item);
    });

    return inst
  };

  nsCtSil.fn_displayInfoMessage = function(title, desc) {
    nsCtSil.fn_displayMessage('info', title, desc);
  };
  nsCtSil.fn_displayOkMessage = function(title, desc) {
    nsCtSil.fn_displayMessage('success', title, desc);
  };
  nsCtSil.fn_displayErrorMessage = function(title, desc) {
    nsCtSil.fn_displayMessage('danger', title, desc);
  };
  nsCtSil.fn_displayMessage = function(severity, title, desc) {
    nsCtSil.el_silMessage.innerHTML = `<div class="alert alert-${severity} m-0 p-2 mb-3">
  <h4 class="alert-title mb-0">${title}</h4>
  <div class="text-secondary">${desc}</div>
</div>
`;

    nsCtSil.el_streamItemListing.classList.add('sil-listing-offset-message');
  };
  nsCtSil.fn_clearMessage = function() {
    nsCtSil.el_silMessage.innerHTML = '';
    nsCtSil.el_streamItemListing.classList.remove('sil-listing-offset-message');
  };

  nsCtSil.fn_previewStream = function(streamItem) {
    nsCtSil.el_silPreviewTitle.innerHTML = `Stream View<br>Stream Item ${streamItem.code}`;

    const componentId = 'nsCtSilStreamPreview';

    if (_nsDLS.componentExists(componentId)) {
      _nsDLS.changeHlsSource(componentId, streamItem.stream_url);
    } else {
      _nsDLS.newComponent(componentId, nsCtSil.el_silPreviewDls, streamItem.stream_url);
    }
  };
  nsCtSil.fn_resetStreamPreview = function() {
    nsCtSil.el_silPreviewTitle.innerHTML = 'Stream View<br>';

    const componentId = 'nsCtSilStreamPreview';

    if (_nsDLS.componentExists(componentId)) {
      _nsDLS.destroyComponent(componentId);
    }
  };

  nsCtSil.wsReq_streamItemListing = function() {
    nsCtWs.fn_sendMessage({
      _bid: '',
      _brc: nsCtWs.CONSTANT.WS_REQCODE_STREAM_ITEM_LISTING,
    });
  };

  nsCtSil.lapi_removeStreamItem = async function(id) {
    const reqURI = '/stream/local-api/delete-stream-item';
    const reqBody = {'id': id};

    let respBundle = null;
    await nsComm.lapi_defaultPOST(reqURI, JSON.stringify(reqBody))
            .then((_respBundle) => { respBundle = _respBundle })
            .catch((_err) => { console.error(`TODO: deadend ${_err}`); });

    return respBundle
  };

  // temp, use existing endpoint, TODO: reqRep via Ws
  nsCtSil.lapi_editStreamItem = async function(data) {
    const reqURI = '/stream/local-api/edit-stream-item';
    const reqBody = data;

    let respBundle = null;
    await nsComm.lapi_defaultPOST(reqURI, JSON.stringify(reqBody))
            .then((_respBundle) => { respBundle = _respBundle })
            .catch((_err) => { console.error(`TODO: deadend ${_err}`); });

    return respBundle
  };


  nsCtSil.el_silNavTabbing.addEventListener('click', function() {
    nsCtSil.fn_resetStreamPreview();
  });

  nsCtSil.dat_silSearchSubj
    .pipe(
      rxjs.debounceTime(333),
    )
    .subscribe({
      next: searchStr => {
        nsCtSil.fn_searchStreamItemListing(searchStr);
      },
    });

  nsCtSil.el_silSearch.addEventListener('input', function() {
    nsCtSil.dat_silSearchSubj.next(this.value);
  });

  nsCtSil.el_silReload.addEventListener('click', function() {
    nsCtSil.fn_resetSearch();
    nsCtSil.wsReq_streamItemListing();
  });

  nsCtSil.el_silPreviewStop.addEventListener('click', function() {
    nsCtSil.fn_resetStreamPreview();
  });

  nsCtSil.fn_resetStreamPreview();
};

// add / edit stream item - header
const nsCtAesi = {ns: 'nsCtAesi'};
nsCtAesi.init = function() {
  nsCtAesi.el_asiNavTabbing = document.getElementById('asiNavTabbing');
  nsCtAesi.el_asimdNavTabbing = document.getElementById('asimdNavTabbing');
  nsCtAesi.el_asieNavTabbing = document.getElementById('asieNavTabbing');
  nsCtAesi.el_asifNavTabbing = document.getElementById('asifNavTabbing');

  nsCtAesi.fn_editStreamItem = function(streamItem) {
    switch (streamItem.source_type) {
      case 'mod_device':
        nsCtAesi.el_asiNavTabbing.click();
        nsCtAesi.el_asimdNavTabbing.click();
        nsCtAesimd.fn_setFormEditMode(streamItem);
        nsCtAesimd.fn_resetDeviceChannelPreview();
        nsCtAesimd.fn_reverseSelection(streamItem.device_code, streamItem.device_channel_id, streamItem.device_stream_type);
        break;

      case 'external':
        nsCtAesi.el_asiNavTabbing.click();
        nsCtAesi.el_asieNavTabbing.click();
        nsCtAesie.fn_setFormEditMode(streamItem);
        break;

      case 'file':
        nsCtAesi.el_asiNavTabbing.click();
        nsCtAesi.el_asifNavTabbing.click();
        nsCtAesif.fn_setFormEditMode(streamItem);
        break;
    }
  };

  // temp, use existing endpoint, TODO: reqRep via Ws
  nsCtAesi.lapi_addStreamItem = async function(data) {
    const reqURI = '/stream/local-api/add-stream-item';
    const reqBody = data;

    let respBundle = null;
    await nsComm.lapi_defaultPOST(reqURI, JSON.stringify(reqBody))
            .then((_respBundle) => { respBundle = _respBundle })
            .catch((_err) => { console.error(`TODO: deadend ${_err}`); });

    return respBundle
  };
  nsCtAesi.lapi_editStreamItem = async function(data) {
    const reqURI = '/stream/local-api/edit-stream-item';
    const reqBody = data;

    let respBundle = null;
    await nsComm.lapi_defaultPOST(reqURI, JSON.stringify(reqBody))
            .then((_respBundle) => { respBundle = _respBundle })
            .catch((_err) => { console.error(`TODO: deadend ${_err}`); });

    return respBundle
  };
};

// add / edit stream item - mod_device
const nsCtAesimd = {ns: 'nsCtAesimd'};
nsCtAesimd.init = function() {
  nsCtAesimd.tmpl_device = `<div class="list-group-item faux-selection device-item p-3">
  <div class="row">
    <div class="col-12">
      <div class="card-title mb-2 dev-code"></div>
      <div class="card-subtitle mb-2 dev-name"></div>
    </div>
  </div>
  <div class="hr-text hr-text-center hr-text-spaceless mt-1 mb-2">Hardware info</div>
  <div class="row">
    <div class="col-4"><div class="text-secondary">Status</div></div>
    <div class="col-8">
      <div class="text-secondary text-truncate dev-status"></div>
    </div>
  </div>
  <div class="row">
    <div class="col-4"><div class="text-secondary">Brand</div></div>
    <div class="col-8">
      <div class="text-secondary text-truncate dev-brand"></div>
    </div>
  </div>
  <div class="row">
    <div class="col-4"><div class="text-secondary">Name</div></div>
    <div class="col-8">
      <div class="text-secondary text-truncate dev-hwname"></div>
    </div>
  </div>
  <div class="row">
    <div class="col-4"><div class="text-secondary">Model</div></div>
    <div class="col-8">
      <div class="text-secondary text-truncate dev-model"></div>
    </div>
  </div>
  <div class="row">
    <div class="col-4"><div class="text-secondary">Type</div></div>
    <div class="col-8">
      <div class="text-secondary text-truncate dev-type"></div>
    </div>
  </div>
  <div class="row">
    <div class="col-4"><div class="text-secondary">Channels</div></div>
    <div class="col-8">
      <div class="text-secondary text-truncate dev-channels"></div>
    </div>
  </div>
</div>
`;

  nsCtAesimd.tmpl_devChan = `<div class="list-group-item faux-selection channel-item p-3">
  <div class="row">
    <div class="col-3">
      <div class="rounded devchan-thumbnail-container">
        <img class="devchan-thumbnail" src="" alt="">
      </div>
    </div>
    <div class="col-6">
      <div class="row">
        <div class="col">
          <div class="card-title mb-2 devchan-title"></div>
        </div>
      </div>
      <div class="hr-text hr-text-center hr-text-spaceless mt-1 mb-2">Hardware info</div>
      <div class="row">
        <div class="col-4"><div class="text-secondary">Name</div></div>
        <div class="col-8">
          <div class="text-secondary text-truncate devchan-name"></div>
        </div>
      </div>
      <div class="row">
        <div class="col-4"><div class="text-secondary">Status</div></div>
        <div class="col-8">
          <div class="text-secondary text-truncate devchan-status"></div>
        </div>
      </div>
    </div>
    <div class="col-3">
      <div class="col devchan-stream-type-container">
      </div>
    </div>
  </div>
</div>
`;
  nsCtAesimd.tmpl_devChanStreamTypeSet = `
<div class="row">
  <div class="btn-group">
    <button class="btn devchan-select"></button>
    <button class="btn btn-icon devchan-preview">
      <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-eye"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M10 12a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" /><path d="M21 12c-2.4 4 -5.4 6 -9 6c-3.6 0 -6.6 -2 -9 -6c2.4 -4 5.4 -6 9 -6c3.6 0 6.6 2 9 6" /></svg>
    </button>
  </div>
</div>
`;

  nsCtAesimd.const_dahuaStreamTypes = [
    ['main', 'Main Stream'],
    // ['extra1', 'Extra Stream 1'],
    // ['extra2', 'Extra Stream 2'],
    // ['extra3', 'Extra Stream 3'],
  ];
  nsCtAesimd.const_hikStreamTypes = [
    ['main', 'Main Stream'],
    ['sub', 'Sub Stream'],
    ['third', 'Third Stream'],
  ];
  nsCtAesimd.const_panaNetcamStreamTypes = [
    ['stream1', 'Stream(1)'],
    ['stream2', 'Stream(2)'],
    ['stream3', 'Stream(3)'],
    ['stream4', 'Stream(4)'],
  ];

  nsCtAesimd.el_msgContainer = document.getElementById('aesimdMessageContainer')

  nsCtAesimd.el_form = document.getElementById('formStreamItem_modDev');
  nsCtAesimd.el_title = document.getElementById('fsimdTitle');
  nsCtAesimd.el_formReset = document.getElementById('fsimdReset');
  nsCtAesimd.el_formCancel = document.getElementById('fsimdCancel');
  nsCtAesimd.el_formSave = document.getElementById('fsimdSave');

  nsCtAesimd.el_dcpPreviewTitle = document.getElementById('dcpPreviewTitle');
  nsCtAesimd.el_dcpPreviewStop = document.getElementById('dcpPreviewStop');
  nsCtAesimd.el_dcpPreviewDls = document.getElementById('dcpPreviewDls');
  nsCtAesimd.el_dcpPreviewTempFeedback = document.getElementById('dcpPreviewTempFeedback');

  nsCtAesimd.el_dlTitle = document.getElementById('dlTitle');
  nsCtAesimd.el_dlSearch = document.getElementById('dlSearch');
  nsCtAesimd.el_dlReload = document.getElementById('dlReload');
  nsCtAesimd.el_deviceListing = document.getElementById('deviceListing');

  nsCtAesimd.el_dclTitle = document.getElementById('dclTitle');
  nsCtAesimd.el_deviceChannelListing = document.getElementById('deviceChannelListing');

  nsCtAesimd.dat_wsSubs = {};

  nsCtAesimd.dat_formMode = 'add'; // 'add' | 'edit'
  nsCtAesimd.dat_formData = null;

  nsCtAesimd.dat_deviceInstances = new Map();
  nsCtAesimd.dat_devChanInstances = new Map();

  nsCtAesimd.dat_referencedDeviceChannels = new Map();

  nsCtAesimd.dat_selectedDevice = null;
  // nsCtAesimd.dat_selectedDeviceChannel = null;

  nsCtAesimd.dat_dlSearchSubj = new rxjs.Subject();

  nsCtAesimd.fn_setupWs = function() {
    nsCtAesimd.dat_wsSubs['ws:open'] = nsCtWs.dat_wsObss['open']
      .subscribe({
        next: (v) => {
          nsCtAesimd.wsReq_deviceSnapshotListing();
        },
      });

    nsCtAesimd.dat_wsSubs['ws:msg:rr'] = nsCtWs.dat_wsObss['msg:rr:asi']
      .subscribe({
        next: (p_rep) => {
          switch(p_rep._brc) {
            case nsCtWs.CONSTANT.WS_REQCODE_DEVICE_SNAPSHOT_LISTING:
              nsCtAesimd.fn_populateDeviceInstances(p_rep.device_snapshots);
              nsCtAesimd.fn_renderDeviceListing();

              nsCtAesimd.fn_populateDevChanInstances(null);
              nsCtAesimd.fn_renderDevChanListing(null);
              break;
            
            default:
              console.warn(`${nsCtAesimd.ns}`, `no handler for req_code '${p_rep._brc}'`);
              break;
          }
        },
      });
  };

  nsCtAesimd.fn_displayInfoMessage = function(title, desc) { nsCtAesimd.fn_displayMessage('info', title, desc); };
  nsCtAesimd.fn_displayOkMessage = function(title, desc) { nsCtAesimd.fn_displayMessage('success', title, desc); };
  nsCtAesimd.fn_displayErrorMessage = function(title, desc) { nsCtAesimd.fn_displayMessage('danger', title, desc); };
  nsCtAesimd.fn_displayMessage = function(severity, title, desc) {
    nsCtAesimd.el_msgContainer.innerHTML = `<div class="alert alert-${severity} m-0 p-2 mb-3">
  <h4 class="alert-title mb-0">${title}</h4>
  <div class="text-secondary">${desc}</div>
</div>`;
  };
  nsCtAesimd.fn_clearMessage = function() {
    nsCtAesimd.el_msgContainer.innerHTML = '';
  };

  nsCtAesimd.fn_setFormAddMode = function() {
    nsCtAesimd.fn_clearMessage();

    nsCtAesimd.dat_formMode = 'add';
    nsCtAesimd.dat_formData = null;

    nsCtAesimd.el_title.innerHTML = 'New Stream Item';

    nsCtAesimd.el_form.reset();
    nsCtAesimd.el_form.id.value = '';

    nsCtAesimd.el_form.code.disabled = false;
    nsCtAesimd.el_formCancel.style.display = 'none';
  };
  nsCtAesimd.fn_setFormEditMode = function(data) {
    nsCtAesimd.fn_clearMessage();

    nsCtAesimd.dat_formMode = 'edit';
    nsCtAesimd.dat_formData = data;

    nsCtAesimd.el_title.innerHTML = `Edit Stream Item '${nsCtAesimd.dat_formData.code}'`;

    nsCtAesimd.el_form.reset();
    nsCtAesimd.el_form.id.value = nsCtAesimd.dat_formData.id;
    nsCtAesimd.el_form.code.value = nsCtAesimd.dat_formData.code;
    nsCtAesimd.el_form.name.value = nsCtAesimd.dat_formData.name;
    nsCtAesimd.el_form.state.value = nsCtAesimd.dat_formData.state;
    nsCtAesimd.el_form.note.value = nsCtAesimd.dat_formData.note;
    nsCtAesimd.el_form.deviceCode.value = nsCtAesimd.dat_formData.device_code; // TODO: normalize
    nsCtAesimd.el_form.deviceChannelId.value = nsCtAesimd.dat_formData.device_channel_id; // TODO: normalize
    nsCtAesimd.el_form.deviceStreamType.value = nsCtAesimd.dat_formData.device_stream_type; // TODO: normalize

    nsCtAesimd.el_form.code.disabled = true;
    nsCtAesimd.el_formCancel.style.display = '';
  };

  nsCtAesimd.fn_updateDeviceSelection = function(devInst) {
    nsCtAesimd.fn_resetDeviceSelection();

    devInst.dat_isSelected = true;
    devInst.el_root.classList.add('active');

    nsCtAesimd.dat_selectedDevice = devInst;
    nsCtAesimd.el_form.deviceCode.value = (devInst === null) ? '' : devInst.dat_item.persistence.code;

    let tempChannels = null;
    if (devInst.dat_item.hardware.analog_channels !== null) tempChannels = devInst.dat_item.hardware.analog_channels;
    if (tempChannels === null && devInst.dat_item.hardware.digital_channels !== null) tempChannels = devInst.dat_item.hardware.digital_channels;
    if (tempChannels === null) tempChannels = [];
    nsCtAesimd.fn_crossmatchAndCacheReferencedDeviceChannels(devInst.dat_item.persistence.code, tempChannels);

    nsCtAesimd.fn_populateDevChanInstances(devInst.dat_item);
    nsCtAesimd.fn_renderDevChanListing(devInst.dat_item.persistence.code);

    nsCtAesimd.el_deviceChannelListing.scrollTop = 0;
  };
  nsCtAesimd.fn_resetDeviceSelection = function() {
    nsCtAesimd.dat_deviceInstances.forEach((devInst) => {
      devInst.dat_isSelected = false;
      devInst.el_root.classList.remove('active');
    });

    if (nsCtAesimd.dat_selectedDevice !== null) {
      nsCtAesimd.fn_removeCacheForReferencedDeviceChannels(nsCtAesimd.dat_selectedDevice.dat_item.persistence.code);
    }

    nsCtAesimd.dat_selectedDevice = null;
    nsCtAesimd.el_form.deviceCode.value = '';

    nsCtAesimd.fn_populateDevChanInstances(null);
    nsCtAesimd.fn_renderDevChanListing(null);
  };

  nsCtAesimd.fn_updateDeviceChannelSelection = function(devChanInst) {
    nsCtAesimd.fn_resetDeviceChannelSelection();

    devChanInst.dat_isSelected = true;
    devChanInst.el_root.classList.add('active');

    nsCtAesimd.el_form.deviceChannelId.value = (devChanInst == null) ? '' : devChanInst.id;
  };
  nsCtAesimd.fn_resetDeviceChannelSelection = function() {
    nsCtAesimd.dat_devChanInstances.forEach((devChanInst) => {
      if (!devChanInst.dat_isSelected) return;

      devChanInst.dat_isSelected = false; 
      devChanInst.el_root.classList.remove('active');
    });
  };

  nsCtAesimd.fn_updateDeviceStreamTypeSelection = function(devChanInst, _el_select) {
    nsCtAesimd.fn_resetDeviceStreamTypeSelection(devChanInst);

    _el_select.classList.add('btn-primary');
    devChanInst.dat_currSelectedButtonElement = _el_select;

    nsCtAesimd.el_form.deviceStreamType.value = _el_select.getAttribute('data-fu-streamType');
  };
  nsCtAesimd.fn_resetDeviceStreamTypeSelection = function(devChanInst) {
    nsCtAesimd.dat_devChanInstances.forEach((_devChanInst) => {
      if (_devChanInst.dat_currSelectedButtonElement === null) return;

      _devChanInst.dat_currSelectedButtonElement.classList.remove('btn-primary');
      _devChanInst.dat_currSelectedButtonElement = null;
    });
  };

  nsCtAesimd.fn_reverseSelection = function(deviceCode, deviceChannelId, deviceStreamType) {
    let devIndex = -1;
    let devChanIndex = -1;

    nsCtAesimd.dat_deviceInstances.forEach((devInst) => {
      devIndex += 1;

      if (deviceCode === devInst.dat_item.persistence.code) {
        devInst.fn_select();
        nsCtAesimd.el_deviceListing.scrollTop = (devIndex * devInst.el_root.offsetHeight);

        nsCtAesimd.dat_devChanInstances.forEach((devChanInst) => {
          devChanIndex += 1;

          if (deviceChannelId === devChanInst.dat_item.channel_id) {
            devChanInst.fn_select();
            nsCtAesimd.el_deviceChannelListing.scrollTop = (devChanIndex * devChanInst.el_root.offsetHeight);

            for (let i = 0; i < devChanInst.el_streamTypeSelects.length; i++) {
              let _el_select = devChanInst.el_streamTypeSelects.item(i);

              if (deviceStreamType === _el_select.getAttribute('data-fu-streamType')) {
                nsCtAesimd.fn_updateDeviceStreamTypeSelection(devChanInst, _el_select);
              }
            }

            return;
          }
        });

        return;
      }
    });
  };

  nsCtAesimd.fn_formReset = function() {
    switch (nsCtAesimd.dat_formMode) {
      case 'add':
        nsCtAesimd.fn_setFormAddMode();
        nsCtAesimd.fn_resetDeviceChannelPreview();
        nsCtAesimd.fn_resetDeviceSelection();
        break;

      case 'edit':
        nsCtAesimd.fn_setFormEditMode(nsCtAesimd.dat_formData);
        nsCtAesimd.fn_resetDeviceChannelPreview();
        nsCtAesimd.fn_reverseSelection(nsCtAesimd.dat_formData.device_code, nsCtAesimd.dat_formData.device_channel_id, nsCtAesimd.dat_formData.device_stream_type);
        break;
    }
  };

  nsCtAesimd.fn_formValidate = function() {
    // TODO
    return true;
  };

  nsCtAesimd.fn_formSaveNewData = async function() {
    const data = {
      'stream_group_id': nsCtAesimd.el_form.streamGroupId.value,

      'id': nsCtAesimd.el_form.id.value,
      'code': nsCtAesimd.el_form.code.value,
      'name': nsCtAesimd.el_form.name.value,
      'state': nsCtAesimd.el_form.state.value,
      'note': nsCtAesimd.el_form.note.value,

      'source_type': 'mod_device',
      'device_code': nsCtAesimd.el_form.deviceCode.value,
      'device_channel_id': nsCtAesimd.el_form.deviceChannelId.value,
      'device_stream_type': nsCtAesimd.el_form.deviceStreamType.value,
    };

    let respBundle = await nsCtAesi.lapi_addStreamItem(data);

    if (respBundle.isOk()) {
      nsCtSil.wsReq_streamItemListing();
      nsCtSil.fn_clearMessage();

      nsCtAesimd.fn_setFormAddMode();
      nsCtAesimd.fn_resetDeviceChannelPreview();
      nsCtAesimd.fn_resetDeviceSelection();
      nsCtAesimd.fn_displayOkMessage('Add ok', `${respBundle.defResp.message}`);

    } else {
      nsCtAesimd.fn_displayErrorMessage('Add error', `${respBundle.defResp.message}`);
    }
  };

  nsCtAesimd.fn_formSaveExistingData = async function() {
    const data = {
      'stream_group_id': nsCtAesimd.el_form.streamGroupId.value,

      'id': nsCtAesimd.el_form.id.value,
      'code': nsCtAesimd.el_form.code.value,
      'name': nsCtAesimd.el_form.name.value,
      'state': nsCtAesimd.el_form.state.value,
      'note': nsCtAesimd.el_form.note.value,

      'source_type': 'mod_device',
      'device_code': nsCtAesimd.el_form.deviceCode.value,
      'device_channel_id': nsCtAesimd.el_form.deviceChannelId.value,
      'device_stream_type': nsCtAesimd.el_form.deviceStreamType.value,
    };

    let respBundle = await nsCtAesi.lapi_editStreamItem(data);

    if (respBundle.isOk()) {
      nsCtSil.wsReq_streamItemListing();
      nsCtSil.fn_displayOkMessage('Update ok', `${respBundle.defResp.message}`);
      nsCtSil.el_silNavTabbing.click();

      nsCtAesimd.fn_setFormAddMode();
      nsCtAesimd.fn_resetDeviceChannelPreview();
      nsCtAesimd.fn_resetDeviceSelection();

    } else {
      nsCtAesimd.fn_displayErrorMessage('Update error', `${respBundle.defResp.message}`);
    }
  };

  nsCtAesimd.el_formReset.addEventListener('click', function(ev) {
    ev.preventDefault();
    nsCtAesimd.fn_formReset();
  }); 
  nsCtAesimd.el_formCancel.addEventListener('click', function(ev) {
    ev.preventDefault();
    nsCtAesimd.fn_setFormAddMode();
    nsCtAesimd.fn_resetDeviceChannelPreview();
    nsCtAesimd.fn_resetDeviceSelection();
    nsCtSil.el_silNavTabbing.click();
  });
  nsCtAesimd.el_formSave.addEventListener('click', function(ev) {
    ev.preventDefault();

    switch (nsCtAesimd.dat_formMode) {
      case 'add': nsCtAesimd.fn_formSaveNewData(); break;
      case 'edit': nsCtAesimd.fn_formSaveExistingData(); break;
      default: return;
    }
  });


  nsCtAesimd.fn_previewDeviceChannel = async function(deviceCode, deviceChannelId, streamType) {
    nsCtAesimd.el_dcpPreviewTitle.innerHTML = `Device ${deviceCode}, channel ${deviceChannelId} - ${streamType}`;
    nsCtAesimd.el_dcpPreviewTempFeedback.innerHTML = ``;

    const componentId = 'nsCtAesimdDeviceChannelPreview';

    const reqURI = '/stream/local-api/device-channel-preview';
    const reqBody = {'device_code': deviceCode, 'device_channel_id': deviceChannelId, 'device_stream_type': streamType};

    let respBundle = null;
    await nsComm.lapi_defaultPOST(reqURI, JSON.stringify(reqBody))
            .then((_respBundle) => { respBundle = _respBundle })
            .catch((_err) => { console.error(`TODO: deadend ${_err}`); });

    let previewUrl = null;
    if (respBundle.isOk()) {
      if (respBundle.defResp.isOk()) {
        previewUrl = respBundle.defResp.data.preview_url;

      } else {
        nsCtAesimd.el_dcpPreviewTempFeedback.innerHTML = `${respBundle.defResp.message}`;
        console.error('fn_previewDeviceChannel:', respBundle.defResp.message);
      }
    } else {
      nsCtAesimd.el_dcpPreviewTempFeedback.innerHTML = `${respBundle.defResp.message}`;
      console.error('fn_previewDeviceChannel:', respBundle.defResp.message);
    }

    if (previewUrl === null) return

    if (_nsDLS.componentExists(componentId)) {
      _nsDLS.changeHlsSource(componentId, previewUrl);
    } else {
      _nsDLS.newComponent(componentId, nsCtAesimd.el_dcpPreviewDls, previewUrl);
    }
  };
  nsCtAesimd.fn_resetDeviceChannelPreview = function() {
    nsCtAesimd.el_dcpPreviewTitle.innerHTML = 'Device Channel Preview';
    nsCtAesimd.el_dcpPreviewTempFeedback.innerHTML = ``;

    const componentId = 'nsCtAesimdDeviceChannelPreview';

    if (_nsDLS.componentExists(componentId)) {
      _nsDLS.destroyComponent(componentId);
    }
  };

  nsCtAesimd.fn_populateDeviceInstances = function(deviceSnapshots) {
    nsCtAesimd.dat_deviceInstances.clear();

    for (let i = 0; i < deviceSnapshots.length; i++) {
      const deviceInstance = nsCtAesimd.fn_newDeviceInstance(i, deviceSnapshots[i]);

      nsCtAesimd.dat_deviceInstances.set(deviceInstance.id, deviceInstance);
    }
  };

  nsCtAesimd.fn_renderDeviceListing = function() {
    nsCtAesimd.el_dlTitle.innerHTML = `Devices (${nsCtAesimd.dat_deviceInstances.size} items)`;
    nsCtAesimd.el_deviceListing.innerHTML = '';

    nsCtAesimd.dat_deviceInstances.forEach((inst) => {
      nsCtAesimd.el_deviceListing.appendChild(inst.el_root);

      inst.fn_renderContent();
    });
  };

  nsCtAesimd.fn_newDeviceInstance = function(idx, deviceSnapshot) {
    const inst = {};
    inst.id = deviceSnapshot.persistence.id;
    // inst.guid = `dl-${deviceSnapshot.persistence.id}`;
    inst.dat_item = deviceSnapshot;

    const tempRoot = document.createElement('div');
    tempRoot.innerHTML = nsCtAesimd.tmpl_device;
    inst.el_root = tempRoot.children[0];
    // inst.el_root.id = inst.guid;

    inst.el_code = inst.el_root.getElementsByClassName('dev-code')[0];
    inst.el_name = inst.el_root.getElementsByClassName('dev-name')[0];
    inst.el_status = inst.el_root.getElementsByClassName('dev-status')[0];
    inst.el_brand = inst.el_root.getElementsByClassName('dev-brand')[0];
    inst.el_hwname = inst.el_root.getElementsByClassName('dev-hwname')[0];
    inst.el_model = inst.el_root.getElementsByClassName('dev-model')[0];
    inst.el_type = inst.el_root.getElementsByClassName('dev-type')[0];
    inst.el_channels = inst.el_root.getElementsByClassName('dev-channels')[0];

    inst.dat_isSelected = false;

    inst.fn_renderContent = function() {
      let statusText = '';
      switch(inst.dat_item.live.conn_state) {
        case 'never':
          statusText = '<span class="badge bg-secondary me-1"></span> Never';
          break;
        case 'alive':
          statusText = '<span class="badge bg-success me-1"></span> OK';
          break;
        case 'lost':
          statusText = `<span class="badge bg-danger me-1"></span> LOST. Last seen at ${inst.dat_item.live.last_seen}.`;
          break;
        default:
          break;
      }

      let analogChannelCount = 0;
      let digitalChannelCount = 0;

      if (inst.dat_item.hardware.analog_channels !== null) {
        analogChannelCount = inst.dat_item.hardware.analog_channels.length;
      }
      if (inst.dat_item.hardware.digital_channels !== null) {
        digitalChannelCount = inst.dat_item.hardware.digital_channels.length;
      }

      inst.el_code.innerHTML = inst.dat_item.persistence.code;
      inst.el_name.innerHTML = inst.dat_item.persistence.name;
      inst.el_status.innerHTML = statusText;
      inst.el_brand.innerHTML = inst.dat_item.persistence.brand;
      inst.el_hwname.innerHTML = inst.dat_item.hardware.device_name;
      inst.el_model.innerHTML = inst.dat_item.hardware.model;
      inst.el_type.innerHTML = inst.dat_item.hardware.device_type;
      inst.el_channels.innerHTML = `Analog: ${analogChannelCount} | Digital: ${digitalChannelCount}`;
    };

    inst.fn_select = function() {
      nsCtAesimd.fn_updateDeviceSelection(inst);
    };

    inst.fn_show = function() { inst.el_root.style.display = ''; };
    inst.fn_hide = function() { inst.el_root.style.display = 'none'; };

    inst.el_root.addEventListener('click', inst.fn_select);

    return inst;
  };

  nsCtAesimd.fn_searchDeviceListing = function(searchStr) {
    let hasSearchStr = true;
    let shownCount = 0;

    if (searchStr === undefined || searchStr === null || searchStr === '') {
      hasSearchStr = false;
    } else {
      searchStr = searchStr.toLowerCase();
    }

    nsCtAesimd.dat_deviceInstances.forEach((inst) => {
      let doShow = true;

      if (hasSearchStr) {
        doShow = (
          inst.dat_item.persistence.code.toLowerCase().includes(searchStr)
          || inst.dat_item.persistence.name.toLowerCase().includes(searchStr)
        );
      }

      if (doShow) {
        inst.fn_show();
        shownCount += 1;
      } else {
        inst.fn_hide();
      }
    });

    nsCtAesimd.el_dlTitle.innerHTML = `Devices (${shownCount} / ${nsCtAesimd.dat_deviceInstances.size} items)`;
  };
  nsCtAesimd.fn_resetSearch = function() {
    nsCtAesimd.el_dlSearch.value = '';
  };

  nsCtAesimd.fn_populateDevChanInstances = function(deviceSnapshot) {
    nsCtAesimd.dat_devChanInstances.clear();

    if (deviceSnapshot === null) return; // on reset

    let tempChannels;

    tempChannels = deviceSnapshot.hardware.analog_channels;
    if (tempChannels !== null) {
      for (let i = 0; i < tempChannels.length; i++) {
        const devChanInstance = nsCtAesimd.fn_newDeviceChannelInstance(
          i,
          deviceSnapshot.persistence.brand,
          tempChannels[i],
        ); 
  
        nsCtAesimd.dat_devChanInstances.set(devChanInstance.id, devChanInstance);
      }
    }

    // TODO: targeted list to show (on observation and iteration of channel stuffs)
    //       for now, use faux handling. assuming if this is an analog hardware
    //       then, the analog_channels will be present, else try digital_channels

    if (nsCtAesimd.dat_devChanInstances.size === 0) {
      tempChannels = deviceSnapshot.hardware.digital_channels;
      if (tempChannels !== null) {
        for (let i = 0; i < tempChannels.length; i++) {
          const devChanInstance = nsCtAesimd.fn_newDeviceChannelInstance(
            i,
            deviceSnapshot.persistence.brand,
            tempChannels[i],
          ); 
    
          nsCtAesimd.dat_devChanInstances.set(devChanInstance.id, devChanInstance);
        }
      }
    }
  };

  nsCtAesimd.fn_renderDevChanListing = function(deviceCode) {
    if (deviceCode === null) {
      nsCtAesimd.el_dclTitle.innerHTML = 'Channels';
    } else {
      nsCtAesimd.el_dclTitle.innerHTML = `Channels of Device '${deviceCode}' (${nsCtAesimd.dat_devChanInstances.size} items)`;
    }

    nsCtAesimd.el_deviceChannelListing.innerHTML = '';

    nsCtAesimd.dat_devChanInstances.forEach((inst) => {
      nsCtAesimd.el_deviceChannelListing.appendChild(inst.el_root);

      inst.fn_renderContent();
    });
  };

  nsCtAesimd.fn_newDeviceChannelInstance = function(idx, deviceBrand, deviceChannel) {
    const inst = {};
    inst.id = deviceChannel.channel_id;
    inst.dat_item = deviceChannel;

    const tempRoot = document.createElement('div');
    tempRoot.innerHTML = nsCtAesimd.tmpl_devChan;
    inst.el_root = tempRoot.children[0];

    inst.el_thumbnail = inst.el_root.getElementsByClassName('devchan-thumbnail')[0];
    inst.el_title = inst.el_root.getElementsByClassName('devchan-title')[0];
    inst.el_streamTypeContainer = inst.el_root.getElementsByClassName('devchan-stream-type-container')[0];
    inst.el_streamTypeSelects = []; // expect array of html elements
    inst.el_streamTypePreviews = []; // expect array of html elements
    inst.el_name = inst.el_root.getElementsByClassName('devchan-name')[0];
    inst.el_status = inst.el_root.getElementsByClassName('devchan-status')[0];

    {
      let streamSets = [];
      switch (deviceBrand) {
        case 'dahua': streamSets = nsCtAesimd.const_dahuaStreamTypes; break;
        case 'hikvision': streamSets = nsCtAesimd.const_hikStreamTypes; break;
        case 'panasonic-netcam': streamSets = nsCtAesimd.const_panaNetcamStreamTypes; break;
        default: break;
      }

      // optz, these can be cached
      for (let ii = 0; ii < streamSets.length; ii++) {
        let currSet = streamSets[ii];
        
        let _tempRoot = document.createElement('div');
        _tempRoot.innerHTML = nsCtAesimd.tmpl_devChanStreamTypeSet;

        let _tempSelectButton = _tempRoot.getElementsByClassName('devchan-select')[0];
        let _tempPreviewButton = _tempRoot.getElementsByClassName('devchan-preview')[0];

        _tempSelectButton.innerHTML = currSet[1];
        _tempSelectButton.setAttribute('data-fu-streamType', currSet[0]);
        _tempPreviewButton.setAttribute('data-fu-streamType', currSet[0]);

        inst.el_streamTypeContainer.appendChild(_tempRoot.children[0]);
      }

      inst.el_streamTypeSelects = inst.el_root.getElementsByClassName('devchan-select');
      inst.el_streamTypePreviews = inst.el_root.getElementsByClassName('devchan-preview');
    }


    inst.dat_isSelected = false;
    inst.dat_currSelectedButtonElement = null;


    inst.fn_renderContent = function() {
      let titleText = `Channel ${inst.id}`;
      if (nsCtAesimd.dat_selectedDevice !== null) {
        const currSelectedDeviceCode = nsCtAesimd.dat_selectedDevice.dat_item.persistence.code;
        const devCache = nsCtAesimd.dat_referencedDeviceChannels.get(currSelectedDeviceCode);

        if (devCache[inst.id] === 0) {
        } else {
          titleText += `<span class="text-green"> - already in this stream group (${devCache[inst.id]}x)</span>`;
        }
      }

      let statusText = inst.dat_item.enabled
        ? '<span class="badge bg-success devchan-status-text"></span> Enabled'
        : '<span class="badge bg-danger devchan-status-text"></span> Not enabled';
      
      inst.el_thumbnail.alt = `thumbnail - channel ${inst.id}`;
      inst.el_title.innerHTML = titleText;
      inst.el_name.innerHTML = inst.dat_item.channel_name;
      inst.el_status.innerHTML = statusText;

      // TODO: available stream type matching with data from hw
    };

    inst.fn_select = function() {
      nsCtAesimd.fn_updateDeviceChannelSelection(inst);

      if (inst.dat_currSelectedButtonElement === null) {
        nsCtAesimd.fn_updateDeviceStreamTypeSelection(inst, inst.el_streamTypeSelects[0]);
      }
    };


    inst.el_root.addEventListener('click', inst.fn_select);

    for (let i = 0; i < inst.el_streamTypeSelects.length; i++) {
      let _el_select = inst.el_streamTypeSelects.item(i);

      _el_select.addEventListener('click', function() {
        nsCtAesimd.fn_updateDeviceStreamTypeSelection(inst, _el_select);
      });
    }

    for (let i = 0; i < inst.el_streamTypePreviews.length; i++) {
      let _el_preview = inst.el_streamTypePreviews.item(i);
      let _el_select = inst.el_streamTypeSelects.item(i);

      _el_preview.addEventListener('click', function() {
        nsCtAesimd.fn_updateDeviceStreamTypeSelection(inst, _el_select);

        if (nsCtAesimd.dat_selectedDevice === null) return;

        nsCtAesimd.fn_resetDeviceChannelPreview();
        nsCtAesimd.fn_previewDeviceChannel(
          nsCtAesimd.dat_selectedDevice.dat_item.persistence.code,
          inst.id,
          _el_preview.getAttribute('data-fu-streamType'),
        );
      });
    }

    return inst;
  };

  nsCtAesimd.fn_crossmatchAndCacheReferencedDeviceChannels = function(deviceCode, deviceChannels) {
    const devCache = nsCtAesimd.dat_referencedDeviceChannels.get(deviceCode);

    if (devCache === undefined) {
      const devchanCache = {};

      for (let i = 0; i < deviceChannels.length; i++) {
        const devchanId = deviceChannels[i].channel_id;

        devchanCache[devchanId] = 0;
      }

      nsCtSil.dat_streamItemInstances.forEach((streamItem) => {
        if (streamItem.dat_item.source_type !== 'mod_device') return;
        if (deviceCode !== streamItem.dat_item.device_code) return;

        if (devchanCache[streamItem.dat_item.device_channel_id] === undefined) {
          devchanCache[streamItem.dat_item.device_channel_id] = 0;
        }

        devchanCache[streamItem.dat_item.device_channel_id] = devchanCache[streamItem.dat_item.device_channel_id] + 1;
      });

      nsCtAesimd.dat_referencedDeviceChannels.set(deviceCode, devchanCache);
    }
  };
  // this is wasteful, but too lazy to implement fine grained handling...
  nsCtAesimd.fn_removeCacheForReferencedDeviceChannels = function(deviceCode) {
    nsCtAesimd.dat_referencedDeviceChannels.delete(deviceCode);
  };

  nsCtAesimd.wsReq_deviceSnapshotListing = function() {
    nsCtWs.fn_sendMessage({
      _bid: '',
      _brc: nsCtWs.CONSTANT.WS_REQCODE_DEVICE_SNAPSHOT_LISTING,
    });
  };


  nsCtAesimd.el_dcpPreviewStop.addEventListener('click', function() {
    nsCtAesimd.fn_resetDeviceChannelPreview();
  });

  nsCtAesimd.dat_dlSearchSubj
    .pipe(
      rxjs.debounceTime(333),
    )
    .subscribe({
      next: searchStr => {
        nsCtAesimd.fn_searchDeviceListing(searchStr);
      },
    });

  nsCtAesimd.el_dlSearch.addEventListener('input', function() {
    nsCtAesimd.dat_dlSearchSubj.next(this.value);
  });

  nsCtAesimd.el_dlReload.addEventListener('click', function() {
    nsCtAesimd.fn_resetSearch();
    nsCtAesimd.wsReq_deviceSnapshotListing();
  });

  nsCtAesimd.fn_setFormAddMode();
  nsCtAesimd.fn_resetDeviceChannelPreview();
};

// add / edit stream item - external
const nsCtAesie = {ns: 'nsCtAesie'};
nsCtAesie.init = function() {
  nsCtAesie.el_msgContainer = document.getElementById('asieMessageContainer')

  nsCtAesie.el_form = document.getElementById('formStreamItem_external');
  nsCtAesie.el_title = document.getElementById('fsieTitle');
  nsCtAesie.el_formReset = document.getElementById('fsieReset');
  nsCtAesie.el_formCancel = document.getElementById('fsieCancel');
  nsCtAesie.el_formSave = document.getElementById('fsieSave');

  nsCtAesie.dat_formMode = 'add'; // 'add' | 'edit'
  nsCtAesie.dat_formData = null;

  nsCtAesie.fn_displayInfoMessage = function(title, desc) { nsCtAesie.fn_displayMessage('info', title, desc); };
  nsCtAesie.fn_displayOkMessage = function(title, desc) { nsCtAesie.fn_displayMessage('success', title, desc); };
  nsCtAesie.fn_displayErrorMessage = function(title, desc) { nsCtAesie.fn_displayMessage('danger', title, desc); };
  nsCtAesie.fn_displayMessage = function(severity, title, desc) {
    nsCtAesie.el_msgContainer.innerHTML = `
<div class="alert alert-${severity} m-0 p-2 mb-3">
  <h4 class="alert-title mb-0">${title}</h4>
  <div class="text-secondary">${desc}</div>
</div>`;
  };
  nsCtAesie.fn_clearMessage = function() {
    nsCtAesie.el_msgContainer.innerHTML = '';
  };

  nsCtAesie.fn_setFormAddMode = function() {
    nsCtAesie.fn_clearMessage();

    nsCtAesie.dat_formMode = 'add';
    nsCtAesie.dat_formData = null;

    nsCtAesie.el_title.innerHTML = 'New Stream Item';

    nsCtAesie.el_form.reset();
    nsCtAesie.el_form.id.value = '';

    nsCtAesie.el_form.code.disabled = false;
    nsCtAesie.el_formCancel.style.display = 'none';
  };
  nsCtAesie.fn_setFormEditMode = function(data) {
    nsCtAesie.fn_clearMessage();

    nsCtAesie.dat_formMode = 'edit';
    nsCtAesie.dat_formData = data;

    nsCtAesie.el_title.innerHTML = `Edit Stream Item '${nsCtAesie.dat_formData.code}'`;

    nsCtAesie.el_form.reset();
    nsCtAesie.el_form.id.value = nsCtAesie.dat_formData.id;
    nsCtAesie.el_form.code.value = nsCtAesie.dat_formData.code;
    nsCtAesie.el_form.name.value = nsCtAesie.dat_formData.name;
    nsCtAesie.el_form.state.value = nsCtAesie.dat_formData.state;
    nsCtAesie.el_form.note.value = nsCtAesie.dat_formData.note;
    nsCtAesie.el_form.external_url.value = nsCtAesie.dat_formData.external_url;

    nsCtAesie.el_form.code.disabled = true;
    nsCtAesie.el_formCancel.style.display = '';
  };

  nsCtAesie.fn_formReset = function() {
    switch (nsCtAesie.dat_formMode) {
      case 'add':
        nsCtAesie.fn_setFormAddMode();
        break;

      case 'edit':
        nsCtAesie.fn_setFormEditMode(nsCtAesie.dat_formData);
        break;
    }
  };

  nsCtAesie.fn_formValidate = function() {
    // TODO
    return true;
  };

  nsCtAesie.fn_formSaveNewData = async function() {
    const data = {
      'stream_group_id': nsCtAesie.el_form.streamGroupId.value,

      'id': nsCtAesie.el_form.id.value,
      'code': nsCtAesie.el_form.code.value,
      'name': nsCtAesie.el_form.name.value,
      'state': nsCtAesie.el_form.state.value,
      'note': nsCtAesie.el_form.note.value,

      'source_type': 'external',
      'external_url': nsCtAesie.el_form.external_url.value,
    };

    let respBundle = await nsCtAesi.lapi_addStreamItem(data);

    if (respBundle.isOk()) {
      nsCtSil.wsReq_streamItemListing();
      nsCtSil.fn_clearMessage();

      const addedCode = nsCtAesie.el_form.code.value;
      nsCtAesie.fn_setFormAddMode();
      nsCtAesie.fn_displayOkMessage('Add ok', `Stream item '${addedCode}' added.`);

    } else {
      nsCtAesie.fn_displayErrorMessage('Add error', `${respBundle.defResp.message}`);
    }
  };

  nsCtAesie.fn_formSaveExistingData = async function() {
    const data = {
      'stream_group_id': nsCtAesie.el_form.streamGroupId.value,

      'id': nsCtAesie.el_form.id.value,
      'code': nsCtAesie.el_form.code.value,
      'name': nsCtAesie.el_form.name.value,
      'state': nsCtAesie.el_form.state.value,
      'note': nsCtAesie.el_form.note.value,

      'source_type': 'external',
      'external_url': nsCtAesie.el_form.external_url.value,
    };

    let respBundle = await nsCtAesi.lapi_editStreamItem(data);

    if (respBundle.isOk()) {
      nsCtSil.wsReq_streamItemListing();
      nsCtSil.fn_displayOkMessage('Update ok', `Stream item '${nsCtAesie.el_form.code.value}' updated.`);
      nsCtSil.el_silNavTabbing.click();

      nsCtAesie.fn_setFormAddMode();

    } else {
      nsCtAesie.fn_displayErrorMessage('Update error', `${respBundle.defResp.message}`);
    }
  };

  nsCtAesie.el_formReset.addEventListener('click', function(ev) {
    ev.preventDefault();
    nsCtAesie.fn_formReset();
  }); 
  nsCtAesie.el_formCancel.addEventListener('click', function(ev) {
    ev.preventDefault();
    nsCtAesie.fn_setFormAddMode();
    nsCtSil.el_silNavTabbing.click();
  });
  nsCtAesie.el_formSave.addEventListener('click', function(ev) {
    ev.preventDefault();

    switch (nsCtAesie.dat_formMode) {
      case 'add': nsCtAesie.fn_formSaveNewData(); break;
      case 'edit': nsCtAesie.fn_formSaveExistingData(); break;
      default: return;
    }
  });

  nsCtAesie.fn_setFormAddMode();
};

// add / edit stream item - file
const nsCtAesif = {ns: 'nsCtAesif'};
nsCtAesif.init = function() {
  nsCtAesif.el_msgContainer = document.getElementById('asifMessageContainer')

  nsCtAesif.el_form = document.getElementById('formStreamItem_file');
  nsCtAesif.el_title = document.getElementById('fsifTitle');
  nsCtAesif.el_formReset = document.getElementById('fsifReset');
  nsCtAesif.el_formCancel = document.getElementById('fsifCancel');
  nsCtAesif.el_formSave = document.getElementById('fsifSave');

  nsCtAesif.dat_formMode = 'add'; // 'add' | 'edit'
  nsCtAesif.dat_formData = null;

  nsCtAesif.fn_displayInfoMessage = function(title, desc) { nsCtAesif.fn_displayMessage('info', title, desc); };
  nsCtAesif.fn_displayOkMessage = function(title, desc) { nsCtAesif.fn_displayMessage('success', title, desc); };
  nsCtAesif.fn_displayErrorMessage = function(title, desc) { nsCtAesif.fn_displayMessage('danger', title, desc); };
  nsCtAesif.fn_displayMessage = function(severity, title, desc) {
    nsCtAesif.el_msgContainer.innerHTML = `<div class="alert alert-${severity} m-0 p-2 mb-3">
  <h4 class="alert-title mb-0">${title}</h4>
  <div class="text-secondary">${desc}</div>
</div>`;
  };
  nsCtAesif.fn_clearMessage = function() {
    nsCtAesif.el_msgContainer.innerHTML = '';
  };

  nsCtAesif.fn_setFormAddMode = function() {
    nsCtAesif.fn_clearMessage();

    nsCtAesif.dat_formMode = 'add';
    nsCtAesif.dat_formData = null;

    nsCtAesif.el_title.innerHTML = 'New Stream Item';

    nsCtAesif.el_form.reset();
    nsCtAesif.el_form.id.value = '';

    nsCtAesif.el_form.code.disabled = false;
    nsCtAesif.el_formCancel.style.display = 'none';
  };
  nsCtAesif.fn_setFormEditMode = function(data) {
    nsCtAesif.fn_clearMessage();

    nsCtAesif.dat_formMode = 'edit';
    nsCtAesif.dat_formData = data;

    nsCtAesif.el_title.innerHTML = `Edit Stream Item '${nsCtAesif.dat_formData.code}'`;

    nsCtAesif.el_form.reset();
    nsCtAesif.el_form.id.value = nsCtAesif.dat_formData.id;
    nsCtAesif.el_form.code.value = nsCtAesif.dat_formData.code;
    nsCtAesif.el_form.name.value = nsCtAesif.dat_formData.name;
    nsCtAesif.el_form.state.value = nsCtAesif.dat_formData.state;
    nsCtAesif.el_form.note.value = nsCtAesif.dat_formData.note;
    nsCtAesif.el_form.filepath.value = nsCtAesif.dat_formData.filepath;

    nsCtAesif.el_form.code.disabled = true;
    nsCtAesif.el_formCancel.style.display = '';
  };

  nsCtAesif.fn_formReset = function() {
    switch (nsCtAesif.dat_formMode) {
      case 'add':
        nsCtAesif.fn_setFormAddMode();
        break;

      case 'edit':
        nsCtAesif.fn_setFormEditMode(nsCtAesif.dat_formData);
        break;
    }
  };

  nsCtAesif.fn_formValidate = function() {
    // TODO
    return true;
  };

  nsCtAesif.fn_formSaveNewData = async function() {
    const data = {
      'stream_group_id': nsCtAesif.el_form.streamGroupId.value,

      'id': nsCtAesif.el_form.id.value,
      'code': nsCtAesif.el_form.code.value,
      'name': nsCtAesif.el_form.name.value,
      'state': nsCtAesif.el_form.state.value,
      'note': nsCtAesif.el_form.note.value,

      'source_type': 'file',
      'filepath': nsCtAesif.el_form.filepath.value,
    };

    let respBundle = await nsCtAesi.lapi_addStreamItem(data);

    if (respBundle.isOk()) {
      nsCtSil.wsReq_streamItemListing();
      nsCtSil.fn_clearMessage();

      const addedCode = nsCtAesif.el_form.code.value;
      nsCtAesif.fn_setFormAddMode();
      nsCtAesif.fn_displayOkMessage('Add ok', `Stream item '${addedCode}' added.`);

    } else {
      nsCtAesif.fn_displayErrorMessage('Add error', `${respBundle.defResp.message}`);
    }
  };

  nsCtAesif.fn_formSaveExistingData = async function() {
    const data = {
      'stream_group_id': nsCtAesif.el_form.streamGroupId.value,

      'id': nsCtAesif.el_form.id.value,
      'code': nsCtAesif.el_form.code.value,
      'name': nsCtAesif.el_form.name.value,
      'state': nsCtAesif.el_form.state.value,
      'note': nsCtAesif.el_form.note.value,

      'source_type': 'file',
      'filepath': nsCtAesif.el_form.filepath.value,
    };

    let respBundle = await nsCtAesi.lapi_editStreamItem(data);

    if (respBundle.isOk()) {
      nsCtSil.wsReq_streamItemListing();
      nsCtSil.fn_displayOkMessage('Update ok', `Stream item '${nsCtAesif.el_form.code.value}' updated.`);
      nsCtSil.el_silNavTabbing.click();

      nsCtAesif.fn_setFormAddMode();

    } else {
      nsCtAesif.fn_displayErrorMessage('Update error', `${respBundle.defResp.message}`);
    }
  };

  nsCtAesif.el_formReset.addEventListener('click', function(ev) {
    ev.preventDefault();
    nsCtAesif.fn_formReset();
  }); 
  nsCtAesif.el_formCancel.addEventListener('click', function(ev) {
    ev.preventDefault();
    nsCtAesif.fn_setFormAddMode();
    nsCtSil.el_silNavTabbing.click();
  });
  nsCtAesif.el_formSave.addEventListener('click', function(ev) {
    ev.preventDefault();

    switch (nsCtAesif.dat_formMode) {
      case 'add': nsCtAesif.fn_formSaveNewData(); break;
      case 'edit': nsCtAesif.fn_formSaveExistingData(); break;
      default: return;
    }
  });

  nsCtAesif.fn_setFormAddMode();
};

document.addEventListener('DOMContentLoaded', function(e) {
  _nsWs.init();

  nsCt.init();

  if (nsCD && nsCD.isEditMode) {
    nsCtWs.init();
    nsCtSil.init();
    nsCtAesi.init();
    nsCtAesimd.init();
    nsCtAesie.init();
    nsCtAesif.init();

    nsCtWs.fn_setupWs();
    nsCtSil.fn_setupWs();
    nsCtAesimd.fn_setupWs();
  }

  // if (nsCD && nsCD.isManagingItems) {}
  nsMain.focusOn('code');
});
