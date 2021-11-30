function login_api() {
    if(!$('input[name=email]')[0].checkValidity() || !$('input[name=password]')[0].checkValidity()){
        alert("Please check the input data again.");
        return ;
    }

    var form = new FormData();
    form.append("email", $('input[name=email]').val());
    form.append("password", $('input[name=password]').val());
    
    var settings = {
        "url": "http://localhost:3000/api/login",
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
            location.href = "/"
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
        "url": "http://localhost:3000/user/token",
        "method": "GET",
        "timeout": 0,
        xhrFields: {
            withCredentials: true
        },
    };

    $.ajax(settings).done(function (response) {
        console.log(response);
    });
}