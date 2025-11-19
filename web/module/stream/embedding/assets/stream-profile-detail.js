'use strict'

// nsContent
const nsCt = {};
nsCt.init = function() {
  nsCt.el_btnSubmitDelete = document.getElementById('btnSubmitDelete');

  nsCt.fn_doSubmitDelete = function() {
    nsCt.el_btnSubmitDelete.click();
  };
};

document.addEventListener('DOMContentLoaded', function(e) {
  nsCt.init();

  nsMain.focusOn('code');
});