window.onload = function () {
    let createtheme = document.getElementById("crTheme");
    let modalcrtheme = document.getElementById("createThemeField");
    let createtask = document.getElementsByClassName("crtask");
    let modalcrtask = document.getElementById("createTaskField");

    createtheme.onclick = function () {
        modalcrtheme.style.display = "block";
        modalcrtheme.onclick = function (event) {
            if (!event.target.className.includes("createForm")) {
                modalcrtheme.style.display = "none";
            }
        }
    };
    for (let i = 0; i < createtask.length; i++) {
        createtask[i].onclick = function () {
            let myval = document.getElementById("create");
            myval.setAttribute('value', createtask[i].children[1].attributes[1].value);
            modalcrtask.style.display = "block";
            modalcrtask.onclick = function (event) {
                if (!event.target.className.includes("createForm")) {
                    modalcrtask.style.display = "none";
                }
            }
            console.log(myval.val)
        };
    }
    let del = document.getElementsByClassName("del");
    for (i = 0; i < del.length; i++) {
        del[i].onclick = function (event) {
            console.log(event);
            event.stopPropagation();
        }
    }

    let newtask = document.getElementsByClassName("edit");
    let modalnewtask = document.getElementById("newTaskField");
    let newtheme = document.getElementsByClassName("edittheme");
    let modalnewtheme = document.getElementById("newThemeField");

    for (let i = 0; i < newtask.length; i++) {
        newtask[i].onclick = function () {
            let newtaskid = document.getElementById("newTaskid");
            let listidtask = document.getElementsByClassName("id");
            newtaskid.setAttribute('value', listidtask[i].attributes[1].value);
            modalnewtask.style.display = "block";
            modalnewtask.onclick = function (event) {
                if (!event.target.className.includes("createForm")) {
                    modalnewtask.style.display = "none";
                }
            }
        };
    }

    for (let i = 0; i < newtheme.length; i++) {
        newtheme[i].onclick = function () {
            let newthemeid = document.getElementById("newThemeid");
            let listidtheme = document.getElementsByClassName("idtheme");
            newthemeid.setAttribute('value', listidtheme[i].attributes[1].value);
            modalnewtheme.style.display = "block";
            modalnewtheme.onclick = function (event) {
                if (!event.target.className.includes("createForm")) {
                    modalnewtheme.style.display = "none";
                }
            }
        };
    }
}
