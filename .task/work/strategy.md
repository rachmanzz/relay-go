Before performing any changes (creation, modification, or deletion), you must document your strategy in this file. You must also read the rules in .task/context/readme.md before starting any task; this is a mandatory requirement that must be followed strictly. Once the strategy is documented, allow the user to validate it. You must wait for an explicit command to execute the strategy before proceeding with any implementation.

- If there is a need for API documentation, use **[./api.docs.md](./api.docs.md)**.
- If there is a need for an execution workflow, use **[./workflow.md](./workflow.md)**.

Strategies must be based on facts, not assumptions. Focus on implementing the requested logic rather than correcting the endpoint structure unless specifically asked. Most importantly, do not modify or "correct" the user's Request structure if it already fulfills the requirements; respect the user's defined request fields and types even if you would personally design them differently. If you have a suggestion for improvement (e.g., a more idiomatic or efficient way), provide it as a note in the terminal instead of applying it automatically based on assumptions.

Always perform a final check after execution to ensure no errors were introduced, all dependencies are correctly imported, and the code follows idiomatic Go patterns for standalone modules.


# Note
After the strategy has been successfully executed and the task is complete, you are required to delete the strategy entry from this file.

# Strategies
[ do not delete this line, please put your stategy down bellow]
