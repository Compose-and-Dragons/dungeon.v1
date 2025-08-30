# Comment réparer votre contexte
## Atténuer et éviter les défaillances de contexte

En complément de notre précédent article, "[Comment les longs contextes échouent](https://www.dbreunig.com/2025/06/22/how-contexts-fail-and-how-to-fix-them.html)", passons en revue les moyens d'atténuer ou d'éviter complètement ces défaillances.

Mais avant cela, récapitulons brièvement quelques-unes des façons dont les longs contextes peuvent échouer :

- **Empoisonnement du contexte** : Quand une hallucination ou une autre erreur s'immisce dans le contexte, où elle est continuellement référencée.
- **Distraction du contexte** : Quand un contexte devient si long que le modèle se concentre excessivement sur le contexte, négligeant ce qu'il a appris pendant l'entraînement.
- **Confusion du contexte** : Quand des informations superflues dans le contexte sont utilisées par le modèle pour générer une réponse de faible qualité.
- **Conflit de contexte** : Quand vous accumulez de nouvelles informations et outils dans votre contexte qui entrent en conflit avec d'autres informations dans le prompt.

Tout ceci concerne la gestion de l'information. Tout ce qui se trouve dans le contexte influence la réponse. Nous revenons au vieil adage de programmation : "[À données erronées, résultats erronés](https://en.wikipedia.org/wiki/Garbage_in,_garbage_out)." Heureusement, il existe de nombreuses options pour traiter les problèmes ci-dessus.

## Tactiques de gestion du contexte

- [RAG : Ajouter sélectivement des informations pertinentes pour aider le LLM à générer une meilleure réponse](#rag)
- [Configuration d'outils : Sélectionner uniquement les définitions d'outils pertinentes à ajouter à votre contexte](#configuration-doutils)
- [Quarantaine de contexte : Isoler les contextes dans leurs propres threads dédiés](#quarantaine-de-contexte)
- [Élagage de contexte : Supprimer les informations non pertinentes ou non nécessaires du contexte](#élagage-de-contexte)
- [Résumé de contexte : Condenser un contexte accumulé en un résumé concis](#résumé-de-contexte)
- [Déchargement de contexte : Stocker des informations en dehors du contexte du LLM, généralement via un outil qui stocke et gère les données](#déchargement-de-contexte)

### RAG

La Génération Augmentée par Récupération (RAG) est l'acte d'ajouter sélectivement des informations pertinentes pour aider le LLM à générer une meilleure réponse.

Tant de choses ont été écrites sur le RAG que nous n'allons pas le couvrir aujourd'hui au-delà de dire : il est bien vivant.

Chaque fois qu'un modèle augmente la mise de la fenêtre de contexte, un nouveau débat "RAG is Dead" naît. Le dernier événement significatif fut quand Llama 4 Scout est arrivé avec une fenêtre de 10 millions de tokens. À cette taille, il est vraiment tentant de penser : "Au diable, jetons tout dedans", et de s'en tenir là.

Mais, comme nous l'avons couvert la dernière fois : si vous traitez votre contexte comme un tiroir à bazar, le bazar [influencera votre réponse](https://www.dbreunig.com/2025/06/22/how-contexts-fail-and-how-to-fix-them.html#context-confusion). Si vous voulez en apprendre davantage, voici un [nouveau cours qui semble excellent](https://maven.com/p/569540/i-don-t-use-rag-i-just-retrieve-documents).

### Configuration d'outils

La Configuration d'outils est l'acte de sélectionner uniquement les définitions d'outils pertinentes à ajouter à votre contexte.

Le terme "loadout" (configuration) est un terme de jeu vidéo qui fait référence à la combinaison spécifique d'capacités, d'armes et d'équipement que vous sélectionnez avant un niveau, un match ou une manche. Habituellement, votre configuration est adaptée au contexte – le personnage, le niveau, la composition du reste de votre équipe, et votre propre ensemble de compétences.

Ici, nous empruntons le terme pour décrire la sélection des outils les plus pertinents pour une tâche donnée.

Peut-être que la façon la plus simple de sélectionner les outils est d'appliquer le RAG à vos descriptions d'outils. C'est exactement ce qu'ont fait Tiantian Gan et Qiyao Sun, qu'ils détaillent dans leur article, "[RAG MCP](https://arxiv.org/abs/2505.03275)." En stockant leurs descriptions d'outils dans une base de données vectorielle, ils sont capables de sélectionner les outils les plus pertinents étant donné un prompt d'entrée.

Lors du prompting de DeepSeek-v3, l'équipe a trouvé que la sélection des bons outils devient critique quand vous avez plus de 30 outils. Au-dessus de 30, les descriptions des outils commencent à se chevaucher, créant de la confusion. Au-delà de 100 outils, le modèle était virtuellement garanti d'échouer à leur test. Utiliser des techniques RAG pour sélectionner moins de 30 outils a donné des prompts dramatiquement plus courts et a résulté en une précision de sélection d'outils jusqu'à 3 fois meilleure.

Pour des modèles plus petits, les problèmes commencent bien avant d'atteindre 30 outils. Un article que nous avons touché dans le dernier post, "[Less is More](https://arxiv.org/abs/2411.15399)," a démontré que Llama 3.1 8b échoue à un benchmark quand on lui donne 46 outils, mais réussit quand on ne lui donne que 19 outils. Le problème est la confusion de contexte, pas les limitations de fenêtre de contexte.

Pour traiter ce problème, l'équipe derrière "Less is More" a développé un moyen de sélectionner dynamiquement les outils en utilisant un recommandeur d'outils alimenté par LLM. Le LLM était invité à raisonner sur le "nombre et type d'outils qu'il 'croit' nécessaires pour répondre à la requête de l'utilisateur." Cette sortie était alors recherchée sémantiquement (RAG d'outils, encore) pour déterminer la configuration finale. Ils ont testé cette méthode avec le [Berkeley Function Calling Leaderboard](https://gorilla.cs.berkeley.edu/leaderboard.html), trouvant que la performance de Llama 3.1 8b s'améliorait de 44%.

L'article "Less is More" note deux autres bénéfices aux contextes plus petits : réduction de la consommation d'énergie et vitesse, des métriques cruciales lors du fonctionnement à la périphérie (c'est-à-dire, faire tourner un LLM sur votre téléphone ou PC, pas sur un serveur spécialisé). Même quand leur méthode de sélection dynamique d'outils échouait à améliorer le résultat d'un modèle, les économies d'énergie et les gains de vitesse valaient l'effort, donnant des économies de 18% et 77%, respectivement.

Heureusement, la plupart des agents ont des surfaces plus petites qui ne nécessitent que quelques outils curatés à la main. Mais si l'étendue des fonctions ou le nombre d'intégrations doit s'étendre, considérez toujours votre configuration.

### Quarantaine de contexte

La Quarantaine de contexte est l'acte d'isoler les contextes dans leurs propres threads dédiés, chacun utilisé séparément par un ou plusieurs LLMs.

Nous voyons de meilleurs résultats quand nos contextes ne sont pas trop longs et ne contiennent pas de contenu non pertinent. Une façon d'atteindre cela est de diviser nos tâches en travaux plus petits et isolés – chacun avec son propre contexte.

Il y a [de](https://arxiv.org/abs/2402.14207) [nombreux](https://www.microsoft.com/en-us/research/articles/magentic-one-a-generalist-multi-agent-system-for-solving-complex-tasks/) exemples de cette tactique, mais un résumé accessible de cette stratégie est l'[article de blog d'Anthropic détaillant leur système de recherche multi-agents](https://www.anthropic.com/engineering/built-multi-agent-research-system). Ils écrivent :

> L'essence de la recherche est la compression : distiller des insights d'un vaste corpus. Les sous-agents facilitent la compression en opérant en parallèle avec leurs propres fenêtres de contexte, explorant différents aspects de la question simultanément avant de condenser les tokens les plus importants pour l'agent de recherche principal. Chaque sous-agent fournit aussi une séparation des préoccupations—outils distincts, prompts, et trajectoires d'exploration—ce qui réduit la dépendance au chemin et permet des investigations approfondies et indépendantes.

La recherche se prête à ce modèle de conception. Quand on donne une question, plusieurs sous-questions ou zones d'exploration peuvent être identifiées et promptées séparément en utilisant plusieurs agents. Cela accélère non seulement la collecte et distillation d'informations (s'il y a du calcul disponible), mais cela empêche chaque contexte d'accumuler trop d'informations ou d'informations non pertinentes à un prompt donné, livrant des résultats de meilleure qualité :

> Nos évaluations internes montrent que les systèmes de recherche multi-agents excellent particulièrement pour les requêtes en largeur d'abord qui impliquent de poursuivre plusieurs directions indépendantes simultanément. Nous avons trouvé qu'un système multi-agents avec Claude Opus 4 comme agent principal et des sous-agents Claude Sonnet 4 surpassait Claude Opus 4 mono-agent de 90,2% sur notre évaluation de recherche interne. Par exemple, quand on lui demande d'identifier tous les membres du conseil d'administration des entreprises du S&P 500 Technologies de l'Information, le système multi-agents a trouvé les bonnes réponses en décomposant cela en tâches pour les sous-agents, tandis que le système mono-agent a échoué à trouver la réponse avec des recherches lentes et séquentielles.

Cette approche aide aussi avec les configurations d'outils, car le concepteur d'agent peut créer plusieurs archétypes d'agents avec leur propre configuration dédiée et instructions sur comment utiliser chaque outil.

Le défi pour les constructeurs d'agents, alors, est de trouver des opportunités pour des tâches isolées à faire tourner sur des threads séparés. Les problèmes qui nécessitent le partage de contexte entre plusieurs agents ne sont pas particulièrement adaptés à cette tactique.

Si le domaine de votre agent se prête d'une manière ou d'une autre à la parallélisation, assurez-vous de [lire tout l'article d'Anthropic](https://www.anthropic.com/engineering/built-multi-agent-research-system). Il est excellent.

### Élagage de contexte

L'Élagage de contexte est l'acte de supprimer des informations non pertinentes ou non nécessaires du contexte.

Les agents accumulent du contexte alors qu'ils déclenchent des outils et assemblent des documents. Parfois, il vaut la peine de faire une pause pour évaluer ce qui a été assemblé et supprimer les déchets. Cela pourrait être quelque chose dont vous chargez votre LLM principal, ou vous pourriez concevoir un outil séparé alimenté par LLM pour réviser et éditer le contexte. Ou vous pourriez choisir quelque chose de plus adapté à la tâche d'élagage.

L'élagage de contexte a une histoire (relativement) longue, car les longueurs de contexte étaient un goulot d'étranglement plus problématique dans le domaine du traitement du langage naturel (NLP), avant ChatGPT. Une méthode d'élagage actuelle, s'appuyant sur cette histoire, est [Provence](https://arxiv.org/abs/2501.16214), "un élagueur de contexte efficace et robuste pour les Questions-Réponses."

Provence est rapide, précis, simple à utiliser, et relativement petit – seulement 1,75 GB. Vous pouvez l'appeler en quelques lignes, comme ceci :

```python
from transformers import AutoModel
provence = AutoModel.from_pretrained("naver/provence-reranker-debertav3-v1", trust_remote_code=True)

# Lire une version markdown de l'entrée Wikipedia pour Alameda, CA
with open('alameda_wiki.md', 'r', encoding='utf-8') as f:
    alameda_wiki = f.read()

# Élaguer l'article, étant donnée une question
question = 'Quelles sont mes options pour quitter Alameda ?'
provence_output = provence.process(question, alameda_wiki)
```

Provence a édité l'article, coupant 95% du contenu, ne me laissant qu'avec [ce sous-ensemble pertinent](https://gist.github.com/dbreunig/b3bdd9eb34bc264574954b2b954ebe83). Il a fait mouche.

On pourrait employer Provence ou une fonction similaire pour abattre des documents ou le contexte entier. De plus, ce modèle est un argument fort pour maintenir une version structurée¹ de votre contexte dans un dictionnaire ou autre forme, à partir de laquelle vous assemblez une chaîne compilée avant chaque appel LLM. Cette structure serait utile lors de l'élagage, permettant de s'assurer que les instructions principales et les objectifs sont préservés tandis que les sections de documents ou d'historique peuvent être élaguées ou résumées.

### Résumé de contexte

Le Résumé de contexte est l'acte de condenser un contexte accumulé en un résumé condensé.

Le Résumé de contexte est d'abord apparu comme un outil pour traiter les fenêtres de contexte plus petites. Alors que votre session de chat s'approchait de dépasser la longueur de contexte maximale, un résumé serait généré et un nouveau thread commencerait. Les utilisateurs de chatbot faisaient cela manuellement, dans ChatGPT ou Claude, demandant au bot de générer un récapitulatif court qui serait alors collé dans une nouvelle session.

Cependant, alors que les fenêtres de contexte augmentaient, les constructeurs d'agents ont découvert qu'il y a des bénéfices à la résumé au-delà de rester dans la limite de contexte totale. Alors que le contexte grandit, il devient distrayant et cause le modèle à moins s'appuyer sur ce qu'il a appris pendant l'entraînement. Nous avons appelé cela [Distraction de Contexte](https://www.dbreunig.com/2025/06/22/how-contexts-fail-and-how-to-fix-them.html#context-distraction). L'équipe derrière l'agent Gemini jouant au Pokémon a découvert que tout au-delà de 100 000 tokens déclenchait ce comportement :

> Bien que Gemini 2.5 Pro supporte un contexte de plus d'1M+ tokens, en faire un usage effectif pour les agents présente une nouvelle frontière de recherche. Dans cette configuration agentique, il a été observé que lorsque le contexte croissait significativement au-delà de 100k tokens, l'agent montrait une tendance à favoriser la répétition d'actions de son vaste historique plutôt que de synthétiser de nouveaux plans. Ce phénomène, bien qu'anecdotique, souligne une distinction importante entre le long contexte pour la récupération et le long contexte pour le raisonnement génératif multi-étapes.

Résumer votre contexte est facile à faire, mais difficile à perfectionner pour un agent donné. Savoir quelles informations doivent être préservées, et détailler cela à une étape de compression alimentée par LLM, est critique pour les constructeurs d'agents. Il vaut la peine de séparer cette fonction comme sa propre étape ou app alimentée par LLM, ce qui vous permet de collecter des données d'évaluation qui peuvent informer et optimiser cette tâche directement.

### Déchargement de contexte

Le Déchargement de contexte est l'acte de stocker des informations en dehors du contexte du LLM, généralement via un outil qui stocke et gère les données.

Cela pourrait être ma tactique favorite, ne serait-ce que parce qu'elle est si simple que vous ne croyez pas qu'elle fonctionnera.

Encore une fois, [Anthropic a un bon résumé de la technique](https://www.anthropic.com/engineering/claude-think-tool), qui détaille leur outil "think", qui est fondamentalement un brouillon :

> Avec l'outil "think", nous donnons à Claude la capacité d'inclure une étape de réflexion supplémentaire—complète avec son propre espace désigné—comme partie de l'obtention de sa réponse finale… C'est particulièrement utile lors de l'exécution de longues chaînes d'appels d'outils ou dans de longues conversations multi-étapes avec l'utilisateur.

J'apprécie vraiment la recherche et autres écrits qu'Anthropic publie, mais je ne suis pas fan du nom de cet outil. Si cet outil était appelé `scratchpad` (brouillon), vous connaîtriez sa fonction immédiatement. C'est un endroit pour que le modèle écrive des notes qui n'encombrent pas son contexte, qui sont disponibles pour référence ultérieure. Le nom "think" entre en conflit avec "[pensée étendue](https://www.anthropic.com/news/visible-extended-thinking)" et anthropomorphise inutilement le modèle...mais je digresse.

Avoir un espace pour noter des notes et des progrès fonctionne. Anthropic montre que coupler l'outil "think" avec un prompt spécifique au domaine (ce que vous feriez de toute façon dans un agent) donne des gains significatifs, jusqu'à 54% d'amélioration contre un benchmark pour des agents spécialisés.

Anthropic a identifié trois scénarios où le modèle de déchargement de contexte est utile :

- **Analyse de sortie d'outils.** Quand Claude doit soigneusement traiter la sortie d'appels d'outils précédents avant d'agir et pourrait avoir besoin de revenir en arrière dans son approche ;
- **Environnements lourds en politiques.** Quand Claude doit suivre des directives détaillées et vérifier la conformité ; et
- **Prise de décision séquentielle.** Quand chaque action s'appuie sur les précédentes et les erreurs coûtent cher (souvent trouvé dans les domaines multi-étapes).

La gestion du contexte est généralement la partie la plus difficile de la construction d'un agent. Programmer le LLM pour, comme dit Karpathy, "[emballer les fenêtres de contexte juste bien](https://x.com/karpathy/status/1937902205765607626)", déployer intelligemment les outils, informations, et maintenance régulière du contexte est le travail du concepteur d'agent.

L'insight clé à travers toutes les tactiques ci-dessus est que le contexte n'est pas gratuit. Chaque token dans le contexte influence le comportement du modèle, pour le meilleur ou pour le pire. Les fenêtres de contexte massives des LLMs modernes sont une capacité puissante, mais elles ne sont pas une excuse pour être négligent avec la gestion de l'information.

Alors que vous construisez votre prochain agent ou optimisez un existant, demandez-vous : Est-ce que tout dans ce contexte mérite sa place ? Si ce n'est pas le cas, vous avez maintenant six façons de le réparer.

---

¹ Au diable, cette liste entière de tactiques est un argument fort [vous devriez programmer vos contextes](https://www.dbreunig.com/2025/06/10/let-the-model-write-the-prompt.html).