# Agentic Workflow Service

## Overview

The **Agentic Workflow Service** is a lightweight Go-based backend that demonstrates how to build an **agentic AI system** capable of:

* Understanding user intent
* Invoking tools (semantic search, summarization)
* Evaluating results for grounding
* Producing structured, traceable outputs

This service acts as an **orchestration layer**, coordinating multiple AI microservices rather than embedding all logic in a single endpoint.

It is designed to reflect how **production AI backends** are built: modular, observable, and explainable.

---

## Why This Service Exists

Many AI backends either:

* Call an LLM directly with no structure, or
* Hard-code a single RAG flow without reasoning or validation

This service demonstrates a **balanced, production-minded approach**:

* Simple enough to understand end-to-end
* Structured enough to support real-world AI workflows

It specifically showcases:

* Agent workflows
* Tool calling
* Evaluation / grounding checks
* Request-level tracing
* Service orchestration

---

## High-Level Architecture

```
Client
  |
  | POST /run
  v
Agentic Workflow Service
  |
  |-- LLM (intent + final answer)
  |-- Tool: Semantic Search (RAG Notes API)
  |-- Tool: Summarization (Summary Service)
  |-- Evaluation (grounding check)
```

Each request is treated as an **agent run**, with its own execution context and trace ID.

---

## Agent Execution Flow

Each `/run` request follows a deterministic workflow:

1. **Intent Detection (LLM)**
   Determines what the user is trying to do and extracts entities (e.g. topic).

2. **Tool Invocation**
   Tools are executed according to the workflow definition (e.g. semantic search against stored notes).

3. **Evaluation**
   The agent checks whether retrieved information supports a grounded response, helping prevent hallucinations.

4. **Final Response (LLM)**
   Generates a structured answer using retrieved evidence as context.

The workflow is explicit and readable to keep reasoning and debugging straightforward.

---

## Workflow Definition

The agent workflow is defined as an ordered sequence of steps:

```go
var Workflow = []Step{
    {Type: StepLLM, Prompt: "intent_decider"},
    {Type: StepTool, Tool: "search"},
    {Type: StepEvaluate},
    {Type: StepLLM, Prompt: "final_answer"},
}
```

This design makes the agent:

* Deterministic
* Inspectable
* Easy to extend (additional tools, retries, or branching)

---

## Tools

### Semantic Search Tool (RAG)

The agent integrates with an external **RAG Notes API** that performs:

* Embedding-based semantic search
* Top-K retrieval of relevant notes

The agent does **not** rely on keyword search.
Instead, it retrieves candidate information semantically and reasons over the results afterward.

---

## Evaluation & Grounding

Before producing a final answer, the agent runs an evaluation step that:

* Inspects retrieved search results
* Checks whether the final answer is grounded in retrieved evidence
* Produces an evaluation score

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
[agent] run=8e29f3fb... step_type=evaluate
[agent] run=8e29f3fb... evaluation_score=1.00 pass=true
[agent] run=8e29f3fb... step_type=llm step=final_answer
```

This allows each response to be traced and debugged deterministically.

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
    "evaluation_score": 1,
    "final_answer": { ... }
  }
}
```

---

## What This Project Demonstrates

* Agentic AI workflows
* Tool orchestration
* RAG integration
* Grounding & evaluation
* Lightweight observability
* Clean Go backend design

This service is intentionally focused and composable, designed to integrate with other AI microservices rather than replace them.

---

## Planned Extensions

* Centralized Generative AI service for LLM execution
* Retry logic based on evaluation results
* Additional tools (summarization, external APIs)
* Expanded tracing and metrics

---

## Tech Stack

* Go
* Fiber
* OpenAI API
* External RAG & Summarization services
* JSON-based agent state

---

## Final Note

This project prioritizes **clarity, correctness, and explainability** over complexity.
It reflects how real AI backend systems evolve incrementally in production environments.
