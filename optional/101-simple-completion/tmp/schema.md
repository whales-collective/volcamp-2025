# Simple OpenAI Completion Schema

```mermaid
flowchart TD
    X[Simple Completion App]:::main
    X --> Initialize(Initialize):::process

    Initialize --> EnvVars(<a href="/101-simple-completion/main.go#L15">Load Environment Variables</a>):::config
    EnvVars -->|MODEL_RUNNER_BASE_URL| BaseURL(<a href="/101-simple-completion/main.go#L15">Base URL</a>):::env
    EnvVars -->|COOK_MODEL| Model(<a href="/101-simple-completion/main.go#L16">Model Name</a>):::env

    Initialize --> OpenAIClient(<a href="/101-simple-completion/main.go#L18">Create OpenAI Client</a>):::client
    BaseURL --> OpenAIClient

    OpenAIClient -->|option.WithBaseURL| ClientConfig(<a href="/101-simple-completion/main.go#L19">Configure Base URL</a>):::config
    OpenAIClient -->|option.WithAPIKey| EmptyAPIKey(<a href="/101-simple-completion/main.go#L20">Empty API Key</a>):::config

    Initialize --> Context(<a href="/101-simple-completion/main.go#L23">Create Context</a>):::context

    Initialize --> Messages(<a href="/101-simple-completion/main.go#L25">Prepare Messages</a>):::message
    Messages -->|openai.UserMessage| UserMessage(<a href="/101-simple-completion/main.go#L27">User Message: Quel est ton nom?</a>):::message

    Messages --> Parameters(<a href="/101-simple-completion/main.go#L30">Chat Completion Parameters</a>):::params
    Model --> Parameters
    Parameters -->|Temperature| Temperature(<a href="/101-simple-completion/main.go#L33">Temperature: 0.5</a>):::params

    OpenAIClient --> ChatCompletion(<a href="/101-simple-completion/main.go#L36">Execute Chat Completion</a>):::api
    Context --> ChatCompletion
    Parameters --> ChatCompletion

    ChatCompletion --> ErrorHandling(<a href="/101-simple-completion/main.go#L38">Error Handling</a>):::error
    ErrorHandling -->|if err != nil| FatalError(<a href="/101-simple-completion/main.go#L39">log.Fatalln</a>):::error

    ChatCompletion --> Response(<a href="/101-simple-completion/main.go#L41">Extract Response Content</a>):::output
    Response -->|completion.Choices.Message.Content| PrintResult(<a href="/101-simple-completion/main.go#L41">Print Result</a>):::final

    classDef main fill:#e1f5fe,stroke:#01579b,stroke-width:3px,color:#000
    classDef process fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef config fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px,color:#000
    classDef env fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef client fill:#fce4ec,stroke:#880e4f,stroke-width:2px,color:#000
    classDef context fill:#f1f8e9,stroke:#33691e,stroke-width:2px,color:#000
    classDef message fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef params fill:#fff8e1,stroke:#f57f17,stroke-width:2px,color:#000
    classDef api fill:#e8eaf6,stroke:#283593,stroke-width:2px,color:#000
    classDef error fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#000
    classDef output fill:#f9fbe7,stroke:#827717,stroke-width:2px,color:#000
    classDef final fill:#e4f7ff,stroke:#0277bd,stroke-width:3px,color:#000
```

**Key Components:**

- **Environment Variables**:
  - `MODEL_RUNNER_BASE_URL`: Base URL for the Docker Model Runner
  - `COOK_MODEL`: Model identifier to use for completion

- **OpenAI Client Configuration**:
  - Base URL configured from environment variable
  - Empty API key (using local model runner)

- **Chat Completion Parameters**:
  - Message: "Quel est ton nom?" (French for "What is your name?")
  - Model: From environment variable
  - Temperature: 0.5

- **Execution Flow**:
  1. Load environment variables
  2. Create OpenAI client with custom base URL
  3. Prepare chat completion parameters
  4. Execute completion request
  5. Handle errors and print response