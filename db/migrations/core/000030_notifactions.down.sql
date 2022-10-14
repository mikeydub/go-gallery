DROP INDEX IF EXISTS notification_owner_id_idx;
DROP TABLE IF EXISTS notifications;

ALTER TABLE users DROP COLUMN IF EXISTS notification_settings;
ALTER TABLE events DROP COLUMN IF EXISTS gallery_id;
ALTER TABLE events DROP COLUMN IF EXISTS comment_id;
ALTER TABLE events DROP COLUMN IF EXISTS admire_id;
ALTER TABLE events DROP COLUMN IF EXISTS feed_event_id;
ALTER TABLE events DROP COLUMN IF EXISTS fingerprint;