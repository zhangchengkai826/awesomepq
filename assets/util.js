function createSelector(arr) {
    selector = document.createElement("SELECT");
    for (var val of arr) {
        option = document.createElement("OPTION");
        option.setAttribute("value", val);
        option.innerText = val;
        selector.appendChild(option);
    }
    return selector;
}

function prepareInputs() {
    var varOptionsDiv = document.getElementById("var options");
    while(varOptionsDiv.firstChild) {
        varOptionsDiv.removeChild(varOptionsDiv.firstChild);
    }

    var cmdCategorySelector = document.getElementById("cmd category selector");
    var cmdCategoryId = cmdCategorySelector.selectedIndex;
    var cmdCategoryName = cmdCategorySelector.options[cmdCategoryId].getAttribute("value");

    switch(cmdCategoryName) {
        case "CREATE TABLE":
            var types = ['(type)', 'bigint', 'bigserial', 'bit', 'bit varying', 'boolean',
            'box', 'bytea', 'character', 'character varying', 'cidr', 'circle', 'date',
            'double precision', 'inet', 'integer', 'interval', 'json', 'jsonb', 'line',
            'lseg', 'macaddr', 'money', 'numeric', 'path', 'pg_lsn', 'point', 'polygon',
            'real', 'smallint', 'smallserial', 'serial', 'text', 'time',
            'timestamp', 'tsquery', 'tsvector', 'txid_snapshot', 'uuid', 'xml'];
            var typeSelector = createSelector(types);
            varOptionsDiv.appendChild(typeSelector);

            var nameInput = document.createElement("INPUT");
            nameInput.setAttribute("value", "(name)");
            varOptionsDiv.appendChild(nameInput);

            varOptionsDiv.appendChild(createSelector(['(options)', 'NOT NULL', 'DEFAULT']));
            varOptionsDiv.appendChild(document.createElement("BR"));
            var addButton = document.createElement("BUTTON");
            addButton.setAttribute("type", "button");
            addButton.innerText = "Add";
            varOptionsDiv.appendChild(addButton);
            break;
        default:
            break;
    }
}