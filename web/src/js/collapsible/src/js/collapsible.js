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
