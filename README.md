# Agentic Workflow Service

## Overview

The **Agentic Workflow Service** is a production-deployed, Go-based backend that demonstrates how to build an **agentic AI system** capable of:

* Understanding user intent
* Orchestrating multiple AI tools
* Coordinating retrieval and summarization microservices
* Evaluating results for grounding
* Producing structured, traceable outputs

This service acts as an **agent orchestration layer**, coordinating multiple internal AI microservices rather than embedding all AI logic in a single endpoint.

It is designed to reflect how **real production AI backends** are built: modular, observable, explainable, and incrementally extensible.

---

## Why This Service Exists

Many AI backends either:

* Call an LLM directly with no structure, or
* Hard-code a single RAG flow without reasoning, validation, or orchestration

This service demonstrates a **production-minded agentic approach** that is:

* Simple enough to understand end-to-end
* Structured enough to support real-world AI workflows

Specifically, it showcases:

* Explicit agent workflows
* Tool orchestration across services
* Grounding and evaluation checks
* Request-level tracing
* Service-to-service AI architecture

---

## High-Level Architecture

```
Client
  |
  | POST /run
  v
Agentic Workflow Service
  |
  |-- LLM (intent detection)
  |-- Tool: Semantic Search (RAG Notes API)
  |-- Tool: Summarization (AI Summary Service)
  |-- LLM (final answer synthesis)
  |-- Evaluation (grounding check)
```

Each request is treated as an **agent run**, with its own execution context and trace ID.

---

## Agent Execution Flow

Each `/run` request follows a deterministic, inspectable workflow:

1. **Intent Detection (LLM)**
   Determines what the user is trying to do and extracts entities (e.g. topic).

2. **Retrieval (Semantic Search Tool)**
   Executes semantic search against an external RAG Notes API to retrieve relevant information.

3. **Summarization (Summary Tool)**
   Passes retrieved context to a dedicated AI Summary microservice to condense and normalize information before reasoning.

4. **Final Answer Synthesis (LLM)**
   Generates a grounded, structured response using summarized and retrieved context.

5. **Evaluation**
   Scores the output for grounding and observability.

The workflow is explicit and readable to keep reasoning, debugging, and extension straightforward.

---

## Workflow Definition

The agent workflow is defined as an ordered sequence of steps:

```go
var Workflow = []Step{
    {Type: StepLLM, Prompt: "intent_decider"},
    {Type: StepTool, Tool: "search"},
    {Type: StepTool, Tool: "summarize"},
    {Type: StepLLM, Prompt: "final_answer"},
    {Type: StepEvaluate},
}
```

This design makes the agent:

* Deterministic
* Inspectable
* Easy to extend (additional tools, retries, or conditional logic)

---

## Tools

### Semantic Search Tool (RAG Notes API)

The agent integrates with an external **RAG Notes API** that provides:

* Embedding-based semantic search
* Top-K retrieval of relevant notes

The agent does **not** rely on keyword search.
Instead, it retrieves candidate information semantically and reasons over the results afterward.

---

### Summarization Tool (AI Summary Service)

Retrieved context is passed to a dedicated **AI Summary microservice**, which:

* Condenses multiple retrieved documents
* Normalizes noisy or overlapping information
* Produces a focused summary for downstream reasoning

This keeps the agent lightweight and avoids coupling summarization logic directly into the orchestration layer.

---

## Evaluation & Grounding

Before producing a final response, the agent runs an evaluation step that:

* Inspects retrieved and summarized context
* Scores how well the final answer is grounded in evidence
* Emits an evaluation score for observability

This reflects a production mindset where **LLM outputs are treated as unreliable by default** and must be validated.

---

## Observability & Tracing

Each agent run is assigned a unique `run_id`.

The service logs:

* Step-level execution
* Tool invocations
* Evaluation outcomes

Example log output:

```
[agent] run=8e29f3fb... step_type=llm step=intent_decider
[agent] run=8e29f3fb... step_type=tool step=search
[agent] run=8e29f3fb... step_type=tool step=summarize
[agent] run=8e29f3fb... step_type=llm step=final_answer
[agent] run=8e29f3fb... evaluation_score=0.60 pass=true
```

This allows each request to be traced and debugged deterministically.

---

## API

### POST /run

Runs the agent workflow for a single input.

**Request**

```json
{
  "input": "What notes do I have about Docker?"
}
```

**Response**

```json
{
  "status": "ok",
  "run_id": "8e29f3fb-c38e-41dd-a560-9a9f844cb9aa",
  "result": {
    "intent_decider": { ... },
    "search": { ... },
    "summarize": { ... },
    "evaluation_score": 0.6,
    "final_answer": { ... }
  }
}
```

---

## What This Project Demonstrates

* Agentic AI workflows
* Tool orchestration across microservices
* RAG-based retrieval
* Summarization as a first-class tool
* Grounding and evaluation
* Lightweight observability
* Clean, production-oriented Go backend design

This service is intentionally focused and composable, designed to **orchestrate AI capabilities**, not replace them.

---

## Tech Stack

* Go
* Fiber
* Docker
* Fly.io
* OpenAI API
* External RAG & Summarization services
* JSON-based agent state

---

## Final Note

This project prioritizes **clarity, correctness, and explainability** over complexity.
It reflects how real AI backend systems evol
