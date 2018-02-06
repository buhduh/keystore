(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);var f=new Error("Cannot find module '"+o+"'");throw f.code="MODULE_NOT_FOUND",f}var l=n[o]={exports:{}};t[o][0].call(l.exports,function(e){var n=t[o][1][e];return s(n?n:e)},l,l.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){
module.exports = require('./src/password');

},{"./src/password":2}],2:[function(require,module,exports){
'use-strict';
var charRange = require('./util/genAsciiCharRangeArr');

// predefined charSet constants
var UPPERCASE = 'UPPERCASE',
    LOWERCASE = 'LOWERCASE',
    DIGIT = 'DIGIT',
    SPECIAL_CHAR = 'SPECIAL_CHAR';

var DEFAULT_OPTIONS = {
    // ASCII character decimal range, i.e. all printable chars excluding Space and Delete
    charMin: 33,
    charMax: 126,
    length: 32,
    // What chars and how many of each set should be in password
    exclusions: [],
    inclusionRules: [
        {
            minNumChars: 3,
            charSet: UPPERCASE
        },
        {
            minNumChars: 3,
            charSet: LOWERCASE
        },
        {
            minNumChars: 3,
            charSet: SPECIAL_CHAR
        },
        {
            minNumChars: 3,
            charSet: DIGIT
        }
    ]
};

function password(options) {
    var possibleCharsAfterExcl,
        necessaryChars = 0,
        charArr = [];

    var options = Object.assign({}, DEFAULT_OPTIONS, options);

    var firstPosCharBeforeExcl = String.fromCharCode(options.charMin),
        lastPosCharBeforeExcl = String.fromCharCode(options.charMax);

    if (typeof options.length !== 'number') {
        console.error('password: invalidParameter in options: "length" must be of type number');
        return false;
    }
    if (options.length < 14) {
        console.warn('password: passwords of character length less than 14 are not recommended');
    }
    if (options.exclusions && options.exclusions.constructor !== Array) {
        console.error('password: invalidParameter in options: "exclusions" must be an Array');
        return false;
    }
    if (typeof options.inclusionRules !== 'undefined' && options.inclusionRules.constructor === Array) {
        // final inclusions, exclusions taking precedent
        options.inclusionRules.forEach(function(rule) {
            // charSet can be represented as a string, not an array. If the array is a constant, provide cooresponding array, else use the string
            if (typeof rule.charSet === 'string' && rule.charSet.length > 1) {
                rule.charSet = _getCharSetFromConstant(rule.charSet);
            }
            rule.charSet = _arrDiff(rule.charSet, options.exclusions);
            // check that after exclusions there are still some characters to pull from this inclusionRule charSet
            if (rule.charSet.length <= 0) {
                console.error('password: invalidParamer in options: one of your inclusionRules were negated completetly by your exclusions');
                return false;
            }
            rule.finalChars = [];
            for (var i = 0; i < rule.minNumChars; i++) {
                necessaryChars += 1;
                rule.finalChars.push(rule.charSet[_getRandomIntInclusive(0, rule.charSet.length - 1)]);
            }
        });
        if (options.length < necessaryChars) {
            console.error('password: invalidParameter in options: "length" and ' +
                '"inclusionRules." inclusionRules character minimum cannot exceed length');
            return false;
        }
        options.inclusionRules.forEach(function(rule) {
            rule.finalChars.forEach(function(char) {
                charArr.push(char);
            });
        });
    }

    possibleCharsAfterExcl = _arrDiff(charRange(firstPosCharBeforeExcl, lastPosCharBeforeExcl), options.exclusions);

    for (i = 0; i < (options.length - necessaryChars); i++) {
        charArr.push(possibleCharsAfterExcl[_getRandomIntInclusive(0, possibleCharsAfterExcl.length - 1)]);
    }

    charArr = _arrShuffle(charArr);
    return charArr.join('');
}

function _getCharSetFromConstant(charSetStr) {
    var charSet;
    switch(charSetStr) {
        case UPPERCASE:
            charSet = charRange('A', 'Z');
            break;
        case LOWERCASE:
            charSet = charRange('a', 'z');
            break;
        case DIGIT:
            charSet = charRange('0', '9');
            break;
        case SPECIAL_CHAR:
            charSet = charRange('!', '/').concat(charRange(':', '@')).concat(charRange('[', '`')).concat(charRange('{', '~'));
            break;
        default:
            charSet = charSetStr;
    }
    return charSet;
}

// substract arr2 from arr1
function _arrDiff(arr1, arr2) {
    return arr1.filter(function(itemA) {
        var pass = true;
        arr2.forEach(function(itemB) {
            if (itemB === itemA) {
                pass = false;
            }
        });
        return pass;
    });
}

// random permutation of character array set
function _arrShuffle(arr) {
    var currentIndex = arr.length,
        temporaryValue, randomIndex;
    while (0 !== currentIndex) {
        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex -= 1;

        temporaryValue = arr[currentIndex];
        arr[currentIndex] = arr[randomIndex];
        arr[randomIndex] = temporaryValue;
    }
    return arr;
}

// get random integer from min to max including min and max
function _getRandomIntInclusive(min, max) {
    var min = Math.ceil(min),
        max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

module.exports = password;

},{"./util/genAsciiCharRangeArr":3}],3:[function(require,module,exports){
module.exports = function(charFirst, charLast) {
    var arr = [],
        i = charFirst.charCodeAt(0),
        j = charLast.charCodeAt(0);
    for (; i <= j; ++i) {
        arr.push(String.fromCharCode(i));
    }
    return arr;
};

},{}],4:[function(require,module,exports){
module.exports = require('./src/js/collapsible');

},{"./src/js/collapsible":5}],5:[function(require,module,exports){
'use-strict';
var COLLAPSED_STATE = 'collapsed';
var EXPANDED_STATE = 'expanded';

var Collapsible = function(targetElem, triggerElem, initState) {
    this.state = initState || COLLAPSED_STATE;
    this.targetElem = targetElem;
    this.triggerElem = triggerElem;
    this.isExpanding = false;
    this.isCollapsing = false;
};

Collapsible.prototype = {
    init: function() {
        this.initState();
        this.bindUi();
    },
    initState: function() {
        this.height = this.getHeight();
        if (this.state === COLLAPSED_STATE) {
            this.targetElem.classList.remove('expanded');
            this.targetElem.classList.add('collapsed');
            this.triggerElem.classList.remove('expanded');
            this.triggerElem.classList.add('collapsed');
            this.setHeight(0);
        } else {
            this.targetElem.classList.remove('collapsed');
            this.targetElem.classList.add('expanded');
            this.triggerElem.classList.remove('collapsed');
            this.triggerElem.classList.add('expanded');
            this.setHeight(this.height);
        }
    },
    setHeight: function(height) {
        this.targetElem.style.height = height + 'px';
    },
    getHeight: function() {
        return this.targetElem.clientHeight;
    },
    bindUi: function() {
        this.triggerElem.addEventListener('click', function(e) {
            this.toggle();
        }.bind(this));
        var expandingComplete = function() {
            this.targetElem.classList.remove('expanding');
            this.targetElem.classList.add('expanded');
            this.triggerElem.classList.remove('expanding');
            this.triggerElem.classList.add('expanded');
            this.targetElem.removeAttribute('style');
            this.isExpanding = false;
            this.isCollapsing = false;
        }.bind(this);
        var collapsingComplete = function() {
            this.targetElem.classList.remove('collapsing');
            this.targetElem.classList.add('collapsed');
            this.triggerElem.classList.remove('collapsing');
            this.triggerElem.classList.add('collapsed');
            this.isCollapsing = false;
            this.isExpanding = false;
        }.bind(this);
        regTransEndEvent(this.targetElem, function() {
            if (this.state === EXPANDED_STATE) {
                expandingComplete();
            } else if (this.state === COLLAPSED_STATE) {
                collapsingComplete();
            }
        }.bind(this));
    },
    show: function() {
        if (!this.isCollapsing && !this.isExpanding) {
            this.isExpanding = true;
            this.setHeight(this.height);
            this.targetElem.classList.remove('collapsed');
            this.targetElem.classList.add('expanding');
            this.triggerElem.classList.remove('collapsed');
            this.triggerElem.classList.add('expanding');
            this.state = EXPANDED_STATE;
        }
    },
    hide: function() {
        if (!this.isCollapsing && !this.isExpanding) {
            this.isCollapsing = true;
            this.setHeight(this.height);
            var timer;
            window.clearTimeout(timer);
            timer = window.setTimeout(function() {
                this.setHeight(0);
            }.bind(this), 10);
            this.targetElem.classList.remove('expanded');
            this.targetElem.classList.add('collapsing');
            this.triggerElem.classList.remove('expanded');
            this.triggerElem.classList.add('collapsing');
            this.state = COLLAPSED_STATE;
        }
    },
    toggle: function() {
        this.state === EXPANDED_STATE ? this.hide() : this.show();
    }
};

function regTransEndEvent(elems, callback) {
    var transitions = {
        'WebkitTransition': 'webkitTransitionEnd',
        'MozTransition': 'transitionend',
        'MSTransition': 'msTransitionEnd',
        'OTransition': 'oTransitionEnd',
        'transition': 'transitionend'
    };
    if (elems.constructor === Array) {
        for (var t in transitions) {
            for (var i = 0; i < elems.length; i++) {
                if (elems[i].style[t] !== undefined) {
                    elems[i].addEventListener(transitions[t], function(e) {
                        callback();
                    }, false);
                }
            }
        }
    } else {
        for (var t in transitions) {
            if (elems.style[t] !== undefined) {
                elems.addEventListener(transitions[t], function(e) {
                    callback();
                }, false);
            }
        }
    }
};

module.exports = Collapsible;

},{}],6:[function(require,module,exports){
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

},{"./collapsible":4,"rand-password-gen":1}]},{},[6]);
