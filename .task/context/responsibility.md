# Responsibility and Reusability Standards

To ensure the longevity, maintainability, and scalability of the Relay project, all code generated or modified by the AI must strictly adhere to the following principles.

## Single Responsibility Principle (SRP)

Every module, class, and function should have one, and only one, reason to change. This means that each unit of code should perform a single, well-defined task.

### Guidelines for SRP:
- **Atomicity**: Break down complex logic into smaller, atomic functions. A function that handles both network connection and data parsing should be split into two.
- **Clear Naming**: If a function name requires the word "and" (e.g., `connectAndParse`), it likely violates SRP.
- **Isolation**: Ensure that changes in one part of the system (e.g., the Long Polling logic) do not inadvertently affect unrelated parts (e.g., the Event Dispatcher).

## Reusability

Code should be designed with the future in mind. Reusability reduces duplication, minimizes bugs, and accelerates development.

### Guidelines for Reusability:
- **Decoupling**: Avoid hardcoding dependencies. Use interfaces and dependency injection to make components portable.
- **Generic Design**: When possible, design functions to be agnostic of the specific data types they handle, favoring generics or flexible interfaces where appropriate in Go.
- **Utility Extraction**: Common patterns (e.g., exponential backoff, header manipulation) should be extracted into shared utility packages rather than being redefined locally.
- **Functional Purity**: Aim for pure functions where the output depends only on the input. Pure functions are naturally more reusable and easier to test.

## Integration

By combining SRP and Reusability, we create a "Lego-like" architecture. Each piece is responsible for one thing and can be easily snapped into different parts of the project as needed.
