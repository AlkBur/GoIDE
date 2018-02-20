function submit(){
    webix.ajax().post(window.conf.Context+"/login", $$("loginForm").getValues(), function(text, data, XmlHttpRequest){
        var dataJSON = data.json();;
        console.log(dataJSON);
        if (!dataJSON.result) {
            /*Ошибка авторизации*/
            webix.message({type:"error", text:"Username or password is Invalid"});
            return;
        }
        window.location.href = window.conf.Context+"/";
    });
}


webix.ready(function(){
    webix.ui({
        view:"form",
        id:"loginForm",
        width:300,
        elements:[
            { view:"text", label:"Username", name:"username"},
            { view:"text", type:"password", label:"Password", name:"password"},
            { margin:5, cols:[
                    { view:"button", value:"Login" , type:"form", click:submit}
            ]}
        ],
        on: {
            onSubmit: function (view, ev) {
              console.log(ev);
            }
        }
    });
   // $$("loginForm").submit(function (event) {
   //     /* stop form from submitting normally */
   //     event.preventDefault();
   //     webix.ajax().post(window.conf.Context+"/login", form.getValues(), function(text, data, xhr){
   //         if (!data.result) {
   //             /*Ошибка авторизации*/
   //             webix.message({type:"error", text:"Username or password is Invalid"});
   //             return;
   //         }
   //         window.location.href = window.conf.Context+"/";
   //     });
   // });





    // $('#loginForm').submit(function (event) {
    //     /* stop form from submitting normally */
    //     event.preventDefault();
    //
    //     console.log(window.conf.Context+"/login");
    //
    //     $.post( window.conf.Context+"/login", $( "#loginForm" ).serialize() )
    //         .done(function( data ) {
    //             if (!data.result) {
    //                 console.log(data);
    //
    //
    //                 $.notify({
    //                     title: '<strong>LOGIN ERROR</strong>',
    //                     icon: 'glyphicon glyphicon-warning-sign',
    //                     message: "Check your username and password and try again."
    //                 },{
    //                     type: 'pastel-error',
    //                     animate: {
    //                         enter: 'animated lightSpeedIn',
    //                         exit: 'animated lightSpeedOut'
    //                     },
    //                     template: '<div data-notify="container" class="col-xs-11 col-sm-3 alert alert-{0}" role="alert">' +
    //                     '<span data-notify="title">{1}</span>' +
    //                     '<span data-notify="message">{2}</span>' +
    //                     '</div>'
    //
    //                 });
    //                 // $("#loginBtn").val('[[.i18n.login_error]]');
    //                 // $('#loginBtn').addClass("denied");
    //                 //
    //                 // setInterval(function () {
    //                 //     $("#loginBtn").val('[[.i18n.login]]');
    //                 //     $('#loginBtn').removeClass("denied");
    //                 // }, 3000);
    //                 return;
    //             }
    //             window.location.href = window.conf.Context+"[[.conf.Context]]/";
    //         });
    // });
});