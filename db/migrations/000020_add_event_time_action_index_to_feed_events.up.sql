CREATE INDEX IF NOT EXISTS feed_events_event_time_action_idx ON feed_events (event_time, action);
DROP INDEX IF EXISTS feeds_event_timestamp_idx;