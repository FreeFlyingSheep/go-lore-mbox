function addButton() {

}

function fold(id) {
    buttonID = 'button-' + id;
    foldID = id;
    button = document.getElementById(buttonID).getElementsByTagName('a')[0];
    content = document.getElementById(foldID);
    if (button.innerHTML == '[-]') {
        button.innerHTML = '[+]'
        content.style.display = 'none';
    } else {
        button.innerHTML = '[-]'
        content.style.display = 'block';
    }
}
