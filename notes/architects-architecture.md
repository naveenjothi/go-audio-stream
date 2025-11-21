# Software Architecture Summary & Checklist

## Summary

Architecture begins by noticing three things:

- **Boundaries** — where responsibilities separate.
- **Flows** — how data moves and transforms.
- **Tension** — what changes often vs. what must stay stable.

You don't begin with the code.  
You begin with:

- The **users**
- The **business**
- The **constraints**
- The **future** the system must endure

Good architecture isn't about perfection —  
It's about designing a system that can survive reality.

---

## Software Architecture Checklist

### 1. Identify Boundaries

- [ ] Have you mapped clear responsibility boundaries?
- [ ] Are modules/services cohesive with a single reason to change?
- [ ] Do boundaries align with the business domain (DDD principles)?

### 2. Understand Flows

- [ ] Is the flow of data defined end-to-end?
- [ ] Are data transformations happening in the correct layers?
- [ ] Are communication patterns intentional (sync, async, events, RPC)?

### 3. Detect Tension Points

- [ ] What components change frequently?
- [ ] What components must remain stable?
- [ ] Does the architecture isolate volatility?

### 4. Look Beyond Code

- [ ] Have you considered user behavior and expectations?
- [ ] Do you fully understand business rules and constraints?
- [ ] Are compliance, security, and legal boundaries accounted for?
- [ ] Is long-term evolution part of the design?

### 5. Ask the Core Questions

- [ ] If this component fails, what breaks?
- [ ] What needs to scale first and fastest?
- [ ] What will change the most in the next 2 years?
- [ ] Who will maintain the system, and is it maintainable?
- [ ] Are the tradeoffs explicit (latency vs. consistency, speed vs. safety)?

### 6. Test for Real-World Survivability

- [ ] Does the system degrade gracefully under failure?
- [ ] Can the design adapt to new requirements without major rewrites?
- [ ] Is technical debt intentional and acknowledged?

---
