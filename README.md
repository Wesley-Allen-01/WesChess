# WesChess


# Directory Structure with Annotations

## **1. `backend/`**
Contains the Go backend code responsible for API functionality, database interaction, and static file serving.

### **a. `cmd/`**
Entry point for the Go application.
- **`server/`**:  
  Houses the main server logic for running the backend application.
  - **`main.go`**:  
    The main entry point of the backend application. Initializes the Gin server, sets up routes (e.g., `/register`), and starts serving API requests and static files.

### **b. `data/`**
Directory for storing backend-related persistent data.
- **`chess_game.db`**:  
  SQLite database file used to manage user data, such as credentials and other application data.

### **c. `internal/`**
Internal package housing core backend functionality.
- **`db/`**:  
  - **`connection.go`**:  
    Handles database connections and initialization, allowing other parts of the backend to query the SQLite database.
- **`game/`**:  
  Placeholder for chess game logic. This could handle validating moves, managing game states, and synchronizing gameplay.
- **`handlers/`**:  
  Contains HTTP handler functions for different backend routes.
  - **`auth.go`**:  
    Implements user authentication logic, such as login and registration, using the database.

### **d. `Dockerfile`**
Dockerfile for building and running the backend service. Configures the Go environment, builds the server, and exposes the backend port.

### **e. `go.mod`**
Go module file listing dependencies for the project.

### **f. `go.sum`**
Dependency checksum file ensuring consistent builds by locking dependency versions.

---

## **2. `frontend/`**
Contains the static files and JavaScript assets used for the application's frontend.

### **a. `node_modules/`**
Directory containing npm dependencies for JavaScript modules required by the frontend.

### **b. `public/`**
Folder for all static assets and HTML files served by the backend or a static file server.
- **`chessboard/`**:  
  Static assets from the `chessboard.js` library.
  - **`chessboard-1.0.0.css`**:  
    Stylesheet for rendering the chessboard.
  - **`chessboard-1.0.0.js`**:  
    JavaScript for initializing and interacting with the chessboard.
- **`img/`**:  
  - **`chesspieces/`**:  
    Directory for chess piece images.
    - **`wikipedia/`**:  
      Subdirectory containing PNG images of all 12 chess pieces (white/black pawns, rooks, knights, bishops, queens, and kings).
- **`favicon.ico`**:  
  Small icon displayed in the browser tab. 
- **`index.html`**:  
  Main page for rendering the chessboard using `chessboard.js`. Includes JavaScript for interactive gameplay.
- **`register.html`**:  
  Registration page with a form for users to create accounts.

### **c. `scripts/`**
(Placeholder for JavaScript code specific to frontend logic.)

### **d. `styles/`**
(Placeholder for CSS files to style application.)

### **e. `package-lock.json`**
Dependency lock file for npm, ensuring consistent frontend builds.

### **f. `package.json`**
Lists frontend dependencies and scripts for managing the project with npm.

---

## **3. `.gitignore`**
Specifies files and directories to be ignored by Git, such as `node_modules/`, `*.db` files, or Docker build artifacts.

---

## **4. `docker-compose.yml`**
Defines the configuration for running multiple services (frontend and backend) in Docker. Coordinates their ports, volumes, and network settings.

---

## **5. `p.md`**
(Placeholder for a file. Could be project-specific notes, documentation, or planning details.)

---

## **6. `package-lock.json`**
(Same as above, likely duplicated.)

---

## **7. `README.md`**
Main documentation file for project. Describes the purpose of the application, setup instructions, and additional project details.
