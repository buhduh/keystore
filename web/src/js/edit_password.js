var pw = require('rand-password-gen'),
  pwField = document.getElementsByClassName('password');

function initPasswordViewToggle() {
  var iconEye = document.getElementsByClassName('iconEye');
  iconEye[0].addEventListener('click', function() {
    pwField[0].setAttribute('type', pwField[0].getAttribute('type') === 'password' ? 'text' : 'password');
  });
}

function initPasswordGen() {

}

function init() {
  initPasswordViewToggle();
  initPasswordGen();
}

function ready(fn) {
  if (document.readyState != 'loading') {
    fn();
  } else {
    document.addEventListener('DOMContentLoaded', fn);
  }
}

ready(init);
