---
marp: true
html: true
theme: default
paginate: true
mermaid: true
---
<style>
.dodgerblue {
  color: dodgerblue;
}
</style>
# Send `{Messages}` to the LLM

<div class="mermaid">
graph TD
    A[LLM Chat Completion] --> B[Messages Array]

    B --> C[System Message]
    B --> D[User Message]
    B --> E[Assistant Message]
    B --> F[Tool Message]

    style C fill:#e1f5fe
    style D fill:#f3e5f5
    style E fill:#e8f5e8
    style F fill:#fff3e0
</div>
