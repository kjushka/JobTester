
window.onload = function(){
    let login = document.getElementById("loginopen");
    let modallog = document.getElementById("loginField");
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
    let btnlog = document.getElementById("btn-login");
    let btnreg = document.getElementById("btn-reg");
    
}