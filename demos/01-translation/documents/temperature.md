# Temperature et top_p

**Réglage de la “temperature” et de “top_p” pour les modèles OpenAI**

La génération de langage naturel a beaucoup progressé grâce aux modèles d’OpenAI comme GPT-3 et GPT-4. Ajuster les paramètres “temperature” et “top_p” est essentiel pour exploiter au mieux ces modèles. Ces réglages permettent de façonner la génération de texte, influençant à la fois la prévisibilité et la créativité des réponses.

**Qu’est-ce que la “temperature” ?**

La “temperature” est un réglage qui contrôle l’aléa lors du choix des mots pendant la création de texte. Des valeurs faibles rendent le texte plus prévisible et cohérent ; des valeurs élevées permettent plus de liberté et de créativité, mais peuvent aussi produire des réponses moins cohérentes.

**Qu’est-ce que “top_p” ?**

“Top_p”, ou échantillonnage par noyau (“nucleus sampling”), est un réglage qui décide combien de mots possibles le modèle va considérer. Une valeur élevée signifie que le modèle examine plus de mots, incluant des options moins probables, ce qui diversifie le texte généré.