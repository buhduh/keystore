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

function deletePw(token) {
  console.log("not implemented");
  console.log("token: " + token);
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
        deletePw(this.dataset.deleteToken);
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
