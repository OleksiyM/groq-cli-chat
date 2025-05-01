# groq-chat

**groq-chat** is a minimalistic, high-performance command-line interface (CLI) for interacting with Groq AI models. Written in Go, it provides a fast, resource-efficient one-shot chat experience ideal for developers and AI enthusiasts.

- **Tiny Memory usage**: 11–14MB
- **Small Docker image size**: \~9–11MB
- **Cross-platform**: Windows, Linux, macOS

---

## ⚡ Quick Overview

- **One-shot prompts only** — no conversation context is preserved.
- **Fast & lightweight** — runs even on older hardware.
- **Interactive CLI** — switch models, view history, and more.

---

## 🚀 Getting Started

### 🔑 Requirements

- **Groq API Key**: [Sign up here](https://groq.com)
- (Optional) **Go 1.18+** if building from source
- (Optional) **Docker** for containerized usage

### 📦 Installation via GitHub Releases

1. Visit [Releases](https://github.com/OleksiyM/groq-cli-chat/releases)
2. Download your platform-specific archive
3. Extract and move the binary into your `$PATH`
   ```bash
   tar -xzf groq-chat-linux-amd64.tar.gz
   mv groq-chat /usr/local/bin/
   ```

### 🐳 Run via Docker

- **Provide GROQ_API_KEY as an value**
```bash
docker run --rm -it \
  -e GROQ_API_KEY=your_api_key \
  -v $HOME/.groq-chat:/root/.groq-chat \
  oleksiyml/groq-chat
```

- **Provide GROQ_API_KEY as an env variable**
```bash
export GROQ_API_KEY=your_api_key

docker run --rm -it \
  -e GROQ_API_KEY=$GROQ_API_KEY \
  -v $HOME/.groq-chat:/root/.groq-chat \
  oleksiyml/groq-chat
```
See details at [Docker hub](https://hub.docker.com/r/oleksiyml/groq-chat).

---

## 💬 Usage

```bash
groq-chat
```

Interactive commands:

- `[i]` — Show model info
- `[m]` — Switch model (up to 10 shown)
- `[h]` — Show history
- `[q]` — Quit

Example one-shot prompt:

```text
You are a medical researcher. Summarize current evidence on the effectiveness of mRNA vaccines in 200-300 words.
```

> Remember: Each prompt must be fully self-contained.

---

## ⚙️ Configuration

### `~/.groq-chat/config.yaml`

Auto-generated on first run:

```yaml
base_url: https://api.groq.com/openai/v1
default_model: llama3-8b-8192
models:
  - llama3-8b-8192
  - mixtral-8x7b-32768
```

- Use `GROQ_API_KEY` ans env variable.
- History is saved as Markdown in `~/.groq-chat/history/` (One answer - new chat)

---

## 🔨 Building from Source

```bash
git clone https://github.com/OleksiyM/groq-cli-chat
cd groq-cli-chat
make build
```

Output will be in `bin/groq-chat`

Optional:

```bash
make clean        # Clean build artifacts
make release      # Build .zip/.tar.gz archives and binaries
make docker       # Docker image (Debian base) - not implemented yet
```

---

## 📁 Project Structure

```
internal/
├── chat/       # Chat loop, history
├── config/     # Config management
├── groq/       # API client
resources/      # UI messages, defaults
Dockerfile      # distroless Debian image (~9MB)
Dockerfile.rhel # scratch-based RHEL image
bin/            # Compiled binaries
```

---

## ✅ Features

- Minimal memory & disk usage
- Clean, interactive CLI
- Configurable via YAML and env vars
- Secure (no API key in source)
- Chat history in Markdown
- Blasing fast
- Crossplatform (Linix, macOS, Windows)

---

## 🔭 Roadmap
	- [x]	Create a CLI interface in Go with an interactive shell
	- [x]	Save chat history as Markdown files 
	- [x]	Support model selection from a list (up to 10 for now)
	- [x]	Built-in model information viewer
	- [x]	Cross-platform releases (`.tar.gz`, `.zip`) for `Linux`, `macOS`, and `Windows`
	- [x]	Compact Docker images based on Debian and RHEL
	- [x]	Configuration via `YAML` and `GROQ_API_KEY` environment variable
	- [x]	GitHub Actions: automatic build and release publishing
	- [ ]	CI: add unit tests and run go test in GitHub Actions
	- [ ]	Multi-turn chat support
	- [ ]	ANSI color output
	- [ ]	Expand CLI functionality: --config flag, auto-loading models
	- [ ]	Generate simple HTML version of chat history (from .md files)
	- [ ]	CLI command autocompletion for [i], [m], [h], etc., in TUI
	- [ ]	Documentation: add demos (GIFs/SVGs), example prompts
	- [ ]	Publish Docker images to Docker Hub (automatically from CI)
	- [ ]	Separate internal API and UI (prep for SDK/library)

---

## 🤝 Contributing

Contributions, ideas, and bug reports are welcome. Stay tuned for issue templates and guidelines.

## 📜 License

[MIT License](LICENSE)

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/), [Cobra](https://github.com/spf13/cobra), and [Viper](https://github.com/spf13/viper)
- Powered by [Groq](https://groq.com)
- Inspired by minimalistic CLI tools

---

*Made with ❤️ for hackers, devs, and prompt-crafters.*

