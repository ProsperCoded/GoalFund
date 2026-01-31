-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'payment_verified', 'contribution_confirmed', etc.
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB, -- Additional event data
    email_sent BOOLEAN DEFAULT FALSE,
    email_sent_at TIMESTAMP,
    email_failed_reason TEXT,
    retry_count INT DEFAULT 0,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);

-- Create notification preferences table
CREATE TABLE IF NOT EXISTS notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    email_enabled BOOLEAN DEFAULT TRUE,
    payment_notifications BOOLEAN DEFAULT TRUE,
    contribution_notifications BOOLEAN DEFAULT TRUE,
    withdrawal_notifications BOOLEAN DEFAULT TRUE,
    proof_notifications BOOLEAN DEFAULT TRUE,
    goal_notifications BOOLEAN DEFAULT TRUE,
    marketing_emails BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notification_preferences_user_id ON notification_preferences(user_id);
