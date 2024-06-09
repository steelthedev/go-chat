# Chat App with Go

A chat application made using Go for the backend logics, algorithms, server building. HTMX for the frontend.

## Debugging/Deployment

> [!NOTE]  
> Running the main file (in this case, it is `main.go`) using the command `go run main.go` would not run the script. Debugger have to use the command `go run .` to debug or run the script. This isn't understood why the `go run main.go` command's not working. Possible theory is that, the whole codebase (./handlers/, ./static/, ./views/, ./go.mod, ./go.sum, ./main.go, ./message.go and ./websocket.go) needs to be run at once.
