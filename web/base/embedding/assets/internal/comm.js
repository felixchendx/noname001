'use strict';

const nsComm = {};
nsComm.init = () => {
  nsComm.verbose = false;

  nsComm.lapi_defaultReqInit = function(_method) {
    const reqHeaders = new Headers();
    reqHeaders.append('Content-Type', 'application/json; charset=utf-8');

    const reqInit = {
      method: _method,
      credentials: 'same-origin',
      headers: reqHeaders,
      body: '{}',
    };

    return reqInit
  };

  nsComm.lapi_newDefaultResponse = function() {
    return {
      status: null, // 'ok' | 'error'
      message: null,
      data: null,

      isError: function() {
        if (this.status == null || this.status == 'error') return true;
        return false;
      },
      isOk: function() { return !this.isError(); },
    };
  };
  nsComm.lapi_newGenericError = function() {
    return {
      errtype: null, // 'server' | 'http' | 'js'
      code: null,
      error: null,
      message: null,
    }
  };
  nsComm.lapi_newResponseBundle = function() {
    return {
      defResp: null,
      genErr: null,

      isError: function() {
        if (this.genErr != null) return true;
        if (this.defResp != null && this.defResp.isError()) return true;
        return false;
      },
      isOk: function() { return !this.isError(); },
    }
  }

  nsComm.lapi_defaultPOST = async function(_uri, _body) {
    const reqInit = nsComm.lapi_defaultReqInit('POST');
    reqInit.body = _body;
  
    const respBundle = nsComm.lapi_newResponseBundle();

    try {
      const resp = await fetch(_uri, reqInit);
      const respStatusCode = resp.status;

      if (respStatusCode == 200 || respStatusCode == 400) {
        const respJson = await resp.json();

        respBundle.defResp = nsComm.lapi_newDefaultResponse();
        respBundle.defResp.status = respJson.status;
        respBundle.defResp.message = respJson.message;
        respBundle.defResp.data = respJson.data;

      // } else if (respStatusCode >= 500 && respStatusCode <= 599) {
      } else {
        // TODO:
        respBundle.genErr = nsComm.lapi_newGenericError();
        respBundle.errtype = 'http';
        respBundle.code = respStatusCode;

        if (nsComm.verbose) console.warn(_uri, respBundle);
      }

    } catch(err) {
      // TODO:
      respBundle.genErr = nsComm.lapi_newGenericError();
      respBundle.errtype = 'js';
      respBundle.error = err;

      if (nsComm.verbose) console.error(_uri, respBundle);
    }

    return respBundle
  };
}
nsComm.init();
