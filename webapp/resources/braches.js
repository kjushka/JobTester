
window.onload = function(){
   
    let create = document.getElementById("createBranch");
    let modalcr = document.getElementById("loginField");

     create.onclick = function(){
        modalcr.style.display="block";
        modalcr.onclick = function(event) {
            if(!event.target.className.includes("createForm")){
            modalcr.style.display = "none";  }      
        }
    };   
    
}