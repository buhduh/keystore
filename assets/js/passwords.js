function getID(elem) {
  return elem.parentNode.parentNode.className.trim()
}

function editPW(elem) {
  var id = getID(elem);
  if(!id || id == "0") {
    //can't really do anything....
    return
  }
  window.location = "/passwords/new?id=" + id;
}

function copyPW(elem) {
  console.log("not implemented")
}

function deletePW(elem) {
  console.log("not implemented")
}

(function() {
}())
