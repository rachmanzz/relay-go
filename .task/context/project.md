# Relay: SSE Driver & Long Polling Module

Relay is a robust communication library designed to provide reliable real-time data streaming by combining Server-Sent Events (SSE) with a sophisticated Long Polling fallback mechanism. It acts as a high-level driver that ensures continuous connectivity between client and server, even in challenging network environments.

## Core Architecture

Relay is built on two primary communication pillars:

1.  **SSE Driver (Primary):** By default, Relay utilizes Server-Sent Events (SSE) to establish a unidirectional, persistent connection. This allows the server to push real-time updates to the client with low overhead and minimal latency.
2.  **Long Polling Module (Secondary/Fallback):** In scenarios where SSE is restricted or unstable, Relay transitions to its Long Polling module. This module acts as the backend for the SSE interface, ensuring that the application remains functional even when persistent HTTP connections are interrupted or blocked by proxies, firewalls, or legacy infrastructure.

## Intelligent Connection Management

The hallmark of Relay is its ability to maintain a stable connection through adaptive switching:

*   **SSE as the Preferred Path:** Relay always attempts to establish and maintain an SSE connection first, as it is the most efficient method for real-time updates.
*   **Automatic Fallback (Auto-Switch):** By default, Relay monitors the health of the SSE stream. If it detects that the connection is frequently dropping, failing to handshake, or unable to remain stable, it automatically switches the underlying transport to Long Polling. This process is seamless and ensures no data is lost during the transition.
*   **Manual Control:** While the auto-switch logic is highly optimized, developers have the flexibility to manually define the preferred transport. This is useful for specific environments where the network characteristics are known in advance (e.g., restricted corporate environments where only standard HTTP requests are allowed).

## Why Long Polling?

In Relay, Long Polling is not just a legacy fallback; it is a meticulously engineered backend for the SSE interface. It is specifically triggered only when SSE cannot maintain a stable and connected state. This hybrid approach guarantees the "always-on" feel of real-time applications while providing the maximum possible compatibility across various network configurations.
