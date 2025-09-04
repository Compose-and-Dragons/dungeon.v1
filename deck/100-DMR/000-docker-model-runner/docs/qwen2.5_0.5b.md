La limite de tokens pour le modèle **Qwen2.5:0.5B** est de 32 768 tokens pour le contexte d’entrée et jusqu’à 8 192 tokens pour la génération de texte en sortie[1][2][3][4].

## Détails techniques

- **Contexte maximal** : 32 768 tokens (entrée)
- **Génération maximale** : 8 192 tokens (sortie)
- Cette limite concerne la version 0.5B du modèle Qwen2.5, qui est basée sur l’architecture transformer et bénéficie d’un entraînement multilingue avancé[1][2].
- Pour les modèles Qwen2.5 de plus grande taille (7B, 14B, 32B, 72B), le contexte maximal peut monter à 128 000 tokens, mais pour la version 0.5B, la limite reste à 32K tokens[3][4].

## Citations directes

- "Supports input sequences of up to 32,768 tokens and can generate output sequences of up to 8,192 tokens"[2].
- "Context Length: Full 32,768 tokens and generation 8192 tokens"[1].

## Comparatif avec autres tailles de Qwen2.5

| Modèle        | Contexte max (tokens) | Génération max (tokens) |
|---------------|----------------------|-------------------------|
| Qwen2.5-0.5B  | 32,768 [3]       | 8,192 [3]           |
| Qwen2.5-7B+   | 128,000 [3]      | 8,192 [3]           |

La version 0.5B est donc conçue pour des usages avec contextes modérés, tout en gardant une capacité d’entrée et de génération bien supérieure à de nombreux modèles de même taille[4][2].

Sources
[1] Qwen/Qwen2.5-0.5B-Instruct https://huggingface.co/Qwen/Qwen2.5-0.5B-Instruct
[2] Qwen2.5 0.5B Instruct · Models https://dataloop.ai/library/model/qwen_qwen25-05b-instruct/
[3] Qwen2.5-LLM: Extending the boundary of LLMs | Qwen https://qwenlm.github.io/blog/qwen2.5-llm/
[4] Qwen2.5-0.5B: Specifications and GPU VRAM Requirements https://apxml.com/models/qwen2-5-0-5b
[5] Qwen/Qwen2.5-0.5B https://huggingface.co/Qwen/Qwen2.5-0.5B
[6] Qwen2.5: A Party of Foundation Models! | Qwen https://qwenlm.github.io/blog/qwen2.5/
[7] Why does Qwen 2.5 support 128k context length, but ... https://www.reddit.com/r/LocalLLaMA/comments/1huvdq4/why_does_qwen_25_support_128k_context_length_but/
[8] Qwen/Qwen2.5-14B-Instruct - Hugging Face https://huggingface.co/Qwen/Qwen2.5-14B-Instruct
[9] Qwen2.5-0.5B-Instruct - m5-docs https://docs.m5stack.com/en/stackflow/models/qwen2.5-0.5b-instruct
[10] Ollama - fmenguy.fr https://fmenguy.fr/ollama/
