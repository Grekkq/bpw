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
            displaySnackbar(`Cannot add new entry\n${error}`, "red");
            return;
        }
        displaySnackbar(`Entry ${formDataJsonString} added succesfully`, "green");
        form.reset();
    });
};

function displaySnackbar(message: string, backgroundColor: string) {
    const x = document.getElementById("snackbar");
    if (!(x instanceof HTMLDivElement)) {
        throw new Error("Missing snackbar div");
    }
    x.textContent = message;
    x.style.backgroundColor = backgroundColor;
    x.className = "show";
    setTimeout(function () { x.className = x.className.replace("show", ""); }, 5000);
}