function post(params) {
    let form = document.createElement("form");
    form.setAttribute("method", "POST");

    for(let key in params) {
        let hiddenField = document.createElement("input");
        hiddenField.setAttribute("type", "hidden");
        hiddenField.setAttribute("name", key);
        hiddenField.setAttribute("value", params[key]);
        form.appendChild(hiddenField);
    }
    document.body.appendChild(form);
    form.submit();
}

/* e: td element */
function update(e) {
    let inputBox = document.createElement("INPUT");
    inputBox.onkeypress = function (event) {
        switch (event.key) {
            case "Enter":
                let inputBox = event.srcElement;
                let newValue = inputBox.value;
                if (newValue === "")
                    newValue = "NULL";
                else
                    newValue = "\'" + newValue + "\'";

                let trNode = inputBox.parentElement.parentElement;
                let colId = Array.prototype.indexOf.call(trNode.children, inputBox.parentNode);
                let colVals = [];
                for (let td of trNode.children)
                    colVals.push(td.innerText);

                let trColInfos = document.getElementById("trColInfo");
                let colNames = [];
                for (let th of trColInfos.children)
                    colNames.push(th.innerHTML.split('<br>')[0]);

                let whereClause = "";
                for (let i = 0; i < colVals.length; i++) {
                    let val;
                    if(colVals[i] === "")
                        val = " IS NULL";
                    else
                        val = "=\'" + colVals[i] + "\'";
                    whereClause += colNames[i] + val + " AND ";
                }
                whereClause = whereClause.slice(0, -5);

                let tbName = window.location.href.split("/").reverse()[0];

                let sql = "UPDATE " + tbName + " SET " + colNames[colId] + "=" + newValue + " WHERE "
                    + whereClause + ";";

                post({"cmd": sql, "exec": true});

                inputBox.parentElement.removeChild(inputBox);
                break;
        }
    };
    e.appendChild(inputBox)
}