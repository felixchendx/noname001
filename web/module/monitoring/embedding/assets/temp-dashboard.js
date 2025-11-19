'use strict';

const nsCtWs = {ns: 'nsCtWs'};
nsCtWs.CONSTANT = {
  WS_REQCODE__NODE_INFO_LISTING: '/ni/l',
  WS_REQCODE__NODE_INFO_ITEM: '/ni/i',
  WS_REQCODE__DEVICE_INFO_LISTING: '/di/l',
  WS_REQCODE__DEVICE_INFO_ITEM: '/di/i',
  WS_REQCODE__STREAM_INFO_LISTING: '/si/l',
  WS_REQCODE__STREAM_INFO_ITEM: '/si/i',
};
nsCtWs.init = function() {
  nsCtWs.dat_wsInstance = null;
  nsCtWs.dat_wsObss = {};
  nsCtWs.dat_wsSubs = {};

  nsCtWs.fn_setupWs = function() {
    const wsId = 'temp-dashboard'
    const wsUri = '/monitoring/temp-dashboard/ws';
    nsCtWs.dat_wsInstance = _nsWs.newWsInstance(wsId, wsUri);

    nsCtWs.dat_wsObss['open'] = nsCtWs.dat_wsInstance.onopenSubj.asObservable();
    nsCtWs.dat_wsObss['close'] = nsCtWs.dat_wsInstance.oncloseSubj.asObservable();
    nsCtWs.dat_wsObss['msg'] = nsCtWs.dat_wsInstance.onmessageSubj.asObservable().pipe(_nsWs.plainJsonParser);

    nsCtWs.dat_wsObss['msg:rr:node'] = new rxjs.Subject();
    nsCtWs.dat_wsObss['msg:rr:device'] = new rxjs.Subject();
    nsCtWs.dat_wsObss['msg:rr:stream'] = new rxjs.Subject();
    nsCtWs.dat_wsObss['msg:df:node'] = new rxjs.ReplaySubject(5);
    nsCtWs.dat_wsObss['msg:df:device'] = new rxjs.ReplaySubject(10);
    nsCtWs.dat_wsObss['msg:df:stream'] = new rxjs.ReplaySubject(10);

    // nsCtWs.dat_wsSubs['open'] = nsCtWs.dat_wsObss['open']
    //   .subscribe({
    //     next: (v) => {},
    //   });

    nsCtWs.dat_wsSubs['close'] = nsCtWs.dat_wsObss['close']
      .subscribe({
        next: (v) => {
          // TODO: reconnect control + feedback
          setTimeout(() => {
            _nsWs.reconnect(wsId);
          }, 10000);
        },
      });

    nsCtWs.dat_wsSubs['msg:router'] = nsCtWs.dat_wsObss['msg']
      .pipe(
        rxjs.tap((msg) => {
          if (msg._bh === undefined || msg._bh === null) {
            console.warn(nsCtWs.ns, 'msg_header not found');
            return
          }

          switch (msg._bh._bt) {
            case '_brr2':
              switch (msg._bh._brc) {
                case nsCtWs.CONSTANT.WS_REQCODE__NODE_INFO_LISTING:
                case nsCtWs.CONSTANT.WS_REQCODE__NODE_INFO_ITEM:
                  nsCtWs.dat_wsObss['msg:rr:node'].next(msg);
                  break

                case nsCtWs.CONSTANT.WS_REQCODE__DEVICE_INFO_LISTING:
                case nsCtWs.CONSTANT.WS_REQCODE__DEVICE_INFO_ITEM:
                  nsCtWs.dat_wsObss['msg:rr:device'].next(msg);
                  break

                case nsCtWs.CONSTANT.WS_REQCODE__STREAM_INFO_LISTING:
                case nsCtWs.CONSTANT.WS_REQCODE__STREAM_INFO_ITEM:
                  nsCtWs.dat_wsObss['msg:rr:stream'].next(msg);
                  break

                default:
                  console.warn(nsCtWs.ns, `no route for req_code '${msg._bh._brc}'`);
                  break
              }
              break

            case '_bdf':
              switch (msg._bh._bto) {
                case 'n':
                  nsCtWs.dat_wsObss['msg:df:node'].next(msg);
                  break

                case 'd':
                  nsCtWs.dat_wsObss['msg:df:device'].next(msg);
                  break

                case 's':
                  nsCtWs.dat_wsObss['msg:df:stream'].next(msg);
                  break

                default:
                  console.warn(nsCtWs.ns, `no route for recipient '${msg._bh._bto}'`);
                  break
              }
              break

            default:
              console.warn(nsCtWs.ns, `no route for msg_type '${msg._bh._bt}'`);
              break
          }
        })
      )
      .subscribe();
  };

  nsCtWs.fn_sendBasicRequest = function(reqCode, payload) {
    nsCtWs.dat_wsInstance.wsock.send(JSON.stringify({
      _bh: {_bt: '_brr2', _bid: '', _brc: reqCode},
      _bp: payload,
    }));
  };
};


const nsCtNode = {};
nsCtNode.TMPL = {};
nsCtNode.TMPL.NODE_ITEM = `
<div class="list-group-item">
  <div class="row">
    <div class="col-auto">
      <div class="status-indicator iderVisualStatus">
        <span class="status-indicator-circle"></span>
        <span class="status-indicator-circle"></span>
        <span class="status-indicator-circle"></span>
      </div>
    </div>
    <div class="col">
      <div class="text-truncate iderNodeId">AHIE</div>
    </div>
  </div>
  <div class="hr-text hr-text-center hr-text-spaceless mt-1 mb-2">Status</div>
  <div class="row">
    <div class="col-auto">
      <div class="faux-space"></div>
    </div>
    <div class="col">
      <div class="row">
        <div class="col-3">
          <div class="text-secondary">Status</div>
        </div>
        <div class="col-9">
          <div class="text-secondary iderTextualStatus"></div>
        </div>
      </div>
      <div class="row">
        <div class="col-3">
          <div class="text-secondary">Last seen at</div>
        </div>
        <div class="col-9">
          <div class="text-secondary iderLastSeen"></div>
        </div>
      </div>
    </div>
  </div>
  <div class="hr-text hr-text-center hr-text-spaceless mt-1 mb-2">Details</div>
  <div class="row">
    <div class="col-auto">
      <div class="faux-space"></div>
    </div>
    <div class="col">
      <div class="row">
        <div class="col-3">
          <div class="text-secondary">CPU</div>
          <div class="text-secondary">MEM</div>
        </div>
        <div class="col-9">
          <div class="text-secondary iderCpu"></div>
          <div class="text-secondary iderMem"></div>
        </div>
      </div>
    </div>
  </div>
</div>
`;
nsCtNode.init = function() {
  nsCtNode.el_nodeSummaryTitle = document.getElementById('nodeSummaryTitle');
  nsCtNode.el_nodeSummarySubtitle = document.getElementById('nodeSummarySubtitle');
  nsCtNode.el_nodeListing = document.getElementById('nodeListing');

  nsCtNode.dat_nodeInstances = new Map();

  nsCtNode.dat_wsSubs = {};

  nsCtNode.fn_linkWs = function() {
    nsCtNode.fn_renderNodeSummary(true);

    nsCtNode.dat_wsSubs['ws:open'] = nsCtWs.dat_wsObss['open']
      .subscribe({
        next: (v) => {
          nsCtWs.fn_sendBasicRequest(nsCtWs.CONSTANT.WS_REQCODE__NODE_INFO_LISTING, null);
        }
      });

    nsCtNode.dat_wsSubs['msg:rr:node'] = nsCtWs.dat_wsObss['msg:rr:node']
      .subscribe({
        next: (msg) => {
          switch (msg._bh._brc) {
            case nsCtWs.CONSTANT.WS_REQCODE__NODE_INFO_LISTING:
              nsCtNode.fn_populateNodeInstances(msg._bp.nil);
              nsCtNode.fn_renderNodeSummary(false);
              nsCtNode.fn_renderNodeListing();
              break
          }
        }
      });

    nsCtNode.dat_wsSubs['msg:df:node'] = nsCtWs.dat_wsObss['msg:df:node']
      .subscribe({
        next: (msg) => {
          switch (msg._bh._btopic) {
            case 'res':
              const nodeResourceInfo = msg._bp;

              const inst = nsCtNode.dat_nodeInstances.get(nodeResourceInfo.id);
              if (inst === undefined || inst === null) { break }

              inst.fn_renderResource(nodeResourceInfo);
              break

            case 'stat':
              const nodeStatusInfo = msg._bp.nsi;

              const inst2 = nsCtNode.dat_nodeInstances.get(nodeStatusInfo.id);
              if (inst2 === undefined || inst2 === null) { break }

              inst2.fn_queueStatusFeed(nodeStatusInfo);
              break
          }
        }
      });
  };

  nsCtNode.fn_populateNodeInstances = function(nodeInfoListing) {
    nsCtNode.dat_nodeInstances.clear();

    for (let i = 0; i < nodeInfoListing.length; i++) {
      const nodeInstance = nsCtNode.fn_newNodeInstance(i, nodeInfoListing[i]);

      nsCtNode.dat_nodeInstances.set(nodeInstance.id, nodeInstance);
    }
  };

  nsCtNode.fn_renderNodeSummary = function(isLoading) {
    if (isLoading) {
      nsCtNode.el_nodeSummaryTitle.innerHTML = '? nodes';
      nsCtNode.el_nodeSummarySubtitle.innerHTML = 'loading...';
      return
    }

    nsCtNode.el_nodeSummaryTitle.innerHTML = `${nsCtNode.dat_nodeInstances.size} node(s)`;
    nsCtNode.el_nodeSummarySubtitle.innerHTML = '_';
  };

  nsCtNode.fn_renderNodeListing = function() {
    nsCtNode.el_nodeListing.innerHTML = '';

    nsCtNode.dat_nodeInstances.forEach((inst) => {
      nsCtNode.el_nodeListing.appendChild(inst.el_root);

      inst.fn_renderContent();
    });
  };

  nsCtNode.fn_newNodeInstance = function(idx, nodeInfo) {
    const inst = {};
    inst.id = nodeInfo.id;
    inst.dat_item = nodeInfo;
    inst.dat_statusFeedSubj = new rxjs.Subject();

    const tempRoot = document.createElement('div');
    tempRoot.innerHTML = nsCtNode.TMPL.NODE_ITEM;
    inst.el_root = tempRoot.children[0];

    inst.el_visualStatus = inst.el_root.getElementsByClassName('iderVisualStatus')[0];
    inst.el_nodeId = inst.el_root.getElementsByClassName('iderNodeId')[0];
    inst.el_lastSeen = inst.el_root.getElementsByClassName('iderLastSeen')[0];
    inst.el_textualStatus = inst.el_root.getElementsByClassName('iderTextualStatus')[0];
    inst.el_cpu = inst.el_root.getElementsByClassName('iderCpu')[0];
    inst.el_mem = inst.el_root.getElementsByClassName('iderMem')[0];
    // inst.el_disk = inst.el_root.getElementsByClassName('iderDisk')[0];

    inst.dat_statusFeedSubj
      .pipe(
        rxjs.concatMap((v) => rxjs.of(v).pipe(rxjs.delay(333))),
      )
      .subscribe((nodeStatusInfo) => {
        inst.fn_renderStatus(nodeStatusInfo);
      });

    inst.fn_renderContent = function() {
      inst.el_nodeId.innerHTML = inst.dat_item.id;
      inst.el_lastSeen.innerHTML = _nsMain.formatTimestamp01(inst.dat_item.laa);

      inst.fn_renderStatus(inst.dat_item.nsi);
      inst.fn_renderResource(inst.dat_item.nri);
    };

    inst.fn_renderStatus = function(nodeStatusInfo) {
      inst.el_lastSeen.innerHTML = _nsMain.formatTimestamp01(nodeStatusInfo.laa);

      inst.el_textualStatus.innerHTML = nodeStatusInfo.txid;

      // 'status-indicator-animated'
      // ^^^ seems to trigger something that cause high CPU usage ?
      inst.el_visualStatus.classList.remove(
        'status-off',
        'status-green', 'status-yellow', 'status-red',
      );

      switch(nodeStatusInfo.vsid) {
        case 'o'  : inst.el_visualStatus.classList.add('status-off'); break
        case 'g_s': inst.el_visualStatus.classList.add('status-green'); break
        case 'g_b': inst.el_visualStatus.classList.add('status-green'); break
        case 'y_s': inst.el_visualStatus.classList.add('status-yellow'); break
        case 'y_b': inst.el_visualStatus.classList.add('status-yellow'); break
        case 'r_s': inst.el_visualStatus.classList.add('status-red'); break
        case 'r_b': inst.el_visualStatus.classList.add('status-red'); break
        default   : inst.el_visualStatus.classList.add('status-off'); break
      }
    };

    inst.fn_renderResource = function(nodeResourceInfo) {
      const cpuText = nodeResourceInfo.dcp;
      const memText = `${nodeResourceInfo.dmup} (${nodeResourceInfo.dmu} / ${nodeResourceInfo.dmt})`;

      inst.el_cpu.innerHTML = cpuText;
      inst.el_mem.innerHTML = memText;
    };

    inst.fn_queueStatusFeed = function(nodeStatusInfo) {
      inst.dat_statusFeedSubj.next(nodeStatusInfo);
    };

    return inst;
  };
};

const nsCtDevice = {};
nsCtDevice.TMPL = {};
nsCtDevice.TMPL.DEVICE_ITEM = `
<div class="list-group-item">
  <div class="row">
    <div class="col-auto">
      <div class="status-indicator iderVisualStatus">
        <span class="status-indicator-circle"></span>
        <span class="status-indicator-circle"></span>
        <span class="status-indicator-circle"></span>
      </div>
    </div>
    <div class="col">
      <div class="text-truncate iderDeviceCode"></div>
      <div class="text-secondary iderDeviceName"></div>
    </div>
  </div>
  <div class="hr-text hr-text-center hr-text-spaceless mt-1 mb-2">Status</div>
  <div class="row">
    <div class="col-auto">
      <div class="faux-space"></div>
    </div>
    <div class="col">
      <div class="row">
        <div class="col-3">
          <div class="text-secondary">Status</div>
        </div>
        <div class="col-9">
          <div class="text-secondary iderTextualStatus"></div>
        </div>
      </div>
      <div class="row">
        <div class="col-3">
          <div class="text-secondary">Last seen at</div>
        </div>
        <div class="col-9">
          <div class="text-secondary iderLastSeen"></div>
        </div>
      </div>
    </div>
  </div>
</div>
`;
nsCtDevice.init = function() {
  nsCtDevice.el_deviceSummaryTitle = document.getElementById('deviceSummaryTitle');
  nsCtDevice.el_deviceSummarySubtitle = document.getElementById('deviceSummarySubtitle');
  nsCtDevice.el_deviceListing = document.getElementById('deviceListing');

  nsCtDevice.dat_deviceInstances = new Map();

  nsCtDevice.dat_wsSubs = {};

  nsCtDevice.fn_linkWs = function() {
    nsCtDevice.fn_renderDeviceSummary(true);

    nsCtDevice.dat_wsSubs['ws:open'] = nsCtWs.dat_wsObss['open']
      .subscribe({
        next: (v) => {
          nsCtWs.fn_sendBasicRequest(nsCtWs.CONSTANT.WS_REQCODE__DEVICE_INFO_LISTING, null);
        },
      });

    nsCtDevice.dat_wsSubs['msg:rr:device'] = nsCtWs.dat_wsObss['msg:rr:device']
      .subscribe({
        next: (msg) => {
          switch (msg._bh._brc) {
            case nsCtWs.CONSTANT.WS_REQCODE__DEVICE_INFO_LISTING:
              nsCtDevice.fn_populateDeviceInstances(msg._bp.dil);
              nsCtDevice.fn_renderDeviceSummary(false);
              nsCtDevice.fn_renderDeviceListing();
              break
          }
        }
      });

    nsCtDevice.dat_wsSubs['msg:df:device'] = nsCtWs.dat_wsObss['msg:df:device']
      .subscribe({
        next: (msg) => {
          switch (msg._bh._btopic) {
            case 'stat':
              const deviceStatusInfo = msg._bp.dsi;

              const inst = nsCtDevice.dat_deviceInstances.get(deviceStatusInfo.id);
              if (inst === undefined || inst === null) { break }

              inst.fn_queueStatusFeed(deviceStatusInfo);
              break
          }
        }
      });
  };

  nsCtDevice.fn_populateDeviceInstances = function(deviceInfoListing) {
    nsCtDevice.dat_deviceInstances.clear();

    for (let i = 0; i < deviceInfoListing.length; i++) {
      const deviceInstance = nsCtDevice.fn_newDeviceInstance(i, deviceInfoListing[i]);

      nsCtDevice.dat_deviceInstances.set(deviceInstance.id, deviceInstance);
    }
  };

  nsCtDevice.fn_renderDeviceSummary = function(isLoading) {
    if (isLoading) {
      nsCtDevice.el_deviceSummaryTitle.innerHTML = '? devices';
      nsCtDevice.el_deviceSummarySubtitle.innerHTML = 'loading...';
      return
    }

    nsCtDevice.el_deviceSummaryTitle.innerHTML = `${nsCtDevice.dat_deviceInstances.size} device(s)`;
    nsCtDevice.el_deviceSummarySubtitle.innerHTML = '_';
  };

  nsCtDevice.fn_renderDeviceListing = function() {
    nsCtDevice.el_deviceListing.innerHTML = '';

    nsCtDevice.dat_deviceInstances.forEach((inst) => {
      nsCtDevice.el_deviceListing.appendChild(inst.el_root);

      inst.fn_renderContent();
    });
  };

  nsCtDevice.fn_newDeviceInstance = function(idx, deviceInfo) {
    const inst = {};
    inst.id = deviceInfo.id;
    inst.dat_item = deviceInfo;
    inst.dat_statusFeedSubj = new rxjs.Subject();

    const tempRoot = document.createElement('div');
    tempRoot.innerHTML = nsCtDevice.TMPL.DEVICE_ITEM;
    inst.el_root = tempRoot.children[0];

    inst.el_visualStatus = inst.el_root.getElementsByClassName('iderVisualStatus')[0];
    inst.el_deviceCode = inst.el_root.getElementsByClassName('iderDeviceCode')[0];
    inst.el_deviceName = inst.el_root.getElementsByClassName('iderDeviceName')[0];
    inst.el_lastSeen = inst.el_root.getElementsByClassName('iderLastSeen')[0];
    inst.el_textualStatus = inst.el_root.getElementsByClassName('iderTextualStatus')[0];

    inst.dat_statusFeedSubj
      .pipe(
        rxjs.concatMap((v) => rxjs.of(v).pipe(rxjs.delay(333))),
      )
      .subscribe((deviceStatusInfo) => {
        inst.fn_renderStatus(deviceStatusInfo);
      });

    inst.fn_renderContent = function() {
      inst.el_deviceCode.innerHTML = inst.dat_item.code;
      inst.el_deviceName.innerHTML = inst.dat_item.name;

      inst.fn_renderStatus(inst.dat_item.dsi);
    };

    inst.fn_renderStatus = function(deviceStatusInfo) {
      inst.el_lastSeen.innerHTML = _nsMain.formatTimestamp01(deviceStatusInfo.laa);

      inst.el_textualStatus.innerHTML = deviceStatusInfo.txid;

      inst.el_visualStatus.classList.remove(
        'status-off',
        'status-green', 'status-yellow', 'status-red',
      );

      switch(deviceStatusInfo.vsid) {
        case 'o'  : inst.el_visualStatus.classList.add('status-off'); break
        case 'g_s': inst.el_visualStatus.classList.add('status-green'); break
        case 'g_b': inst.el_visualStatus.classList.add('status-green'); break
        case 'y_s': inst.el_visualStatus.classList.add('status-yellow'); break
        case 'y_b': inst.el_visualStatus.classList.add('status-yellow'); break
        case 'r_s': inst.el_visualStatus.classList.add('status-red'); break
        case 'r_b': inst.el_visualStatus.classList.add('status-red'); break
        default   : inst.el_visualStatus.classList.add('status-off'); break
      }
    };

    inst.fn_queueStatusFeed = function(devicemStatusInfo) {
      inst.dat_statusFeedSubj.next(devicemStatusInfo);
    };

    return inst;
  };
};

const nsCtStream = {};
nsCtStream.TMPL = {};
nsCtStream.TMPL.STREAM_ITEM = `
<div class="list-group-item">
  <div class="row">
    <div class="col-auto">
      <div class="status-indicator iderVisualStatus">
        <span class="status-indicator-circle"></span>
        <span class="status-indicator-circle"></span>
        <span class="status-indicator-circle"></span>
      </div>
    </div>
    <div class="col">
      <div class="text-truncate iderStreamCode"></div>
    </div>
  </div>
  <div class="hr-text hr-text-center hr-text-spaceless mt-1 mb-2">Status</div>
  <div class="row">
    <div class="col-auto">
      <div class="faux-space"></div>
    </div>
    <div class="col">
      <div class="row">
        <div class="col-3">
          <div class="text-secondary">Status</div>
        </div>
        <div class="col-9">
          <div class="text-secondary iderTextualStatus"></div>
        </div>
      </div>
      <div class="row">
        <div class="col-3">
          <div class="text-secondary">Last seen at</div>
        </div>
        <div class="col-9">
          <div class="text-secondary iderLastSeen"></div>
        </div>
      </div>
    </div>
  </div>
</div>
`;
nsCtStream.init = function() {
  nsCtStream.el_streamSummaryTitle = document.getElementById('streamSummaryTitle');
  nsCtStream.el_streamSummarySubtitle = document.getElementById('streamSummarySubtitle');
  nsCtStream.el_streamListing = document.getElementById('streamListing');

  nsCtStream.dat_streamInstances = new Map();

  nsCtStream.dat_wsSubs = {};

  nsCtStream.fn_linkWs = function() {
    nsCtStream.fn_renderStreamSummary(true);

    nsCtStream.dat_wsSubs['ws:open'] = nsCtWs.dat_wsObss['open']
      .subscribe({
        next: (v) => {
          nsCtWs.fn_sendBasicRequest(nsCtWs.CONSTANT.WS_REQCODE__STREAM_INFO_LISTING, null);
        },
      });

    nsCtStream.dat_wsSubs['msg:rr:stream'] = nsCtWs.dat_wsObss['msg:rr:stream']
      .subscribe({
        next: (msg) => {
          switch (msg._bh._brc) {
            case nsCtWs.CONSTANT.WS_REQCODE__STREAM_INFO_LISTING:
              nsCtStream.fn_populateStreamInstances(msg._bp.sil);
              nsCtStream.fn_renderStreamSummary(false);
              nsCtStream.fn_renderStreamListing();
              break

            case nsCtWs.CONSTANT.WS_REQCODE__STREAM_INFO_ITEM:
              console.log(msg._bp.sii);
              break
          }
        }
      });

    nsCtStream.dat_wsSubs['msg:df:stream'] = nsCtWs.dat_wsObss['msg:df:stream']
      .subscribe({
        next: (msg) => {
          switch (msg._bh._btopic) {
            case 'stat':
              const streamStatusInfo = msg._bp.ssi;

              const inst = nsCtStream.dat_streamInstances.get(streamStatusInfo.id);
              if (inst === undefined || inst === null) { break }

              inst.fn_queueStatusFeed(streamStatusInfo);
              break
          }
        }
      });
  };

  nsCtStream.fn_populateStreamInstances = function(streamInfoListing) {
    nsCtStream.dat_streamInstances.clear();

    for (let i = 0; i < streamInfoListing.length; i++) {
      const streamInstance = nsCtStream.fn_newStreamInstance(i, streamInfoListing[i]);

      nsCtStream.dat_streamInstances.set(streamInstance.id, streamInstance);
    }
  };

  nsCtStream.fn_renderStreamSummary = function(isLoading) {
    if (isLoading) {
      nsCtStream.el_streamSummaryTitle.innerHTML = '? streams';
      nsCtStream.el_streamSummarySubtitle.innerHTML = 'loading...';
      return
    }

    nsCtStream.el_streamSummaryTitle.innerHTML = `${nsCtStream.dat_streamInstances.size} stream(s)`;
    nsCtStream.el_streamSummarySubtitle.innerHTML = '_';
  };

  nsCtStream.fn_renderStreamListing = function() {
    nsCtStream.el_streamListing.innerHTML = '';

    nsCtStream.dat_streamInstances.forEach((inst) => {
      nsCtStream.el_streamListing.appendChild(inst.el_root);

      inst.fn_renderContent();
    });
  };

  nsCtStream.fn_newStreamInstance = function(idx, streamInfo) {
    const inst = {};
    inst.id = streamInfo.id;
    inst.dat_item = streamInfo;
    inst.dat_statusFeedSubj = new rxjs.Subject();

    const tempRoot = document.createElement('div');
    tempRoot.innerHTML = nsCtStream.TMPL.STREAM_ITEM;
    inst.el_root = tempRoot.children[0];

    inst.el_visualStatus = inst.el_root.getElementsByClassName('iderVisualStatus')[0];
    inst.el_streamCode = inst.el_root.getElementsByClassName('iderStreamCode')[0];
    inst.el_lastSeen = inst.el_root.getElementsByClassName('iderLastSeen')[0];
    inst.el_textualStatus = inst.el_root.getElementsByClassName('iderTextualStatus')[0];

    inst.dat_statusFeedSubj
      .pipe(
        rxjs.concatMap((v) => rxjs.of(v).pipe(rxjs.delay(333))),
      )
      .subscribe((streamStatusInfo) => {
        inst.fn_renderStatus(streamStatusInfo);
      });

    inst.fn_renderContent = function() {
      inst.el_streamCode.innerHTML = inst.dat_item.code;

      inst.fn_renderStatus(inst.dat_item.ssi);
    };

    inst.fn_renderStatus = function(streamStatusInfo) {
      inst.el_lastSeen.innerHTML = _nsMain.formatTimestamp01(streamStatusInfo.laa);

      inst.el_textualStatus.innerHTML = streamStatusInfo.txid;

      inst.el_visualStatus.classList.remove(
        'status-off',
        'status-green', 'status-yellow', 'status-red',
      );

      switch(streamStatusInfo.vsid) {
        case 'o'  : inst.el_visualStatus.classList.add('status-off'); break
        case 'g_s': inst.el_visualStatus.classList.add('status-green'); break
        case 'g_b': inst.el_visualStatus.classList.add('status-green'); break
        case 'y_s': inst.el_visualStatus.classList.add('status-yellow'); break
        case 'y_b': inst.el_visualStatus.classList.add('status-yellow'); break
        case 'r_s': inst.el_visualStatus.classList.add('status-red'); break
        case 'r_b': inst.el_visualStatus.classList.add('status-red'); break
        default   : inst.el_visualStatus.classList.add('status-off'); break
      }
    };

    inst.fn_queueStatusFeed = function(streamStatusInfo) {
      inst.dat_statusFeedSubj.next(streamStatusInfo);
    };

    return inst;
  };
};

window.addEventListener("DOMContentLoaded", function() {
  _nsWs.init();

  nsCtWs.init();
  nsCtWs.fn_setupWs();

  nsCtNode.init();
  nsCtNode.fn_linkWs();

  nsCtDevice.init();
  nsCtDevice.fn_linkWs();

  nsCtStream.init();
  nsCtStream.fn_linkWs();
});
