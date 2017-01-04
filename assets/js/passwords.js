//TODO do i need to check headers?
//spinny icon and shiz
//There might be a better way to do this..., good enough for now
function copyPW(token) {
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

function deletePW(token) {
  console.log("not implemented");
}
