---
marp: true
theme: default
paginate: true
---
<style>
.dodgerblue {
  color: dodgerblue;
}
.indianred {
  color: indianred;
}
.forestgreen {
  color: forestgreen;
}
</style>
## Parameters 📊


- **Model Capacity**: More parameters = <span class="dodgerblue">**ability to capture more complexity in data**</span>, but also require <span class="indianred">**more computational resources**</span>
- **Size Impact**: Parameter count directly affects model file size, <span class="indianred">**memory usage**</span>, and <span class="indianred">**inference speed**</span>
- **Trade-offs**: Higher parameter counts <span class="dodgerblue">**improve capability**</span> but <span class="indianred">increase hardware requirements</span> and <span class="indianred">**processing time**</span>


<!--
Les paramètres d'un modèle sont les valeurs numériques (poids et biais) que
   le modèle apprend pendant l'entraînement pour effectuer ses prédictions.

  Dans un réseau de neurones :
  - Poids : déterminent l'importance des connexions entre neurones
  - Biais : permettent d'ajuster les seuils d'activation

  Plus un modèle a de paramètres, plus il peut capturer de complexité dans
  les données, mais plus il consomme de ressources (mémoire, calcul).

  Exemples :
  - GPT-3 : 175 milliards de paramètres
  - Llama 2 7B : 7 milliards de paramètres
  - Un petit modèle : quelques millions

  Le nombre de paramètres influence directement :
  - La capacité du modèle (ce qu'il peut apprendre)
  - Sa taille sur disque
  - Le temps d'inférence
  - Les ressources nécessaires pour l'utiliser
-->

[← Previous](000-small-models-are-dumb.md) | [Next →](002-temperature.md)
