# Rapport Technique Jan-nano

**Auteurs :** Alan Dao (Gia Tuan Dao), Dinh Bach Vu

**Date de soumission :** 28 juin 2025

**ArXiv ID :** 2506.22760v2

**Classification :** cs.CL (Computation and Language)

---

## Résumé

La plupart des modèles de langage font face à un compromis fondamental où des capacités puissantes nécessitent des ressources computationnelles substantielles. Nous brisons cette contrainte avec Jan-nano, un modèle de langage de 4 milliards de paramètres qui redéfinit l'efficacité grâce à une spécialisation radicale : au lieu d'essayer de tout connaître, il maîtrise l'art de trouver n'importe quoi instantanément.

Affiné à partir de Qwen3-4B en utilisant notre nouveau système multi-étapes d'Apprentissage par Renforcement avec Récompenses Vérifiables (RLVR) qui élimine complètement la dépendance à l'entraînement de prédiction du prochain token (SFT), Jan-nano atteint 83,2% sur le benchmark SimpleQA avec l'intégration MCP tout en fonctionnant sur du matériel grand public. Avec une longueur de contexte de 128K, Jan-nano prouve que l'intelligence n'est pas une question d'échelle, mais de stratégie.

---

## Introduction

### Le Défi des Modèles de Langage Modernes

Les modèles de langage contemporains sont confrontés à un dilemme persistant : l'équilibre entre performance et efficacité computationnelle. Les modèles les plus performants nécessitent généralement des ressources considérables, limitant leur accessibilité et leur déploiement pratique.

### Notre Approche : Spécialisation Intelligente

Jan-nano représente un changement de paradigme dans la conception des modèles de langage. Plutôt que de poursuivre une approche généraliste coûteuse en ressources, nous avons développé un modèle spécialisé dans la recherche et la récupération d'informations avec une efficacité remarquable.

---

## Architecture et Méthodologie

### Modèle de Base

- **Paramètres :** 4 milliards
- **Architecture de base :** Qwen3-4B
- **Longueur de contexte :** 128K tokens natifs
- **Optimisé pour :** Tâches de recherche approfondie

### Innovation : Système RLVR (Reinforcement Learning with Verifiable Rewards)

Jan-nano utilise une approche révolutionnaire d'entraînement qui abandonne complètement l'entraînement de prédiction du prochain token (SFT) traditionnel au profit d'un système multi-étapes d'apprentissage par renforcement avec récompenses vérifiables.

#### Caractéristiques Clés du Système RLVR :

1. **Élimination de la Dépendance SFT** : Suppression complète de l'entraînement basé sur la prédiction du prochain token
2. **Récompenses Vérifiables** : Système de validation en temps réel des performances
3. **Approche Multi-étapes** : Entraînement progressif et adaptatif

### Intégration MCP (Model Context Protocol)

Jan-nano est optimisé pour fonctionner de manière transparente avec les serveurs MCP, permettant une intégration efficace avec divers outils de recherche et sources de données.

---

## Performance et Évaluation

### Benchmark SimpleQA

**Résultat Principal :** 83,2% de réussite sur SimpleQA avec intégration MCP

Cette performance exceptionnelle démontre la capacité du modèle à :
- Traiter efficacement les requêtes de recherche
- Intégrer des sources d'information externes
- Maintenir une précision élevée malgré sa taille compacte

### Efficacité Computationnelle

- **Compatibilité :** Matériel grand public
- **VRAM :** Optimisé pour un déploiement local efficient
- **Performance :** Maintien des performances sur toute la longueur du contexte 128K

---

## Capacités Spécialisées

### 1. Recherche Approfondie

Jan-nano excelle dans le traitement de documents de recherche entiers, d'articles académiques longs et de conversations multi-tours complexes grâce à sa fenêtre de contexte étendue.

### 2. Intégration d'Outils

- **Appel de Fonctions** : Excellent support pour l'appel de fonctions
- **Intégration d'Outils** : Capacité native d'intégration avec diverses sources de données
- **Protocole MCP** : Optimisation spécifique pour les serveurs MCP

### 3. Déploiement Local

Le modèle est conçu pour être efficace en termes de mémoire, permettant un déploiement dans des environnements locaux ou embarqués.

---

## Versions et Variantes

### Jan-nano Standard
- **Paramètres :** 4,02B
- **Type de Tenseur :** BF16
- **Contexte :** Standard

### Jan-nano-128k
- **Fenêtre de Contexte Native :** 128K tokens
- **Performance :** Sans dégradation sur toute la longueur du contexte
- **Application :** Recherche approfondie et analyse de documents longs

### Formats Disponibles
- **Safetensors** : Format principal
- **GGUF** : Pour déploiement optimisé
- **Quantifications** : Multiples niveaux (3-bit, 4-bit, 5-bit, 6-bit, 8-bit)

---

## Implications et Impact

### Redéfinition de l'Intelligence Artificielle

Jan-nano démontre qu'**"l'intelligence n'est pas une question d'échelle, mais de stratégie"**. Cette approche ouvre de nouvelles perspectives pour :

1. **Démocratisation de l'IA** : Accès aux capacités avancées sans ressources massives
2. **Efficacité Énergétique** : Réduction significative de l'empreinte carbone
3. **Innovation Architecturale** : Nouveau paradigme de conception de modèles

### Applications Pratiques

- **Recherche Académique** : Assistant pour l'analyse de littérature
- **Développement Logiciel** : Intégration dans des systèmes de développement
- **Applications Embarquées** : Déploiement sur dispositifs avec ressources limitées
- **Analyse de Documents** : Traitement de documents longs et complexes

---

## Méthodologie d'Évaluation

L'évaluation de Jan-nano utilise une approche basée sur MCP qui évalue :
- **Précision Factuelle** : Validation des informations récupérées
- **Efficacité d'Intégration** : Performance avec les serveurs MCP
- **Performance en Monde Réel** : Tests dans des conditions d'utilisation réelles

---

## Conclusion

Jan-nano représente une avancée significative dans la conception de modèles de langage compacts mais puissants. En abandonnant l'approche traditionnelle "tout savoir" au profit d'une spécialisation intelligente "trouver tout", ce modèle ouvre la voie à une nouvelle génération d'assistants IA efficaces et accessibles.

### Points Clés à Retenir

1. **Innovation Méthodologique** : Le système RLVR élimine la dépendance au SFT traditionnel
2. **Performance Exceptionnelle** : 83,2% sur SimpleQA avec seulement 4B paramètres
3. **Efficacité Remarquable** : Fonctionnement sur matériel grand public
4. **Vision Stratégique** : L'intelligence par la spécialisation plutôt que par l'échelle

---

## Références et Citation

```bibtex
@misc{dao2025jannanotechnicalreport,
  title={Jan-nano Technical Report}, 
  author={Alan Dao and Dinh Bach Vu},
  year={2025},
  eprint={2506.22760},
  archivePrefix={arXiv},
  primaryClass={cs.CL},
  url={https://arxiv.org/abs/2506.22760}
}
```

---

## Disponibilité

- **Modèle Principal** : [Hugging Face - Menlo/Jan-nano](https://huggingface.co/Menlo/Jan-nano)
- **Version 128K** : [Hugging Face - Menlo/Jan-nano-128k](https://huggingface.co/Menlo/Jan-nano-128k)
- **Format GGUF** : [Hugging Face - Menlo/Jan-nano-gguf](https://huggingface.co/Menlo/Jan-nano-gguf)
- **Article Original** : [ArXiv:2506.22760](https://arxiv.org/abs/2506.22760)

---

*Note : Ce document constitue une traduction et synthèse basée sur les informations publiquement disponibles sur l'article Jan-nano Technical Report. Pour les détails techniques complets, veuillez consulter l'article original sur ArXiv.*