const worker = new Worker('worker.js')
function convert(file) {
    var reader = new FileReader();
    reader.readAsDataURL(file);
}

document.addEventListener("DOMContentLoaded", function(){
    // file selected
    document.getElementById("file").addEventListener("change", function(evt) {
        convert(evt.target.files[0]);
    }, false);
});