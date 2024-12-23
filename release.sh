cd server

echo "Building server"
go build -o ../vscode_client/lsp-server main.go


cd ../vscode_client

echo "Building and publishing VSCode client"
npm run publish

# if personal access token is expired, you can create a new one and replace the old one.
# 1. create pat on https://dev.azure.com/<pulisher>/_usersSettings/tokens
# 2. vsce login <publisher>
