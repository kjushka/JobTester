window.onload = function(){
    let create = document.getElementById("createBranch");
    let modalcr = document.getElementById("createField");
    create.onclick = function(){
        modalcr.style.display="block";
        modalcr.onclick = function(event) {
            if(!event.target.className.includes("createForm")){
                modalcr.style.display = "none";  }
        }
    };
    let del = document.getElementById("del");
    del.onclick = function(event) {
        event.stopPropagation();
    }
}