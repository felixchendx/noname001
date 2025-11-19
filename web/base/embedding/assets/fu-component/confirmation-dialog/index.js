'use strict';

const nsFuConfirmationDialog = {};

function initFuConfirmationDialog() {
  nsFuConfirmationDialog.template = `
  `;

  nsFuConfirmationDialog.rendered = () => {};
  nsFuConfirmationDialog.renderAll = () => {
    let components = document.getElementsByTagName('fu-confirmation-dialog');
    for (let i = 0; i < components.length; i++) {
      // const currMarkup = (new DOMParser()).parseFromString(nsFuConfirmationDialog.template, 'text/html');
      let currComponent = components.item(i);
      let currMarkup = nsFuConfirmationDialog.template;
      currComponent.innerHTML = currMarkup;

      // let currComponentID = currComponent.getAttribute('id');
      // let innerEl = document.getElementById('fuConfirmationDialog');
      // innerEl.setAttribute('id', currComponentID)

      // currComponent.removeAttribute('id');
    }
  };
}
