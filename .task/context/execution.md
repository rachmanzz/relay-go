# Execution and Implementation Workflow

This document establishes the mandatory protocol for moving from a task request to its technical implementation.

## Strategy-First Requirement

A comprehensive strategy must be documented before any file execution (creating, modifying, or deleting files).

1.  **Mandatory Documentation**: All proposed changes must be written to `.task/work/strategy.md`.
2.  **No Immediate Execution**: After documenting the strategy, stop and wait for a specific command (prompt) to proceed.
3.  **Exception for Inquiries**: A strategy is NOT required if the prompt is purely for informational purposes, answering questions, or tasks that do not involve modifying the codebase.

## Implementation Boundaries and Standards

To ensure the technical integrity of the Relay module, these boundaries must be strictly observed:

- **Module-Centric Architecture**: This project is a standalone module, not a framework-dependent application. Ensure code is idiomatic, lightweight, and focused on its specific functionality as an SSE/Long Polling driver.
- **Fact-Based Implementation**: Ground every change in the current codebase state. Verify paths, imports, and dependencies before applying changes.
- **Dependency Integrity**: Ensure all imports are necessary and compatible. Minimize external dependencies to keep the module portable.
- **Zero Assumption Policy**: If a requirement is ambiguous, ask for clarification instead of guessing.
- **Post-Execution Verification**: Confirm all changes solve the requested problem without introducing regressions.
- **Cleanup Requirement**: Delete the strategy entry from `.task/work/strategy.md` once the task is complete.

By following this workflow, we maintain a controlled, transparent, and high-quality development process that prioritizes user intent and system stability.
