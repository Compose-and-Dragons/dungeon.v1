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
## Quantization 🧮

- **Precision Reduction**: Technique that reduces numerical precision of model weights from 32-bit floats to lower bit representations (16-bit, 8-bit, or even 4-bit)
- **Size Benefits**: Converting from float32 to int8 can reduce model size by 4x, <span class="dodgerblue">**making large models fit on limited hardware**</span>
- **Performance Gains**: Lower precision computations are <span class="dodgerblue">**faster and use less memory**</span>, enabling quicker inference times
- **Quality Trade-off**: Slight accuracy loss in exchange for dramatic efficiency improvements, but modern techniques (GPTQ, GGML) minimize this impact
- **Deployment Enabler**: Allows running larger, more capable models on consumer hardware that couldn't handle full-precision versions


<!--
La quantization d'un modèle est une technique qui consiste à réduire la
  précision numérique des poids et activations du modèle pour diminuer sa
  taille et accélérer l'inférence.

  Au lieu d'utiliser des nombres à virgule flottante sur 32 bits (float32),
  on utilise des représentations plus compactes comme 16 bits, 8 bits, ou
  même 4 bits. Par exemple, passer de float32 à int8 divise par 4 la taille
  du modèle.

  Cette compression permet :
  - De réduire l'utilisation mémoire
  - D'accélérer les calculs
  - De faire tourner des modèles plus gros sur du matériel limité

  Le compromis est une légère perte de précision, mais les techniques
  modernes (comme GPTQ, GGML) maintiennent généralement de bonnes
  performances.
-->

[← Previous](../002-temperature.md) | [Next →](../../001-context-this-is-the-way/demos/0-bot-draft/000-personality-and-soul.md)
