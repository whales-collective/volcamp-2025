# Embeddings and Cosine Similarity Demo Schema

```mermaid
flowchart TD
    X[Embeddings Demo]:::main
    X --> Initialize(Initialize):::process

    Initialize --> OpenAIClient(<a href="/200-first-let-s-talk-about-rag/1-embeddings-distances/main.go#L29-L32">Create OpenAI Client</a>):::client

    Initialize --> Chunks(<a href="/200-first-let-s-talk-about-rag/1-embeddings-distances/main.go#L13-L18">Define Text Chunks</a>):::data
    Chunks --> ChunkAnimals[🐿️ Écureuils grimpent<br/>🐟 Truites nagent<br/>🐸 Grenouilles nagent<br/>🐰 Lapins courent]:::chunks

    Initialize --> UserQuestion(<a href="/200-first-let-s-talk-about-rag/1-embeddings-distances/main.go#L38">User Question</a>):::question
    UserQuestion --> Question["Quels sont les animaux qui nagent ?"]:::questiontext

    OpenAIClient --> EmbedQuestion(<a href="/200-first-let-s-talk-about-rag/1-embeddings-distances/main.go#L42-L50">Generate Embeddings from Question</a>):::embed
    Question --> EmbedQuestion

    OpenAIClient --> EmbedChunks(<a href="/200-first-let-s-talk-about-rag/1-embeddings-distances/main.go#L57-L74">Generate Embeddings from Chunks</a>):::embed
    ChunkAnimals --> EmbedChunks

    EmbedQuestion --> QuestionVector[Question Vector]:::vector
    EmbedChunks --> ChunkVectors[Chunk Vectors]:::vector

    QuestionVector --> CosineSim(<a href="/200-first-let-s-talk-about-rag/1-embeddings-distances/main.go#L68-L71">Calculate Cosine Similarity</a>):::calc
    ChunkVectors --> CosineSim

    CosineSim --> RagFunction[<a href="/200-first-let-s-talk-about-rag/1-embeddings-distances/main.go#L68">rag.CosineSimilarity</a>]:::ragfunc

    CosineSim --> Evaluation(<a href="/200-first-let-s-talk-about-rag/1-embeddings-distances/main.go#L77-L83">IsGoodCosineSimilarity</a>):::eval
    Evaluation --> Threshold[<a href="/200-first-let-s-talk-about-rag/1-embeddings-distances/main.go#L78">Threshold > 0.65</a>]:::threshold
    Threshold --> Result[✅ or ❌]:::result

    CosineSim --> Output(<a href="/200-first-let-s-talk-about-rag/1-embeddings-distances/main.go#L72">Display Results</a>):::output

    classDef main fill:#e1f5fe,stroke:#01579b,stroke-width:3px,color:#000
    classDef process fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef config fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px,color:#000
    classDef client fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef data fill:#fce4ec,stroke:#880e4f,stroke-width:2px,color:#000
    classDef question fill:#f1f8e9,stroke:#33691e,stroke-width:2px,color:#000
    classDef embed fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef vector fill:#fff8e1,stroke:#f57f17,stroke-width:2px,color:#000
    classDef calc fill:#e8eaf6,stroke:#283593,stroke-width:2px,color:#000
    classDef eval fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#000
    classDef output fill:#e4f7ff,stroke:#0277bd,stroke-width:3px,color:#000
    classDef chunks fill:#fde7f3,stroke:#ad1457,stroke-width:2px,color:#000
    classDef questiontext fill:#e3f2fd,stroke:#0277bd,stroke-width:2px,color:#000
    classDef ragfunc fill:#fff9c4,stroke:#f9a825,stroke-width:2px,color:#000
    classDef threshold fill:#f3e5ab,stroke:#f57f17,stroke-width:2px,color:#000
    classDef result fill:#c8e6c9,stroke:#388e3c,stroke-width:2px,color:#000
```

## Processus détaillé

### 1. Initialisation
- **Client OpenAI** : Création d'un client OpenAI avec l'URL de base personnalisée et une clé API vide

### 2. Données de test
- **Chunks prédéfinis** : 4 phrases décrivant des animaux et leurs actions :
  - Écureuils qui grimpent
  - Truites qui nagent
  - Grenouilles qui nagent
  - Lapins qui courent
- **Question utilisateur** : "Quels sont les animaux qui nagent ?"

### 3. Génération des embeddings
- **Embedding de la question** : Conversion de la question en vecteur numérique
- **Embeddings des chunks** : Conversion de chaque chunk en vecteur numérique

### 4. Calcul de similarité
- **Similarité cosinus** : Utilise `rag.CosineSimilarity()` pour comparer le vecteur de la question avec chaque vecteur de chunk
- **Évaluation** : Fonction `IsGoodCosineSimilarity()` qui retourne ✅ si similarité > 0.65, sinon ❌

### 5. Résultats attendus
Les chunks contenant "nagent" (truites et grenouilles) devraient avoir une similarité cosinus élevée avec la question sur les animaux qui nagent.

**Fonctions clés :**
- `main()` : Point d'entrée principal
- `IsGoodCosineSimilarity()` : Évaluation du seuil de similarité
- `rag.CosineSimilarity()` : Calcul de la similarité cosinus (fonction externe)