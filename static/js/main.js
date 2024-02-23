addEventListener("click", function(e) {
  if (e.target.className == "option") {
    console.log("clicked", e.target.innerHTML);
    document.getElementById("submission").value = e.target.innerHTML;
  }
  // console.log(e.target)
  if (e.target.className == "album_input") {
    document.getElementById("searchresults").innerHTML = "";
  }
})
