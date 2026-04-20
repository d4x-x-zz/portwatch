// Throttle middleware for the notify pipeline.
//
// ThrottledAlerter suppresses repeated notifications within a sliding
// time window. The first notification in a window is always forwarded.
// Subsequent notifications within the same window are held as "pending".
//
// Call Flush to drain any pending notification immediately, e.g. on
// daemon shutdown so no change is silently dropped.
package notify
