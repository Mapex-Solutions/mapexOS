// Package runtime holds the cross-service contract types and constants for
// the workflow service runtime module.
//
// These subject constants are the wire-level contract for messages published
// by the workflow service and consumed by other services:
//   - mapexos.trigger.workflow.execute → triggers MS (TRIGGERS stream) —
//     unified workflow trigger/plugin execution dispatch.
//
// Ownership: workflow service (publisher).
// Consumers (Go): triggers.
// Reciprocity: mirrored by workspace_js/packages/schemas/src/services/workflow/runtime.
//
// Contracts stay leaf-level — no imports from services/.
package runtime

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// SubjectTriggerWorkflowExecute is the NATS subject for workflow-originated
// trigger/plugin execution requests dispatched to the triggers service.
// Resolved at package init — e.g. "dev.mapexos.trigger.workflow.execute".
var SubjectTriggerWorkflowExecute = config.Subject("trigger", "workflow.execute")

// StreamWorkflowResume is the JetStream stream for resume/callback messages
// targeting the workflow Runtime. Cross-service: published by the JS Workflow
// Executor (TS) and the Triggers service to deliver execution results back to
// the Runtime. Resolved at package init — e.g. "DEV-MAPEXOS-WORKFLOW-RESUME".
var StreamWorkflowResume = config.StreamName("WORKFLOW", "RESUME")

// SubjectWorkflowResume is the wildcard subject pattern for the workflow
// resume stream. Resolved at package init — e.g. "dev.mapexos.workflow.resume.>".
var SubjectWorkflowResume = config.Subject("workflow", "resume") + ".>"

// StreamWorkflowJSCode is the JetStream stream for JS code execution
// requests dispatched by the workflow Runtime to the JS Workflow Executor.
// Resolved at package init — e.g. "DEV-MAPEXOS-JSWORKFLOWEXECUTOR-CODE".
var StreamWorkflowJSCode = config.StreamName("JSWORKFLOWEXECUTOR", "CODE")

// SubjectWorkflowJSCode is the fixed subject for JS code execution requests
// on the workflow JS code stream. Resolved at package init — e.g.
// "dev.mapexos.workflow.js.code".
var SubjectWorkflowJSCode = config.Subject("workflow", "js.code")

// StreamWorkflowExecution is the JetStream stream for workflow execution
// commands. Cross-service: published by the router service (router-execute
// subject) and the workflow Runtime (per-instance subworkflow subjects).
// Consumed exclusively by the workflow Runtime. Resolved at package init —
// e.g. "DEV-MAPEXOS-WORKFLOW-EXECUTION".
var StreamWorkflowExecution = config.StreamName("WORKFLOW", "EXECUTION")

// SubjectWorkflowExecution is the wildcard subject pattern for the workflow
// execution stream. Resolved at package init — e.g.
// "dev.mapexos.workflow.execution.>".
var SubjectWorkflowExecution = config.Subject("workflow", "execution") + ".>"
