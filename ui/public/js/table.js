"use strict";
async function loadData() {
    const userIdField = document.querySelector("#userId");
    if (!(userIdField instanceof HTMLInputElement)) {
        throw new Error("User ID form is not present");
    }
    const userIdValue = userIdField.value;
    const data = await getDataFromApi(userIdValue);
    if (data.length == 0) {
        errorToast(`No data present in db for this userId: ${userIdValue}`);
    }
    loadDataIntoTable(data);
}
async function getDataFromApi(userIdValue) {
    const formDataJsonString = JSON.stringify({ userId: userIdValue });
    const fetchOptions = {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: formDataJsonString
    };
    const response = await fetch("api/listEntries", fetchOptions);
    const json = await response.json();
    const entries = [];
    for (const entry of json) {
        entries.push(await entry);
    }
    return entries;
}
function loadDataIntoTable(data) {
    const template = document.querySelector("#row");
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
