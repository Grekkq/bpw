export default class AddEntry extends HTMLElement {
    template!: HTMLTemplateElement;
    connectedCallback() {
        const templateNull = document.querySelector("#add-entry");
        if (!(templateNull instanceof HTMLTemplateElement)) {
            throw new Error("add-entry template tag is missing.");
        }
        this.template = templateNull;
        const content = this.template.content.cloneNode(true);
        if (!(content instanceof Node)) {
            throw new Error("AddEntry failed to clone content of template.");
        }
        this.appendChild(content);
        this.submitFormEventHandler();
    }

    private submitFormEventHandler() {
        const formAddEntry = this.querySelector("#addEntryForm");
        if (!(formAddEntry instanceof HTMLFormElement)) {
            throw new Error("addEntryForm form tag is missing.");
        }
        formAddEntry.addEventListener("submit", async (e: SubmitEvent) => {
            e.preventDefault();
            const form = e.currentTarget;
            if (!(form instanceof HTMLFormElement)) {
                throw new Error("Element is not a form");
            }
            const formData = new FormData(form);
            const url = "/api/addEntry";
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
                errorToast(`Cannot add new entry\n${error}`);
                return;
            }
            successToast(`Entry ${formDataJsonString} added succesfully`);
            form.reset();
        });
    }
}


