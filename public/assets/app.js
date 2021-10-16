"use strict";

//q: как по человечески модуль инпортунть?
//q: как быть с console.log на проде?

const searchApiEndpoint = './mock-api/search.json';
const neighboursApiEndpoint = './mock-api/neighbours.json';

window.addEventListener('load', function () {
    console.log('app started');

    //graph init
    const graphContainerId = 'graph-container';
    let graph = new GraphEngine(neighboursApiEndpoint, graphContainerId);

    // search input
    init_search_app(graph);


}, false);


function init_search_app(graph) {
    const minSearchQueryLength = 3;
    const nodeIdAttr = 'data-id';
    const searchInput = document.getElementById('search-form-input');
    const searchResultContainer = document.getElementById('search-form-results');
    const searchInputClear = document.getElementById('search-form-clear');

    let activeRequestController;
    let debounce;

    let search_handler = (e) => {
        let searchString = e.target.value.trim().toLowerCase();
        console.log(searchString);

        if (!!searchString) {
            searchResultContainer.hidden = true;
            graph.displayNotInitState();
        }

        if (searchString.length < minSearchQueryLength) {
            console.log('skip search event', searchString)
            return false
        }

        if (!!debounce) {
            clearTimeout(debounce);
        }

        debounce = setTimeout(() => {
            let url = searchApiEndpoint + '?' + new URLSearchParams({search: searchString,});
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

                    if (!data.results.length) {
                        searchResultContainer.hidden = true;
                        graph.displayNotFoundState(searchString);
                        return
                    }

                    searchResultContainer.innerHTML = '';
                    data.results.forEach((node) => {
                        console.log('display node = ', node);
                        let createdNode = document.createElement('ul');
                        createdNode.setAttribute(nodeIdAttr, node.id);
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
        graph.displayNotInitState();
    }


    let select_handler = (e) => {
        console.log('select root node handler', e.target);
        //todo lock search-input for new chars
        searchResultContainer.hidden = true;
        searchInput.value = e.target.innerText;
        graph.initRoot(e.target.getAttribute(nodeIdAttr));
    }

    searchInput.addEventListener("keyup", search_handler, false);
    searchInputClear.addEventListener("click", clear_handler, false);
    searchResultContainer.addEventListener("click", select_handler, false);
}


class GraphEngine {
    root_node_id;
    graph;
    api;

    #graphContainerId;
    #container;

    constructor(apiEndpoint, graphContainerId) {
        this.#graphContainerId = graphContainerId;
        this.#container = document.getElementById(this.#graphContainerId);
        this.api = apiEndpoint;
        this.displayNotInitState();
    }

    initRoot(root_node_id) {
        this.root_node_id = root_node_id;
        console.log('init graph engine for', root_node_id)

        let relations = this.fetchNodeRelations(root_node_id);

        // init graph
        this.graph = new sigma(this.#graphContainerId);

        // todo init graph-controls

        this.appendRelations(relations);
    }

    fetchNodeRelations() {
        //todo impl
        return {
            nodes: [],
            edges: []
        };
    }

    appendRelations(relations) {
        //todo impl
    }

    displayNotInitState() {
        this.#container.innerText = 'Use search to explore a couscous';
    }

    displayNotFoundState(searchString) {
        this.#container.innerText = 'No results for ' + searchString;
    }
}