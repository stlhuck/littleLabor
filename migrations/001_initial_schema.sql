-- 001_initial_schema.sql
-- Initial database schema for littleLabor

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    points INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Chores table
CREATE TABLE chores (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    frequency VARCHAR(50), -- daily, weekly, monthly, etc.
    points_value INT NOT NULL DEFAULT 10,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User-Chore assignments (many-to-many with metadata)
CREATE TABLE user_chores (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chore_id INT NOT NULL REFERENCES chores(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, chore_id)
);

-- Chore completions tracking
CREATE TABLE completions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chore_id INT NOT NULL REFERENCES chores(id) ON DELETE CASCADE,
    completed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    points_earned INT NOT NULL
);

-- Rewards table
CREATE TABLE rewards (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    point_cost INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User rewards (earned/redeemed)
CREATE TABLE user_rewards (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reward_id INT NOT NULL REFERENCES rewards(id) ON DELETE CASCADE,
    earned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    redeemed_at TIMESTAMP,
    status VARCHAR(50) DEFAULT 'earned' -- earned, redeemed
);

-- Create indexes for common queries
CREATE INDEX idx_user_chores_user_id ON user_chores(user_id);
CREATE INDEX idx_user_chores_chore_id ON user_chores(chore_id);
CREATE INDEX idx_completions_user_id ON completions(user_id);
CREATE INDEX idx_completions_chore_id ON completions(chore_id);
CREATE INDEX idx_completions_completed_at ON completions(completed_at);
CREATE INDEX idx_user_rewards_user_id ON user_rewards(user_id);
CREATE INDEX idx_user_rewards_status ON user_rewards(status);
