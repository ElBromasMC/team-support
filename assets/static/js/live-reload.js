(() => {
    let socket
    let retryInterval = 2000

    function connectWebSocket() {
        socket = new WebSocket('ws://localhost:8010')

        socket.onopen = _ => {
            console.log("Connected to the Live server!")
        }

        socket.onmessage = function (event) {
            if (event.data === 'reload') {
                window.location.reload()
            }
        }

        socket.onerror = error => {
            console.log("WebSocket Error: ", error);
        }

        socket.onclose = event => {
            console.log("Live server connection closed.", event)
            setTimeout(() => {
                console.log("Attempting to reconnect...")
                connectWebSocket()
            }, retryInterval)
        }
    }

    connectWebSocket();
})()