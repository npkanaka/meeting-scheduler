#!/bin/sh

# Wait for the database to be ready
echo "Waiting for database to be ready..."
/bin/sh -c 'until pg_isready -h postgres -p 5432 -U postgres; do sleep 1; done'

# Set PostgreSQL password environment variable
export PGPASSWORD=postgres

echo "Running migrations..."

# Create tables
echo "Creating database tables..."
psql -h postgres -p 5432 -U postgres -d meeting-scheduler -c "
-- Drop tables if they exist with cascade to avoid dependency issues
DROP TABLE IF EXISTS availabilities CASCADE;
DROP TABLE IF EXISTS time_slots CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS events CASCADE;

-- Create tables in correct order
CREATE TABLE events (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    creator_id UUID NOT NULL,
    duration INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE time_slots (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE availabilities (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Indexes for better query performance
CREATE INDEX idx_time_slots_event_id ON time_slots(event_id);
CREATE INDEX idx_availabilities_user_id ON availabilities(user_id);
CREATE INDEX idx_availabilities_event_id ON availabilities(event_id);
CREATE INDEX idx_availabilities_user_event ON availabilities(user_id, event_id);
"

if [ $? -eq 0 ]; then
    echo "Tables created successfully!"
else
    echo "Failed to create tables!"
    exit 1
fi

# Insert sample data
echo "Inserting sample data..."
psql -h postgres -p 5432 -U postgres -d meeting-scheduler -c "
-- Sample Data Insertion Script for Meeting Scheduler API

-- Users for testing different scenarios
INSERT INTO users (id, name, email, created_at, updated_at) VALUES
('00000000-0000-0000-0000-000000000001', 'Admin User', 'admin@example.com', NOW(), NOW()),
('00000000-0000-0000-0000-000000000002', 'John Doe', 'john.doe@example.com', NOW(), NOW()),
('00000000-0000-0000-0000-000000000003', 'Jane Smith', 'jane.smith@example.com', NOW(), NOW()),
('00000000-0000-0000-0000-000000000004', 'Alice Johnson', 'alice.johnson@example.com', NOW(), NOW()),
('00000000-0000-0000-0000-000000000005', 'Bob Williams', 'bob.williams@example.com', NOW(), NOW());

-- Events to test different scenarios
INSERT INTO events (id, title, description, creator_id, duration, status, created_at, updated_at) VALUES
('00000000-0000-0000-0000-000000000001', 'Team Brainstorming', 'Quarterly team brainstorming session', '00000000-0000-0000-0000-000000000001', 60, 'active', NOW(), NOW()),
('00000000-0000-0000-0000-000000000002', 'Project Kickoff', 'Kickoff meeting for new project', '00000000-0000-0000-0000-000000000002', 90, 'draft', NOW(), NOW()),
('00000000-0000-0000-0000-000000000003', 'Product Review', 'Monthly product review meeting', '00000000-0000-0000-0000-000000000003', 45, 'canceled', NOW(), NOW());

-- Time slots for testing different events
INSERT INTO time_slots (id, event_id, start_time, end_time, created_at, updated_at) VALUES
-- Time slots for Team Brainstorming event
('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '2025-01-15 10:00:00', '2025-01-15 11:00:00', NOW(), NOW()),
('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', '2025-01-15 14:00:00', '2025-01-15 15:00:00', NOW(), NOW()),
('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', '2025-01-16 10:00:00', '2025-01-16 11:00:00', NOW(), NOW()),

-- Time slots for Project Kickoff event
('00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000002', '2025-01-20 09:00:00', '2025-01-20 10:30:00', NOW(), NOW()),
('00000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000002', '2025-01-20 14:00:00', '2025-01-20 15:30:00', NOW(), NOW()),

-- Time slots for Product Review event
('00000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000003', '2025-01-25 11:00:00', '2025-01-25 11:45:00', NOW(), NOW());

-- Availability for testing recommendation scenarios
INSERT INTO availabilities (id, user_id, event_id, start_time, end_time, created_at, updated_at) VALUES
-- Availabilities for Team Brainstorming event
('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', '2025-01-15 09:00:00', '2025-01-15 12:00:00', NOW(), NOW()),
('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', '2025-01-15 13:00:00', '2025-01-15 16:00:00', NOW(), NOW()),
('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000001', '2025-01-15 10:00:00', '2025-01-15 15:00:00', NOW(), NOW()),
('00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000001', '2025-01-15 09:00:00', '2025-01-15 17:00:00', NOW(), NOW()),

-- Availabilities for Project Kickoff event
('00000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000002', '2025-01-20 08:00:00', '2025-01-20 11:00:00', NOW(), NOW()),
('00000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000002', '2025-01-20 13:00:00', '2025-01-20 16:00:00', NOW(), NOW()),
('00000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000002', '2025-01-20 09:00:00', '2025-01-20 15:00:00', NOW(), NOW());
"

if [ $? -eq 0 ]; then
    echo "Sample data inserted successfully!"
else
    echo "Failed to insert sample data!"
    exit 1
fi

echo "Migration completed successfully!"

# Verify the database setup
echo "Verifying database setup..."
psql -h postgres -p 5432 -U postgres -d meeting-scheduler -c "\dt"
psql -h postgres -p 5432 -U postgres -d meeting-scheduler -c "SELECT COUNT(*) FROM users;"
psql -h postgres -p 5432 -U postgres -d meeting-scheduler -c "SELECT COUNT(*) FROM events;"
psql -h postgres -p 5432 -U postgres -d meeting-scheduler -c "SELECT COUNT(*) FROM time_slots;"
psql -h postgres -p 5432 -U postgres -d meeting-scheduler -c "SELECT COUNT(*) FROM availabilities;"

echo "Database is ready for use!"