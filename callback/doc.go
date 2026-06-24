// Package callback handles WorkTool server-to-server callbacks.
//
// WorkTool exposes two callback protocols:
//
//   - QA message callback: WorkTool POSTs user chat messages to your URL;
//     use ParseQARequest and QAAck / QATextReply to respond.
//     Configure with Robot.SetQACallback.
//
//   - Event callback: WorkTool POSTs events (command results, online/offline, etc.)
//     to your URL; use NewEventParser to decode plain JSON.
//     Configure with Robot.SetEventCallback.
package callback
