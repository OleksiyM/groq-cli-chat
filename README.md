# groq-chat

**groq-chat** is a lightweight, high-performance command-line interface (CLI) for interacting with Groq AI models. Built in Go, it offers a one-shot chat experience with minimal resource usage (11-14MB memory) and blazing-fast performance, even on older hardware. With a sleek, interactive interface, it supports model selection, chat history, and model information retrieval, all packaged in tiny Docker images (~9-11MB).

Tested on Windows 11, Fedora 42, Ubuntu 24.04 and macOS, `groq-chat` is perfect for developers, AI enthusiasts, and anyone seeking a fast, efficient CLI to chat with Groq’s powerful language models.

## Important Note: One-Shot Prompting
The `groq-chat` CLI does **not** remember previous messages or maintain conversation context. Each question or input is treated as a **new, standalone prompt** sent to the Groq model. This design choice ensures simplicity and accommodates models with low context windows, such as `llama3-8b-8192` (8k tokens).

To get high-quality answers, craft **effective one-shot prompts** that include all necessary context and instructions in a single message. A good one-shot prompt is clear, specific, and self-contained.

**Example One-Shot Prompt**:
```
You are an expert historian. Provide a detailed summary of the key events in the American Revolution (1775-1783), including major battles, political developments, and the role of key figures like George Washington and Thomas Jefferson, in 300-500 words.
```
This prompt specifies the role, task, scope, and desired length, ensuring a comprehensive response without relying on prior messages.

## Configuration
- **Config File**: On first start, `groq-chat` creates `~/.groq-chat/config.yaml` by fetching all available models from the Groq API (`https://api.groq.com/openai/v1/models`). The first model in the list is set as the default.
  ```yaml
  base_url: https://api.groq.com/openai/v1
  api_key: ""
  default_model: llama3-8b-8192
  models:
    - llama3-8b-8192
    - mixtral-8x7b-32768
    # ... (all available models)
  ```
- **Model Management**: It’s recommended to edit `config.yaml` to remove unused models, keeping the list to 10 or fewer for better usability in the `[m]` model selection menu (model is selected by number from 0 to 9).
- **Environment Variable**: Set `GROQ_API_KEY` to authenticate API requests.
  ```bash
  export GROQ_API_KEY=gsk_...
  ```
- **History**: Chat logs are saved as Markdown files in `~/.groq-chat/history/`.

## Features
- **Minimal Footprint**: Uses 11-14MB memory and runs in ~9-11MB Docker images.
- **Cross-Platform**: Works on Windows, Linux, and macOS.
- **Interactive CLI**: Supports `[i]`nfo, `[m]`odel selection, `[h]`istory, and `[q]`uit commands.
- **Chat History**: Saves conversations as Markdown files in `~/.groq-chat/history/`.
- **Configurable**: Uses `~/.groq-chat/config.yaml` for model preferences.
- **Secure**: Leverages environment variables (`GROQ_API_KEY`) for API authentication.
- **Versioned**: Run `groq-chat --version` to check the current version (0.1.0).

## Project Structure
The project is organized for clarity and maintainability, with Go source files and Docker configurations:

```
groq-chat/
├── Dockerfile              # Debian-based Docker image (9MB)
├── Dockerfile.rhel         # RHEL-based Docker image (~9-11MB)
├── main.go                 # CLI entry point with Cobra command setup
├── internal/
│   ├── chat/               # Chat logic and history management
│   │   ├── chat.go         # Main chat loop and UI
│   │   ├── history.go      # Chat history listing
│   ├── config/             # Configuration loading and validation
│   │   ├── config.go       # Config file handling with Viper
│   │   ├── models.go       # Model validation utilities
│   ├── groq/               # Groq API client
│   │   ├── client.go       # HTTP client for API requests
│   │   ├── types.go        # API response structs
├── resources/              # Constants and messages
│   ├── messages.go         # Error and UI messages
│   ├── defaults.go         # Default constants (empty for now)
├── bin/                    # Build output
│   ├── groq-chat           # Local binary
│   ├── release/            # Release binaries and archives
│       ├── groq-chat-linux-amd64.tar.gz
│       ├── groq-chat-windows-amd64.zip
│       └── ...
```

- **Go Files**: Modular design with separate packages for chat, config, and API interactions.
- **Dockerfiles**: Two minimal images (Debian: `gcr.io/distroless/static-debian12`, RHEL: `scratch` with CA certificates).
- **Releases**: Cross-platform binaries (`tar.gz`, `zip`) for Linux, Windows, and macOS.

## Prerequisites
- **Go**: For building from source (optional, version 1.18+ recommended).
- **Docker**: For running containerized images (optional).
- **GROQ_API_KEY**: Mandatory. Obtain from [Groq](https://groq.com) to authenticate API requests.

## Installation
Currently, `groq-chat` is distributed via GitHub Releases. Cloning and GitHub Actions setup will be available in the next version.

### Download from GitHub Releases
1. Visit the [GitHub Releases page](https://github.com/OleksiyM/groq-cli-chat/releases).
2. Download the appropriate binary for your platform:
   - Linux: `groq-chat-linux-amd64.tar.gz`
   - Windows: `groq-chat-windows-amd64.zip`
   - macOS: `groq-chat-macos-amd64.tar.gz`
3. Extract the archive:
   ```bash
   tar -xzf groq-chat-linux-amd64.tar.gz
   # or
   unzip groq-chat-windows-amd64.zip
   ```
4. Move the binary to a directory in your `PATH` (optional):
   ```bash
   mv groq-chat /usr/local/bin/
   ```

### Run with Docker
1. Download the Docker image from GitHub Releases or build locally.
2. Run the Debian-based image (9MB):
   ```bash
   docker run --rm -it -e GROQ_API_KEY=$GROQ_API_KEY -v $HOME/.groq-chat:/root/.groq-chat groq-chat:debian
   ```
3. Or use the RHEL-based image (~9-11MB):
   ```bash
   docker run --rm -it -e GROQ_API_KEY=$GROQ_API_KEY -v $HOME/.groq-chat:/root/.groq-chat groq-chat:rhel
   ```

**Notes**:
- Groq API key (export it: `export GROQ_API_KEY=gsk_...`).
- The `-v $HOME/.groq-chat:/root/.groq-chat` mount persists config (`config.yaml`) and history (`.md` files).
- Ensure `~/.groq-chat/` is writable:
  ```bash
  mkdir -p ~/.groq-chat
  chmod -R 777 ~/.groq-chat
  ```

## Usage
Run `groq-chat` to start the interactive CLI:
```bash
./groq-chat
# or
groq-chat.exe
# or
docker run --rm -it -e GROQ_API_KEY=$GROQ_API_KEY -v $HOME/.groq-chat:/root/.groq-chat groq-chat:debian
```

### Commands
- `[i]`: Display current model information (ID, owner, context window).
- `[m]`: Select a different model (lists up to 10 available models).
- `[h]`: List chat history (Markdown files in `~/.groq-chat/history/`).
- `[q]`: Quit the CLI.
- Any other input: Send a chat message to the Groq API.


## Building from Source
1. Clone the repository (available in next version).
2. Install Go (1.18+).
3. Build the binary:
   ```bash
   make build
   # Output: bin/groq-chat
   ```
4. Build release binaries and Docker images:
   ```bash
   make buid
   make release
   # make docker-build-debian (not implemented yet)
   # make docker-build-rhel (not implemented yet)
   ```

## Roadmap
- [x] Create GitHub repository with cloning instructions.
- [ ] Add GitHub Actions for automated builds, tests, and Docker Hub pushes.
- [ ] **May be** Enhance UI with ANSI colors for better readability.
- [ ] **May be** Support additional platforms and distributions.
- [ ] **May be** Restore advanced command-line flags (e.g., custom config paths).

## Contributing
Contributions are welcome! Stay tuned for the GitHub repository to submit issues, feature requests, or pull requests.

## License
[MIT License](LICENSE).

## Acknowledgments
- Built with [Go](https://golang.org), [Cobra](https://github.com/spf13/cobra), and [Viper](https://github.com/spf13/viper).
- Powered by [Groq](https://groq.com) AI models.
- Inspired by minimal, high-performance CLI tools.

---
*Created with ❤️ for AI enthusiasts and CLI lovers.*