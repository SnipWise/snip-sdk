# Compression de contexte

La **compression de contexte** est une technique qui permet de réduire la taille des informations (tokens) envoyées à un modèle de langage tout en préservant l'essentiel du sens.

**Pourquoi c'est utile :**
- Les LLMs ont une fenêtre de contexte limitée (ex: 8k, 128k tokens)
- Réduire les coûts (facturation au token)
- Accélérer l'inférence

**Principales approches :**

1. **Résumé/distillation** : condenser un long texte en version plus courte
2. **Sélection sélective** : ne garder que les passages pertinents (via embeddings, reranking)
3. **Compression par prompts** : entraîner un modèle à compresser l'information en "soft tokens" ou représentations denses
4. **Élagage d'attention** : supprimer les tokens à faible importance selon les scores d'attention

C'est particulièrement utile pour le RAG (Retrieval-Augmented Generation) et les conversations longues.