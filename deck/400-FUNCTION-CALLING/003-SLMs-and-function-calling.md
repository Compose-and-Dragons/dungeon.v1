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
## Small Language Models (SLMs) and Function Calling

**Challenges:**
- SLMs are generally less capable at function calling
- Require careful <span class="dodgerblue">model selection</span> and <span class="dodgerblue">precise function definitions</span>

**Success Factors:**
- Clear, detailed function descriptions
- Precise parameter specifications
- => Well-structured tool definitions guide the LLM's prompt construction
