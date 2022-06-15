export default class ListEntries extends HTMLElement {
    connectedCallback() {
        const templateNull = document.querySelector("#list-entries");
        if (!(templateNull instanceof HTMLTemplateElement)) {
            throw new Error("list-entries template tag is missing.");
        }
        this.template = templateNull;
        const content = this.template.content.cloneNode(true);
        if (!(content instanceof Node)) {
            throw new Error("ListEntries failed to clone content of template.");
        }
        this.appendChild(content);
        const loadButton = this.querySelector("#loadEntries");
        if (!(loadButton instanceof HTMLButtonElement)) {
            throw new Error("loadEntries button is missing.");
        }
        loadButton.addEventListener("click", () => this.loadData());
    }
    async loadData() {
        const userIdField = this.querySelector("#userId");
        if (!(userIdField instanceof HTMLInputElement)) {
            throw new Error("User ID form is not present");
        }
        const userIdValue = userIdField.value;
        const data = await this.getDataFromApi(userIdValue);
        if (data.length == 0) {
            errorToast(`No data present in db for this userId: ${userIdValue}`);
        }
        this.loadDataIntoTable(data);
    }
    async getDataFromApi(userIdValue) {
        const formDataJsonString = JSON.stringify({ userId: userIdValue });
        const fetchOptions = {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: formDataJsonString
        };
        const response = await fetch("api/listEntries", fetchOptions);
        if (response.status != 200) {
            errorToast("Failed to get responese from server");
        }
        const json = await response.json();
        const entries = [];
        for (const entry of json) {
            entries.push(await entry);
        }
        return entries;
    }
    loadDataIntoTable(data) {
        const template = this.querySelector("#row");
        if (!(template instanceof HTMLTemplateElement)) {
            throw new Error("Template table row not present");
        }
        template.parentNode?.replaceChildren(template);
        for (const item of data) {
            const clone = template.content.cloneNode(true);
            const row = clone.querySelectorAll("td");
            if (!row[0] || !row[1] || !row[2] || !row[3] || !row[4]) {
                throw new Error("Template doesn't have enough <td> tags");
            }
            row[0].textContent = item.timestamp.toFixed();
            row[1].textContent = item.sys.toFixed();
            row[2].textContent = item.dia.toFixed();
            row[3].textContent = item.pulse.toFixed();
            row[4].textContent = item.comment;
            template.parentNode?.appendChild(clone);
        }
    }
}
