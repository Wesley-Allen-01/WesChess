const chess = new Chess();

const roomID = window.location.pathname.split("/")[2];
document.getElementById("roomID").innerText = roomID;

const user_color = document.getElementById('board').dataset.playerColor;

if (user_color === "w") {
    document.getElementById("playerColor").innerText = "Player Color: White";
} else {
    document.getElementById("playerColor").innerText = "Player Color: Black";
}

const fullColor = {
    w: "white",
    b: "black"
};

function updateMoveStatus() {
    if (user_color === chess.turn()) {
        document.getElementById("whoseMove").innerText = "Your Move";
    } else {
        document.getElementById("whoseMove").innerText = "Opponent's Move";
    }
}

updateMoveStatus();

const ws = new WebSocket(`ws://localhost:8080/ws/${roomID}`);

ws.onopen = () => {
    console.log("WebSocket connected");
    document.getElementById("status").innerText = "Connected to game.";
};

ws.onmessage = (event) => {
    const move = JSON.parse(event.data); 
    console.log("Received move:", move);

    
    const result = chess.move({ from: move.source, to: move.target, promotion: 'q' });
    if (result) {
        board.move(`${move.source}-${move.target}`); 
        if (chess.in_checkmate()) {
            document.getElementById("status").innerText = "Checkmate!";
            if (user_color === chess.turn()) {
                alert("You suck!");
            } else {
                alert("You are the 🐐!");
            }
        } else {
            updateMoveStatus();
        }
    } else {
        console.error("Received invalid move from server");
    }
};

ws.onclose = () => {
    console.log("WebSocket disconnected");
    document.getElementById("status").innerText = "Disconnected from game.";
};

ws.onerror = (error) => {
    console.error("WebSocket error:", error);
    document.getElementById("status").innerText = "WebSocket error.";
};


function highlightLegalMoves(moves) {
    moves.forEach(move => {
        const square = document.querySelector(`.square-${move.to}`);
        if (square) {
            square.style.backgroundColor = "rgba(0, 255, 0, 0.5)"; // Green for legal moves
        }
    })
}

function clearHighlights() {
    document.querySelectorAll(".square-55d63").forEach(square => {
        square.style.backgroundColor = ""; // Reset square backgrounds
    });
}

function handleDragStart(source, piece, position, orientation) {
    console.log(`Drag started: ${source}`);
    if (chess.turn() !== user_color) {
        return false;
    }

    const legal_moves = chess.moves({ square: source });
    if (legal_moves.length === 0) {
        return false;
    }
    console.log("we gonna call the highlighta")
    highlightLegalMoves(legal_moves);
    return true;
}

function handleDrop(source, target) {
    console.log(`Attempting move: ${source} -> ${target}`)
    if (chess.turn() !== user_color) {
        console.log("Not your turn");
        return "snapback"; 
    }

    const move = chess.move({ from: source, to: target, promotion: 'q' });
    if (move === null) {
        console.log("Invalid move, snapping back");
        return "snapback"; 
    }

    ws.send(JSON.stringify({ source, target }));

    clearHighlights();

    if (chess.in_checkmate()) {
        document.getElementById("status").innerText = "Checkmate!";
        if (user_color === chess.turn()) {
            alert("You suck!");
        } else {
            alert("You are the 🐐!");
        }
    }
    else {
        updateMoveStatus();
    }
    console.log("Move sent to server:", { source, target });
}




var config = {
    position: "start",
    draggable: true,
    pieceTheme: '/static/img/chesspieces/wikipedia/{piece}.png',
    onDrop: handleDrop,
    onDragStart: handleDragStart,
    orientation: fullColor[user_color],
};


const board = Chessboard("board", config);
