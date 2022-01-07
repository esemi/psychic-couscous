"use strict";

const searchApiEndpoint = './mock-api/search.json';
const nodeInfoApiEndpoint = './mock-api/node-info.json';

window.addEventListener('load', function () {
    console.log('app started');

    // init graph app
    let graphApp = new GraphEngine(nodeInfoApiEndpoint);

    // init search app
    init_search_app(graphApp, searchApiEndpoint);

}, false);


function init_search_app(graph, searchApi) {
    const minSearchQueryLength = 3;
    const nodeIdAttr = 'data-id';
    const searchInput = document.getElementById('search-form-input');
    const searchResultContainer = document.getElementById('search-form-results');
    const searchInputClear = document.getElementById('search-form-clear');

    let activeRequestController;
    let debounce;

    let search_handler = (e) => {
        let searchString = e.target.value.trim().toLowerCase();

        if (!!searchString) {
            searchResultContainer.hidden = true;
            graph.showPreInitState();
        }

        if (searchString.length < minSearchQueryLength) {
            console.log('skip search event', searchString)
            return false
        }

        if (!!debounce) {
            clearTimeout(debounce);
        }

        debounce = setTimeout(() => {
            let url = searchApi + '?' + new URLSearchParams({search: searchString,});
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

                    if (!data.hasOwnProperty('results') || !data.results.length) {
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
        graph.init();
    }


    let select_handler = (e) => {
        console.log('select root node handler', e.target);
        searchResultContainer.hidden = true;
        searchInput.value = e.target.innerText;
        graph.addRootNode(e.target.getAttribute(nodeIdAttr));
    }

    searchInput.addEventListener("keyup", search_handler, false);
    searchInputClear.addEventListener("click", clear_handler, false);
    searchResultContainer.addEventListener("click", select_handler, false);
}


class GraphEngine {
    rootNodeId;
    graph;
    nodeInfoApiEndpoint;
    activeRequestController;

    static graphContainerId = 'graph-container';
    static stateNotInitId = 'graph-state-not-init';
    static stateNotFoundId = 'graph-state-not-found';

    #containerStateNotInit;
    #containerStateNotFound;
    #containerForGraph;

    constructor(nodeInfoApiEndpoint) {
        this.nodeInfoApiEndpoint = nodeInfoApiEndpoint;

        this.#containerStateNotInit = document.getElementById(GraphEngine.stateNotInitId);
        this.#containerStateNotFound = document.getElementById(GraphEngine.stateNotFoundId);
        this.#containerForGraph = document.getElementById(GraphEngine.graphContainerId);

        this.init();
    }

    addRootNode(root_node_id) {
        // load and display root node with relations
        console.log('init graph engine for ', root_node_id)
        this.rootNodeId = root_node_id;

        this._fetchNodeRelations(this.rootNodeId).then((nodeInfo) => {
            this.appendRelations(nodeInfo.nodes, nodeInfo.edges);
            this.showGraph();
        });
    }

    _fetchNodeRelations(nodeId) {
        // fetch all relations by nodeId
        let url = this.nodeInfoApiEndpoint + '?' + new URLSearchParams({node_id: nodeId,});
        console.log('send node info request', nodeId, url);

        let nodeInfoRequest = new Request(url, {
            headers: {'Accept': 'application/json'},
            method: 'GET',
            signal: this.activeRequestController.signal
        });

        return fetch(nodeInfoRequest)
            .then((response) => {
                if (!response.ok) {
                    throw new Error("Request error, status = " + response.status);
                }
                return response.json()
            })
            .then((responseJson) => {
                console.log('fetch node info response', responseJson)
                return {
                    nodes: responseJson.nodes,
                    edges: responseJson.relations,
                }
            });
    }

    appendRelations(nodes, edges) {
        // append nodes and edges to graph
        this.graph.addNodes(
            nodes.map(node => {
                return {
                    id: node.id,
                    attributes: {
                        text: node.name,
                    },
                    data: {
                        name: node.name,
                        type: node.type,
                    }
                }
            })
        );
        this.graph.addEdges(
            edges.map(edge => {
                return {
                    id: `${edge.from_id}:${edge.to_id}`,
                    source: edge.from_id,
                    target: edge.to_id,
                    attributes: {
                        text: `${edge.type}[${edge.name}]`
                    },
                    data: {
                        name: edge.name,
                        type: edge.type,
                    }
                }
            })
        );

        console.log('debug: edges=', this.graph.getNodes().size, '; nodes=', this.graph.getEdges().size)
    }

    init() {
        if (!!this.activeRequestController) {
            this.activeRequestController.abort();
        }
        this.activeRequestController = new AbortController();

        // @see <https://doc.linkurio.us/ogma/3.2.0/quickstart.html>
        this.graph = new Ogma({
            container: this.#containerForGraph.id,
            options: {
                directedEdges: true,
                edgesAlwaysCurvy: true,
            }
        });
        this.graph.styles.addNodeRule({
            color: this.graph.rules.map({
                field: 'type',
                values: {
                    'movie': 'green',
                    'name': 'yellow',
                },
                fallback: 'black'
            })
        });

        //
        this.graph.events.onClick(event => {
            const target = event.target;
            if (!target || !target.isNode) {
                return
            }
            console.log(`Node ${target.getId()} was clicked`);
        });

        this.showPreInitState();
    }

    drawGraph() {
        // apply graph layout and draw results
        this.graph.layouts.force({
            incremental: true,
        })
    }

    showGraph() {
        this.hideAll();
        this.#containerForGraph.hidden = false;
        this.drawGraph();
    }

    showPreInitState() {
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