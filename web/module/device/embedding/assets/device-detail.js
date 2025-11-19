'use strict';

// nsContent
const nsCt = {};
nsCt.init = function() {
  nsCt.el_btnSubmitDelete = document.getElementById('btnSubmitDelete');
  nsCt.el_togglePassword = document.getElementById('togglePassword');
  nsCt.el_passwordInput = document.getElementById('passwordInput');

  nsCt.fn_doSubmitDelete = function() {
    nsCt.el_btnSubmitDelete.click();
  };

  nsCt.el_togglePassword.addEventListener('click', function(ev) {
    const iconDontShowPassword = `<svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-eye-off"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M10.585 10.587a2 2 0 0 0 2.829 2.828" /><path d="M16.681 16.673a8.717 8.717 0 0 1 -4.681 1.327c-3.6 0 -6.6 -2 -9 -6c1.272 -2.12 2.712 -3.678 4.32 -4.674m2.86 -1.146a9.055 9.055 0 0 1 1.82 -.18c3.6 0 6.6 2 9 6c-.666 1.11 -1.379 2.067 -2.138 2.87" /><path d="M3 3l18 18" /></svg>`;
    const iconShowPassword = `<svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-eye"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M10 12a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" /><path d="M21 12c-2.4 4 -5.4 6 -9 6c-3.6 0 -6.6 -2 -9 -6c2.4 -4 5.4 -6 9 -6c3.6 0 6.6 2 9 6" /></svg>`;

    let typePassword = nsCt.el_passwordInput.getAttribute('type');

    if (typePassword == 'password') {
      nsCt.el_passwordInput.setAttribute('type', 'text');
      nsCt.el_togglePassword.innerHTML = iconDontShowPassword;
    } else {
      nsCt.el_passwordInput.setAttribute('type', 'password');
      nsCt.el_togglePassword.innerHTML = iconShowPassword;
    }
  });
};

// websocket
const nsCtWs = { ns: 'nsCtWs' };
nsCtWs.CONSTANT = {
  WS_REQCODE_DEVICE_SNAPSHOT: '/device/snapshot',
  WS_REQCODE_TEMP_ERROR_DETAILS: '/device/temp-error-details',

  WS_REQCODE_DEVICE_RELOAD: '/device/reload',
};
nsCtWs.init = function() {
  nsCtWs.dat_wsInstance = null;
  nsCtWs.dat_wsObss = {};
  nsCtWs.dat_wsSubs = {};

  nsCtWs.fn_wsSetup = function() {
    const wsId = 'device-detail';
    const wsUri = '/device/device/detail/ws' + window.location.search;

    nsCtWs.dat_wsInstance = _nsWs.newWsInstance(wsId, wsUri);

    nsCtWs.dat_wsObss['open'] = nsCtWs.dat_wsInstance.onopenSubj.asObservable();
    nsCtWs.dat_wsObss['msg'] = nsCtWs.dat_wsInstance.onmessageSubj.asObservable().pipe(_nsWs.plainJsonParser);

    nsCtWs.dat_wsObss['msg:ev'] = new rxjs.ReplaySubject(3);
    nsCtWs.dat_wsObss['msg:rr'] = new rxjs.Subject();

    nsCtWs.dat_wsSubs['msg:router'] = nsCtWs.dat_wsObss['msg']
      .pipe(
        rxjs.tap((msg) => {
          switch(msg._bt) {
            case '_bev': nsCtWs.dat_wsObss['msg:ev'].next(msg._bp); break;
            case '_brr': nsCtWs.dat_wsObss['msg:rr'].next(msg._bp); break;

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

  nsCtWs.ws_reqDeviceSnapshot = function() {
    nsCtWs.fn_sendMessage({
      _bid: '',
      _brc: nsCtWs.CONSTANT.WS_REQCODE_DEVICE_SNAPSHOT,
    });
  };
  nsCtWs.ws_reqTempErrorDetails = function() {
    nsCtWs.fn_sendMessage({
      _bid: '',
      _brc: nsCtWs.CONSTANT.WS_REQCODE_TEMP_ERROR_DETAILS,
    });
  };

  nsCtWs.ws_reqDeviceReload = function() {
    nsCtWs.fn_sendMessage({
      _bid: '',
      _brc: nsCtWs.CONSTANT.WS_REQCODE_DEVICE_RELOAD,
    });
  };
};

// additional content partition - panel info
const nsCtPi = { ns: 'nsCtPi' };
nsCtPi.CONSTANT = {
  // copied as is from source
  DEVICE_LIVE_STATE__NEW: 'dls:new',
  DEVICE_LIVE_STATE__INACTIVE: 'dls:inactive',
  DEVICE_LIVE_STATE__INIT_BEGIN: 'dls:init:begin',
  DEVICE_LIVE_STATE__INIT_FAIL: 'dls:init:fail',
  DEVICE_LIVE_STATE__INIT_OK: 'dls:init:ok',
  DEVICE_LIVE_STATE__DISCONNECTED: 'dls:disconnected',
  DEVICE_LIVE_STATE__RELOAD_BEGIN: 'dls:reload:begin',
  DEVICE_LIVE_STATE__RELOAD_FAIL: 'dls:reload:fail',
  DEVICE_LIVE_STATE__RELOAD_OK: 'dls:reload:ok',
  DEVICE_LIVE_STATE__DESTROY: 'dls:destroy',
};
nsCtPi.init = function() {
  nsCtPi.tmpl_opcapInit = `
<div class="text-muted">
  <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="currentColor"  class="icon icon-tabler icons-tabler-filled icon-tabler-square-minus"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M19 2a3 3 0 0 1 3 3v14a3 3 0 0 1 -3 3h-14a3 3 0 0 1 -3 -3v-14a3 3 0 0 1 3 -3zm-4 9h-6l-.117 .007a1 1 0 0 0 .117 1.993h6l.117 -.007a1 1 0 0 0 -.117 -1.993z" /></svg>
<div>`;
  nsCtPi.tmpl_opcapUnknown = `
<div class="text-warning" title="Unknown" data-bs-toggle="tooltip" data-bs-placement="top">
  <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="currentColor"  class="icon icon-tabler icons-tabler-filled icon-tabler-help-square"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M19 2a3 3 0 0 1 2.995 2.824l.005 .176v14a3 3 0 0 1 -2.824 2.995l-.176 .005h-14a3 3 0 0 1 -2.995 -2.824l-.005 -.176v-14a3 3 0 0 1 2.824 -2.995l.176 -.005h14zm-7 13a1 1 0 0 0 -.993 .883l-.007 .117l.007 .127a1 1 0 0 0 1.986 0l.007 -.117l-.007 -.127a1 1 0 0 0 -.993 -.883zm1.368 -6.673a2.98 2.98 0 0 0 -3.631 .728a1 1 0 0 0 1.44 1.383l.171 -.18a.98 .98 0 0 1 1.11 -.15a1 1 0 0 1 -.34 1.886l-.232 .012a1 1 0 0 0 .111 1.994a3 3 0 0 0 1.371 -5.673z" /></svg>
<div>`;
  nsCtPi.tmpl_opcapCanDo = `
<div class="text-green">
  <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="currentColor"  class="icon icon-tabler icons-tabler-filled icon-tabler-square-check"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M18.333 2c1.96 0 3.56 1.537 3.662 3.472l.005 .195v12.666c0 1.96 -1.537 3.56 -3.472 3.662l-.195 .005h-12.666a3.667 3.667 0 0 1 -3.662 -3.472l-.005 -.195v-12.666c0 -1.96 1.537 -3.56 3.472 -3.662l.195 -.005h12.666zm-2.626 7.293a1 1 0 0 0 -1.414 0l-3.293 3.292l-1.293 -1.292l-.094 -.083a1 1 0 0 0 -1.32 1.497l2 2l.094 .083a1 1 0 0 0 1.32 -.083l4 -4l.083 -.094a1 1 0 0 0 -.083 -1.32z" /></svg>
<div>`;
  nsCtPi.tmpl_opcapCannotDo = `
<div class="text-red">
  <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="currentColor"  class="icon icon-tabler icons-tabler-filled icon-tabler-square-x"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M19 2h-14a3 3 0 0 0 -3 3v14a3 3 0 0 0 3 3h14a3 3 0 0 0 3 -3v-14a3 3 0 0 0 -3 -3zm-9.387 6.21l.094 .083l2.293 2.292l2.293 -2.292a1 1 0 0 1 1.497 1.32l-.083 .094l-2.292 2.293l2.292 2.293a1 1 0 0 1 -1.32 1.497l-.094 -.083l-2.293 -2.292l-2.293 2.292a1 1 0 0 1 -1.497 -1.32l.083 -.094l2.292 -2.293l-2.292 -2.293a1 1 0 0 1 1.32 -1.497z" /></svg>
<div>`;

  nsCtPi.tmpl_reloadDevice = `
<svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-reload"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M19.933 13.041a8 8 0 1 1 -9.925 -8.788c3.899 -1 7.935 1.007 9.425 4.747" /><path d="M20 4v5h-5" /></svg>
Reload device`;
  nsCtPi.tmpl_reloadDevicePending = `<div class="spinner-border spinner-border-sm" role="status"></div>`;

  nsCtPi.el_statusText = document.getElementById('statusText');
  nsCtPi.el_connectionText = document.getElementById('connectionText');
  nsCtPi.el_opcapState = document.getElementById('opcapState');
  nsCtPi.el_opcapReadDeviceInfo = document.getElementById('opcapReadDeviceInfo');
  nsCtPi.el_opcapReadRtspStream = document.getElementById('opcapReadRtspStream');
  nsCtPi.el_opcapReadStreamInfo = document.getElementById('opcapReadStreamInfo');
  nsCtPi.el_opcapReadAnalogInputChannels = document.getElementById('opcapReadAnalogInputChannels');
  nsCtPi.el_opcapReadDigitalInputChannels = document.getElementById('opcapReadDigitalInputChannels');
  nsCtPi.el_opcapErrDetailsTitle = document.getElementById('opcapErrDetailsTitle');
  nsCtPi.el_opcapErrDetailsBody = document.getElementById('opcapErrDetailsBody');
  nsCtPi.el_hwBrand = document.getElementById('hwBrand');
  nsCtPi.el_hwName = document.getElementById('hwName');
  nsCtPi.el_hwModel = document.getElementById('hwModel');
  nsCtPi.el_hwType = document.getElementById('hwType');
  nsCtPi.el_hwAnalogChannelsCount = document.getElementById('hwAnalogChannelsCount');
  nsCtPi.el_hwDigitalChannelsCount = document.getElementById('hwDigitalChannelsCount');
  nsCtPi.el_reloadWarning = document.getElementById('reloadWarning');
  nsCtPi.el_reloadDevice = document.getElementById('reloadDevice');

  nsCtPi.dat_wsSubs = {};

  nsCtPi.fn_wsSetup = function() {
    nsCtPi.dat_wsSubs['ws:open'] = nsCtWs.dat_wsObss['open']
      .subscribe({
        next: (v) => {
          nsCtWs.ws_reqDeviceSnapshot();
        },
      });

    nsCtPi.dat_wsSubs['ws:msg:ev'] = nsCtWs.dat_wsObss['msg:ev']
      .pipe(
        rxjs.concatMap((p_ev) => rxjs.of(p_ev).pipe(rxjs.delay(333))),
      )
      .subscribe({
        next: (p_ev) => {
          nsCtWs.ws_reqDeviceSnapshot();
          // switch(p_ev.ev_code) { // TODO
          //   default:
          //     break;
          // }
        },
      });

    nsCtPi.dat_wsSubs['ws:msg:rr'] = nsCtWs.dat_wsObss['msg:rr']
      .subscribe({
        next: (p_rep) => {
          switch(p_rep._brc) {
            case nsCtWs.CONSTANT.WS_REQCODE_DEVICE_SNAPSHOT:
              const ds = p_rep.device_snapshot;
              nsCtPi.fn_renderInfoPanel(ds);
              nsCtWs.ws_reqTempErrorDetails()
              break;

            case nsCtWs.CONSTANT.WS_REQCODE_TEMP_ERROR_DETAILS:
              const tempErrorDetails = p_rep.temp_error_details;
              nsCtPi.fn_updateTempErrorDetails(tempErrorDetails);
              break;

            default:
              console.warn(`${nsCtPi.ns}`, `no handler for req_code '${p_rep._brc}'`);
              break;
          }
        },
      });
  };

  nsCtPi.fn_renderInfoPanel = function(ds) {
    let statusText = '-';
    let connectionText = '-';

    let opcapStateText = '-';
    let opcapReadDeviceInfoText = nsCtPi.tmpl_opcapInit;
    let opcapReadRtspStreamText = nsCtPi.tmpl_opcapInit;
    let opcapReadStreamInfoText = nsCtPi.tmpl_opcapInit;
    let opcapReadAnalogInputChannelsText = nsCtPi.tmpl_opcapInit;
    let opcapReadDigitalInputChannelsText = nsCtPi.tmpl_opcapInit;

    let opcapErrDetailsTitleText = 'Detail messages';
    let opcapErrDetailsBodyText = '-';

    let hwBrandText = '-';
    let hwNameText = '-';
    let hwModelText = '-';
    let hwTypeText = '-';
    let hwAnalogChannelsCountText = '-';
    let hwDigitalChannelsCountText = '-';

    let showReloadWarning = false;
    let reloadDeviceText = nsCtPi.tmpl_reloadDevice;
    let allowReload = false;

    if (ds === undefined || ds === null) {

    } else {
      switch (ds.live.state) {
        case nsCtPi.CONSTANT.DEVICE_LIVE_STATE__NEW:
          statusText = 'waiting for initialization';
          break;

        case nsCtPi.CONSTANT.DEVICE_LIVE_STATE__INACTIVE:
          statusText = 'INACTIVE';
          break;

        case nsCtPi.CONSTANT.DEVICE_LIVE_STATE__INIT_BEGIN:
          statusText = 'initializing<span class="animated-dots"></span>';
          break;

        case nsCtPi.CONSTANT.DEVICE_LIVE_STATE__INIT_FAIL:
          statusText = `init FAIL: <code style="display: block;">${ds.live.conn_state_msg}</code>`;
          allowReload = true;
          break;

        case nsCtPi.CONSTANT.DEVICE_LIVE_STATE__INIT_OK:
          statusText = 'init OK';
          showReloadWarning = true;
          allowReload = true;
          break;

        case nsCtPi.CONSTANT.DEVICE_LIVE_STATE__DISCONNECTED:
          statusText = '<span class="badge bg-danger me-1"></span>DISCONNECTED';
          allowReload = true;
          break;

        case nsCtPi.CONSTANT.DEVICE_LIVE_STATE__RELOAD_BEGIN:
          statusText = 'reloading<span class="animated-dots"></span>';
          reloadDeviceText = nsCtPi.tmpl_reloadDevicePending;
          break;

        case nsCtPi.CONSTANT.DEVICE_LIVE_STATE__RELOAD_FAIL:
          statusText = `reload FAIL: <code style="display: block;">${ds.live.conn_state_msg}</code>`;
          allowReload = true;
          break;

        case nsCtPi.CONSTANT.DEVICE_LIVE_STATE__RELOAD_OK:
          statusText = 'reload OK';
          showReloadWarning = true;
          allowReload = true;
          break;

        case nsCtPi.CONSTANT.DEVICE_LIVE_STATE__DESTROY:
          statusText = 'DESTROY';
          break;

        default:
          console.warn(`${nsCtPi.ns}`, `no handler for DLS '${ds.live.state}'`);
          break;
      }

      // TODO: after rearrange conn stuffs
      switch (ds.live.conn_state) {
        case 'never':
          connectionText = '<span class="badge bg-secondary me-1"></span>Never connected';
          break;

        case 'alive':
          connectionText = '<span class="badge bg-success me-1"></span>OK';
          break;

        case 'lost':
          const fmt1 = nsCtPi.fn_formatTimestamp(ds.live.last_seen);
          connectionText = `LOST / DISCONNECTED, last seen at ${fmt1}`;
          break;
      }

      // TODO: switch case with human wording...
      opcapStateText = ds.op_cap.state;

      // TODO: temp until backend can test rtsp
      opcapReadRtspStreamText = nsCtPi.tmpl_opcapUnknown;

      opcapReadDeviceInfoText = ds.op_cap.can_read_device_info ? nsCtPi.tmpl_opcapCanDo : nsCtPi.tmpl_opcapCannotDo;
      opcapReadStreamInfoText = ds.op_cap.can_read_stream_info ? nsCtPi.tmpl_opcapCanDo : nsCtPi.tmpl_opcapCannotDo;
      opcapReadAnalogInputChannelsText = ds.op_cap.can_read_analog_input_channels ? nsCtPi.tmpl_opcapCanDo : nsCtPi.tmpl_opcapCannotDo;
      opcapReadDigitalInputChannelsText = ds.op_cap.can_read_digital_input_channels ? nsCtPi.tmpl_opcapCanDo : nsCtPi.tmpl_opcapCannotDo;

      hwBrandText = ds.persistence.brand;
      hwNameText = ds.hardware.device_name;
      hwModelText = ds.hardware.model;
      hwTypeText = ds.hardware.device_type;

      hwAnalogChannelsCountText = Array.isArray(ds.hardware.analog_channels) ? ds.hardware.analog_channels.length : '-';
      hwDigitalChannelsCountText = Array.isArray(ds.hardware.digital_channels) ? ds.hardware.digital_channels.length : '-';
    }


    nsCtPi.el_statusText.innerHTML = statusText;
    nsCtPi.el_connectionText.innerHTML = connectionText;

    nsCtPi.el_opcapState.innerHTML = opcapStateText;
    nsCtPi.el_opcapReadDeviceInfo.innerHTML = opcapReadDeviceInfoText;
    nsCtPi.el_opcapReadRtspStream.innerHTML = opcapReadRtspStreamText;
    nsCtPi.el_opcapReadStreamInfo.innerHTML = opcapReadStreamInfoText;
    nsCtPi.el_opcapReadAnalogInputChannels.innerHTML = opcapReadAnalogInputChannelsText;
    nsCtPi.el_opcapReadDigitalInputChannels.innerHTML = opcapReadDigitalInputChannelsText;

    nsCtPi.el_opcapErrDetailsTitle.innerHTML = opcapErrDetailsTitleText;
    nsCtPi.el_opcapErrDetailsBody.innerHTML = opcapErrDetailsBodyText;

    nsCtPi.el_hwBrand.innerHTML = hwBrandText;
    nsCtPi.el_hwName.innerHTML = hwNameText;
    nsCtPi.el_hwModel.innerHTML = hwModelText;
    nsCtPi.el_hwType.innerHTML = hwTypeText;
    nsCtPi.el_hwAnalogChannelsCount.innerHTML = hwAnalogChannelsCountText;
    nsCtPi.el_hwDigitalChannelsCount.innerHTML = hwDigitalChannelsCountText;

    nsCtPi.el_reloadWarning.style.display = showReloadWarning ? '' : 'none';
    nsCtPi.el_reloadDevice.innerHTML = reloadDeviceText;
    nsCtPi.el_reloadDevice.disabled = !allowReload;
  };

  nsCtPi.fn_updateTempErrorDetails = function(tempErrorDetails) {
    let itemCount = 0;
    let tempBody = '';

    console.log(tempErrorDetails); // temp, for easier debugging
    
    for (let pn in tempErrorDetails) {
      itemCount += 1;
      tempBody += `<div class="row mb-1">
  <div class="col-4">${pn}</div>
  <code class="col-8">${tempErrorDetails[pn]}</code>
</div>`;
    }

    nsCtPi.el_opcapErrDetailsTitle.innerHTML = `Detail messages (${itemCount} items)`;
    nsCtPi.el_opcapErrDetailsBody.innerHTML = tempBody;
  };

  nsCtPi.fn_formatTimestamp = function(ts) {
    // temp, make date util stuffs
    // https://moment.github.io/luxon/#/formatting?id=table-of-tokens
    return luxon.DateTime.fromISO(ts).toFormat('dd LLLL yyyy, HH:mm:ss ZZZZ');
  }

  nsCtPi.el_reloadDevice.addEventListener('click', function() {
    nsCtWs.ws_reqDeviceReload();
  });

  nsCtPi.fn_renderInfoPanel(null);
};

document.addEventListener('DOMContentLoaded', function() {
  _nsWs.init();

  nsCt.init();
  nsCtWs.init();
  nsCtPi.init();

  if (nsCD && nsCD.isEditMode) {
    nsCtWs.fn_wsSetup();
    nsCtPi.fn_wsSetup();
  }

  nsMain.focusOn('code');
});
