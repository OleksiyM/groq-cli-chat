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
- **Response statistics** — see token count, completion time, and tokens per second.
- **Support for multiple AI providers** — works with any OpenAI-compatible API.

---

## 🚀 Getting Started

### 🔑 Requirements

- **Groq API Key**: [Sign up here](https://groq.com)
- (Optional) **OpenAI API Key** if you want to use OpenAI models, also you can get xAI, OpenRouter, Mistral, etc.
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
After that you will see:
```txt
🍎 One-shot Groq CLI chat
[i]nfo | select [m]odel | [u]pdate models | [h]istory | change [c]onfig | [q]uit

[model_name] >
```

### Interactive commands:

- `[i]` — Show current model info
- `[m]` — Switch model (up to 20 shown)
- `[u]` — Update models list from the current provider. Excluded models (part of the names) can be configured in `config_<provider>.yaml` (see below)
- `[h]` — Show history - list of the saved answers in Markdown format
- `[c]` — Change config (it should be previously saved in `config_<provider>.yaml`)
- `[q]` — Quit

### One-shot prompts 
<details>
  <summary>Examples of good one-shot prompts</summary>

**Example 1:**
```text
You are a DevOps engineer. Explain the best practices for setting up a CI/CD pipeline using Jenkins in 200-300 words.
```

**Example 2:**
```text
You are a Linux system administrator. Summarize the key steps to secure a Linux server against common vulnerabilities in 200-300 words.
```

**Example 3:**
```text
You are a cloud engineer. Describe the process of deploying a containerized application on Kubernetes, including key commands, in 200-300 words.
```

**Example 4:**
```text
You are a network engineer. Summarize how to configure a firewall using `iptables` to allow HTTP and SSH traffic in 200-300 words.
```

**Example 5:**
```text
You are a DevOps specialist. Explain how to monitor system performance using `top`, `htop`, and `vmstat` commands in 200-300 words.
```

**Example 6:**
```text
You are a Linux user. Summarize how to use `grep` and `awk` together to extract and process log file data in 200-300 words.
```

**Example 7:**
```text
You are a system administrator. Describe the steps to automate backups using `rsync` and cron jobs in 200-300 words.
```
</details>

> Remember: Each prompt must be fully self-contained.

After receiving a response, you'll see statistics in the format:
```text
───┤ Stats: 31 tokens | 0.02 sec | 1911.78 tok/sec ├───
```

This shows:
- Total tokens used
- Completion time in seconds
- Tokens per second processing rate

---

## ⚙️ Configuration

### Auto-generated on first run

`~/.groq-chat/config.yaml`

```yaml
api_key_name: GROQ_API_KEY
app_title: "\U0001F34E One-shot Groq CLI chat"
base_url: https://api.groq.com/openai/v1
default_model: allam-2-7b
excluded_models:
    - whisper
    - playai
models:
    - allam-2-7b
    - compound-beta
    - compound-beta-mini
    - deepseek-r1-distill-llama-70b
    - gemma2-9b-it
    - llama-3.1-8b-instant
    - llama-3.3-70b-versatile
    - llama-guard-3-8b
    - llama3-70b-8192
    - llama3-8b-8192
    - meta-llama/llama-4-maverick-17b-128e-instruct
    - meta-llama/llama-4-scout-17b-16e-instruct
    - mistral-saba-24b
    - qwen-qwq-32b
provider_name: Groq
```

- Use `GROQ_API_KEY` as env variable.
- You canchange any of the values.
- Models list will be updated automatically via `[u]` command. 
- You can exclude some models from the list when updating. Each line is part of the excluded model(s) name
- You can create your own config file and use it with `[c]` command. (see below)

### Custom providers configurations

<details>
  <summary>Examples of config files</summary>

#### xAI (Grok)

**config_grok.yaml**
```yaml
api_key_name: XAI_API_KEY
app_title: "\U0001F680 One-shot Grok-3 CLI chat"
base_url: https://api.x.ai/v1
default_model: grok-3-mini-beta
excluded_models:
    - fast
    - playai
models:
    - grok-3-beta
    - grok-3-mini-beta
provider_name: Grok
```

#### OpenAI

**config_openai.yaml**
```yaml
api_key_name: OPENAI_API_KEY
app_title: "One-shot OpenAI CLI chat"
base_url: https://api.openai.com
default_model: gpt-4.1-mini
excluded_models:
    - fast
    - playai
    - llama
models:
    - gpt-4.1-mini
    - gpt-4.1
provider_name: OpenAI
```
#### OpenRouter

**config_openrouter.yaml**
```yaml
api_key_name: OPENROUTER_API_KEY
app_title: "One-shot OpenRouter CLI chat"
base_url: https://openrouter.ai/api/v1
default_model: nvidia/llama-3.3-nemotron-super-49b-v1:free
excluded_models:
    - fast
    - playai
    - llama
models:
    - deepseek/deepseek-chat-v3-0324:free
    - deepseek/deepseek-v3-base:free
    - deepseek/deepseek-r1:free
    - nvidia/llama-3.1-nemotron-ultra-253b-v1:free
    - nvidia/llama-3.3-nemotron-super-49b-v1:free
    - qwen/qwen3-235b-a22b:free
provider_name: OpenRouter
```
</details>


- Use `api_key_name` as env variable.

## Chat history

- The chat history is saved as Markdown files in `~/.groq-chat/history/`. Each chat is a separate file, named after the timestamp of the chat creation.
- It can be viewed in any Markdown viewer.
- You can delete them manually

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
make build        # Build the binary
make clean        # Clean build artifacts
make release      # Build .zip/.tar.gz archives and binaries
make docker       # Docker image (Debian base) - not implemented yet
```

---

## 📁 Project Structure

```
cmd/
└── main.go/    # Main CLI entry point
internal/
├── chat/       # Chat loop, history
├── config/     # Config management
├── groq/       # API client
resources/      # UI messages, defaults
Dockerfile      # distroless Debian image (~9MB)
Dockerfile.rhel # scratch-based RHEL image
bin/            # Compiled binaries
go.mod          # Go module file
go.sum          # Go module checksum file
Makefile        # Build automation
README.md       # This file
```

---

## ✅ Features

- Minimal memory & disk usage
- Clean, interactive CLI
- Configurable via YAML and env vars
- Secure (no API key in source)
- Chat history in Markdown
- Blasing fast
- Crossplatform (Linux, macOS, Windows)
- Response statistics (tokens, time, tokens/sec)
- Support for multiple AI providers (OpenAI compartible APIs)

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
  - [x]	Response statistics (tokens, time, tokens/sec)
  - [x]	Support for multiple AI providers (OpenAI-compatible APIs)
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

*Made with ❤️ for hackers, devs, DevOps, SysAdmins and prompt-crafters.*
