"use strict";

//q: как по человечески модуль инпортунть?
//q: как быть с console.log на проде?

const searchApiEndpoint = './mock-api/search.json';
const neighboursApiEndpoint = './mock-api/neighbours.json';

window.addEventListener('load', function () {
    console.log('app started');

    // init graph app
    let graphApp = new GraphEngine(neighboursApiEndpoint);
    graphApp.graph.addNode(
        {id: 'nn000001', title: 'First test node'},
        {id: 'nn000002', title: 'Second test node'},
    );
    graphApp.appendRelations([]);
    // // init search app
    // init_search_app(graphApp);


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
            graph.showNotInitState();
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
                        graph.showNotFoundState(searchString);
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
        graph.reset();
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

    static graphContainerId = 'graph-container';
    static stateNotInitId = 'graph-state-not-init';
    static stateNotFoundId = 'graph-state-not-found';
    #containerStateNotInit;
    #containerStateNotFound;
    #containerForGraph;

    constructor(apiEndpoint) {
        this.api = apiEndpoint;

        this.#containerStateNotInit = document.getElementById(GraphEngine.stateNotInitId);
        this.#containerStateNotFound = document.getElementById(GraphEngine.stateNotFoundId);
        this.#containerForGraph = document.getElementById(GraphEngine.graphContainerId);

        // init graph
        this.graph = new sigma(GraphEngine.graphContainerId);
        // todo init graph-controls

        this.reset();
    }

    initRoot(root_node_id) {
        console.log('init graph engine for', root_node_id)
        this.root_node_id = root_node_id;

        let relations = this.fetchNodeRelations(this.root_node_id);
        this.appendRelations(relations);
        this.showGraph();
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
        //todo add new nodes
        //todo add new edges

        // fit nodes size by new degrees
        this.graph.nodes.forEach((node, ) => {
            // node size depends on its degree
            console.log('fit node size = ', node);
            console.log('fit node size = ', this.graph.degree(node));

            // atts.size = Math.sqrt(this.graph.degree(node)) / 2;
        });
    }

    showGraph() {
        this.hideAll();
        this.#containerForGraph.hidden = false;
        this.graph.refresh();
    }

    reset() {
        this.graph.clear();
        this.showNotInitState();
    }

    showNotInitState() {
        this.hideAll();
        this.#containerStateNotInit.hidden = false
    }

    showNotFoundState(searchString) {
        this.hideAll();
        this.#containerStateNotFound.innerText = 'No results for ' + searchString;
        this.#containerStateNotFound.hidden = false
    }

    hideAll() {
        this.#containerStateNotInit.hidden = true;
        this.#containerStateNotFound.hidden = true;
        this.#containerForGraph.hidden = true;
    }
}