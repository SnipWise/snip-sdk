Le « context shift » de llama.cpp est le mécanisme qui permet de continuer une génération quand la fenêtre de contexte est pleine, en supprimant les tokens les plus anciens et en « faisant glisser » le reste du contexte pour libérer de la place.[3][4]

## Mécanisme de base

Quand le nombre total de tokens (prompt + historique + sortie déjà générée) atteint la taille max de contexte \(n_{\text{ctx}}\), llama.cpp va:[2][3]
- Garder un certain nombre de tokens au début (souvent le système + le début de la conversation, contrôlé par `--keep`).[4][3]
- Jeter une partie des tokens les plus anciens après cette zone « protégée » et décaler le reste, de façon à ce que de nouveaux tokens puissent être ajoutés sans dépasser \(n_{\text{ctx}}\).[7][2]

## Objectif et impact

L’objectif est de permettre des conversations « infinies » sans réévaluer tout le contexte à chaque fois, en continuant la génération même après avoir rempli la fenêtre.[6][3]
En pratique, cela signifie que le modèle perd progressivement les messages les plus anciens, ce qui peut dégrader la cohérence si la partie supprimée contenait des instructions importantes non protégées par `--keep`.[2][4]

## Activation et désactivation

En ligne de commande, il existe des options comme `--context-shift` et `--no-context-shift` (et la variable d’environnement associée) pour activer ou désactiver ce comportement.[5][8][6]
Si le context shift est désactivé, lorsque le contexte est plein, la génération s’arrête ou retourne une erreur quand la requête dépasse la fenêtre, plutôt que de faire tourner le contexte.[8][5]

## Différence avec d’autres techniques

Le context shift est un « sliding window » au niveau de la séquence de tokens de la conversation, pas une modification de l’architecture du modèle.[7][2]
C’est différent de la sliding window attention (SWA) intégrée dans certains modèles, qui limite l’attention à une fenêtre récente et peut être incompatible avec le context shift dans certains cas d’usage.[1]

[1](https://www.reddit.com/r/LocalLLaMA/comments/1nkvkle/gemma_3_27b_context_shifting_not_supported_in/)
[2](https://github.com/abetlen/llama-cpp-python/discussions/1394)
[3](https://qwen.readthedocs.io/en/latest/run_locally/llama.cpp.html)
[4](https://steelph0enix.github.io/posts/llama-cpp-guide/)
[5](https://github.com/ggml-org/llama.cpp/issues/9390)
[6](https://manpages.debian.org/unstable/llama.cpp-tools/llama-server.1.en.html)
[7](https://github.com/ggml-org/llama.cpp/issues/3969)
[8](https://github.com/ggml-org/llama.cpp/issues/12038)
[9](https://www.jan.ai/docs/desktop/llama-cpp)
[10](https://www.ajeetraina.com/how-to-increase-context-window-size-in-docker-model-runner-with-llama-cpp/)