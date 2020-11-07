window.onload = function(){
    let createtheme = document.getElementById("crTheme");
    let modalcrtheme = document.getElementById("createThemeField");
    let createtask = document.getElementById("crTask");
    let modalcrtask = document.getElementById("createTaskField");
    
     createtheme.onclick = function(){
        modalcrtheme.style.display="block";
        modalcrtheme.onclick = function(event) {
            if(!event.target.className.includes("createForm")){
            modalcrtheme.style.display = "none";  }      
        }
    };  
    createtask.onclick = function(){
        modalcrtask.style.display="block";
        modalcrtask.onclick = function(event) {
            if(!event.target.className.includes("createForm")){
            modalcrtask.style.display = "none";  }      
        }
    };  
    let del = document.getElementById("deltask");
    del.onclick = function(event) {
        event.stopPropagation();
    }
    
}