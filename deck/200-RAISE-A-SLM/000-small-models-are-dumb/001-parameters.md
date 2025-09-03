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
## Parameters 📊

- **Weights and Biases**: Numerical values learned during training that determine how neurons connect and activate in the neural network
- **Model Capacity**: More parameters = ability to capture more complexity in data, but also require more computational resources
- **Size Impact**: Parameter count directly affects model file size, memory usage, and inference speed
- **Scale Examples**: GPT-3 has 175B parameters, Llama 2 7B has 7B, <span class="dodgerblue">**while small models have just millions**</span>
- **Trade-offs**: Higher parameter counts improve capability but increase hardware requirements and processing time


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