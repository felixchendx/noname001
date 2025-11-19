'use strict';

// TODO: spinner for loading data thru ws

// nsContent
const nsCt = {};
nsCt.init = function() {
  nsCt.el_btnSubmitDelete = document.getElementById('btnSubmitDelete');

  nsCt.fn_doSubmitDelete = function() {
    nsCt.el_btnSubmitDelete.click();
  };
};

// nsDlgStreamItem
const nsDSI = {};
nsDSI.init = function() {
  nsDSI.el_title         = document.getElementById('dlgStreamItem_title');

  nsDSI.el_id            = document.getElementById('dlgStreamItem_ID');
  nsDSI.el_code          = document.getElementById('dlgStreamItem_Code');
  nsDSI.el_name          = document.getElementById('dlgStreamItem_Name');
  nsDSI.el_state         = document.getElementById('dlgStreamItem_State');
  nsDSI.el_note          = document.getElementById('dlgStreamItem_Note');

  nsDSI.el_streamGroupID = document.getElementById('dlgStreamItem_StreamGroupID');

  nsDSI.el_sourceType      = document.getElementById('dlgStreamItem_SourceType');
  nsDSI.el_deviceCode      = document.getElementById('dlgStreamItem_DeviceCode');
  nsDSI.el_deviceChannelID = document.getElementById('dlgStreamItem_DeviceChannelID');
  nsDSI.el_externalURL     = document.getElementById('dlgStreamItem_ExternalURL');
  nsDSI.el_filePath        = document.getElementById('dlgStreamItem_Filepath');

  nsDSI.el_btnSaveStreamItem = document.getElementById('btnSaveStreamItem');

  // === TEMP CHANNEL PREVIEW ===
  nsDSI.el_containerDCP = document.getElementById('dcpContainer')
  nsDSI.el_buttonDCP = document.getElementById('dcpButton')
  nsDSI.el_vidContainerDCP = document.getElementById('dcpVideoContainer')
  nsDSI.el_vidDCP = document.getElementById('dcpVideo')
  nsDSI.el_feedbackDCP = document.getElementById('dcpPreviewFeedback')


  nsDSI.fn_onOpenDialog = function(item, isAddMode) {
    nsMain.focusOn('dlgStreamItem_Code');

    nsDSI.el_id.value         = isAddMode ? "" : item.ID;
    nsDSI.el_code.value       = isAddMode ? "" : item.Code;
    nsDSI.el_name.value       = isAddMode ? "" : item.Name;
    nsDSI.el_state.value      = isAddMode ? "" : item.State;
    nsDSI.el_note.value       = isAddMode ? "" : item.Note;

    nsDSI.el_sourceType.value      = isAddMode ? "" : item.SourceType;
    nsDSI.el_deviceCode.value      = isAddMode ? "" : item.DeviceCode;
    nsDSI.el_deviceChannelID.value = isAddMode ? "" : item.DeviceChannelID;
    nsDSI.el_externalURL.value     = isAddMode ? "" : item.ExternalURL;
    nsDSI.el_filePath.value        = isAddMode ? "" : item.Filepath;

    if (isAddMode == true) {
      nsDSI.el_title.innerText = "New Stream item"
      nsDSI.el_btnSaveStreamItem.setAttribute("onclick", "nsDSI.fn_saveStreamItem(true);");
    } else {
      nsDSI.el_title.innerText = `Stream Item ${nsDSI.el_code.value}`
      nsDSI.el_btnSaveStreamItem.setAttribute("onclick", "nsDSI.fn_saveStreamItem(false);");
    }

    nsDSI.fn_visibilitySubSourceType();
    nsDSI.el_sourceType.addEventListener("change", nsDSI.fn_visibilitySubSourceType);
  };

  nsDSI.fn_visibilitySubSourceType = function() {
    nsDSI.el_subSourceType            = document.getElementById('dlgStreamItem_subSourceTpeRow');
    nsDSI.el_ContainerDeviceCode      = document.getElementById('dlgStreamItem_DeviceCodeContainer');
    nsDSI.el_ContainerDeviceChannelID = document.getElementById('dlgStreamItem_DeviceChannelIDContainer');
    nsDSI.el_ContainerExternalURL     = document.getElementById('dlgStreamItem_ExternalURLContainer');
    nsDSI.el_ContainerfilePath        = document.getElementById('dlgStreamItem_FilepathContainer');


    nsDSI.el_subSourceType.style.display = "none";

    nsDSI.el_ContainerDeviceCode.style.display      = "none";
    nsDSI.el_ContainerDeviceChannelID.style.display = "none";
    nsDSI.el_containerDCP.style.display = "none";
    nsDSI.el_vidContainerDCP.style.display = "none";
    nsDSI.el_ContainerExternalURL.style.display     = "none";
    nsDSI.el_ContainerfilePath.style.display        = "none";

    switch (nsDSI.el_sourceType.value) {
      case "mod_device":
        nsDSI.el_subSourceType.style.display = "block";
        nsDSI.el_ContainerDeviceCode.style.display = "inline-block";
        nsDSI.el_ContainerDeviceChannelID.style.display = "inline-block";

        nsDSI.el_containerDCP.style.display = "inline-block";
        nsDSI.el_vidContainerDCP.style.display = "";

        nsDSI.el_filePath.value = "";
        nsDSI.el_externalURL.value = "";
        break;

      case "external":
        nsDSI.el_subSourceType.style.display = "block";
        nsDSI.el_ContainerExternalURL.style.display = "block";

        nsDSI.el_deviceCode.value = "";
        nsDSI.el_deviceChannelID.value = "";
        nsDSI.el_filePath.value = "";
        break;

      case "file":
        nsDSI.el_subSourceType.style.display = "block";
        nsDSI.el_ContainerfilePath.style.display = "block";

        nsDSI.el_deviceCode.value = "";
        nsDSI.el_deviceChannelID.value = "";
        nsDSI.el_externalURL.value = "";
        break;

      default:
        nsDSI.el_deviceCode.value = "";
        nsDSI.el_deviceChannelID.value = "";
        nsDSI.el_externalURL.value = "";
        nsDSI.el_filePath.value = "";
        break;
    }

    nsDSI.fn_dcpReset();
  };

  nsDSI.fn_saveStreamItem = async function(isAddMode) {
    const reqBody = {
      "id"               : nsDSI.el_id.value,
      "code"             : nsDSI.el_code.value,
      "name"             : nsDSI.el_name.value,
      "state"            : nsDSI.el_state.value,
      "note"             : nsDSI.el_note.value,

      "stream_group_id"  : nsDSI.el_streamGroupID.value,

      "source_type"      : nsDSI.el_sourceType.value,
      "device_code"      : nsDSI.el_deviceCode.value,
      "device_channel_id": nsDSI.el_deviceChannelID.value,
      "external_url"     : nsDSI.el_externalURL.value,
      "filepath"         : nsDSI.el_filePath.value,
    };

    const apiFunction = isAddMode ? nsDSI.fn_addStreamItem : nsDSI.fn_editStreamItem;
    const respBundle = await apiFunction(reqBody);
    if (respBundle.isOk()) {
      location.reload();
    } else {
      console.error('TODO: show feedback somewhere ?', respBundle);
    }
  };

  nsDSI.fn_addStreamItem = async function(reqBody) {
    const reqURI = '/stream/local-api/add-stream-item';
    let respBundle = null;

    await nsComm.lapi_defaultPOST(reqURI, JSON.stringify(reqBody))
            .then((_respBundle) => { respBundle = _respBundle })
            .catch((_err) => { console.error(`TODO: deadend ${_err}`); });
  
    return respBundle
  };

  nsDSI.fn_editStreamItem = async function(reqBody) {
    const reqURI = '/stream/local-api/edit-stream-item';
    let respBundle = null;

    await nsComm.lapi_defaultPOST(reqURI, JSON.stringify(reqBody))
            .then((_respBundle) => { respBundle = _respBundle })
            .catch((_err) => { console.error(`TODO: deadend ${_err}`); });
  
    return respBundle
  };

  nsDSI.fn_onRemoveStreamItem = async function(id) {
    const reqURI = '/stream/local-api/delete-stream-item';
    const reqBody = {
      "id": id,
    };

    let respBundle = null;

    await nsComm.lapi_defaultPOST(reqURI, JSON.stringify(reqBody))
            .then((_respBundle) => { respBundle = _respBundle })
            .catch((_err) => { console.error(`TODO: deadend ${_err}`); });

    if (respBundle.isOk()) {
      location.reload();
    } else {
      console.error('TODO: show feedback somewhere ?', respBundle);
    }
  };


  // === TEMP CHANNEL PREVIEW ===
  nsDSI.dat_currentHLS = null;

  // TODO: add reset to event dialog close
  nsDSI.fn_dcpReset = function() {
    if (nsDSI.dat_currentHLS != null) {
      nsDSI.dat_currentHLS.stopLoad()
      nsDSI.dat_currentHLS = null;
    }

    nsDSI.el_buttonDCP.disabled = false;
    nsDSI.el_feedbackDCP.innerHTML = "";
  };

  nsDSI.fn_deviceChannelPreview = async function() {
    nsDSI.fn_dcpReset();
    
    nsDSI.el_buttonDCP.disabled = true;
    nsDSI.el_feedbackDCP.innerHTML = `
      <div class="progress">
        <div class="progress-bar progress-bar-indeterminate"></div>
      </div>
    `;

    const deviceCode = nsDSI.el_deviceCode.value;
    const channelID = nsDSI.el_deviceChannelID.value;

    const reqURI = '/stream/local-api/device-channel-preview';
    const reqBody = {"device_code": deviceCode, "device_channel_id": channelID};

    let respBundle = null;

    await nsComm.lapi_defaultPOST(reqURI, JSON.stringify(reqBody))
            .then((_respBundle) => { respBundle = _respBundle })
            .catch((_err) => { console.error(`TODO: deadend ${_err}`); });

    if (respBundle.isOk) {
      nsDSI.el_buttonDCP.disabled = false;

      if (respBundle.defResp.isOk()) {
        nsDSI.el_feedbackDCP.innerHTML = "";

        const vidEl = nsDSI.el_vidDCP;
        const vidSrc = respBundle.defResp.data.preview_url;
  
        if (Hls.isSupported()) {
          var hls = new Hls();
          hls.loadSource(vidSrc);
          hls.attachMedia(vidEl);
          hls.on(Hls.Events.MANIFEST_PARSED, function() {
            vidEl.play();
          });
  
          nsDSI.dat_currentHLS = hls;
  
        } else if (vidEl.canPlayType('application/vnd.apple.mpegurl')) {
          // TODO: branching check
          vidEl.src = vidSrc;
          vidEl.addEventListener('loadedmetadata', function() {
            vidEl.play();
          });
        }
      } else {
        nsDSI.dat_currentHLS = null;
        nsDSI.el_feedbackDCP.innerHTML = `
          <div class="alert alert-danger m-0 p-2">
            <h4 class="alert-title mb-0">Error</h4>
            <div class="text-secondary">${respBundle.defResp.message}</div>
          </div>
        `;
      }

    } else {
      console.error('err1', respBundle); // for quick check without server log access

      nsDSI.dat_currentHLS = null;
      nsDSI.el_feedbackDCP.innerHTML = `
        <div class="alert alert-danger m-0 p-2">
          <h4 class="alert-title mb-0">Error</h4>
          <div class="text-secondary">Something went wrong...</div>
        </div>
      `;
    }
  };
  nsDSI.el_buttonDCP.addEventListener('click', nsDSI.fn_deviceChannelPreview);

};

document.addEventListener('DOMContentLoaded', function(e) {
  nsCt.init();
  nsDSI.init();

  nsMain.focusOn('code');
});
