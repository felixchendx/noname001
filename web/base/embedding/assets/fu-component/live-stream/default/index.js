'use strict';

// requires: hls.js

// TODO: loading bar styling and message structure + styling
// TODO: hls loader policy + error handling 

// TODO: _nsLogger in it's own package
//       either ^^^ or immediate tampering of window.console object

const _nsDLS = {}; // aw shiet...
_nsDLS.debug = false;
_nsDLS.instances = {}; // map[data-fu-id] = dlsInstance

_nsDLS.logger = {
  debug: function(id, msg) { console.debug('DLS', `${id}: ${msg}`); },
  info : function(id, msg) { console.info('DLS', `${id}: ${msg}`); },
  log  : function(id, msg) { console.log('DLS', `${id}: ${msg}`); },
  warn : function(id, msg) { console.warn('DLS', `${id}: ${msg}`); },
  error: function(id, msg) { console.error('DLS', `${id}: ${msg}`); },

  // fatal: function(id, msg) { console.error('DLS', `${id} [FATAL]: ${msg}`); },
};

_nsDLS.template = `
<div class="dls-loading">
  <div class="progress">
    <div class="progress-bar progress-bar-indeterminate"></div>
  </div>
</div>
<div class="dls-message">
</div>
<video class=dls-video></video>
<canvas class="dls-overlay"></canvas>
`;

_nsDLS.defaultLiveStreamConf = function() {
  // var lsConf = {
  //   autoStartLoad: true,
  //   startPosition: -1,
  //   debug: false,
  //   capLevelOnFPSDrop: false,
  //   capLevelToPlayerSize: false,
  //   defaultAudioCodec: undefined,
  //   initialLiveManifestSize: 1,
  //   maxBufferLength: 30,
  //   maxMaxBufferLength: 600,
  //   backBufferLength: Infinity,
  //   frontBufferFlushThreshold: Infinity,
  //   maxBufferSize: 60 * 1000 * 1000,
  //   maxBufferHole: 0.1,
  //   highBufferWatchdogPeriod: 2,
  //   nudgeOffset: 0.1,
  //   nudgeMaxRetry: 3,
  //   maxFragLookUpTolerance: 0.25,
  //   liveSyncDurationCount: 3,
  //   liveSyncOnStallIncrease: 1,
  //   liveMaxLatencyDurationCount: Infinity,
  //   liveDurationInfinity: false,
  //   preferManagedMediaSource: false,
  //   enableWorker: true,
  //   enableSoftwareAES: true,
  //   fragLoadPolicy: {
  //     default: {
  //       maxTimeToFirstByteMs: 9000,
  //       maxLoadTimeMs: 100000,
  //       timeoutRetry: {
  //         maxNumRetry: 2,
  //         retryDelayMs: 0,
  //         maxRetryDelayMs: 0,
  //       },
  //       errorRetry: {
  //         maxNumRetry: 5,
  //         retryDelayMs: 3000,
  //         maxRetryDelayMs: 15000,
  //         backoff: 'linear',
  //       },
  //     },
  //   },
  //   startLevel: undefined,
  //   audioPreference: {
  //     characteristics: 'public.accessibility.describes-video',
  //   },
  //   subtitlePreference: {
  //     lang: 'en-US',
  //   },
  //   startFragPrefetch: false,
  //   testBandwidth: true,
  //   progressive: false,
  //   lowLatencyMode: true,
  //   fpsDroppedMonitoringPeriod: 5000,
  //   fpsDroppedMonitoringThreshold: 0.2,
  //   appendErrorMaxRetry: 3,
  //   loader: customLoader,
  //   fLoader: customFragmentLoader,
  //   pLoader: customPlaylistLoader,
  //   xhrSetup: XMLHttpRequestSetupCallback,
  //   fetchSetup: FetchSetupCallback,
  //   abrController: AbrController,
  //   bufferController: BufferController,
  //   capLevelController: CapLevelController,
  //   fpsController: FPSController,
  //   timelineController: TimelineController,
  //   enableDateRangeMetadataCues: true,
  //   enableMetadataCues: true,
  //   enableID3MetadataCues: true,
  //   enableWebVTT: true,
  //   enableIMSC1: true,
  //   enableCEA708Captions: true,
  //   stretchShortVideoTrack: false,
  //   maxAudioFramesDrift: 1,
  //   forceKeyFrameOnDiscontinuity: true,
  //   abrEwmaFastLive: 3.0,
  //   abrEwmaSlowLive: 9.0,
  //   abrEwmaFastVoD: 3.0,
  //   abrEwmaSlowVoD: 9.0,
  //   abrEwmaDefaultEstimate: 500000,
  //   abrEwmaDefaultEstimateMax: 5000000,
  //   abrBandWidthFactor: 0.95,
  //   abrBandWidthUpFactor: 0.7,
  //   abrMaxWithRealBitrate: false,
  //   maxStarvationDelay: 4,
  //   maxLoadingDelay: 4,
  //   minAutoBitrate: 0,
  //   emeEnabled: false,
  //   licenseXhrSetup: undefined,
  //   drmSystems: {},
  //   drmSystemOptions: {},
  //   requestMediaKeySystemAccessFunc: requestMediaKeySystemAccess,
  //   cmcd: {
  //     sessionId: uuid(),
  //     contentId: hash(contentURL),
  //     useHeaders: false,
  //   },
  // };

  // this conf set only works smoothly on 1.8.3 (pre-auth), choppy on 1.12.2
  // const lsConf001 = {
  //   maxBufferLength: 5, // seconds
  //   maxMaxBufferLength: 10, // seconds
  //   backBufferLength: 0,
  //   frontBufferFlushThreshold: Infinity,
  //   maxBufferSize: 60 * 1000 * 1000,

  //   highBufferWatchdogPeriod: 2, // seconds

  //   liveSyncDurationCount: 0.1,
  //   liveSyncOnStallIncrease: 0.1,
  //   liveMaxLatencyDurationCount: 0.3,
  //   maxLiveSyncPlaybackRate: 2,
  //   liveDurationInfinity: true,
  // };

  // https://github.com/video-dev/hls.js/blob/master/docs/API.md#fine-tuning
  const lsConf002 = {
    autoStartLoad: true,
    startPosition: -1,

    maxBufferLength: 10, // seconds
    maxMaxBufferLength: 60, // seconds
    backBufferLength: 5,
    frontBufferFlushThreshold: Infinity,
    maxBufferSize: 50 * 1000 * 1000,

    highBufferWatchdogPeriod: 2, // seconds

    // TODO: revisit later
    // - periodically auto adjust hls.targetLatency for each instance...
    // - also this, abrEwmaDefaultEstimate

    // hls.targetLatency = (liveSyncDurationCount * EXT-X-TARGETDURATION) + (liveSyncOnStallIncrease * n-stalls)
    liveSyncMode: 'edge',
    liveSyncDurationCount: 1,
    liveSyncOnStallIncrease: 0.5,
    liveMaxLatencyDurationCount: 4,
    maxLiveSyncPlaybackRate: 1.1, // min 1, max 2
    liveDurationInfinity: true,

    enableWorker: true, // TODO

    startFragPrefetch: true,
    lowLatencyMode: true,

    progressive: true,
    // xhrSetup: function(xhr, url) {
    //   xhr.open('GET', url, true);
    //   xhr.withCredentials = true;
    //   xhr.setRequestHeader('Authorization', 'Basic dmlld2VyMDAxOmJpd2VyMDAx');
    //   xhr.send();
    // },
    fetchSetup: function(context, initParams) {
      initParams.credentials = 'same-origin';
      // TODO: dynamic authn
      initParams.headers.append('Authorization', 'Basic dmlld2VyMDAxOmJpd2VyMDAx'); // btoa()
      return new Request(context.url, initParams);
    },

    // loader: customLoader, // TODO: explore, might be easier to access fields
  };

  return lsConf002;
};

_nsDLS.componentExists = function(componentId) {
  const dls = _nsDLS.instances[componentId];

  if (dls == undefined || dls == null) {
    return false;
  }

  return true;
};

// @params componentId: string
// @params componentEl: HTMLElement
// @params hlsSourceUrl: HTTP URL
_nsDLS.newComponent = function(componentId, componentEl, hlsSourceUrl) {
  const dls = {};
  dls.id = componentId;
  dls.compEl = componentEl; // TODO: validate
  dls.sourceURL = '';

  dls.compEl.innerHTML = _nsDLS.template;

  dls.loadingEl = dls.compEl.getElementsByClassName('dls-loading')[0];
  dls.messageEl = dls.compEl.getElementsByClassName('dls-message')[0];
  dls.videoEl   = dls.compEl.getElementsByClassName('dls-video')[0];


  // note: sourceUrl from js take precedence over sourceUrl from tag
  const sourceUrlFromTag = dls.compEl.getAttribute('data-fu-src');

  if (hlsSourceUrl == undefined || hlsSourceUrl == null) {
    dls.sourceURL = sourceUrlFromTag;
  } else {
    dls.sourceURL = hlsSourceUrl;
  }

  if (dls.sourceURL == undefined || dls.sourceURL == null) {
    dls.sourceURL = '';
  }

  dls.compEl.setAttribute('data-fu-id', dls.id);
  dls.compEl.setAttribute('data-fu-src', dls.sourceURL);


  if (dls.sourceURL == undefined || dls.sourceURL == null || dls.sourceURL == '') {
    const errMessage = '[FATAL] data source required!';
    _nsDLS.logger.error(dls.id, errMessage);

    dls.messageEl.innerHTML = `<p>${errMessage}</p>`;

    dls.loadingEl.style.display = 'none';
    dls.messageEl.style.display = '';
    dls.videoEl.style.display = 'none';

  } else {
    dls.loadingEl.style.display = '';
    dls.messageEl.style.display = 'none';
    dls.videoEl.style.display = 'none';

    _nsDLS._injectHlsInstance(dls);

    // TODO: poster / thumbnail and stuffs
    dls.hlsInstance.loadSource(dls.sourceURL);
  }

  _nsDLS.instances[dls.id] = dls;
};

// @params componentId: string
_nsDLS.destroyComponent = function(componentId) {
  const dls = _nsDLS.instances[componentId];

  if (dls == undefined || dls == null) {
    _nsDLS.logger.warn(componentId, `dls instance not found`);
    return
  }

  dls.compEl.innerHTML = '';
  dls.compEl.removeAttribute('data-fu-id');

  if (dls.hlsInstance == undefined || dls.hlsInstance == null) {
  } else {
    dls.hlsInstance.stopLoad();
    dls.hlsInstance.detachMedia();
    dls.hlsInstance.destroy();
    dls.hlsInstance = null;
  }

  _nsDLS.instances[componentId] = null;
};


// intended to be used only once for sweeping start
_nsDLS.renderAllUsed = false;
_nsDLS.renderAll = function() {
  if (!Hls.isSupported()) {
    _nsDLS.logger.error('ALL', 'lib hls not supported?!');
    return
  }

  if (_nsDLS.renderAllUsed) {
    _nsDLS.logger.warn('ALL', 'renderAll already used!');
    return
  }

  _nsDLS.renderAllUsed = true;

  const components = document.getElementsByTagName('fu-default-live-stream');
  for (let i = 0; i < components.length; i++) {
    _nsDLS.newComponent(i, components.item(i), null);
  }
}

// intended for internal use
// @params dls: dlsInstance
_nsDLS._injectHlsInstance = function(dls) {
  dls.hlsInstance = new Hls(_nsDLS.defaultLiveStreamConf());

  // dls.hlsInstance.on(Hls.Events.FRAG_BUFFERED, function(ev, data) {
  //   // console.log('Quality switch API', {
  //   //   bwEstimate: dls.hlsInstance.bandwidthEstimate,
  //   // });
  //   console.log('live stream API', {
  //     liveSyncPos: dls.hlsInstance.liveSyncPosition,
  //     latency: dls.hlsInstance.latency,
  //     maxLatency: dls.hlsInstance.maxLatency,
  //     targetLatency: dls.hlsInstance.targetLatency,
  //     drift: dls.hlsInstance.drift,
  //     playingDate: dls.hlsInstance.playingDate,
  //   });
  //   // console.log('playback rate', vidEl.playbackRate)
  // });

  // dls.hlsInstance.on(Hls.Events.FPS_DROP, function(ev, data) {
  //   console.log('FPS_DROP', data);
  // });

  dls.hlsInstance.on(Hls.Events.MANIFEST_PARSED, function(ev, data) {
    dls.loadingEl.style.display = 'none';
    dls.messageEl.style.display = 'none';
    dls.videoEl.style.display = '';

    // TODO: attach after play
    dls.hlsInstance.attachMedia(dls.videoEl);

    // dls.hlsInstance.startLoad(startPosition=-1,skipSeekToStartPosition=false); // ??? silent error ?

    dls.videoEl.setAttribute('controls', '');
    dls.videoEl.setAttribute('disablepictureinpicture', ''); // test later, seems like my browser does not support that attr yet

    dls.videoEl.muted = true; // cannot use setAttribute to alter this...

    dls.videoEl.play();
  });

  // TODO: better targeted handling
  dls.hlsInstance.on(Hls.Events.ERROR, function(ev, data) {
    const errorType = data.type;
    const errorDetails = data.details;
    const errorFatal = data.fatal;

    if (errorFatal) {
      switch (errorType) {
        case Hls.ErrorTypes.NETWORK_ERROR:
          // _nsDLS.resetHls(dls.id);
          break;
        case Hls.ErrorTypes.MEDIA_ERROR:
          _nsDLS.logger.warn(dls.id, 'fatal media error encountered, trying to recover...');
          dls.hlsInstance.recoverMediaError();
          break;
        case Hls.ErrorTypes.KEY_SYSTEM_ERROR:
        case Hls.ErrorTypes.MUX_ERROR:
        case Hls.ErrorTypes.OTHER_ERROR:
        default:
          break;
      }

      const errMessage = `[FATAL] ${errorType} - ${errorDetails}`;
      _nsDLS.logger.error(dls.id, errMessage);

      dls.messageEl.innerHTML = `<p>${errMessage}</p>`;

      dls.loadingEl.style.display = 'none';
      dls.messageEl.style.display = '';
      dls.videoEl.style.display = 'none';
    }

    if (_nsDLS.debug) {
      switch (errorType) {
        case Hls.ErrorTypes.NETWORK_ERROR:
          // https://github.com/video-dev/hls.js/blob/master/docs/API.md#network-errors
          switch (errorDetails) {
            case Hls.ErrorDetails.MANIFEST_LOAD_ERROR:
            // ...
          }
        case Hls.ErrorTypes.MEDIA_ERROR:
        case Hls.ErrorTypes.KEY_SYSTEM_ERROR:
        case Hls.ErrorTypes.MUX_ERROR:
        case Hls.ErrorTypes.OTHER_ERROR:
        default:
          _nsDLS.logger.debug(dls.id, `${errorType} - ${errorDetails} - ${errorFatal}`);
          break;
      }
    }
  });
};

// @params componentId: string
// @params hlsSourceUrl: HTTP URL
_nsDLS.changeHlsSource = function(componentId, hlsSourceUrl) {
  const dls = _nsDLS.instances[componentId];

  if (dls == undefined || dls == null) {
    _nsDLS.logger.error(componentId, 'dls instance not found!');
    return
  }

  if (dls.hlsInstance == undefined || dls.hlsInstance == null) {
  } else {
    dls.hlsInstance.stopLoad();
    dls.hlsInstance.detachMedia();
    dls.hlsInstance.destroy();
    dls.hlsInstance = null;
  }

  if (hlsSourceUrl == undefined || hlsSourceUrl == null) {
    dls.sourceURL = '';
  } else {
    dls.sourceURL = hlsSourceUrl;
  }
  dls.compEl.setAttribute('data-fu-src', dls.sourceURL);


  if (dls.sourceURL == undefined || dls.sourceURL == null || dls.sourceURL == '') {
    const errMessage = '[FATAL] data source required!';
    _nsDLS.logger.error(dls.id, errMessage);

    dls.messageEl.innerHTML = `<p>${errMessage}</p>`;

    dls.loadingEl.style.display = 'none';
    dls.messageEl.style.display = '';
    dls.videoEl.style.display = 'none';

  } else {
    dls.loadingEl.style.display = '';
    dls.messageEl.style.display = 'none';
    dls.videoEl.style.display = 'none';

    _nsDLS._injectHlsInstance(dls);
  
    // TODO: poster / thumbnail
    dls.hlsInstance.loadSource(dls.sourceURL);
  }
};

_nsDLS.reloadStream = function(componentId) {
  const dls = _nsDLS.instances[componentId];

  if (dls == undefined || dls == null) {
    _nsDLS.logger.error(componentId, 'dls instance not found!');
    return
  }

  if (dls.hlsInstance == undefined || dls.hlsInstance == null) {
  } else {
    dls.hlsInstance.stopLoad();
    dls.hlsInstance.detachMedia();
    dls.hlsInstance.destroy();
    dls.hlsInstance = null;
  }

  if (dls.sourceURL == undefined || dls.sourceURL == null || dls.sourceURL == '') {
    _nsDLS.logger.error(componentId, 'cannot reload without data source. use changeHlsSource instead.');
    return
  }
  dls.compEl.setAttribute('data-fu-src', dls.sourceURL);

  dls.loadingEl.style.display = '';
  dls.messageEl.style.display = 'none';
  dls.videoEl.style.display = 'none';

  _nsDLS._injectHlsInstance(dls);

  // TODO: poster / thumbnail
  dls.hlsInstance.loadSource(dls.sourceURL);
};

// TODO: stopStream
