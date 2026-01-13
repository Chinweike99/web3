### What We’re Building (Job Queue System – CLI)

### A job queue where:
Users submit jobs from the CLI
Jobs go into a queue
Multiple workers process jobs concurrently
The system shuts down cleanly
No goroutine leaks
No shared-state chaos

## Core Mental Model (read this twice)
Jobs are data
Queue owns job flow
Workers consume jobs
Channels move data
Goroutines do work
Context tells everyone when to stop