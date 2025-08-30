# Comment les longs contextes échouent

## Gérer votre contexte est la clé du succès des agents

Alors que les fenêtres de contexte des modèles de pointe continuent de croître¹, avec de nombreux modèles supportant jusqu'à 1 million de tokens, je vois beaucoup de discussions enthousiastes sur la façon dont les longues fenêtres de contexte vont débloquer les agents de nos rêves. Après tout, avec une fenêtre suffisamment large, vous pouvez simplement tout jeter dans un prompt dont vous pourriez avoir besoin - outils, documents, instructions, et plus - et laisser le modèle s'occuper du reste.

Les longs contextes ont paralysé l'enthousiasme pour le RAG (pas besoin de trouver le meilleur document quand vous pouvez tout faire rentrer dans le prompt !), ont permis l'engouement pour MCP (connectez-vous à tous les outils et les modèles peuvent faire n'importe quel travail !), et ont alimenté l'enthousiasme pour les agents².

Mais en réalité, des contextes plus longs ne génèrent pas de meilleures réponses. Surcharger votre contexte peut faire échouer vos agents et applications de manière surprenante. Les contextes peuvent devenir empoisonnés, distrayants, confus ou conflictuels. C'est particulièrement problématique pour les agents, qui dépendent du contexte pour rassembler des informations, synthétiser les résultats et coordonner les actions.

Passons en revue les façons dont les contextes peuvent devenir incontrôlables, puis examinons les méthodes pour atténuer ou éviter complètement les échecs de contexte.

## Échecs de contexte

### Empoisonnement du contexte

L'empoisonnement du contexte se produit lorsqu'une hallucination ou une autre erreur se glisse dans le contexte, où elle est référencée de manière répétée.

L'équipe DeepMind a souligné l'empoisonnement du contexte dans le [rapport technique de Gemini 2.5](https://storage.googleapis.com/deepmind-media/gemini/gemini_v2_5_report.pdf), que [nous avons analysé la semaine dernière](https://www.dbreunig.com/2025/06/17/an-agentic-case-study-playing-pokémon-with-gemini.html). En jouant à Pokémon, l'agent Gemini hallucinait parfois pendant qu'il jouait, empoisonnant son contexte :

> Une forme particulièrement flagrante de ce problème peut se produire avec "l'empoisonnement du contexte" - où de nombreuses parties du contexte (objectifs, résumé) sont "empoisonnées" avec de la désinformation sur l'état du jeu, ce qui peut souvent prendre très longtemps à annuler. En conséquence, le modèle peut devenir obsédé par l'atteinte d'objectifs impossibles ou non pertinents.

Si la section "objectifs" de son contexte était empoisonnée, l'agent développait des stratégies insensées et répétait des comportements en poursuivant un objectif qui ne pouvait pas être atteint.

### Distraction du contexte

La distraction du contexte se produit lorsqu'un contexte devient si long que le modèle se concentre excessivement sur le contexte, négligeant ce qu'il a appris pendant l'entraînement.

Alors que le contexte grandit pendant un flux de travail agentique - tandis que le modèle rassemble plus d'informations et construit un historique - ce contexte accumulé peut devenir distrayant plutôt qu'utile. L'agent Gemini jouant à Pokémon a clairement démontré ce problème :

> Bien que Gemini 2.5 Pro supporte un contexte de plus de 1M tokens, en faire un usage efficace pour les agents présente une nouvelle frontière de recherche. Dans cette configuration agentique, il a été observé qu'à mesure que le contexte croissait significativement au-delà de 100k tokens, l'agent montrait une tendance à favoriser la répétition d'actions de son vaste historique plutôt que de synthétiser de nouveaux plans. Ce phénomène, bien qu'anecdotique, souligne une distinction importante entre le long contexte pour la récupération et le long contexte pour le raisonnement génératif multi-étapes.

Au lieu d'utiliser son entraînement pour développer de nouvelles stratégies, l'agent devenait obsédé par la répétition d'actions passées de son historique de contexte étendu.

Pour les modèles plus petits, le plafond de distraction est beaucoup plus bas. Une [étude de Databricks](https://www.databricks.com/blog/long-context-rag-performance-llms) a trouvé que la justesse du modèle commençait à chuter autour de 32k pour Llama 3.1 405b et plus tôt pour les modèles plus petits.

Si les modèles commencent à mal se comporter bien avant que leurs fenêtres de contexte ne soient remplies, à quoi bon les très grandes fenêtres de contexte ? En résumé : la synthèse³ et la récupération de faits. Si vous ne faites ni l'un ni l'autre, méfiez-vous du plafond de distraction de votre modèle choisi.

### Confusion du contexte

La confusion du contexte se produit lorsque du contenu superflu dans le contexte est utilisé par le modèle pour générer une réponse de faible qualité.

Pendant un moment, il semblait vraiment que tout le monde allait livrer un [MCP](https://www.dbreunig.com/2025/03/18/mcps-are-apis-for-llms.html). Le rêve d'un modèle puissant, connecté à tous vos services et affaires, faisant toutes vos tâches banales semblait à portée de main. Il suffisait de jeter toutes les descriptions d'outils dans le prompt et d'y aller. [Le prompt système de Claude](https://www.dbreunig.com/2025/05/07/claude-s-system-prompt-chatbots-are-more-than-just-models.html) nous a montré la voie, car il s'agit principalement de définitions d'outils ou d'instructions pour utiliser les outils.

Mais même si [la consolidation et la concurrence ne ralentissent pas les MCP](https://www.dbreunig.com/2025/06/16/drawbridges-go-up.html), la confusion du contexte le fera. Il s'avère qu'il peut y avoir une limite au nombre d'outils.

Le [Berkeley Function-Calling Leaderboard](https://gorilla.cs.berkeley.edu/leaderboard.html) est un benchmark d'utilisation d'outils qui évalue la capacité des modèles à utiliser efficacement les outils pour répondre aux prompts. Maintenant dans sa 3e version, le classement montre que chaque modèle performe moins bien quand il dispose de plus d'un outil⁴. De plus, l'équipe de Berkeley a "conçu des scénarios où aucune des fonctions fournies n'est pertinente... nous nous attendons à ce que la sortie du modèle soit aucun appel de fonction." Pourtant, tous les modèles appelleront occasionnellement des outils qui ne sont pas pertinents.

En parcourant le classement des appels de fonction, vous pouvez voir le problème s'aggraver à mesure que les modèles deviennent plus petits :

Un exemple frappant de confusion du contexte peut être vu dans un [article récent](https://arxiv.org/pdf/2411.15399) qui a évalué la performance de petits modèles sur le [benchmark GeoEngine](https://arxiv.org/abs/2404.15500), un essai qui présente 46 outils différents. Quand l'équipe a donné à un Llama 3.1 8b quantifié (compressé) une requête avec tous les 46 outils, il a échoué, même si le contexte était bien dans la fenêtre de contexte de 16k. Mais quand ils n'ont donné au modèle que 19 outils, il a réussi.

Le problème est : si vous mettez quelque chose dans le contexte, le modèle doit y faire attention. Il peut s'agir d'informations non pertinentes ou de définitions d'outils inutiles, mais le modèle en tiendra compte. Les grands modèles, en particulier les modèles de raisonnement, deviennent meilleurs pour ignorer ou rejeter le contexte superflu, mais nous voyons continuellement des informations inutiles faire trébucher les agents. Les contextes plus longs nous permettent d'y fourrer plus d'informations, mais cette capacité s'accompagne d'inconvénients.

### Conflit du contexte

Le conflit du contexte se produit lorsque vous accumulez de nouvelles informations et outils dans votre contexte qui entrent en conflit avec d'autres informations du contexte.

C'est une version plus problématique de la confusion du contexte : le mauvais contexte ici n'est pas non pertinent, il entre directement en conflit avec d'autres informations du prompt.

Une équipe de Microsoft et Salesforce a brillamment documenté cela dans un [article récent](https://arxiv.org/pdf/2505.06120). L'équipe a pris des prompts de plusieurs benchmarks et a "fragmenté" leurs informations à travers plusieurs prompts. Pensez-y de cette façon : parfois, vous pourriez vous asseoir et taper des paragraphes dans ChatGPT ou Claude avant d'appuyer sur entrée, considérant chaque détail nécessaire. D'autres fois, vous pourriez commencer avec un prompt simple, puis ajouter d'autres détails quand la réponse du chatbot n'est pas satisfaisante. L'équipe Microsoft/Salesforce a modifié les prompts de benchmark pour ressembler à ces échanges multi-étapes :

Toutes les informations du prompt sur le côté gauche sont contenues dans les plusieurs messages sur le côté droit, qui seraient joués en plusieurs tours de chat.

Les prompts fragmentés ont donné des résultats dramatiquement pires, avec une chute moyenne de 39%. Et l'équipe a testé une gamme de modèles - le score du vanté o3 d'OpenAI a chuté de 98.1 à 64.1.

Que se passe-t-il ? Pourquoi les modèles performent-ils moins bien si l'information est rassemblée par étapes plutôt que tout d'un coup ?

La réponse est la confusion du contexte : le contexte assemblé, contenant l'intégralité de l'échange de chat, contient les premières tentatives du modèle de répondre au défi avant qu'il ait toutes les informations. Ces réponses incorrectes restent présentes dans le contexte et influencent le modèle quand il génère sa réponse finale. L'équipe écrit :

> Nous trouvons que les LLM font souvent des suppositions dans les premiers tours et tentent prématurément de générer des solutions finales, sur lesquelles ils s'appuient excessivement. En termes plus simples, nous découvrons que quand les LLM prennent un mauvais virage dans une conversation, ils se perdent et ne récupèrent pas.

Cela n'augure rien de bon pour les constructeurs d'agents. Les agents assemblent le contexte à partir de documents, d'appels d'outils et d'autres modèles chargés de sous-problèmes. Tout ce contexte, tiré de sources diverses, a le potentiel d'être en désaccord avec lui-même. De plus, quand vous vous connectez aux outils MCP que vous n'avez pas créés, il y a une plus grande chance que leurs descriptions et instructions entrent en conflit avec le reste de votre prompt.

## Conclusion

L'arrivée des fenêtres de contexte d'un million de tokens semblait transformatrice. La capacité de jeter tout ce dont un agent pourrait avoir besoin dans le prompt a inspiré des visions d'assistants superintelligents qui pourraient accéder à n'importe quel document, se connecter à chaque outil et maintenir une mémoire parfaite.

Mais comme nous l'avons vu, des contextes plus grands créent de nouveaux modes d'échec. L'empoisonnement du contexte intègre des erreurs qui se composent au fil du temps. La distraction du contexte fait que les agents s'appuient fortement sur leur contexte et répètent des actions passées plutôt que d'avancer. La confusion du contexte mène à l'utilisation d'outils ou de documents non pertinents. Le conflit du contexte crée des contradictions internes qui font dérailler le raisonnement.

Ces échecs touchent le plus durement les agents parce que les agents opèrent exactement dans les scénarios où les contextes gonflent : rassembler des informations de sources multiples, faire des appels d'outils séquentiels, s'engager dans un raisonnement multi-tours et accumuler des historiques étendus.

Heureusement, il y a des solutions ! Dans un prochain article, nous couvrirons les techniques pour atténuer ou éviter ces problèmes, des méthodes pour charger dynamiquement les outils jusqu'à la mise en place de quarantaines de contexte.

Lisez l'article de suivi, "[Comment corriger votre contexte](/2025/06/26/how-to-fix-your-context.html)"

---

### Notes de bas de page

¹ Gemini 2.5 et GPT-4.1 ont des fenêtres de contexte d'1 million de tokens, assez larges pour y jeter [Infinite Jest](https://en.wikipedia.org/wiki/Infinite_Jest), avec beaucoup de place en rab.

² La section "[Texte de forme longue](https://ai.google.dev/gemini-api/docs/long-context#long-form-text)" dans les docs Gemini résume bien cet optimisme.

³ En fait, dans l'étude Databricks citée ci-dessus, une façon fréquente dont les modèles échouaient quand on leur donnait de longs contextes était qu'ils retournaient des résumés du contexte fourni, tout en ignorant toutes les instructions contenues dans le prompt.

⁴ Si vous êtes sur le classement, faites attention aux colonnes "Live (AST)". [Ces métriques utilisent des définitions d'outils du monde réel contribuées au produit par l'entreprise](https://gorilla.cs.berkeley.edu/blogs/12_bfcl_v2_live.html), "évitant les inconvénients de la contamination de jeux de données et des benchmarks biaisés."