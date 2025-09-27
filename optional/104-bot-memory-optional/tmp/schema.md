# Interactive Bot with Memory Schema

```mermaid
flowchart TD
    X[Interactive Bot with Memory]:::main
    X --> Initialize(Initialize):::process

    Initialize --> EnvVars(<a href="/104-bot-memory/main.go#L19">Load Environment Variables</a>):::config
    EnvVars -->|MODEL_RUNNER_BASE_URL| BaseURL(<a href="/104-bot-memory/main.go#L19">Base URL</a>):::env
    EnvVars -->|COOK_MODEL| Model(<a href="/104-bot-memory/main.go#L20">Model Name</a>):::env
    EnvVars -->|TEMPERATURE| Temperature(<a href="/104-bot-memory/main.go#L30">Temperature Setting</a>):::env
    EnvVars -->|TOP_P| TopP(<a href="/104-bot-memory/main.go#L31">Top P Setting</a>):::env
    EnvVars -->|AGENT_NAME| AgentName(<a href="/104-bot-memory/main.go#L32">Agent Name</a>):::env
    EnvVars -->|SYSTEM_INSTRUCTIONS| SystemInstructions(<a href="/104-bot-memory/main.go#L33">System Instructions</a>):::env

    Initialize --> OpenAIClient(<a href="/104-bot-memory/main.go#L22">Create OpenAI Client</a>):::client
    BaseURL --> OpenAIClient

    OpenAIClient -->|option.WithBaseURL| ClientConfig(<a href="/104-bot-memory/main.go#L23">Configure Base URL</a>):::config
    OpenAIClient -->|option.WithAPIKey| EmptyAPIKey(<a href="/104-bot-memory/main.go#L24">Empty API Key</a>):::config

    Initialize --> Context(<a href="/104-bot-memory/main.go#L27">Create Context</a>):::context

    Initialize --> InitMessages(<a href="/104-bot-memory/main.go#L36">Initialize Messages with System</a>):::memory
    SystemInstructions --> InitMessages

    InitMessages --> MainLoop(<a href="/104-bot-memory/main.go#L40">Start Main Chat Loop</a>):::loop

    MainLoop --> Reader(<a href="/104-bot-memory/main.go#L41">Create Stdin Reader</a>):::input
    Reader --> Prompt(<a href="/104-bot-memory/main.go#L42">Display Bot Prompt with Memory Icon</a>):::prompt
    Prompt --> ReadInput(<a href="/104-bot-memory/main.go#L43">Read User Input</a>):::input

    ReadInput --> CheckExit(<a href="/104-bot-memory/main.go#L45">Check for /bye Command</a>):::check
    CheckExit -->|/bye| ExitMessage(<a href="/104-bot-memory/main.go#L46">Display Goodbye</a>):::exit
    ExitMessage --> EndLoop(<a href="/104-bot-memory/main.go#L47">Break Loop</a>):::exit

    CheckExit -->|Continue| CheckMemory(<a href="/104-bot-memory/main.go#L50">Check for /memory Command</a>):::check
    CheckMemory -->|/memory| DisplayMemory(<a href="/104-bot-memory/main.go#L51">Display Conversational Memory</a>):::memory
    DisplayMemory --> MainLoop

    CheckMemory -->|Regular Message| AppendUser(<a href="/104-bot-memory/main.go#L56">Append User Message to Memory</a>):::memory

    AppendUser --> Parameters(<a href="/104-bot-memory/main.go#L58">Chat Completion Parameters</a>):::params
    Model --> Parameters
    Temperature --> Parameters
    TopP --> Parameters
    Parameters -->|All Messages| MessageHistory(<a href="/104-bot-memory/main.go#L59">Include Message History</a>):::memory
    Parameters -->|Temperature| TempParam(<a href="/104-bot-memory/main.go#L61">openai.Opt temperature</a>):::params
    Parameters -->|TopP| TopPParam(<a href="/104-bot-memory/main.go#L62">openai.Opt topP</a>):::params

    OpenAIClient --> StreamCompletion(<a href="/104-bot-memory/main.go#L65">Create Streaming Completion</a>):::stream
    Context --> StreamCompletion
    Parameters --> StreamCompletion

    StreamCompletion --> NewLine(<a href="/104-bot-memory/main.go#L67">Print New Line</a>):::output
    NewLine --> InitAnswer(<a href="/104-bot-memory/main.go#L70">Initialize Answer String</a>):::accumulator
    InitAnswer --> StreamLoop(<a href="/104-bot-memory/main.go#L71">Stream Processing Loop</a>):::streamloop
    StreamLoop -->|stream.Next| NextChunk(<a href="/104-bot-memory/main.go#L72">Get Next Chunk</a>):::chunk
    NextChunk --> CheckContent(<a href="/104-bot-memory/main.go#L74">Check Content Available</a>):::check
    CheckContent -->|Has Content| ExtractContent(<a href="/104-bot-memory/main.go#L75">Extract Delta Content</a>):::chunk
    ExtractContent --> AccumulateAnswer(<a href="/104-bot-memory/main.go#L77">Accumulate to Answer</a>):::accumulator
    AccumulateAnswer --> PrintChunk(<a href="/104-bot-memory/main.go#L78">Print Content</a>):::output
    PrintChunk --> StreamLoop

    StreamLoop --> ErrorCheck(<a href="/104-bot-memory/main.go#L82">Check Stream Error</a>):::error
    ErrorCheck -->|if err != nil| FatalError(<a href="/104-bot-memory/main.go#L83">log.Fatalln</a>):::error

    ErrorCheck -->|No Error| AppendAssistant(<a href="/104-bot-memory/main.go#L87">Append Assistant Response to Memory</a>):::memory
    AppendAssistant --> Separator(<a href="/104-bot-memory/main.go#L89">Print Separator Line</a>):::output
    Separator --> MainLoop

    DisplayMemory --> MemoryFunction[<a href="/104-bot-memory/main.go#L117">DisplayConversationalMemory Function</a>]:::function
    MemoryFunction --> MessageConverter[<a href="/104-bot-memory/main.go#L96">MessageToMap Function</a>]:::function
    MessageConverter --> JSONMarshal(<a href="/104-bot-memory/main.go#L97">MarshalJSON</a>):::conversion
    JSONMarshal --> JSONUnmarshal(<a href="/104-bot-memory/main.go#L103">Unmarshal to map</a>):::conversion
    JSONUnmarshal --> StringMap(<a href="/104-bot-memory/main.go#L107">Convert to string map</a>):::conversion

    classDef main fill:#e1f5fe,stroke:#01579b,stroke-width:3px,color:#000
    classDef process fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef config fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px,color:#000
    classDef env fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef client fill:#fce4ec,stroke:#880e4f,stroke-width:2px,color:#000
    classDef context fill:#f1f8e9,stroke:#33691e,stroke-width:2px,color:#000
    classDef memory fill:#fff3e0,stroke:#ff6f00,stroke-width:3px,color:#000
    classDef loop fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef input fill:#fff8e1,stroke:#f57f17,stroke-width:2px,color:#000
    classDef prompt fill:#e8eaf6,stroke:#283593,stroke-width:2px,color:#000
    classDef check fill:#f9fbe7,stroke:#827717,stroke-width:2px,color:#000
    classDef exit fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#000
    classDef params fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#000
    classDef stream fill:#e3f2fd,stroke:#1976d2,stroke-width:2px,color:#000
    classDef streamloop fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px,color:#000
    classDef chunk fill:#e8f5e8,stroke:#388e3c,stroke-width:2px,color:#000
    classDef accumulator fill:#e1f5fe,stroke:#0277bd,stroke-width:2px,color:#000
    classDef output fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#000
    classDef error fill:#ffebee,stroke:#d32f2f,stroke-width:2px,color:#000
    classDef function fill:#f1f8e9,stroke:#689f38,stroke-width:2px,color:#000
    classDef conversion fill:#fce4ec,stroke:#ad1457,stroke-width:2px,color:#000
```

**Key Components:**

- **Environment Variables**:
  - `MODEL_RUNNER_BASE_URL`: Base URL for the Docker Model Runner
  - `COOK_MODEL`: Model identifier to use for completion
  - `TEMPERATURE`: Controls randomness in responses (parsed as float64)
  - `TOP_P`: Controls nucleus sampling (parsed as float64)
  - `AGENT_NAME`: Display name for the bot in prompts
  - `SYSTEM_INSTRUCTIONS`: System message defining bot behavior

- **Memory Management**:
  - **Persistent Message History**: Messages slice maintains entire conversation
  - **System Message**: Initial system instructions stored in memory
  - **User Messages**: Each user input appended to conversation history
  - **Assistant Messages**: Complete responses accumulated and stored
  - **Memory Display**: `/memory` command shows full conversation history

- **OpenAI Client Configuration**:
  - Base URL configured from environment variable
  - Empty API key (using local model runner)

- **Interactive Chat Loop**:
  - Continuous loop for user interaction with memory persistence
  - Enhanced prompt with brain emoji indicating memory capability
  - Memory inspection via `/memory` command
  - Exit command `/bye` to terminate session

- **Response Processing**:
  - **Content Accumulation**: Full assistant response collected during streaming
  - **Memory Storage**: Complete response added to conversation history
  - **Real-time Display**: Content streamed to user as generated

- **Memory Functions**:
  - **`MessageToMap`**: Converts OpenAI messages to readable map format
  - **`DisplayConversationalMemory`**: Shows formatted conversation history
  - **JSON Conversion**: Marshal/unmarshal for message processing

**Memory Features**:
1. **Conversation Persistence**: All messages (system, user, assistant) stored
2. **Context Awareness**: Each request includes full conversation history
3. **Memory Inspection**: `/memory` command displays conversation log
4. **Response Accumulation**: Streaming responses fully captured for storage
5. **Continuous Context**: Bot remembers previous exchanges throughout session

**Interactive Flow**:
1. Initialize with system message in memory
2. Enter continuous chat loop with memory persistence
3. Handle special commands (`/bye`, `/memory`)
4. Append user messages to conversation history
5. Include full message history in API requests
6. Stream and accumulate assistant responses
7. Store complete assistant responses in memory
8. Maintain conversation context across interactions