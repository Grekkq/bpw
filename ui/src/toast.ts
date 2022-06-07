function errorToast(message: string, visibleFor = 3) {
    const successColor = window.getComputedStyle(document.body).getPropertyValue("--error-color");
    displayToast(message, successColor, visibleFor);
}

function successToast(message: string, visibleFor = 3) {
    const successColor = window.getComputedStyle(document.body).getPropertyValue("--success-color");
    displayToast(message, successColor, visibleFor);
}

function displayToast(message: string, color: string, visibleFor = 3) {
    const toast = document.getElementById("toast");
    if (!(toast instanceof HTMLDivElement)) {
        throw new Error("Missing toast div");
    }
    toast.textContent = message;
    toast.style.borderColor = color;
    animate(toast, visibleFor);    
}

function animate(toast: HTMLDivElement, visibleFor: number) {
    toast.className = "fadeIn";
    setTimeout(function () {
        toast.className = toast.className.replace("fadeIn", "fadeOut");
        setTimeout(() => {
            toast.className = toast.className.replace("fadeOut", "");
        }, 450);
    }, visibleFor * 1000);
}