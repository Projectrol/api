--- every noti type has same 3 settings: add, remove, mention. Next settings if specific per types
-- task: "1111" (add_to: 1, remove_from: 1, mention_in: 1, due: 1)
-- project: "111" (add_to: 1, remove_from: 1, mention_in: 1)
-- event: "111" (add_to: 1, remove_from: 1, menition_in: 1, notice_before: 1)
-- event_notice_before: notice before event occured (in minutes)
--
CREATE TABLE IF NOT EXISTS user_notifications_settings (
    id serial primary key,
    user_id serial NOT NULL UNIQUE,
    is_via_inbox boolean DEFAULT 'true',
    is_via_email boolean DEFAULT 'true',
    task_noti_settings text NOT NULL DEFAULT '1111',
    project_noti_settings text NOT NULL DEFAULT '111',
    event_noti_settings text NOT NULL DEFAULT '1111',
    event_notice_before integer DEFAULT 30,
    created_at timestamp NOT NULL DEFAULT(NOW() at time zone 'utc'),
    updated_at timestamp,
    FOREIGN KEY (user_id) REFERENCES users(id)
)