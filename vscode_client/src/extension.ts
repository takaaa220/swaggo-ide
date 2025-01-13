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
import { homedir } from "os";
import { existsSync } from "fs";
import { promisify } from "util";

let client: LanguageClient;

export async function activate(context: ExtensionContext) {
  const binaryPath = await getLanguageServerBinaryPath();
  if (!binaryPath) {
    console.log("swaggo-language-server binary not found");
    return;
  }

  const serverOptions: ServerOptions = {
    run: { command: binaryPath, transport: TransportKind.stdio },
    debug: { command: binaryPath, transport: TransportKind.stdio },
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

async function getLanguageServerBinaryPath(): Promise<string | undefined> {
  const devPath = process.env.SWAGGO_LS_DEV_PATH;

  const getGoBinPath = () => {
    const gopath = process.env.GOPATH ?? join(homedir(), "go");
    return join(gopath, "bin");
  };

  const binaryPath = devPath ?? join(getGoBinPath(), "swaggo-language-server");
  if (existsSync(binaryPath)) {
    console.log("swaggo-language-server binary found", binaryPath);
    return binaryPath;
  }

  const installed = await installBinary();
  return installed ? binaryPath : undefined;
}

async function installBinary(): Promise<boolean> {
  const res = await Window.showInformationMessage(
    "Install swaggo-language-server binary?",
    "Yes"
  );
  if (res !== "Yes") {
    return false;
  }

  try {
    await promisify(exec)(
      "go install github.com/takaaa220/swaggo-ide/swaggo-language-server@latest"
    );
    Window.showInformationMessage(
      "swaggo-language-server installed successfully"
    );
    return true;
  } catch (e) {
    Window.showErrorMessage(`Failed to install swaggo-language-server: ${e}`);
    return false;
  }
}
