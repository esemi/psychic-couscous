"use strict";

window.addEventListener('load', function () {
    console.log('app started');

    const searchApiEndpoint = '/mock-api/search.json';

    // search input
    init_search_app(searchApiEndpoint);


}, false);


function init_search_app(apiEndpoint) {
    const minSearchQueryLength = 3;
    const searchInput = document.getElementById('search-form-input');
    const searchResultContainer = document.getElementById('search-form-results');
    const searchInputClear = document.getElementById('search-form-clear');

    let activeRequestController;
    let debounce;

    let search_handler = (e) => {
        let searchString = e.target.value;
        console.log(searchString);

        if (!!searchString) {
            searchResultContainer.hidden = true;
        }

        if (searchString.length < minSearchQueryLength) {
            console.log('skip search event', searchString)
            return false
        }

        if (!!debounce) {
            clearTimeout(debounce);
        }

        debounce = setTimeout(() => {
            let url = apiEndpoint + '?' + new URLSearchParams({search: searchString,});
            console.log('send search request', searchString, url);

            //abort existing request
            if (!!activeRequestController) {
                activeRequestController.abort();
            }

            activeRequestController = new AbortController();

            let searchRequest = new Request(url, {
                headers: {'Accept': 'application/json'},
                method: 'GET',
                signal: activeRequestController.signal
            });

            fetch(searchRequest)
                .then((response) => {
                    if (!response.ok) {
                        throw new Error("Search error, status = " + response.status);
                    }
                    return response.json()
                })
                .then((data) => {
                    console.log('fetch search result', data)
                    //clear search result
                    searchResultContainer.innerHTML = '';
                    data.results.forEach((node) => {
                        console.log('display node = ', node);
                        let createdNode = document.createElement('ul');
                        createdNode.setAttribute('data-id', node.id);
                        createdNode.setAttribute('data-type', node.type);
                        createdNode.appendChild(document.createTextNode(node.name));
                        searchResultContainer.appendChild(createdNode);
                    });
                    searchResultContainer.hidden = false;
                });

        }, 350);
    };


    let clear_handler = () => {
        searchResultContainer.hidden = true;
        searchInput.value = '';
        searchInput.focus();
    }

    searchInput.addEventListener("keyup", search_handler, false);
    searchInputClear.addEventListener("click", clear_handler, false);
}
