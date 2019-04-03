function createSelector(arr) {
    let selector = document.createElement("SELECT");
    for (let val of arr) {
        let option = document.createElement("OPTION");
        option.setAttribute("value", val);
        option.innerText = val;
        selector.appendChild(option);
    }
    return selector;
}

function createColumnDefinitionDiv() {
    let columnDefinitionDiv = document.createElement("DIV");

    let types = ['(type)', 'bigint', 'bigserial', 'bit', 'bit varying', 'boolean',
        'box', 'bytea', 'character', 'character varying', 'cidr', 'circle', 'date',
        'double precision', 'inet', 'integer', 'interval', 'json', 'jsonb', 'line',
        'lseg', 'macaddr', 'money', 'numeric', 'path', 'pg_lsn', 'point', 'polygon',
        'real', 'smallint', 'smallserial', 'serial', 'text', 'time',
        'timestamp', 'tsquery', 'tsvector', 'txid_snapshot', 'uuid', 'xml'];
    let typeSelector = createSelector(types);
    columnDefinitionDiv.appendChild(typeSelector);

    let nameInput = document.createElement("INPUT");
    nameInput.setAttribute("value", "(name)");
    columnDefinitionDiv.appendChild(nameInput);

    columnDefinitionDiv.appendChild(createSelector(['(options)', 'NOT NULL', 'DEFAULT']));

    let additionalInput = document.createElement("INPUT");
    additionalInput.setAttribute("value", "(additional value)");
    columnDefinitionDiv.appendChild(additionalInput);

    let addButton = document.createElement("BUTTON");
    addButton.setAttribute("type", "button");
    addButton.setAttribute("onclick", "addColumnDefinitionDiv()");
    addButton.innerText = "Add";
    columnDefinitionDiv.appendChild(addButton);

    return columnDefinitionDiv;
}

function addColumnDefinitionDiv() {
    let varOptionsDiv = document.getElementById("var options");
    varOptionsDiv.appendChild(createColumnDefinitionDiv());
}

function prepareInputs() {
    let varOptionsDiv = document.getElementById("var options");
    while(varOptionsDiv.firstChild) {
        varOptionsDiv.removeChild(varOptionsDiv.firstChild);
    }

    let cmdCategorySelector = document.getElementById("cmd category selector");
    let cmdCategoryId = cmdCategorySelector.selectedIndex;
    let cmdCategoryName = cmdCategorySelector.options[cmdCategoryId].getAttribute("value");

    switch(cmdCategoryName) {
        case "CREATE TABLE":
            addColumnDefinitionDiv();
            break;
        default:
            break;
    }
}