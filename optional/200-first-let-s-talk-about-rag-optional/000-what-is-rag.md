---
marp: true
theme: default
paginate: true
---
<style>
.dodgerblue {
  color: dodgerblue;
}
</style>

# RAG: retrieval augmented generation
> Breaking information into chunks

## Make context smaller and more relevant

---
## Helping the model remember large amounts of information
- Language models have a **context limit** (e.g., 4K, 16K, [32K tokens](https://qwenlm.github.io/blog/qwen2.5-llm/#model-card))
- By breaking information into **chunks**, we can **store** and **retrieve**
  information **beyond** this limit.
- <span class="dodgerblue">**Better focus**</span>
- <span class="dodgerblue">**Less data to load**</span>


<!-- 
Expliquer
- **Better focus**
- **Less data to load**
-->