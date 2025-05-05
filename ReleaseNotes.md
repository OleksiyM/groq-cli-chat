# Release Notes - groq-chat

## v1.0.0 - Response Statistics & Multi-Provider Support

### New Features
- **Response Statistics**: Added detailed statistics after each response showing:
  - Total tokens used
  - Completion time in seconds
  - Tokens per second processing rate
- **Multi-Provider Support**: Now compatible with any OpenAI-compatible API provider
  - Works with xAI's grok-3-mini-beta and other compatible models
  - Simply configure the base_url in your config.yaml

### Improvements
- **Enhanced Error Handling**: Better error messages when API responses are empty or malformed
- **Response Processing**: Improved handling of API responses with proper validation
- **Documentation**: Updated README with new features and usage examples

### Bug Fixes
- Fixed issue with response not being displayed in the terminal
- Fixed timing calculation for completion time when not provided by the API
- Removed debug output that was appearing in the terminal

### Technical Details
- Refactored the Chat method in the client.go file to properly handle API responses
- Added validation to ensure responses contain valid content
- Implemented fallback timing calculation when the API doesn't provide completion time

---

## v0.1.1
- fixes: improve UI formatting and readability
- 
---

## v0.1.0 - Initial Release

### Features
- Minimalistic CLI interface for interacting with Groq AI models
- One-shot prompts with no conversation context
- Interactive commands: [i], [m], [h], [q]
- Model selection and information viewing
- Chat history saved as Markdown files
- Configuration via YAML and environment variables
- Cross-platform support (Windows, Linux, macOS)
- Docker images (Debian and RHEL-based)
- Low resource usage (11-14MB RAM)
- Small Docker image size (9-11MB)