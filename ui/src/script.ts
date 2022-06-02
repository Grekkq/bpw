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


        const formDataObject = Object.fromEntries(formData.entries()) as any;
        if(formDataObject["sys"]){
            formDataObject["sys"] = parseInt(formDataObject["sys"]);
        }
        if(formDataObject["dia"]){
            formDataObject["dia"] = parseInt(formDataObject["dia"]);
        }
        if(formDataObject["pulse"]){
            formDataObject["pulse"] = parseInt(formDataObject["pulse"]);
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
            throw new Error(error);
        }

        return res.json();
    });
};