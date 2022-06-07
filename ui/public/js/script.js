"use strict";
window.onload = function () {
    const formAddEntry = document.getElementById("test");
    formAddEntry?.addEventListener("submit", async (e) => {
        e.preventDefault();
        const form = e.currentTarget;
        if (!(form instanceof HTMLFormElement)) {
            throw new Error("Element is not a form");
        }
        const formData = new FormData(form);
        const url = form.action;
        const formDataObject = Object.fromEntries(formData.entries());
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
            errorToast(`Cannot add new entry\n${error}`);
            return;
        }
        successToast(`Entry ${formDataJsonString} added succesfully`);
        form.reset();
    });
};
