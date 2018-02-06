(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);var f=new Error("Cannot find module '"+o+"'");throw f.code="MODULE_NOT_FOUND",f}var l=n[o]={exports:{}};t[o][0].call(l.exports,function(e){var n=t[o][1][e];return s(n?n:e)},l,l.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){
module.exports = require('./src/js/collapsible');

},{"./src/js/collapsible":2}],2:[function(require,module,exports){
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

},{}],3:[function(require,module,exports){
var Collapsible = require('./collapsible');

//TODO do i need to check headers?
//spinny icon and shiz
//There might be a better way to do this..., good enough for now
function copyPw(token) {
  var req = new XMLHttpRequest();
  req.onreadystatechange = function() {
    if (req.readyState === XMLHttpRequest.DONE) {
      if (req.status === 200) {
        var data = JSON.parse(req.responseText)
        if (!data.data) {
          return
        }
        var e = document.getElementById("password_container");
        e.value = data.data;
        e.focus();
        e.select()
        document.execCommand("copy")
        e.value = "";
      } else {
        alert('There was a problem with the request.');
      }
    }
  }
  //This has to be synchronous
  req.open("POST", AJAX_EP, false);
  req.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  req.send("payload="+encodeURIComponent(token));
}

function deletePw(token, domain) {
  if(!confirm("Are you sure you wish to delete domain: '" + domain + "'.")) {
    return;
  }
  var req = new XMLHttpRequest();
  req.onreadystatechange = function() {
    if (req.readyState === XMLHttpRequest.DONE) {
      if (req.status === 200) {
        var data = JSON.parse(req.responseText)
        if (!data.data) {
          return
        }
        location.reload();
      } else {
        alert('There was a problem with the request.');
      }
    }
  }
  req.open("POST", AJAX_EP);
  req.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  req.send("payload="+encodeURIComponent(token));
}

function removeClass(el, className) {
  if (el.classList) {
    el.classList.remove(className);
  } else {
    el.className = el.className.replace(new RegExp('(^|\\b)' + className.split(' ').join('|') + '(\\b|$)', 'gi'), ' ');
  }
}

function initCollapsible() {
  var targetElems = document.getElementsByClassName('collapsibleTarget'),
  triggerElems = document.getElementsByClassName('collapsibleTrigger');

  for (var i = 0; i < targetElems.length; i++) {
    new Collapsible(targetElems[i], triggerElems[i], 'collapsed').init();
  }
}

function initActions() {
  var copyBtns = document.getElementsByClassName('copy'),
      deleteBtns = document.getElementsByClassName('delete');

  for (var i = 0; i < copyBtns.length; i++) {
    copyBtns[i].addEventListener('click', function() {
        copyPw(this.dataset.copyToken);
    }, false);
  }

  for (var i = 0; i < deleteBtns.length; i++) {
    deleteBtns[i].addEventListener('click', function() {
        var domain = this.parentElement.parentElement.children[0].innerText;
        deletePw(this.dataset.deleteToken, domain);
    }, false);
  }
}

function init() {
  initCollapsible();
  initActions();
  //TODO create a loading state so flash of expanded content of collapsibleTargets are not seen on page load
}

function ready(fn) {
  if (document.readyState != 'loading') {
    fn();
  } else {
    document.addEventListener('DOMContentLoaded', fn);
  }
}

ready(init);

},{"./collapsible":1}]},{},[3]);
