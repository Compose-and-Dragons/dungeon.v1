# Les Petits Modèles de Langage sont l'Avenir de l'IA Agentique

**arXiv:2506.02153v1 [cs.AI] 2 Jun 2025**

**Auteurs :** Peter Belcak¹, Greg Heinrich¹, Shizhe Diao¹, Yonggan Fu¹, Xin Dong¹, Saurav Muralidharan¹, Yingyan Celine Lin¹,², Pavlo Molchanov¹

¹NVIDIA Research, ²Georgia Institute of Technology

agents@nvidia.com

## Résumé

Les grands modèles de langage (LLM) sont souvent loués pour leur performance quasi-humaine sur une large gamme de tâches et valorisés pour leur capacité à maintenir une conversation générale. L'essor des systèmes d'IA agentique introduit cependant une multitude d'applications dans lesquelles les modèles de langage effectuent un petit nombre de tâches spécialisées de manière répétitive et avec peu de variation.

Nous présentons ici la position selon laquelle les petits modèles de langage (SLM) sont suffisamment puissants, intrinsèquement plus adaptés et nécessairement plus économiques pour de nombreuses invocations dans les systèmes agentiques, et constituent donc l'avenir de l'IA agentique. Notre argumentation s'appuie sur le niveau actuel de capacités des SLM, les architectures communes des systèmes agentiques et l'économie du déploiement des modèles de langage.

Nous soutenons en outre que dans les situations où les capacités conversationnelles généralistes sont essentielles, les systèmes agentiques hétérogènes (c'est-à-dire, des agents invoquant plusieurs modèles différents) constituent le choix naturel. Nous discutons des obstacles potentiels à l'adoption des SLM dans les systèmes agentiques et esquissons un algorithme général de conversion d'agents LLM vers SLM.

Notre position, formulée comme une déclaration de valeur, souligne l'importance de l'impact opérationnel et économique qu'un passage même partiel des LLM vers les SLM aura sur l'industrie des agents IA. Nous visons à stimuler la discussion sur l'utilisation efficace des ressources IA et espérons faire progresser les efforts pour réduire les coûts de l'IA actuelle. En appelant à des contributions et critiques de notre position, nous nous engageons à publier toute correspondance de ce type à l'adresse research.nvidia.com/labs/lpr/slm-agents.

## 1. Introduction

Le déploiement de l'intelligence artificielle agentique connaît une croissance météorique. Des enquêtes récentes montrent que plus de la moitié des grandes entreprises informatiques utilisent activement des agents IA, 21% les ayant adoptés au cours de la seule dernière année. Outre les utilisateurs, les marchés voient également une valeur économique substantielle dans les agents IA : fin 2024, le secteur de l'IA agentique avait reçu plus de 2 milliards USD de financement de startups, était évalué à 5,2 milliards USD et devrait croître à près de 200 milliards USD d'ici 2034. Pour le dire simplement, il y a une attente croissante que les agents IA joueront un rôle substantiel dans l'économie moderne.

Les composants centraux qui alimentent la plupart des agents IA modernes sont des modèles de langage (très) grands. Ce sont les LLM qui fournissent l'intelligence fondamentale permettant aux agents de prendre des décisions stratégiques sur quand et comment utiliser les outils disponibles, contrôler le flux d'opérations nécessaires pour accomplir les tâches, et, si nécessaire, décomposer les tâches complexes en sous-tâches gérables et effectuer des raisonnements pour la planification d'actions et la résolution de problèmes. Un agent IA typique communique alors simplement avec un point de terminaison API LLM choisi en effectuant des requêtes vers une infrastructure cloud centralisée qui héberge ces modèles.

Les points de terminaison API LLM sont spécifiquement conçus pour servir un grand volume de requêtes diverses en utilisant un seul LLM généraliste. Ce modèle opérationnel est profondément ancré dans l'industrie — si profondément ancré, en fait, qu'il forme la base de paris en capital substantiels : alors que le marché du service d'API LLM qui sous-tend les applications agentiques était estimé à 5,6 milliards USD en 2024, l'investissement dans l'infrastructure cloud d'hébergement a bondi à 57 milliards USD la même année. L'écart de 10 fois entre l'investissement et la taille du marché a été accepté, car on suppose que ce modèle opérationnel restera la pierre angulaire de l'industrie sans altérations substantielles, et que le grand investissement initial livrera des retours comparables aux solutions logicielles et internet traditionnelles dans les 3-4 ans.

Dans ce travail, nous reconnaissons la dominance du modèle opérationnel standard mais contestons verbalement un de ses aspects, à savoir la coutume selon laquelle les requêtes des agents pour accéder à l'intelligence linguistique sont - malgré leur simplicité comparative - gérées par des choix singleton de LLM généralistes. Nous déclarons (Section 2), argumentons (Section 3) et défendons (Section 4) la position selon laquelle les petits, plutôt que les grands, modèles de langage sont l'avenir de l'IA agentique. Nous reconnaissons cependant l'engagement commercial et la praxis désormais héritée qui cause l'état contraire du présent (Section 5). En remède, nous fournissons un aperçu d'un algorithme de conversion pour la migration d'applications agentiques des LLM vers les SLM (Section 6), et appelons à une discussion plus large (Section 7).

## 2. Position

### 2.1 Définitions

Pour concrétiser notre position, nous utilisons les définitions de travail suivantes :

**DT1** Un SLM est un modèle de langage qui peut tenir sur un appareil électronique grand public commun et effectuer des inférences avec une latence suffisamment faible pour être pratique lors du service des requêtes agentiques d'un utilisateur.

**DT2** Un LLM est un modèle de langage qui n'est pas un SLM.

Nous justifions le libellé de ces définitions en Annexe A, mais notons que leur choix a peu d'incidence sur l'essence de notre position. Nous notons qu'à partir de 2025, nous serions à l'aise de considérer la plupart des modèles de moins de 10 milliards de paramètres comme des SLM.

### 2.2 Déclaration

Nous soutenons que les SLM sont :

**V1** principalement suffisamment puissants pour gérer les tâches de modélisation linguistique des applications agentiques ;

**V2** intrinsèquement plus adaptés opérationnellement pour l'utilisation dans les systèmes agentiques que les LLM ;

**V3** nécessairement plus économiques pour la grande majorité des utilisations de modèles de langage dans les systèmes agentiques que leurs homologues LLM généralistes par la vertu de leur petite taille ;

et que sur la base des vues V1-V3, les SLM sont l'avenir de l'IA agentique.

Le phrasé de notre position est délibéré. Dans sa déclaration, nous souhaitons transmettre que le développement futur décrit est ultimement une conséquence nécessaire des différences entre les SLM et les LLM si les priorités naturelles sont suivies. Nous ne faisons pas une recommandation ou n'essayons pas d'imposer une obligation — nous faisons une déclaration de ce que nous voyons comme un reflet fidèle des valeurs de la communauté dans ce contexte.

### 2.3 Élaboration

Nous affirmons que la dominance des LLM dans la conception d'agents IA est à la fois excessive et mal alignée avec les demandes fonctionnelles de la plupart des cas d'usage agentiques. Bien que les LLM offrent une généralité impressionnante et une fluidité conversationnelle, la majorité des sous-tâches agentiques dans les systèmes agentiques déployés sont répétitives, délimitées et non-conversationnelles — appelant des modèles efficaces, prévisibles et peu coûteux. Dans ce contexte, les SLM ne suffisent pas seulement, mais sont souvent préférables. Ils offrent plusieurs avantages : latence plus faible, exigences mémoire et computationnelles réduites, et coûts opérationnels significativement plus bas, tout en maintenant une performance de tâche adéquate dans des domaines contraints.

Notre position découle d'une vision pragmatique des modèles d'usage de modèles de langage dans les architectures agentiques. Ces systèmes décomposent typiquement des objectifs complexes en sous-tâches modulaires, chacune pouvant être gérée de manière fiable par des SLM spécialisés ou affinés. Nous soutenons qu'insister sur les LLM pour toutes ces tâches reflète une mauvaise allocation de ressources computationnelles — une qui est économiquement inefficace et environnementalement insoutenable à grande échelle.

De plus, dans les cas où le raisonnement général ou le dialogue en domaine ouvert est essentiel, nous préconisons les systèmes agentiques hétérogènes, où les SLM sont utilisés par défaut et les LLM sont invoqués sélectivement et avec parcimonie. Cette composition modulaire — combinant la précision et l'efficacité des SLM avec la généralité des LLM — permet la construction d'agents à la fois rentables et capables.

Ultimement, nous observons que changer le paradigme d'architectures centrées sur les LLM vers des architectures SLM-d'abord représente pour beaucoup non seulement un raffinement technique mais aussi un devoir moral humain. Alors que la communauté IA fait face à des coûts d'infrastructure croissants et des préoccupations environnementales, adopter et normaliser l'utilisation des SLM dans les flux de travail agentiques peut jouer un rôle crucial dans la promotion d'un déploiement IA responsable et durable.

## 3. Arguments de Position

Nous soutenons les vues V1-V3 par les arguments non-exclusifs suivants.

### 3.1 Les SLM sont déjà suffisamment puissants pour l'utilisation dans les agents

**A1** Les SLM sont suffisamment puissants pour prendre la place des LLM dans les systèmes agentiques. Cet argument soutient la vue V1.

Au cours des dernières années, les capacités des petits modèles de langage ont considérablement progressé. Bien que les lois d'échelle des modèles de langage restent observées, la courbe d'échelle entre la taille du modèle et les capacités devient de plus en plus raide, impliquant que les capacités des nouveaux petits modèles de langage sont beaucoup plus proches de celles des précédents grands modèles de langage. En effet, les avancées récentes montrent que les petits modèles de langage bien conçus peuvent égaler ou dépasser la performance de tâche précédemment attribuée uniquement à des modèles beaucoup plus grands.

Voici des exemples notables :

• **Série Microsoft Phi.** Phi-2 (2,7 milliards) atteint des scores de raisonnement de bon sens et de génération de code à la hauteur des modèles de 30 milliards tout en fonctionnant ~15× plus rapidement. Phi-3 small (7 milliards) atteint une compréhension du langage et un raisonnement de bon sens à la hauteur et des scores de génération de code jusqu'aux modèles de 70 milliards de la même génération.

• **Famille NVIDIA Nemotron-H.** Les modèles hybrides Mamba-Transformer de 2/4,8/9 milliards atteignent une précision de suivi d'instructions et de génération de code comparable aux LLM denses de 30 milliards de la même génération à une fraction d'un ordre de grandeur des FLOP d'inférence.

• **Série Huggingface SmolLM2.** La famille SmolLM2 de modèles de langage compacts avec des tailles allant de 125 millions à 1,7 milliard de paramètres excelle chacun dans leur compréhension du langage, appel d'outils et performance de suivi d'instructions jusqu'aux contemporains de 14 milliards tout en égalant les modèles de 70 milliards d'il y a 2 ans.

• **NVIDIA Hymba-1.5B.** Ce SLM hybride Mamba-attention démontre la meilleure précision d'instruction et un débit de tokens 3,5× supérieur aux modèles transformer de taille comparable.

• **Série DeepSeek-R1-Distill.** DeepSeek-R1-Distill est une famille de modèles de raisonnement avec des tailles de 1,5-8 milliards, entraînés sur des échantillons générés par DeepSeek-R1. Ils démontrent de fortes capacités de raisonnement de bon sens. Notamment, le modèle DeepSeek-R1-Distill-Qwen-7B surpasse de grands modèles propriétaires comme Claude-3.5-Sonnet-1022 et GPT-4o-0513.

### 3.2 Les SLM sont plus économiques dans les systèmes agentiques

**A2** Les SLM sont plus économiques que les LLM dans les systèmes agentiques. Cet argument soutient la vue V3.

Les petits modèles offrent des avantages significatifs en efficacité coût, adaptabilité et flexibilité de déploiement. Ces avantages sont spécifiquement précieux dans les flux de travail agentiques où la spécialisation et le raffinement itératif sont critiques.

• **Efficacité d'inférence.** Servir un SLM de 7 milliards est 10-30× moins cher (en latence, consommation d'énergie et FLOP) qu'un LLM de 70-175 milliards, permettant des réponses agentiques en temps réel à grande échelle.

• **Agilité de fine-tuning.** L'affinement efficace en paramètres (ex. LoRA et DoRA) et l'affinement complet des paramètres pour les SLM nécessitent seulement quelques heures-GPU, permettant d'ajouter, corriger ou spécialiser des comportements du jour au lendemain plutôt qu'en semaines.

• **Déploiement edge.** Les avancées dans les systèmes d'inférence sur appareil comme ChatRTX démontrent l'exécution locale de SLM sur des GPU grand public, montrant une inférence agentique en temps réel, hors ligne avec une latence plus faible et un contrôle de données plus fort.

### 3.3 Les SLM sont plus flexibles

**A3** Les SLM possèdent une plus grande flexibilité opérationnelle comparés aux LLM. Cet argument soutient les vues V2 et V3.

En raison de leur petite taille et de la réduction associée des coûts de pré-entraînement et d'affinement, les SLM sont intrinsèquement plus flexibles que leurs homologues volumineux lorsqu'ils apparaissent dans les systèmes agentiques. Il devient ainsi beaucoup plus abordable et pratique d'entraîner, adapter et déployer plusieurs modèles experts spécialisés pour différentes routines agentiques.

### 3.4 Les agents exposent seulement une fonctionnalité très restreinte des modèles de langage

**A4** Les applications agentiques sont des interfaces vers un sous-ensemble limité de capacités de modèles de langage. Ceci soutient les vues V1 et V2.

Un agent IA est essentiellement une passerelle fortement instruite et orchestrée extérieurement vers un modèle de langage comportant une interface humain-ordinateur et une sélection d'outils qui, lorsqu'engagés correctement, font quelque chose d'utilitaire. De cette perspective, le grand modèle de langage sous-jacent qui était ingéniéré pour être un généraliste puissant est à travers un ensemble de prompts minutieusement écrits et une gestion de contexte méticuleusement orchestrée restreint à opérer dans une petite section de sa palette de compétences autrement large.

### 3.5 Les interactions agentiques nécessitent un alignement comportemental étroit

**A5** Les interactions agentiques nécessitent un alignement comportemental étroit. Ceci s'aligne avec la vue V2.

Un agent IA typique a des interactions fréquentes avec le code, soit par l'appel d'outils de modèles de langage soit en retournant une sortie qui doit être analysée par un morceau de code agentique qui fait l'appel de modèle de langage. Il est essentiel pour le succès de ces interactions que l'appel d'outil généré et la sortie générée se conforment aux exigences strictes de formatage imposées par l'ordre, le typage et la nature des paramètres de l'outil, et l'attente du code invoquant le modèle de langage, respectivement.

### 3.6 Les systèmes agentiques sont naturellement hétérogènes

**A6** Les systèmes agentiques permettent naturellement l'hétérogénéité dans la sélection de modèles qu'ils utilisent. Ceci s'aligne avec la vue V2.

Un modèle de langage peut lui-même être un outil appelé par un autre modèle de langage. De même, chaque fois que le code de l'agent invoque un modèle de langage, il peut, en principe, choisir n'importe quel modèle de langage. Nous soutenons qu'incorporer plusieurs modèles de langage de différentes tailles et capacités pour des requêtes ou opérations de différents niveaux de complexité offre un moyen naturel pour l'introduction des SLM.

### 3.7 Les interactions agentiques sont des voies naturelles pour rassembler des données pour l'amélioration future

**A7** Les interactions agentiques sont une bonne source de données pour l'amélioration future de modèles. Ceci soutient fondamentalement la vue V2.

Comme noté en Section 3.4, les invocations d'outils et de modèles de langage pendant un processus agentique sont souvent accompagnées d'un prompting minutieux qui concentre le modèle de langage sur la livraison de la fonctionnalité étroite requise à ce moment. Chacune de ces invocations est elle-même une source naturelle de données pour l'amélioration future.

## 4. Vues Alternatives

Les vues alternatives significatives suivantes ont été exprimées dans la littérature académique et populaire.

### 4.1 Les généralistes LLM auront toujours l'avantage d'une compréhension du langage plus générale

**VA1** Soit T une tâche unique utilisant le langage général et soient L, S un grand et un petit modèle de langage de la même génération, respectivement. La performance de L sur T l'emportera toujours sur celle de S.

Cette vue alternative dispute la vue V2 et repose sur les contre-arguments suivants :

**CA1** Il existe un corpus non-négligeable de preuves empiriques de la supériorité des grands modèles de langage en compréhension générale du langage sur les petits modèles de langage de la même génération. Les LLM acquièrent leurs capacités de compréhension du langage conformément aux lois d'échelle. Leur plus grande échelle leur permet alors de démontrer de meilleures performances à travers un large éventail de tâches spécialisées de traitement du langage naturel.

**Réfutation.** Nous croyons que le contre-argument CA1 est trop limité pour attaquer la vue V2, notamment parce que :

**A8** Les études populaires de lois d'échelle supposent que l'architecture du modèle soit gardée constante dans la même génération, alors que le travail récent sur l'entraînement de petits modèles de langage démontre qu'il y a des bénéfices de performance distincts à considérer différentes architectures pour différentes tailles de modèles.

**A9** La flexibilité des petits modèles de langage vient au secours. Un petit modèle de langage peut être facilement affiné pour la tâche T pour performer au niveau désiré de fiabilité.

### 4.2 L'inférence LLM sera toujours moins chère à cause de leur centralisation

**VA2** Le bénéfice de coût d'inférence par token de la petitesse des SLM spécialisés dans les applications agentiques est éclipsé par l'économie d'échelle.

**Reconnaissance.** Nous reconnaissons que la vue alternative VA2 est une vue valide, les considérations économiques exactes étant hautement spécifiques au cas. Nous croyons que le jury délibère encore sur la vue alternative VA2, mais que plusieurs facteurs indiquent que la vue V3 pourrait prévaloir.

### 4.3 Mondes également possibles

**VA3** Tant le monde agentique utilisant les SLM que le monde agentique utilisant les LLM sont des mondes également possibles, mais le "monde agentique LLM" a une avance considérable en termes de pratique de déploiement et d'optimisation, et l'inertie industrielle canalise déjà les efforts vers l'innovation uniquement dans cette direction.

**Reconnaissance.** Nous reconnaissons la vue alternative VA3 comme une possibilité distincte, mais maintenons la position que le poids des avantages décrits à travers les arguments A1-A7 peut plausiblement renverser l'état actuel des choses.

## 5. Obstacles à l'Adoption

Il serait prudent de se demander : Si les arguments A1-A7 sont vraiment convaincants, pourquoi les nouvelles générations d'agents ne font-elles que perpétuer le statu quo d'utilisation de LLM généralistes ?

Nous croyons que les suivants sont parmi les principaux obstacles d'aujourd'hui à l'adoption généralisée des SLM :

**B1** Grandes quantités d'investissement initial dans l'infrastructure d'inférence LLM centralisée.

**B2** Utilisation de benchmarks généralistes dans l'entraînement, conception et évaluation des SLM.

**B3** Manque de sensibilisation populaire. Les SLM ne reçoivent souvent pas le niveau d'intensité marketing et d'attention presse que reçoivent les LLM, malgré leur meilleure adéquation dans de nombreux scénarios industriels.

Nous notons que les obstacles B1-B3 sont des obstacles pratiques et loin d'être des défauts fondamentaux de la technologie SLM dans le contexte de l'IA agentique.

## 6. Algorithme de Conversion d'Agents LLM vers SLM

La nature même des applications agentiques leur permet d'éventuellement changer de l'utilisation de généralistes LLM vers l'utilisation de spécialistes SLM à beaucoup de leurs interfaces. Dans les étapes suivantes, nous esquissons un algorithme qui décrit une façon possible d'effectuer le changement du modèle sous-jacent sans douleur.

**S1 Sécuriser la collecte de données d'usage.** L'étape initiale implique de déployer l'instrumentation pour logger tous les appels d'agents non-IHM, capturant les prompts d'entrée, les réponses de sortie, les contenus des appels d'outils individuels, et optionnellement les métriques de latence pour une optimisation ciblée ultérieure.

**S2 Curation et filtrage des données.** Commencer à collecter des données à travers les pipelines de l'étape S1. Une fois qu'une quantité satisfaisante de données a été collectée (10k-100k exemples étant suffisants pour l'affinement de petits modèles en règle générale), il est nécessaire de supprimer toute PII, PHI ou toute autre donnée sensible spécifique à l'application.

**S3 Clustering de tâches.** Employer des techniques de clustering non-supervisées sur les prompts collectés et les actions d'agents pour identifier des motifs récurrents de requêtes ou opérations internes d'agents. Ces clusters aident à définir des tâches candidates pour la spécialisation SLM.

**S4 Sélection SLM.** Pour chaque tâche identifiée, sélectionner un ou plusieurs candidats SLM. Les critères de sélection incluent les capacités inhérentes du SLM, sa performance sur des benchmarks pertinents pour le type de tâche, sa licence et son empreinte de déploiement.

**S5 Affinement SLM spécialisé.** Pour chaque tâche sélectionnée et candidat SLM correspondant, préparer un ensemble de données spécifique à la tâche à partir des données curées collectées dans les étapes S2 et S3. Ensuite, affiner les SLM choisis sur ces ensembles de données spécialisés.

**S6 Itération et raffinement.** On peut réentraîner les SLM et le modèle routeur périodiquement avec de nouvelles données pour maintenir la performance et s'adapter aux motifs d'usage évolutifs.

## 7. Appel à Discussion

L'industrie de l'IA agentique montre des signes de promesse d'avoir un effet transformateur sur le travail de col blanc et au-delà.

C'est l'opinion des auteurs que toute économie de dépenses ou amélioration sur la durabilité de l'infrastructure IA agirait comme un catalyseur pour cette transformation, et qu'il est ainsi éminemment désirable d'explorer toutes les options pour le faire.

Nous appelons donc à des contributions et critiques de notre position, à être dirigées vers agents@nvidia.com, et nous engageons à publier toute correspondance de ce type à research.nvidia.com/labs/lpr/slm-agents.

---

## Annexe A : Définitions

Cette annexe fournit deux justifications pour le choix de définitions en Section 2.1.

### A.1 Argument pragmatique

Il est désirable d'avoir une définition des SLM qui répond à trois critères clés :

• **Intemporalité.** La définition devrait être intemporelle : Elle devrait éviter la dépendance à des métriques spécifiques au matériel comme le nombre de paramètres ou les FLOP, qui deviennent rapidement obsolètes avec l'avancement technologique.

• **Praticité.** La définition est susceptible d'avoir une généralité beaucoup plus large si elle est ancrée dans l'usage pratique, reflétant l'objectif du monde réel de déployer les SLM sur des appareils grand public largement disponibles.

• **Alignement de motivation.** La définition devrait capturer la motivation fondamentale qui motive l'entraînement des SLM en premier lieu, qui est de permettre des modèles de langage capables qui peuvent fonctionner sur l'appareil ou dans des budgets significativement contraints comparés aux LLM.

### A.2 Argument limite

Pour explorer la distinction entre petits et grands modèles de langage dans le contexte de l'IA agentique, adoptons la lentille intransigeante d'un extrémiste, pour qui l'intelligence doit être soit maximalement petite soit maximalement large.

Imaginons un système super-intelligent s'étendant sur des échelles galactiques, mobilisant toute la matière disponible pour optimiser ses computations. Un tel système, bien que théoriquement capable d'adresser des questions profondes, ferait face à des contraintes physiques insurmontables. La vitesse de la lumière limite la communication, avec des délais aller-retour à travers une galaxie pouvant s'étendre sur des dizaines de milliers d'années.

À l'inverse, considérons un système intelligent infiniment petit, réduit au substrat minimal capable de computation. Un tel système, semblable aux organismes biologiques les plus simples, manquerait des senseurs, effecteurs ou capacité computationnelle pour interagir de manière significative avec son environnement.

Par conséquent : Les humains, souvent considérés comme un pinacle d'intelligence, offrent un ancrage utile pour définir les SLM et LLM. Avec un ratio masse cerveau-corps dépassé seulement par de petits mammifères comme les souris, les humains équilibrent l'efficacité computationnelle avec l'embodiment pratique.

## Annexe B : Études de Cas de Remplacement LLM vers SLM

Cette annexe évalue l'étendue potentielle de remplacement des invocations de grands modèles de langage par des petits modèles de langage dans trois agents open source populaires : MetaGPT, Open Operator et Cradle.

### B.1 Étude de cas 1 : MetaGPT

**Nom :** MetaGPT **Licence :** Apache 2.0 **Objectif :** MetaGPT est un framework multi-agents conçu pour émuler une société de logiciels. Il assigne des rôles comme Gestionnaire de Produit, Architecte, Ingénieur et Ingénieur QA pour gérer collaborativement des tâches incluant la rédaction d'exigences, la conception de système, l'implémentation et les tests.

**Invocations LLM :**

• Actions basées sur les rôles. Chaque rôle d'agent invoque les LLM pour remplir ses responsabilités spécialisées (ex. codage, documentation). • Templates de prompts. Prompts structurés utilisés pour des sorties cohérentes. • Intelligence dynamique. Utilisée pour la planification, le raisonnement et l'adaptation. • Génération Augmentée par Récupération (RAG). Récupère des documents pertinents pour améliorer la génération.

**Évaluation pour le Remplacement SLM :** Les SLM seraient bien adaptés pour la génération de code de routine et les tâches de boilerplate, ainsi que pour produire des réponses structurées basées sur des templates prédéfinis. Cependant, ils nécessiteraient des données d'affinement supplémentaires pour performer de manière fiable des tâches plus complexes, comme le raisonnement architectural et la planification adaptative ou le débogage, qui bénéficieraient initialement de la compréhension contextuelle plus large et de la généralité des LLM.

**Conclusion :** Dans le cas de MetaGPT, nous estimons qu'environ 60% de ses requêtes LLM pourraient être gérées de manière fiable par des SLM appropriément spécialisés.

### B.2 Étude de cas 2 : Open Operator

**Nom :** Open Operator **Licence :** Licence MIT **Objectif :** Open Operator est un agent d'automatisation de flux de travail permettant aux utilisateurs de définir des comportements d'agents qui peuvent effectuer des tâches comme des appels API, surveillance et orchestration utilisant des outils et services.

**Invocations LLM :**

• Traitement du langage naturel. Analyse l'intention utilisateur. • Prise de décision. Guide le flux d'exécution. • Génération de contenu. Écrit des résumés, rapports.

**Évaluation pour le Remplacement SLM :** Les SLM seraient bien adaptés pour des tâches comme l'analyse et le routage de commandes simples, ainsi que la génération de messages basés sur des templates prédéfinis. Ils pourraient atteindre leurs limitations en traitant des tâches plus complexes qui nécessiteraient un raisonnement multi-étapes ou la capacité de maintenir un flux de conversation et un contexte dans le temps — domaines où les LLM continueraient d'offrir des avantages significatifs.

**Conclusion :** Dans le cas d'Open Operator, nous estimons qu'environ 40% de ses requêtes LLM pourraient être gérées de manière fiable par des SLM appropriément spécialisés.

### B.3 Étude de cas 3 : Cradle

**Nom :** Cradle **Licence :** Licence MIT **Objectif :** Cradle est conçu pour le Contrôle Général d'Ordinateur (CGO), permettant aux agents d'opérer des applications GUI via l'entrée de captures d'écran et l'interaction utilisateur simulée.

**Invocations LLM :**

• Interprétation d'interface. Comprend le contexte visuel. • Planification d'exécution de tâches. Détermine les séquences d'actions GUI. • Gestion d'erreurs. Diagnostique et réagit aux états logiciels inattendus.

**Évaluation pour le Remplacement SLM :** Les SLM seraient bien adaptés pour gérer les flux de travail d'interaction GUI répétitifs et l'exécution de séquences de clics pré-apprises. Cependant, ils feraient face à des défis concernant les tâches impliquant l'adaptation GUI dynamique ou la résolution d'erreurs non structurées, qui nécessiteraient un degré plus élevé de compréhension contextuelle typiquement fourni par les LLM.

**Conclusion :** Dans le cas de Cradle, nous estimons qu'environ 70% de ses requêtes LLM pourraient être gérées de manière fiable par des SLM appropriément spécialisés.

---

## Références

[Les références originales de l'article académique sont conservées mais non traduites car elles représentent des citations bibliographiques standard]

[1] Aashima. Small language models vs. llms: Finding the right fit for your needs, October 2024. [2] ABBYY. Small language models vs. large language models, November 2024. [3] Marah Abdin, et al. Phi-3 technical report: A highly capable language model locally on your phone. arXiv preprint arXiv:2404.14219, 2024. [4] Adyog. The economics of ai training and inference: How deepseek broke the cost curve, February 2025. [5] Ishika Agarwal, et al. Delift: Data efficient language model instruction fine tuning. arXiv preprint arXiv:2411.04425, 2024.

[...et ainsi de suite pour les 82 références restantes...]

---

_Manuscrit. En cours de révision._