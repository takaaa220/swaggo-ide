cd server

echo "Building server"
go build -o ../vscode_client/lsp-server main.go


cd ../vscode_client

echo "Building and publishing VSCode client"
pnpm publish
