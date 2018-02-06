var pw = require('rand-password-gen'),
  Collapsible = require('./collapsible'),
  pwField = document.getElementsByClassName('password'),
  generatedPwField = document.querySelectorAll('[name=generatedPw]'),
  generatePwBtn = document.getElementsByClassName('generatePw'),
  acceptPwBtn = document.getElementsByClassName('acceptPw'),
  pwLengthField = document.querySelectorAll('[name=pwLength]'),
  pwSpecialCheck = document.querySelectorAll('[name=pwSpecial]'),
  pwSpecialField = document.querySelectorAll('[name=pwSpecialCharNum]'),
  pwNumCheck = document.querySelectorAll('[name=pwNum]'),
  pwNumField = document.querySelectorAll('[name=pwNumCharNum]'),
  pwUppercase = document.querySelectorAll('[name=pwUppercase]'),
  pwUppercaseCharNum = document.querySelectorAll('[name=pwUppercaseCharNum]'),
  pwLowercase = document.querySelectorAll('[name=pwLowercase]'),
  pwLowercaseCharNum = document.querySelectorAll('[name=pwLowercaseCharNum]'),
  excludedCharsField = document.getElementsByClassName('excludedChars'),
  firstPwGenOpened = false;

function initPasswordViewToggle() {
  var iconEye = document.getElementsByClassName('iconEye');
  iconEye[0].addEventListener('click', function() {
    pwField[0].setAttribute('type', pwField[0].getAttribute('type') === 'password' ? 'text' : 'password');
  });
}

function initPasswordGen() {
  var targetElems = document.getElementsByClassName('collapsibleTarget'),
    triggerElems = document.getElementsByClassName('collapsibleTrigger'),
    pwGenTarget = document.getElementsByClassName('pwGenTarget');

  pwGenTarget[0].addEventListener('click', function(event) {
    event.preventDefault();
    if (!firstPwGenOpened) {
      firstPwGenOpened = true;
      generatePw();
    }
  });

  generatePwBtn[0].addEventListener('click', function(event) {
    event.preventDefault();
    generatePw();
  });

  acceptPwBtn[0].addEventListener('click', function(event) {
    event.preventDefault();
    pwField[0].value = generatedPwField[0].value;
  });

  for (var i = 0; i < targetElems.length; i++) {
    new Collapsible(targetElems[i], triggerElems[i], 'collapsed').init();
  }
}

function generatePw() {
  var optsObj = getPwOptionsObj();
  passwordStr = pw(optsObj);
  generatedPwField[0].value = passwordStr;
}

function getPwOptionsObj() {
  var optsObj = {},
    inclusionRules = [];

  var length = parseInt(pwLengthField[0].value),
    incSpecialChars = pwSpecialCheck[0].checked,
    specialCharsCount = parseInt(pwSpecialField[0].value),
    incNumChars = pwNumCheck[0].checked,
    numCharsCount = parseInt(pwNumField[0].value),
    incUppercase = pwUppercase[0].checked,
    uppercaseCharsCount = parseInt(pwUppercaseCharNum[0].value),
    incLowercase = pwLowercase[0].checked,
    lowercaseCharsCount = parseInt(pwLowercaseCharNum[0].value);

  if (length > 0) {
    optsObj = Object.assign(optsObj, {length: length});
  }

  if (incSpecialChars && specialCharsCount > 0) {
    inclusionRules.push({minNumChars: specialCharsCount, charSet: 'SPECIAL_CHAR'});
  }

  if (incNumChars && numCharsCount > 0) {
    inclusionRules.push({minNumChars: numCharsCount, charSet: 'DIGIT'});
  }

  if (incUppercase && uppercaseCharsCount > 0) {
    inclusionRules.push({minNumChars: uppercaseCharsCount, charSet: 'UPPERCASE'});
  }

  if (incLowercase && lowercaseCharsCount > 0) {
    inclusionRules.push({minNumChars: lowercaseCharsCount, charSet: 'LOWERCASE'});
  }

  if (inclusionRules.length) {
    optsObj = Object.assign(optsObj, {inclusionRules: inclusionRules});
  }

  var excludedChars = excludedCharsField[0].value;
  var excludedCharsArr = [];
  if (excludedChars.length > 0) {
    // first remove out all space characters
    excludedChars = excludedChars.replace(/\s/g, '');
    // check if string is comma separated characters
    if (/^(.,)+[^, ]$/.test(excludedChars)) {
      excludedCharsArr = excludedChars.split(',');
      optsObj = Object.assign(optsObj, {exclusions: excludedCharsArr});
    } else {
      console.error('malformed comma separated excluded characters');
    }
  }

  return optsObj;
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
