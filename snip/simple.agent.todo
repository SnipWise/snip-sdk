about the RAG Capabilities

	agent0 := snip.NewAgent(ctx,
		snip.AgentConfig{
			Name:               "Local Agent",
			SystemInstructions: "You are a helpful assistant.",
			ModelID:            chatModelId,
			EngineURL:          engineURL,
		},
		snip.ModelConfig{
			Temperature: 0.5,
			TopP:        0.9,
		},
		snip.EnableChatStreamFlowWithMemory(),
        snip.EnableSimilaritySearch(RagAgent), --> Not Flow ==> Actually, Change the Flow
        snip.EnableToolCalling(ToolAgent), --> Flow or Not Flow ? ==> Actually, Change the Flow
	)

With this we are doing agents composition