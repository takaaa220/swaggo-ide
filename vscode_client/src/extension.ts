import {
  ExtensionContext,
  Uri,
  window as Window,
  workspace,
  commands,
} from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";
import { exec } from "child_process";
import { dirname, join } from "path";
import { readdir } from "fs/promises";

let client: LanguageClient;

export async function activate(context: ExtensionContext) {
  const serverModule = Uri.joinPath(context.extensionUri, "lsp-server").fsPath;

  const serverOptions: ServerOptions = {
    run: { command: serverModule, transport: TransportKind.stdio },
    debug: { command: serverModule, transport: TransportKind.stdio },
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: "go" }],
    synchronize: {
      fileEvents: [workspace.createFileSystemWatcher("**/*.go")],
    },
    outputChannel: Window.createOutputChannel("Go Swag"),
    traceOutputChannel: Window.createOutputChannel("Go Swag Trace"),
  };

  context.subscriptions.push(
    commands.registerCommand("swaggo-language-server-client.format", format)
  );

  try {
    client = new LanguageClient(
      "swaggo",
      "Go Swag",
      serverOptions,
      clientOptions
    );
  } catch (error) {
    void Window.showErrorMessage(`Failed to create language client: ${error}`);
  }

  client
    .start()
    .catch((error) => client.error(`Failed to start the server: ${error}`));
}

export async function deactivate() {
  if (!client) {
    return undefined;
  }
  await client.stop();
}

async function format(filePath: string) {
  if (!client) {
    return;
  }
  if (!filePath) {
    Window.showErrorMessage("Must provide a file path to format");
  }

  const directoryPath = dirname(filePath);

  const entries = await readdir(directoryPath, { withFileTypes: true });
  const excludeFilePaths = entries
    .filter(
      (entry) => entry.isFile() && join(directoryPath, entry.name) !== filePath
    )
    .map((entry) => join(directoryPath, entry.name));

  try {
    const command = `swag fmt --dir ${directoryPath} --exclude ${excludeFilePaths.join(
      ","
    )}`;

    exec(command);
    Window.showInformationMessage("File formatted successfully");
  } catch (e) {
    Window.showErrorMessage(`Failed to format the file: ${e}`);
  }
}
