function fold(id) {
    buttonID = 'button-' + id;
    foldID = 'fold-' + id;
    button = document.getElementById(buttonID).getElementsByTagName('a')[0];
    content = document.getElementById(foldID);
    if (button.innerHTML == '[-] Collapse') {
        button.innerHTML = '[+] Expand'
        content.style.display = 'none';
    } else {
        button.innerHTML = '[-] Collapse'
        content.style.display = 'block';
    }
}
