function login_api() {
    if(!$('input[name=email]')[0].checkValidity() || !$('input[name=password]')[0].checkValidity()){
        alert("Please check the input data again.");
        return ;
    }

    var form = new FormData();
    form.append("email", $('input[name=email]').val());
    form.append("password", $('input[name=password]').val());
    
    var settings = {
        "url": location.protocol+"//"+location.host+"/api/api/login",
        "method": "POST",
        "timeout": 0,
        "processData": false,
        "mimeType": "multipart/form-data",
        "contentType": false,
        "data": form,
        xhrFields: {
            withCredentials: true
        },
    };

    $.ajax(settings).done(function(data, status, xhr) {
        msg = JSON.parse(data);
        if (msg == "OK"){
            location.href = "/";
        } else {
            alert("fail");
        }
    }).fail(function (data, textStatus, errorThrown) {
        if(data){
            msg = JSON.parse(data.responseText).error;
            $('#error')[0].innerHTML = "<p style='color:red;'>"+ msg +"</p>";
        } else {
            alert("fail");
        }
    });
}

function token_api() {
    var settings = {
        "url": location.protocol+"//"+location.host+"/api/user/token",
        "method": "GET",
        "timeout": 0,
        xhrFields: {
            withCredentials: true
        },
    };

    $.ajax(settings).done(function (response, textStatus, xhr) {
        let token = JSON.parse(response);
        $('#token_table')[0].innerHTML = `
        <table>
            <tr><th>Access Token</th><td class="col_a">${token[0]}</td></tr>
            <tr><th>Refresh Token</th><td class="col_a">${token[1]}</td></tr>
        </table>
        `;
    }).fail(function (data, textStatus, errorThrown) {
        if (data.status == "401" || data.status == "403"){
            alert("로그인 해주세요.");
            location.href = "/form/login.html";
        }
    });
}

function admin_api() {
    var settings = {
        "url": location.protocol+"//"+location.host+"/api/user/admin",
        "method": "GET",
        "timeout": 0,
        xhrFields: {
            withCredentials: true
        },
    };

    $.ajax(settings).done(function (response, textStatus, xhr) {
        let users=JSON.parse(response);
        let result = "<table><tr> <th>ID</th> <th>Email</th> <th>Permission</th> </tr>";
        for (i = 0; i < users.length; i++) {
            result += `
            <tr>
                <td>${users[i]['ID']}</td>
                <td>${users[i]['Email']}</td>
                <td>${users[i]['Permission']}</td>
            </tr>`;
        }
        result += "</table>";
        $('#users_table')[0].innerHTML = result;
    }).fail(function (data, textStatus, errorThrown) {
        if (data.status == "401" || data.status == "403"){
            alert("로그인 해주세요.");
            location.href = "/form/login.html";
        }
    });
}

function logout_api() {
    var settings = {
        "url": location.protocol+"//"+location.host+"/api/user/logout",
        "method": "GET",
        "timeout": 0,
        xhrFields: {
            withCredentials: true
        },
    };

    $.ajax(settings).done(function (response, textStatus, xhr) {
        alert("로그아웃 되었습니다.");
        location.href = "/form/login.html";
    }).fail(function (data, textStatus, errorThrown) {
        if (data.status == "401" || data.status == "403"){
            alert("로그인 해주세요.");
            location.href = "/form/login.html";
        }
    });
}
