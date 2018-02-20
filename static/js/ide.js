var show_hint_play = false;

var menu = {
    id: "main_menu",
    view:"menu",
    autowidth: true,
    data: [
        {value: "Файл", submenu:[
             {value: "Сохранить все (Shift+Ctrl+S)", icon:"save"},
                "Закрыть все", { $template:"Separator" },
                {value: "Импорт", icon:"upload"},
                {value: "Экспорт", icon:"download"},
                { $template:"Separator" },
                {value: "Выход", icon:"window-close"}]},
        {value: "Правка", submenu:["Отменить", "Вернуть", { $template:"Separator" },
            "Выделить все", "Выделить идентификатор", { $template:"Separator" },
            "Перейти к строке", "Удалить строку", { $template:"Separator" },
            "Копировать строки выше", "Копировать строки ниже",
            "Переместить строки выше", "Переместить строки ниже"]},
        {value: "Текст", submenu:["Форматировать", { $template:"Separator" },
            "Автозаполнение", "Переход к определению", "Показать информацию выражения",
            "Найти места использования", { $template:"Separator" }, "Включить комментарий"]},
        {value: "Найти", submenu:["Найти в файле", { $template:"Separator" },
            "Найти в файлах", { $template:"Separator" }, "Найти", "Найти следующий",
            "Найти предыдущий", "Заменить", "Заменить все"]},
        {value: "Запустить", submenu:["Построить", "Построить и выполнить", { $template:"Separator" },
            "Тест", { $template:"Separator" }, "go get", "go install", { $template:"Separator" },
                "go vet"]},
        {value: "Настройки"},
        {value: "Помощь", submenu:["IDE документация", "Горячие клавиши", "О программе"]},
        {id: "play", value: "", icon: "play", on:{
            onMouseMove:function (ev) {
                console.log(ev);
            }
            } }
    ],
    on:{
        onMenuItemClick:function(id){
            webix.message("Click: "+this.getMenuItem(id).value);
        },
        onMouseMove:function (id, e, node) {
            if(!show_hint_play && id=="play"){
                show_hint_play = true;
                $$("menu_tooltip").show({id:"tooltip", value:"Построить и выполнить"}, {x:e.pageX-30, y:30});
                console.log(node);
                console.log(e);
            }

        },
        onMouseOut:function (e) {
            if(show_hint_play) {
                show_hint_play = false;
                $$("menu_tooltip").hide();
            }
        }
    },
    type:{
        subsign:true
    }
};

var activUser = {
    view: "icon",
    icon: "user-circle"
};

var treeFiles = {
    view:"tree",
    width: 250,
    data: [
        {id:"root", value:"Cars", open:true, data:[
                { id:"1", open:true, value:"Toyota", data:[
                        { id:"1.1", value:"Avalon" },
                        { id:"1.2", value:"Corolla" },
                        { id:"1.3", value:"Camry" }
                    ]},
                { id:"2", open:true, value:"Skoda", data:[
                        { id:"2.1", value:"Octavia" },
                        { id:"2.2", value:"Superb" }
                    ]}
            ]}
    ]
};

var tabsPanel = {
    view:"menu",
    autowidth: true,
    data: [
        "Результат",
        "Поиск",
        "Сообщения"
    ]
};


var closeMenu = {
    view: "icon",
    icon: "bars"
};



webix.ready(function(){
    webix.ui({
        id: "menu_tooltip",
        view:"tooltip",
        template:"#value#",
        height:40
    });

    webix.ui({
        rows: [
            // { view:"toolbar", elements: menu },
             { view:"toolbar", elements:[ menu, activUser ] },
            { cols:[
                {rows: [
                    {view: "toolbar", id:"toolbar", elements:[
                        {
                            view: "label",
                            label: "Files"
                        },
                        {
                            view: "icon" , icon: "fas fa-bars fa-xs",
                            click: function(){
                                if( $$("menu").config.hidden){
                                    $$("menu").show();
                                } else {
                                    $$("menu").hide();
                                }
                            }
                        }
                    ]},
                    treeFiles
                ]},
                {view:"resizer"},
                { rows: [
                    { cols:[
                        // {view:"tabview", borderless:true, cells:[{ header:"Empty", body:{ } }]},
                        {view:"tamplate", value: ""},
                        {view:"resizer"},
                        {rows: [
                            {view: "toolbar", id:"toolbar", elements:[
                                {
                                    view: "label",
                                    label: "Схема"
                                },
                                {
                                    view: "icon", icon: "bars",
                                    click: function(){
                                        if( $$("menu").config.hidden){
                                            $$("menu").show();
                                        } else {
                                            $$("menu").hide();
                                        }
                                    }
                                }
                            ]},
                            { template:"Колонка 2", width: 250 }
                        ]}
                    ]},
                    {view:"resizer"},
                    { rows: [
                        {view:"toolbar", elements:[tabsPanel, closeMenu]},
                        { template:"Строка 3", height: 150}
                    ]}
                ]}
            ]},
            { template:"Строка 3", autoheight: true }
        ]
    });
});