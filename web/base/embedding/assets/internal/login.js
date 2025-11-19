'use strict';

document.addEventListener("DOMContentLoaded", function() {
  const inputUsername = document.getElementById("inputUsername");
  const inputPassword = document.getElementById("inputPassword");
  const showPassword = document.getElementById("showPassword");
  const hidePassword = document.getElementById("hidePassword");

  inputUsername.focus();
  inputUsername.setSelectionRange(-1, -1);
  hidePassword.style.display = 'none';

  showPassword.addEventListener("click", function(e) {
    e.preventDefault();

    inputPassword.type = 'text';
    showPassword.style.display = 'none';
    hidePassword.style.display = '';

    hidePassword.focus();

  }, false);

  hidePassword.addEventListener("click", function(e) {
    e.preventDefault();

    inputPassword.type = 'password';
    showPassword.style.display = '';
    hidePassword.style.display = 'none';

    inputPassword.focus();
    inputPassword.setSelectionRange(-1, -1);

  }, false);
});
