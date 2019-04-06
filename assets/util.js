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

function getTblName() {
    return window.location.href.split("/").reverse()[0];
}

function getColNames() {
    let trColInfos = document.getElementById("trColInfo");
    let colNames = [];
    for(let i = 0; i < trColInfos.children.length-1; i++)
        colNames.push(trColInfos.children[i].innerHTML.split('<br>')[0]);
    return colNames;
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
                for (let i = 0; i < trNode.children.length-1; i++)
                    colVals.push(trNode.children[i].innerText);

                let colNames = getColNames();

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

                let sql = "UPDATE " + getTblName() + " SET " + colNames[colId] + "=" + newValue + " WHERE "
                    + whereClause + ";";

                post({"cmd": sql, "exec": true});

                inputBox.parentElement.removeChild(inputBox);
                break;
        }
    };
    e.appendChild(inputBox)
}

/* e input element */
function showInsertRow(e) {
    e.setAttribute("value", "OK");
    e.setAttribute("onclick", "insert(this)");
    let tr = document.createElement("TR");
    tr.setAttribute("id", "insertPlaceHolder");
    let trColInfos = document.getElementById("trColInfo");
    for(let i = 0; i < trColInfos.children.length-1; i++) {
        let inputBox = document.createElement("INPUT");
        inputBox.setAttribute("style",
            "width:" + trColInfos.children[i].offsetWidth * 0.8 + "px;");
        let td = document.createElement("TD");
        td.appendChild(inputBox);
        tr.appendChild(td);
    }
    let datatbl = document.getElementById("datatbl");
    datatbl.appendChild(tr);
}

function insert(e) {
    let insertPlaceHolder = document.getElementById("insertPlaceHolder");

    let insertVals = [];
    for(let i = 0; i < insertPlaceHolder.children.length; i++)
        insertVals.push(insertPlaceHolder.children[i].children[0].value);
    let sql = "INSERT INTO " + getTblName() + " VALUES (";
    for(let val of insertVals) {
        if(val === "")
            sql += "NULL,";
        else
            sql += "\'" + val + "\',";
    }
    sql = sql.slice(0, -1);
    sql += ");";

    post({"cmd": sql, "exec": true});

    insertPlaceHolder.parentElement.removeChild(insertPlaceHolder);
    let insertBtn = document.getElementById("insertBtn");
    insertBtn.setAttribute(
        "value", "Insert");
    insertBtn.setAttribute(
        "onclick", "showInsertRow(this)");
}

function del(e) {
    let tr = e.parentElement.parentElement;
    let curRowVals = [];
    for(let i = 0; i < tr.children.length-1; i++)
        curRowVals.push(tr.children[i].innerHTML);

    let colNames = getColNames();

    let sql = "DELETE FROM " + getTblName() + " WHERE ";
    for(let i = 0; i < colNames.length; i++) {
        sql += colNames[i];
        if(curRowVals[i] === "")
            sql += " IS NULL";
        else
            sql += "=\'" + curRowVals[i] + "\'";
        sql += " AND ";
    }
    sql = sql.slice(0, -5);
    sql += ";";

    post({"cmd": sql, "exec": true});
}

function showAddColInputs(e) {
    let th = e.parentElement;
    th.removeChild(e);

    let nameInput = document.createElement("INPUT");
    nameInput.setAttribute("id", "newColName");
    nameInput.setAttribute("value", "(new column name)");
    th.appendChild(nameInput);

    th.appendChild(document.createElement("BR"));

    let typeInput = document.createElement("INPUT");
    typeInput.setAttribute("id", "newColType");
    typeInput.setAttribute("value", "(new column type)");
    th.appendChild(typeInput);

    th.appendChild(document.createElement("BR"));

    let btOk = document.createElement("INPUT");
    btOk.setAttribute("type", "button");
    btOk.setAttribute("onclick", "addCol(this)");
    btOk.setAttribute("value", "Add");
    th.appendChild(btOk);
}

function addCol(e) {
    let newColName = document.getElementById("newColName").value;
    let newColType = document.getElementById("newColType").value;
    let sql = "ALTER TABLE " + getTblName() +
        " ADD COLUMN " + newColName + " " + newColType + ";";

    post({"cmd": sql, "exec": true});

    let th = e.parentElement;
    while(th.firstChild) 
        th.removeChild(th.firstChild);
    let btPlus = document.createElement("INPUT");
    btPlus.setAttribute("type", "button");
    btPlus.setAttribute("onclick", "showAddColInputs(this)");
    btPlus.setAttribute("value", "+");
    th.appendChild(btPlus);
}