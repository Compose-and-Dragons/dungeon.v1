# Rapport Technique Jan-nano

**Auteurs:** Alan Dao (Gia Tuan Dao), Dinh Bach Vu

**ArXiv:** 2506.22760v1 (soumis le 28 juin 2025)

---

## Résumé

La plupart des modèles de langage font face à un compromis fondamental où des capacités puissantes nécessitent des ressources computationnelles substantielles. Nous brisons cette contrainte avec Jan-nano, un modèle de langage de 4 milliards de paramètres qui redéfinit l'efficacité grâce à une spécialisation radicale : au lieu d'essayer de tout savoir, il maîtrise l'art de trouver n'importe quoi instantanément.

Affiné à partir de Qwen3-4B en utilisant notre nouveau système multi-étapes d'Apprentissage par Renforcement avec Récompenses Vérifiables (RLVR) qui élimine complètement la dépendance à l'entraînement de prédiction du token suivant (SFT), Jan-nano atteint 83,2% sur le benchmark SimpleQA avec intégration MCP tout en fonctionnant sur du matériel grand public. Avec une longueur de contexte de 128K, Jan-nano prouve que l'intelligence ne concerne pas l'échelle, mais la stratégie.

---

## Introduction

Les modèles de langage modernes sont confrontés à un défi persistant : l'équilibre entre performance et efficacité computationnelle. Traditionnellement, l'augmentation des capacités d'un modèle nécessitait une augmentation proportionnelle des ressources computationnelles, créant des barrières à l'adoption et limitant l'accessibilité.

Jan-nano représente une approche révolutionnaire qui remet en question ce paradigme établi. Plutôt que de poursuivre une croissance aveugle en taille de modèle, nous avons développé une architecture spécialisée qui optimise l'efficacité de recherche et de récupération d'informations.

---

## Architecture et Innovation Technique

### Spécialisation Radicale

Jan-nano adopte une philosophie de "spécialisation radicale" - au lieu de tenter d'emmagasiner toutes les connaissances possibles, le modèle excelle dans l'art de localiser et récupérer des informations pertinentes de manière quasi-instantanée.

### Système RLVR (Apprentissage par Renforcement avec Récompenses Vérifiables)

L'innovation clé de Jan-nano réside dans son système multi-étapes RLVR qui élimine complètement la dépendance à l'entraînement traditionnel de prédiction du token suivant (SFT). Cette approche permet :

- **Entraînement plus ciblé** : Le modèle apprend directement à optimiser les tâches de recherche
- **Vérification des récompenses** : Les résultats peuvent être validés objectivement
- **Efficacité améliorée** : Réduction significative des ressources nécessaires

### Intégration du Protocole de Contexte de Modèle (MCP)

Jan-nano a été optimisé pour fonctionner de manière transparente avec les serveurs du Protocole de Contexte de Modèle (MCP), permettant une intégration efficace avec divers outils de recherche et sources de données.

---

## Performances et Évaluation

### Benchmark SimpleQA

Jan-nano atteint un score impressionnant de 83,2% sur le benchmark SimpleQA avec intégration MCP, démontrant ses capacités exceptionnelles en matière de :

- **Précision factuelle** : Capacité à fournir des réponses exactes et vérifiables
- **Intégration d'outils** : Utilisation efficace des ressources externes
- **Performance constante** : Maintien de la qualité à travers diverses tâches

### Contexte Étendu (128K)

La version Jan-Nano-128k propose une fenêtre de contexte native de 128k qui permet des capacités de recherche plus profondes et plus complètes sans la dégradation de performance typiquement associée aux méthodes d'extension de contexte.

**Améliorations clés :**
- 🔍 **Recherche plus approfondie** : Le contexte étendu permet le traitement de documents de recherche entiers, de documents longs et de conversations multi-tours complexes
- ⚡ **Fenêtre native 128k** : Construite dès le départ pour gérer efficacement les longs contextes, maintenant les performances sur toute la gamme de contexte

---

## Applications Pratiques

### Tâches de Recherche Profonde

Jan-Nano est un modèle de langage compact de 4 milliards de paramètres spécialement conçu et entraîné pour les tâches de recherche approfondie. Il excelle dans :

- **Analyse de documents** : Traitement et synthèse de longs documents académiques
- **Recherche multi-sources** : Intégration d'informations provenant de diverses sources
- **Conversations complexes** : Maintien du contexte sur de longues interactions

### Déploiement sur Matériel Grand Public

Avec seulement 4 milliards de paramètres, Jan-nano peut fonctionner sur du matériel grand public, rendant les capacités de recherche avancées accessibles à un public plus large.

---

## Caractéristiques Techniques

### Spécifications du Modèle

- **Paramètres** : 4,02 milliards
- **Architecture** : Basée sur Qwen3
- **Contexte** : Jusqu'à 128K tokens
- **Licence** : Apache 2.0
- **Type de tensor** : BF16

### Compatibilité et Formats

Jan-nano est disponible en plusieurs formats optimisés :

- **Quantification 3-bit** : Q3_K_S (1,89 GB), Q3_K_M (2,08 GB)
- **Quantification 4-bit** : Q4_K_S (2,38 GB), Q4_K_M (2,50 GB)
- **Quantification 8-bit** : Q8_0 (4,28 GB)
- **Format GGUF** : Compatible avec les déploiements locaux

---

## Méthodologie d'Évaluation

L'évaluation a été menée en utilisant notre approche de benchmark basée sur MCP, qui évalue les performances du modèle sur les tâches SimpleQA tout en exploitant ses capacités d'intégration native des serveurs MCP. Cette méthodologie reflète mieux les performances réelles de Jan-Nano en tant que modèle de recherche augmenté par des outils, validant à la fois sa précision factuelle et son efficacité dans l'intégration d'outils.

---

## Utilisation et Configuration

### Paramètres Recommandés

Pour une utilisation optimale de Jan-nano :

- **Température** : 0.7
- **Top-p** : 0.8
- **Top-k** : 20
- **Min-p** : 0

### Capacités Principales

Jan-nano excelle particulièrement dans :

- **Utilisation d'outils** : Excellents appels de fonction et intégration d'outils
- **Recherche** : Capacités améliorées de recherche et de traitement d'informations
- **Efficacité** : Déploiement économe en VRAM pour utilisation locale

---

## Conclusion et Impact

Jan-nano prouve que l'intelligence ne concerne pas l'échelle, mais la stratégie. En se concentrant sur la spécialisation plutôt que sur la généralisation massive, ce modèle démontre qu'il est possible d'atteindre des performances exceptionnelles avec des ressources computationnelles modérées.

Cette approche ouvre de nouvelles perspectives pour :

- **Démocratisation de l'IA** : Rendre les capacités avancées accessibles sur du matériel grand public
- **Efficacité énergétique** : Réduction significative de l'empreinte carbone des modèles IA
- **Innovation architecturale** : Nouvelles stratégies d'optimisation pour les modèles futurs

Le succès de Jan-nano suggère un changement de paradigme dans le développement des modèles de langage, privilégiant l'efficacité stratégique à la croissance brute des paramètres.

---

## Références

Pour citer ce travail :

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

**Note** : Jan-Nano est un modèle non-pensant, optimisé pour la récupération et l'intégration d'informations plutôt que pour le raisonnement génératif complexe.

---

*Document traduit de l'anglais vers le français basé sur le rapport technique arXiv:2506.22760v1*