# Interactive Bot Chat Schema

```mermaid
flowchart TD
    X[Interactive Bot Chat App]:::main
    X --> Initialize(Initialize):::process

    Initialize --> EnvVars(<a href="/103-bot-draft/main.go#L18">Load Environment Variables</a>):::config
    EnvVars -->|MODEL_RUNNER_BASE_URL| BaseURL(<a href="/103-bot-draft/main.go#L18">Base URL</a>):::env
    EnvVars -->|COOK_MODEL| Model(<a href="/103-bot-draft/main.go#L19">Model Name</a>):::env
    EnvVars -->|TEMPERATURE| Temperature(<a href="/103-bot-draft/main.go#L29">Temperature Setting</a>):::env
    EnvVars -->|TOP_P| TopP(<a href="/103-bot-draft/main.go#L30">Top P Setting</a>):::env
    EnvVars -->|AGENT_NAME| AgentName(<a href="/103-bot-draft/main.go#L31">Agent Name</a>):::env
    EnvVars -->|SYSTEM_INSTRUCTIONS| SystemInstructions(<a href="/103-bot-draft/main.go#L32">System Instructions</a>):::env

    Initialize --> OpenAIClient(<a href="/103-bot-draft/main.go#L21">Create OpenAI Client</a>):::client
    BaseURL --> OpenAIClient

    OpenAIClient -->|option.WithBaseURL| ClientConfig(<a href="/103-bot-draft/main.go#L22">Configure Base URL</a>):::config
    OpenAIClient -->|option.WithAPIKey| EmptyAPIKey(<a href="/103-bot-draft/main.go#L23">Empty API Key</a>):::config

    Initialize --> Context(<a href="/103-bot-draft/main.go#L26">Create Context</a>):::context

    Initialize --> MainLoop(<a href="/103-bot-draft/main.go#L34">Start Main Chat Loop</a>):::loop

    MainLoop --> Reader(<a href="/103-bot-draft/main.go#L35">Create Stdin Reader</a>):::input
    Reader --> Prompt(<a href="/103-bot-draft/main.go#L36">Display Bot Prompt</a>):::prompt
    Prompt --> ReadInput(<a href="/103-bot-draft/main.go#L37">Read User Input</a>):::input

    ReadInput --> CheckExit(<a href="/103-bot-draft/main.go#L39">Check for /bye Command</a>):::check
    CheckExit -->|/bye| ExitMessage(<a href="/103-bot-draft/main.go#L40">Display Goodbye</a>):::exit
    ExitMessage --> EndLoop(<a href="/103-bot-draft/main.go#L41">Break Loop</a>):::exit

    CheckExit -->|Continue| PrepareMessages(<a href="/103-bot-draft/main.go#L44">Prepare Chat Messages</a>):::message
    SystemInstructions --> PrepareMessages
    PrepareMessages -->|openai.SystemMessage| SystemMessage(<a href="/103-bot-draft/main.go#L45">System Message</a>):::message
    PrepareMessages -->|openai.UserMessage| UserMessage(<a href="/103-bot-draft/main.go#L46">User Message</a>):::message

    PrepareMessages --> Parameters(<a href="/103-bot-draft/main.go#L49">Chat Completion Parameters</a>):::params
    Model --> Parameters
    Temperature --> Parameters
    TopP --> Parameters
    Parameters -->|Temperature| TempParam(<a href="/103-bot-draft/main.go#L52">openai.Opt temperature</a>):::params
    Parameters -->|TopP| TopPParam(<a href="/103-bot-draft/main.go#L53">openai.Opt topP</a>):::params

    OpenAIClient --> StreamCompletion(<a href="/103-bot-draft/main.go#L56">Create Streaming Completion</a>):::stream
    Context --> StreamCompletion
    Parameters --> StreamCompletion

    StreamCompletion --> NewLine(<a href="/103-bot-draft/main.go#L58">Print New Line</a>):::output
    NewLine --> StreamLoop(<a href="/103-bot-draft/main.go#L60">Stream Processing Loop</a>):::streamloop
    StreamLoop -->|stream.Next| NextChunk(<a href="/103-bot-draft/main.go#L61">Get Next Chunk</a>):::chunk
    NextChunk --> CheckContent(<a href="/103-bot-draft/main.go#L63">Check Content Available</a>):::check
    CheckContent -->|Has Content| PrintChunk(<a href="/103-bot-draft/main.go#L64">Print Delta Content</a>):::output
    PrintChunk --> StreamLoop

    StreamLoop --> ErrorCheck(<a href="/103-bot-draft/main.go#L68">Check Stream Error</a>):::error
    ErrorCheck -->|if err != nil| FatalError(<a href="/103-bot-draft/main.go#L69">log.Fatalln</a>):::error

    ErrorCheck -->|No Error| Separator(<a href="/103-bot-draft/main.go#L71">Print Separator Line</a>):::output
    Separator --> MainLoop

    classDef main fill:#e1f5fe,stroke:#01579b,stroke-width:3px,color:#000
    classDef process fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef config fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px,color:#000
    classDef env fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef client fill:#fce4ec,stroke:#880e4f,stroke-width:2px,color:#000
    classDef context fill:#f1f8e9,stroke:#33691e,stroke-width:2px,color:#000
    classDef loop fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef input fill:#fff8e1,stroke:#f57f17,stroke-width:2px,color:#000
    classDef prompt fill:#e8eaf6,stroke:#283593,stroke-width:2px,color:#000
    classDef check fill:#f9fbe7,stroke:#827717,stroke-width:2px,color:#000
    classDef exit fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#000
    classDef message fill:#e4f7ff,stroke:#0277bd,stroke-width:2px,color:#000
    classDef params fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#000
    classDef stream fill:#e3f2fd,stroke:#1976d2,stroke-width:2px,color:#000
    classDef streamloop fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px,color:#000
    classDef chunk fill:#e8f5e8,stroke:#388e3c,stroke-width:2px,color:#000
    classDef output fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#000
    classDef error fill:#ffebee,stroke:#d32f2f,stroke-width:2px,color:#000
```

**Key Components:**

- **Environment Variables**:
  - `MODEL_RUNNER_BASE_URL`: Base URL for the Docker Model Runner
  - `COOK_MODEL`: Model identifier to use for completion
  - `TEMPERATURE`: Controls randomness in responses (parsed as float64)
  - `TOP_P`: Controls nucleus sampling (parsed as float64)
  - `AGENT_NAME`: Display name for the bot in prompts
  - `SYSTEM_INSTRUCTIONS`: System message defining bot behavior

- **OpenAI Client Configuration**:
  - Base URL configured from environment variable
  - Empty API key (using local model runner)

- **Interactive Chat Loop**:
  - Continuous loop for user interaction
  - Stdin reader for user input capture
  - Dynamic prompt showing agent name and model
  - Exit command `/bye` to terminate session

- **Message System**:
  - System Message: From `SYSTEM_INSTRUCTIONS` environment variable
  - User Message: Real-time input from user

- **Streaming Parameters**:
  - Model: From environment variable
  - Temperature: Dynamic from environment
  - TopP: Dynamic from environment

- **Interactive Flow**:
  1. Load all environment variables and parse numeric values
  2. Create OpenAI client with custom base URL
  3. Enter continuous chat loop
  4. Display interactive prompt with agent name and model
  5. Read user input from stdin
  6. Check for exit command (`/bye`)
  7. Prepare system and user messages
  8. Configure streaming completion parameters
  9. Create streaming completion request
  10. Process stream chunks in real-time
  11. Print each content delta as it arrives
  12. Display separator and continue loop
  13. Handle stream errors and exit conditions

**Interactive Features**:
- **Continuous conversation**: Loop allows multiple exchanges
- **Real-time streaming**: Responses appear as they're generated
- **Configurable agent**: Agent name and behavior via environment
- **Clean exit**: `/bye` command for graceful termination
- **Visual separation**: Dashed lines between conversation turns