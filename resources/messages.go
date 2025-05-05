package resources

const (
	MenuOptions = "[i]nfo | select [m]odel | [u]pdate models | [h]istory | change [c]onfig | [q]uit"
	WelcomeMessage = `🍎 One-shot Groq CLI chat
` + MenuOptions + `

`

	Prompt           = "[%s] > "
	InfoModel        = "Current model: %s\n"
	InfoModelDetails = `────────┤ Model Information ├─────────
- ID: %s
- Owned By: %s
- Active: %v
- Context Window: %d tokens
─────────────────────────────────────

`
	SelectModelHeader = `────────┤ Available models ├─────────`
	SelectModelPrompt = `─────────────────────────────────────
Select model (0-%d): `
	GoodbyeMessage     = "Goodbye!"
	InfoModelUnchanged = "Invalid selection, model unchanged: %s\n"

	HistoryFormat = `# Chat History (%s)
**Model**: %s
**User**: %s
**Response**: %s
**Stats**:
- Total Tokens: %d
- Completion Time: %.2f seconds
- Tokens per Second: %.2f
`

	StatsFormat = `───┤ Stats: %d tokens | %.2f sec | %.2f tok/sec ├───

`

	// Error messages
	ErrLoadConfig          = "failed to load config: %v\n"
	ErrExecuteCmd          = "failed to execute command: %v\n"
	ErrHomeDir             = "failed to get home directory: %v"
	ErrConfigDir           = "failed to create config directory: %v"
	ErrCreateConfig        = "failed to create default config: %v"
	ErrReadConfig          = "failed to read config: %v"
	ErrUnmarshalConfig     = "failed to unmarshal config: %v"
	ErrNoAPIKey            = "API key environment variable not set"
	ErrWriteConfig         = "failed to write config: %v"
	ErrCreateClient        = "failed to create Groq client: %v"
	ErrListModels          = "failed to list models: %v"
	ErrInvalidClientParams = "invalid client parameters: baseURL or apiKey is empty"
	ErrCreateRequest       = "failed to create request: %v"
	ErrHTTP                = "HTTP request failed: %v"
	ErrAPI                 = "API error (status %d): %s"
	ErrDecodeResponse      = "failed to decode response: %v"
	ErrEncodePayload       = "failed to encode payload: %v"
	ErrChat                = "chat request failed: %v"
	ErrSelectModel         = "failed to select model: %v"
	ErrReadInput           = "failed to read input"
	ErrInvalidChoice       = "invalid choice: %s"
	ErrChoiceOutOfRange    = "choice %d out of range (0-%d)"
	ErrSaveHistory         = "failed to save chat history: %v"
	ErrCreateHistoryDir    = "failed to create history directory: %v"
	ErrNoModels            = "no models available"
	ErrReadHistoryDir      = "failed to read history directory: %v"
	ErrInvalidConfig       = "invalid configuration: %v"
	ErrInvalidDefaultModel = "invalid default model: %s"
	ErrGetModel            = "failed to retrieve model information: %v"

	// Info messages
	InfoConfigCreated = "Config created at ~/.groq-chat/config.yaml. Please review and adjust models and default model."
)

const DefaultBaseURL = "https://api.groq.com/openai/v1"
