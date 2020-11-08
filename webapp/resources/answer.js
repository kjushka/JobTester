window.onload = function () {
    let rightAns = document.getElementsByClassName("right-ans")
    let wrongAns = document.getElementsByClassName("wrong-ans")
    for (let i = 0; i < rightAns.length; i++) {
        rightAns[i].onclick = function () {
            let response = fetch(
                "http://localhost:8080/answers/"
            )
        };
    }
}