## Edventures: LLM Knowledge Store

The main solution in this project is it provides a knowledge store for Large Language Models (LLMs), leveraging **Qdrant** as the vector database.

---

### Prerequisites

Before you get started, make sure you have the following installed:

* **Docker:** You'll need Docker to run the Qdrant vector database.
* **Go:** This project is written in Go, so ensure you have a Go environment set up.

---

### How to run it

#### 1. Start Qdrant Vector Database

Qdrant is crucial for efficient vector storage and retrieval. The easiest way to get it up and running is by using Docker.

First, pull the Qdrant Docker image:

```bash
docker pull qdrant/qdrant
```


Next, run the Qdrant container. This command maps the necessary ports and mounts a volume to ensure your data is persistent.
Bash

```bash

docker run -p 6333:6333 -p 6334:6334 \
    -v "$(echo $HOME)/qdrant_storage:/qdrant/storage:z" \
    qdrant/qdrant
```

-p 6333:6333: Maps the gRPC port for client communication.
-p 6334:6334: Maps the HTTP port for the REST API and web UI.
-v "$(echo $HOME)/qdrant_storage:/qdrant/storage:z": Creates a persistent volume at ~/qdrant_storage on your host machine. This means your data won't disappear if the container stops or is removed.

#### 3. Environmental Variables

Set your GOOGLE_API_KEY environmental variable. You can do this by:
Creating a .env file in your root directory and adding GOOGLE_API_KEY=your_key_here
Or by exporting it in your terminal: export GOOGLE_API_KEY=your_key_here

have the env file either by create .env file in root director or by exproting it 

#### 2. Run the Application

Once your Qdrant instance is running, navigate to the project's root directory in your terminal and execute the application:

```Bash
go run .
```

This command will start the LLM Knowledge Store application, which will automatically connect to your running Qdrant instance.

---
## Brief Write-up
---
- Its command-line interface (CLI) based AI Coach.
- The core innovation lies in its robust knowledge management system, which dynamically collects, stores, and embeds user-specific information to enhance the AI's contextual understanding and response generation.

#### Architecture Overview

Frontend (CLI): Handles user interaction, input/output, command parsing (for example, .exit, .help, .upload), and displays AI responses and progress. It manages conversational flow and interfaces with core AI and knowledge components.

Backend/Core Logic
* chatclient -  Wraps the Google Generative AI (Gemini) API, initializing the chat model with the "Anna" persona and managing message sending and response retrieval.

* knowledge - This central component manages dynamic knowledge integration.

    * Extractor Uses a separate chatclient to extract key knowledge points from user inputs (questions, discussions, uploaded documents), processing text into manageable chunks for storage.

    * Store Manages persistence and retrieval of extracted knowledge, using Qdrant (a vector database) for efficient storage and similarity search of embeddings, and the Google GenAI embedding model for creating these.

* docUpload - Enables document uploads and reviews, leveraging the knowledge system to extract and store information, and using a dedicated prompt for structured review.

* promptStore - Centralizes and manages all system prompts, personas (for example, AnnaSystemPrompt), and prompt templates (for example, KnowledgeExtractionPrompt) used across the application for consistency.

The flow begins with the main package setting up clients and services. User input first goes to knowledge.AddInputToKnowledge for extraction and storage. For general queries, the input is augmented with retrieved relevant knowledge before being sent to the LLM via chatclient.SendMessage.

#### Design Decisions and Trade-offs
* Knowledge-Centric Approach: The primary design focuses on proactively extracting, embedding, and retrieving relevant data for every query, significantly enhancing the LLM's ability to provide personalized and contextually accurate responses. This introduces increased latency due to extra API calls for embedding and vector database lookups, and adds architectural complexity.

* Separate Knowledge Extraction: LLM A dedicated Extractor component with its own chatclient is used for knowledge extraction, rather than the main conversational LLM. This increases LLM API calls, potentially leading to higher costs and more latency for initial extraction, but keeps the main conversational flow lean and focused.

* Qdrant for Vector Storage: Qdrant was chosen for its performance in similarity search and ease of Go integration. This requires a separate Qdrant instance and introduces a newer technology; PostgreSQL might have been a more mature alternative.

* Generative AI (Gemini): for Embeddings and Chat Leveraging the Google Generative AI suite offers powerful models for conversation and high-quality text embeddings. This creates a dependency on a specific cloud provider and incurs associated API costs, with performance influenced by API rate limits and network latency.


#### Ideas for Improvement if I Had More Time
* Periodically summarize dense knowledge chunks and user histories for more efficient retrieval and a better signal-to-noise ratio.

* Utilize smaller local embedding models and experiment with vector quantization and varying dimensions to reduce API calls, latency, and storage.

* Employ Go's concurrency for parallelizing embedding generation and database operations, especially with large datasets.

* Introduce multiple AI personas (e.g., Legal, Marketing) that can be switched or inferred, each with tailored prompts and knowledge.

* Guide users through complex tasks by breaking them into smaller, guided sub-tasks.

* Enhance the AI's ability to ask clarifying questions when user queries are vague.

* Improve PDF uploading to parse complex documents and extract information based on user-defined guidelines.

* Offer a transparent interface for users to view, search, and manage their stored knowledge.


#### I had too many idea, but only ten fingers :)
