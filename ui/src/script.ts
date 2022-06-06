window.onload = function () {
    const formAddEntry = document.getElementById("test");
    formAddEntry?.addEventListener("submit", async (e: SubmitEvent) => {
        e.preventDefault();

        const form = e.currentTarget;
        if (!(form instanceof HTMLFormElement)) {
            throw new Error("Element is not a form");
        }
        const formData = new FormData(form);
        const url = form.action;

        type FormMap = { [k: string]: string | number; };
        const formDataObject = Object.fromEntries(formData.entries()) as FormMap;
        if (formDataObject["sys"]) {
            if (typeof formDataObject["sys"] == "string") {
                formDataObject["sys"] = parseInt(formDataObject["sys"]);
            }
        }
        if (formDataObject["dia"]) {
            if (typeof formDataObject["dia"] == "string") {
                formDataObject["dia"] = parseInt(formDataObject["dia"]);
            }
        }
        if (formDataObject["pulse"]) {
            if (typeof formDataObject["pulse"] == "string") {
                formDataObject["pulse"] = parseInt(formDataObject["pulse"]);
            }
        }

        const formDataJsonString = JSON.stringify(formDataObject);
        const fetchOptions = {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: formDataJsonString
        };
        const res = await fetch(url, fetchOptions);
        if (!res.ok) {
            const error = await res.text();
            displayToast(`Cannot add new entry\n${error}`, "red");
            return;
        }
        displayToast(`Entry ${formDataJsonString} added succesfully`, "green");
        form.reset();
    });
};

function displayToast(message: string, backgroundColor: string, visibleFor = 3) {
    const toast = document.getElementById("toast");
    if (!(toast instanceof HTMLDivElement)) {
        throw new Error("Missing toast div");
    }
    toast.textContent = message;
    toast.style.backgroundColor = backgroundColor;
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
