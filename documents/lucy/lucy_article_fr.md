# Lucy : recherche web agentique sur mobile avec vecteurs de tâches générés par machine

## Résumé

Les petits modèles de langage (SLM) sont intrinsèquement limités dans les tâches à forte intensité de connaissances en raison de leur capacité contrainte. Bien que le calcul au moment du test offre une voie vers des performances améliorées, la plupart des approches traitent le raisonnement comme un processus fixe ou heuristique. Dans ce travail, nous proposons un nouveau paradigme : considérer le raisonnement interne du modèle, délimité par les balises `<think>` et `</think>`, comme une machine à vecteurs de tâches dynamique. Plutôt que de traiter le contenu à l'intérieur de ces balises comme une simple trace de pensée, nous interprétons le processus de génération lui-même comme un mécanisme par lequel le modèle construit et affine ses propres vecteurs de tâches à la volée. Nous avons développé une méthode pour optimiser cette machine à vecteurs de tâches dynamique via RLVR et entraîné avec succès un modèle de recherche web agentique. Nous présentons Lucy, un SLM de 1,7 milliard de paramètres qui exploite ce mécanisme de raisonnement dynamique avec l'intégration MCP pour atteindre 78,3% de précision sur le benchmark SimpleQA, performant au niveau de modèles beaucoup plus grands comme DeepSeek-V3. Ceci démontre que les petits modèles peuvent rivaliser avec les grands lorsqu'ils sont équipés d'un raisonnement de tâches structuré et auto-construit.

## 1. Introduction

Les grands modèles de langage (LLM) ont démontré des capacités remarquables en compréhension et génération de langage naturel. Cependant, face à des tâches à forte intensité de connaissances nécessitant des informations à jour ou un raisonnement multi-étapes, ils échouent souvent en raison de limites de connaissance inhérentes et de dynamiques d'inférence instables. Une solution naturelle consiste à augmenter les LLM avec des outils externes, en particulier la recherche web, leur permettant de récupérer des faits actuels et de vérifier les affirmations de manière dynamique.

Cela a donné naissance au paradigme de la recherche agentique, où les modèles alternent entre raisonnement et utilisation d'outils pour résoudre des requêtes complexes. Cependant, un défi critique demeure : comment stabiliser le processus de raisonnement qui guide le comportement de recherche.

Bien que des modèles comme ToolFormer et ReAct montrent que les LLM peuvent apprendre à appeler des outils, leurs trajectoires de raisonnement sont souvent incohérentes, redondantes ou divergentes, surtout sur plusieurs tours. Comme le soulignent Wu et al. (2025), le défi central dans la recherche d'informations à long horizon ne vient pas de la capacité d'accéder aux outils de recherche, mais plutôt de la difficulté de construire des trajectoires de raisonnement fiables, c'est-à-dire décomposer efficacement les tâches complexes tout en préservant la cohérence avec l'objectif global.

Dans ce travail, nous soutenons que la base d'un tel raisonnement existe déjà dans le modèle. Les LLM modernes ne sont pas seulement des appariateurs de motifs, ils construisent implicitement des représentations de tâches internes pendant l'inférence. Cette idée est soutenue par le concept de vecteurs de tâches, qui ont montré encoder des informations spécifiques aux tâches dans l'espace latent même avant que la génération de tokens commence.

Nous étendons cette insight : plutôt que de voir le processus de raisonnement comme une trace passive, nous traitons la génération dans les balises `<think>` comme une machine à vecteurs de tâches dynamique, un mécanisme computationnel auto-modifiant par lequel le modèle construit, met à jour et affine activement sa représentation de tâche au moment du test.

## Contributions principales

- **Une recette pour optimiser le vecteur de tâche** : Nous introduisons des stratégies d'entraînement et des contraintes architecturales qui guident le modèle à générer un raisonnement plus stable et auto-cohérent dans la boucle `<think>`, lui permettant efficacement d'affiner sa propre représentation de tâche pendant l'inférence.

- **Analyse théorique et empirique de la stabilité des vecteurs de tâches** : Nous analysons comment un processus de raisonnement bien structuré réduit le bruit et améliore l'attribution du crédit dans l'utilisation d'outils multi-étapes.

- **Preuve que les petits modèles peuvent maîtriser le comportement agentique** : Nous démontrons qu'un modèle de 1,7B paramètres peut atteindre de fortes performances via l'apprentissage par renforcement.

- **Défier l'hypothèse du goulot d'étranglement des données** : Nous démontrons qu'un comportement agentique de haute qualité peut émerger dans un petit modèle de 1,7B entraîné sur un ensemble de données limité et isolé.

## 2. Travaux connexes

Notre travail fait le pont entre trois domaines de recherche critiques pour développer des agents de recherche d'informations robustes, avec un accent particulier sur la stabilisation du processus de raisonnement via l'optimisation des vecteurs de tâches.

### 2.1 Raisonnement dynamique dans les modèles de langage

Des travaux récents ont révélé que les LLM construisent des représentations de tâches internes pendant l'inférence, un phénomène formalisé par les vecteurs de tâches. Ces représentations latentes gouvernent la trajectoire de raisonnement du modèle, particulièrement dans les contextes agentiques où les modèles alternent entre penser et agir. Bien que des approches comme IRCoT et ReAct démontrent la valeur des traces de raisonnement explicites, elles souffrent souvent d'incohérence dans les scénarios multi-tours.

### 2.2 Génération augmentée par recherche

L'intégration de la récupération avec les LLM a évolué à travers deux paradigmes : (1) les systèmes de Génération Augmentée par Récupération (RAG) qui souffrent de limitations de récupération statique, et (2) les approches d'utilisation d'outils dynamiques. Alors que les méthodes RAG luttent avec un contexte non pertinent, les modèles augmentés par outils font face à des défis pour maintenir la cohérence de raisonnement à travers plusieurs étapes de recherche.

### 2.3 Apprentissage par renforcement pour systèmes agentiques

S'appuyant sur les fondations RLHF, les avancées récentes ont simplifié l'optimisation de politique par des méthodes comme DPO et GRPO. Cependant, ces approches optimisent typiquement pour les résultats finaux plutôt que la stabilité de raisonnement.

## 3. Méthodologie

Notre méthodologie est fondée sur les principes architecturaux du projet Jan-Nano, mais diverge significativement dans son approche au processus de raisonnement du modèle. Nous employons un système d'apprentissage par renforcement multi-étapes qui contourne le besoin de fine-tuning supervisé (SFT) et intègre un serveur RAG local pour la récupération d'informations en temps réel.

Là où Jan-Nano offre une approche "sans réflexion" en générant directement des réponses, notre travail se concentre sur la rétention et l'optimisation du raisonnement par chaîne de pensée (CoT) du modèle. Un défi principal dans les modèles de langage utilisant le raisonnement par chaîne de pensée est la tendance à la "surréflexion", où le modèle génère des étapes de raisonnement excessivement verbeuses ou redondantes.

### 3.1 Récompense composite

Notre cadre d'apprentissage par renforcement emploie une fonction de récompense multi-composants qui combine plusieurs termes de récompense spécialisés pour façonner de manière exhaustive le comportement de l'agent.

#### 3.1.1 Composants de récompense fondamentaux

Le signal d'entraînement principal comprend des termes de récompense principaux qui imposent la correction, la conformité structurelle et l'utilisation appropriée des outils :

- **Récompense de correction (R_correct)** : C'est le signal objectif principal. Pour la tâche QA, c'est une récompense binaire (1 ou 0) déterminée par une correspondance de sous-chaîne entre la réponse générée par le modèle et la vérité terrain.

- **Récompense de validité XML (R_xml)** : Pour assurer l'intégrité structurelle, cette fonction vérifie les balises XML bien formées (`<think>`, `<tool_call>`, `<answer>`). Elle pénalise les balises déséquilibrées ou les générations logiquement incohérentes.

- **Récompense d'adhérence au format (R_format)** : Cette fonction évalue la conformité globale de la sortie du modèle avec le schéma de réponse prédéfini.

- **Récompense d'exécution d'outils (R_tool_exec)** : Une récompense est fournie pour exécuter avec succès un outil appelé.

#### 3.1.2 Fonctions de récompense centrées sur le comportement

En plus des objectifs standard de correction et de formatage, nous incorporons des signaux de récompense auxiliaires qui incitent à un comportement de recherche efficace et un raisonnement ciblé :

**Ratio Visite/Recherche** : Pour promouvoir une utilisation judicieuse de l'outil de recherche, nous appliquons une fonction de récompense basée sur le ratio entre les visites de documents et les requêtes de recherche émises.

**Réflexion efficace** : Pour promouvoir un raisonnement concis mais efficace, nous introduisons un signal de récompense basé sur la longueur de la portée de raisonnement interne du modèle, délimitée par les balises `<think>` et `</think>`. La récompense est définie en utilisant une distribution normale asymétrique centrée à 35 tokens avec une asymétrie négative, encourageant un raisonnement informatif mais succinct.

### 3.2 Cadre d'apprentissage par renforcement en deux étapes

Notre entraînement est divisé en deux étapes distinctes, chacune avec un objectif spécifique :

**Étape 1 : Entraînement fondamental** pour la correction et la maîtrise des outils. L'étape initiale est conçue pour enseigner au modèle les compétences fondamentales de correction et d'utilisation efficace des outils tout en façonnant doucement son comportement de raisonnement.

**Étape 2 : Fine-tuning** pour l'adhérence au format et la maximisation de la précision. Après avoir atteint la convergence à l'Étape 1, le modèle entre dans une phase de raffinement visant à assurer la conformité avec les contraintes de sortie et l'optimisation des performances.

## 4. Expériences et résultats

### 4.1 Expériences

Nous avons évalué notre modèle sur l'ensemble de données SimpleQA, suivant le protocole d'évaluation établi dans le projet Jan-Nano. Pour évaluer les capacités pratiques d'utilisation d'outils des modèles, nous avons intégré notre code d'évaluation avec un serveur MCP (Model Context Protocol).

Lucy atteint **78,3% de précision** malgré son nombre compact de 1,7B paramètres. Cela représente une amélioration substantielle de 19,1 points de pourcentage par rapport à la baseline de 4B paramètres et égale les performances du modèle significativement plus grand DeepSeek-67B (78,2%).

| Modèle | SimpleQA | Paramètres |
|--------|----------|------------|
| OpenAI o1 | 42,6% | Inconnu |
| Grok 3 | 44,6% | Inconnu |
| o3 | 49,4% | Inconnu |
| Claude-3.7-Sonnet | 50,0% | Inconnu |
| Gemini-2.5 Pro | 52,9% | Inconnu |
| ChatGPT-4.5 | 62,5% | Inconnu |
| **Avec MCP :** | | |
| DeepSeek-671B | 78,2% | 671B |
| **Lucy (mode think)** | **78,3%** | **1,7B** |
| Jan-nano | 80,7% | 4B |
| Jan-nano-128k | 83,2% | 4B |

Notre découverte clé démontre une efficacité exceptionnelle des paramètres, validant notre approche d'entraînement : en préservant et optimisant le raisonnement par chaîne de pensée via des récompenses ciblées sur le comportement, nous élicitons avec succès un raisonnement complexe et une utilisation d'outils à partir d'une architecture compacte.

## 5. Discussion

### 5.1 Saut émergent de la réflexion redondante

Pendant l'optimisation du processus de raisonnement de Lucy, nous avons observé une adaptation intriguante dans un scénario représentatif impliquant des tâches de recherche multi-étapes. Lors de la pénalisation des portées de raisonnement excessivement longues, le modèle a appris à omettre stratégiquement les étapes de réflexion pendant les opérations à faible décision.

Par exemple, dans une tâche de recherche en 5 étapes, le modèle alterne entre recherche et lecture :

**Comportement de base :** Délibération complète à chaque étape
- `<think>` Rechercher sur le sujet `</think>`
- `<think>` Lire le résultat `</think>`
- `<think>` Lire le résultat `</think>`
- `<think>` Rechercher des détails `</think>`

**Comportement optimisé :** Réflexion supprimée pour la lecture
- `<think>` Rechercher sur le sujet `</think>`
- Lire le résultat (pas de balises de réflexion)
- Lire le résultat (pas de balises de réflexion)
- `<think>` Rechercher des détails `</think>`

Ce comportement spécifique suggère que Lucy peut allouer dynamiquement la capacité de raisonnement basée sur la prévisibilité de l'action. Le modèle réservait la réflexion pour les opérations à haute incertitude (ex. formulation de requête) tout en contournant la délibération pour les actions déterministes (ex. traitement du contenu récupéré).

### 5.2 L'illusion du calcul au moment du test

L'hypothèse commune que "plus de réflexion améliore les performances" s'avère peu fiable pour les petits modèles. Notre benchmarking a révélé des cas d'échec où Lucy ne pouvait pas formuler de requêtes correctes non pas à cause de données internet manquantes, mais parce qu'elle manquait de connaissances fondamentales sur les entités interrogées—démontrant que le calcul au moment du test ne peut pas compenser un ancrage conceptuel manquant.

Les modèles plus grands comme Jan-nano dépassent cette limitation en s'auto-corrigeant pendant la recherche, tandis que Lucy (1,7B) circule souvent dans des hypothèses incorrectes. Le processus d'entraînement a à peine suivi le rythme de la baseline de Jan-nano, suggérant que les petits modèles atteignent des barrières de connaissances fondamentales qu'aucune quantité de raisonnement ne peut surmonter.

## 6. Conclusion

Lucy démontre que les petits modèles de langage, lorsqu'équipés de techniques de raisonnement structuré et d'apprentissage par renforcement, peuvent rivaliser avec des systèmes beaucoup plus grands sur des tâches complexes à forte intensité de connaissances. En recadrant le processus de raisonnement interne du modèle comme une machine à vecteurs de tâches dynamique, nous permettons à Lucy de construire et d'affiner itérativement ses objectifs pendant l'inférence.

Notre utilisation des balises `<think>` à la fois comme échafaudage et cible d'optimisation—combinée avec des récompenses centrées sur le comportement et un format de dialogue XML structuré—résulte en un modèle de recherche web agentique hautement efficace.

Malgré sa taille modeste de 1,7B paramètres, Lucy atteint 78,3% de précision sur le benchmark SimpleQA sous paramètres MCP, égalant des modèles plusieurs centaines de fois plus grands. Ces résultats défient les hypothèses conventionnelles autour de l'échelle, des exigences de données et du raisonnement au moment du test.

Ultimement, Lucy suggère une nouvelle direction pour les petits modèles capables : non pas par l'échelle brutale des paramètres, mais en entraînant les modèles à mieux penser—pas plus longtemps. Nous espérons que ce travail inspirera une exploration plus poussée des systèmes agentiques légers optimisés pour l'interaction réelle et l'utilisation d'outils.