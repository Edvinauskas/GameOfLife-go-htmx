# GameOfLife-go-htmx

## An interactive implementation of [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life), written in HTML/HTMX and Go. 

This interactive implementation of the game uses a websocket to communicate with the Go backend (which serves the HTML too).

To start the Go server, run the command below and go to localhost:3000:
```
go run .
```

The game grid is dynamically generated on load to fill your browser window.

Click any cell to make it 'Alive'. Once you're ready, click START to begin the simulation. You can pause at any time to add or remove cells, then click START to resume.

Honorable mentions: 

CSS Framework: [Bulma](https://bulma.io/)

Go libraries: [gorilla/websocket](https://github.com/gorilla/websocket)

HTMX extention: [htmx-ext-ws](https://htmx.org/extensions/ws/)




