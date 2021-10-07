"use strict";

window.addEventListener('load', function () {
    console.log('app started');

    // search input
    init_search_app();



}, false);


function init_search_app() {
    const searchInput = document.getElementById('search-form-input');

    searchInput.addEventListener("keyup", e => {
        let searchString = e.target.value;
        console.log(searchString);

        

        const filteredCharacters = hpCharacters.filter(character => {
            return (
                character.name.includes(searchString) ||
                character.house.includes(searchString)
            );
        });
        displayCharacters(filteredCharacters);
    });

}

function search_listener() {
    var postSearch = null;
    var lockMicroTime = 0;

    jQuery(".js-search-term").bind("keyup change", function()
    {
        //новый локтайм и внутреннее время
        lockMicroTime = new Date().getTime();
        var innerTime = lockMicroTime;
        var target = this;

        //ждём 200 милисекунд и запускаем запрос с проверками
        //@TODO переписать на замыкание (IE8 bug =( )
        setTimeout(function()
        {
            if(lockMicroTime != innerTime)
                return false;

            var idW = jQuery(target).attr('idW');
            var term = jQuery.trim( jQuery(target).val() );
            var result = jQuery(".js-search-result");
            var colSpan = ( typeof idW === 'undefined') ? 5 : 6;
            result.parent().removeClass('hide');

            if( /^[\wА-ЯЁа-яё\s.-]{3,30}$/.test(term) !== true )
            {
                result.html('<tr><td colspan="'+colSpan+'">Такое условие поиска недопустимо.</td></tr>');
                return false;
            }

            //обрубаем лишние запросы
            if (postSearch != null)
                postSearch.abort();

            result.html('<tr><td colspan="'+colSpan+'">'+loadImage+'</td></tr>');

            postSearch = jQuery.post('/ajax/search/',
                {
                    'idW': ( typeof idW == 'undefined') ? 0 : idW,
                    'term': term,
                    'format': 'html'
                },
                function(data){
                    postSearch = null;
                    result.html(data);
                },
                'html');

        },200);

    });
    // todo search event function
    // todo search call search-API
    // todo cancel request before sending new
}