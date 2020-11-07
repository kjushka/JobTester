
window.onload = function(){
    let login = document.getElementById("loginopen");
    let modallog = document.getElementById("loginField");
    let create = document.getElementById("createBranch");
    let modalcr = document.getElementById("createField");
    let reg = document.getElementById("regopen");
    let modalreg = document.getElementById("registrField");
    login.onclick = function(){
        modallog.style.display="block";
        modallog.onclick = function(event) {
            if(!event.target.className.includes("loginForm")){
            modallog.style.display = "none";  }      
        }
    };
    reg.onclick = function(){
        modalreg.style.display="block";
        modalreg.onclick = function(event) {
            if(!event.target.className.includes("regForm")){
            modalreg.style.display = "none";  }      
        }
    };   
     create.onclick = function(){
        modalcr.style.display="block";
        modalcr.onclick = function(event) {
            if(!event.target.className.includes("createForm")){
            modalcr.style.display = "none";  }      
        }
    };   
    
}