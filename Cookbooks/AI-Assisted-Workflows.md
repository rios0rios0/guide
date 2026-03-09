# AI-Assisted Workflows

> **TL;DR:** Use the RPI methodology (Research, Plan, Implement, Review) for complex AI-assisted engineering tasks. Each phase uses a specialized mindset and produces a traceable artifact. Clear context between phases to prevent the AI from conflating investigation with implementation.

## Overview

AI coding assistants excel at well-scoped, straightforward tasks but often fail on complex, multi-file changes. The root cause is that they conflate investigation with implementation -- when asked for code, they generate plausible-sounding solutions without verifying existing patterns, API contracts, or architectural constraints.

The **RPI methodology** (Research, Plan, Implement, Review) solves this by splitting complex work into four sequential phases. Each phase has a distinct objective, produces a traceable artifact, and constrains the AI away from premature action. This is conceptualized as a **type transformation pipeline**:

```
Uncertainty → Knowledge → Strategy → Working Code → Validated Code
```

## When to Use RPI

| Use RPI When | Skip RPI When |
|---|---|
| Multi-file changes spanning 3+ files | Single-file changes under 50 lines |
| Introducing a new pattern or API | Typo fixes, log message updates |
| Integrating with an external dependency | Renaming a variable or function |
| Requirements are unclear or ambiguous | Adding a test for existing code |
| Changes touch multiple architectural layers | Updating documentation only |

## The Four Phases

### 1. Research (Uncertainty → Knowledge)

**Objective:** Transform uncertainty into verified knowledge by investigating the codebase, external APIs, and documentation.

**Constraints:**
- Do NOT write any production code during this phase
- Focus on verified truth, not plausible assumptions
- Document findings with traceable evidence (file paths, line numbers, API docs)

**Activities:**
- Search the codebase for existing implementations and patterns
- Read external API documentation and verify contracts
- Identify architectural constraints (which layer should this live in?)
- Document a single recommended approach per decision point

**Artifact:** `{{YYYY-MM-DD}}-<topic>-research.md` containing findings and recommendations.

### 2. Plan (Knowledge → Strategy)

**Objective:** Convert verified knowledge into an actionable, step-by-step implementation strategy.

**Constraints:**
- Validate that research artifacts exist before proceeding
- Reference specific files and line numbers for precision targeting
- Include validation checkpoints between steps

**Activities:**
- Break the implementation into ordered tasks with checkboxes
- Specify which files to create, modify, or delete
- Define the test strategy (what tests to write, which patterns to follow)
- Identify potential risks and mitigation steps

**Artifact:** Plan document with implementation tasks, file references, and validation steps.

### 3. Implement (Strategy → Working Code)

**Objective:** Execute the plan tasks sequentially, producing working code with verification at each step.

**Constraints:**
- Follow the plan strictly -- do not improvise or skip steps
- Verify each step before proceeding to the next
- Log all changes for the review phase

**Activities:**
- Execute plan tasks in order
- Run tests after each significant change
- Update the changes log with what was modified and why
- Pause for intermediate review if the plan specified checkpoints

**Artifact:** Modified code plus `{{YYYY-MM-DD}}-<topic>-changes.md` documenting all changes.

### 4. Review (Working Code → Validated Code)

**Objective:** Validate the implementation against research findings, the plan, and project standards.

**Constraints:**
- Compare implementation against research findings -- did we implement what we discovered?
- Check adherence to project conventions (Clean Architecture layers, BDD test structure, naming conventions)
- Run validation tools (`make lint`, `make test`, `make sast`)

**Activities:**
- Verify that all plan tasks were completed
- Run the full validation suite
- Check for standards compliance (architecture, naming, testing patterns)
- Identify follow-up work or iteration requirements

**Artifact:** `{{YYYY-MM-DD}}-<topic>-review.md` with validation results and follow-up items.

## Phase Artifacts

| Phase | Artifact | Purpose |
|---|---|---|
| Research | `*-research.md` | Verified knowledge and recommendations |
| Plan | Plan document | Step-by-step implementation strategy |
| Implement | `*-changes.md` | Change log for traceability |
| Review | `*-review.md` | Validation results and follow-up items |

Store artifacts in a `.copilot-tracking/` directory (or equivalent) at the project root. These files provide an audit trail and can be referenced in pull request descriptions.

## Integration with Team Standards

The RPI methodology complements the existing development standards:

- **Research phase:** Verify that proposed changes follow Clean Architecture layer separation and dependency direction
- **Plan phase:** Ensure the test strategy follows BDD patterns (`// given`, `// when`, `// then`) and uses the correct test doubles
- **Implement phase:** Follow the operations vocabulary (`List`, `Get`, `Insert`, `Update`, `Delete`) and language-specific conventions
- **Review phase:** Run `make lint`, `make test`, and `make sast` as defined in the CI/CD standards

## Tips for Effective Use

1. **Clear context between phases.** Start a fresh conversation or use `/clear` between phases. Research findings persist in artifact files, not in conversation history.
2. **One phase at a time.** Do not combine Research and Implementation. The quality of the output depends on the discipline of separating investigation from action.
3. **Trust the artifacts.** Each phase reads the artifacts from the previous phase. Well-written artifacts make subsequent phases faster and more accurate.
4. **Iterate when needed.** If the Review phase reveals issues, loop back to the appropriate phase rather than patching in place.
