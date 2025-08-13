-- Educational Platform Schema
-- Migration 011: Comprehensive educational content and course management system

-- Instructors table for course creators
CREATE TABLE IF NOT EXISTS instructors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    title VARCHAR(255),
    bio TEXT,
    avatar VARCHAR(500),
    credentials TEXT[],
    experience TEXT,
    rating DECIMAL(3,2) DEFAULT 0.00,
    course_count INTEGER DEFAULT 0,
    total_students INTEGER DEFAULT 0,
    total_revenue DECIMAL(15,2) DEFAULT 0.00,
    social_links JSONB DEFAULT '{}',
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Courses table for educational courses
CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(50) NOT NULL, -- beginner, intermediate, advanced, expert
    level VARCHAR(50) NOT NULL, -- crypto_basics, trading_fundamentals, ai_trading, advanced_strategies
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) DEFAULT 'USD',
    duration_minutes INTEGER DEFAULT 0, -- estimated completion time in minutes
    prerequisites TEXT[], -- required course IDs or skills
    learning_path VARCHAR(100), -- structured learning sequence
    instructor_id UUID NOT NULL REFERENCES instructors(id) ON DELETE CASCADE,
    thumbnail VARCHAR(500),
    trailer_video VARCHAR(500),
    certificate BOOLEAN DEFAULT false,
    status VARCHAR(20) NOT NULL DEFAULT 'draft', -- draft, published, archived
    enrollment_count INTEGER DEFAULT 0,
    completion_count INTEGER DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0.00,
    review_count INTEGER DEFAULT 0,
    total_revenue DECIMAL(15,2) DEFAULT 0.00,
    tags TEXT[],
    language VARCHAR(10) DEFAULT 'en',
    difficulty_level INTEGER DEFAULT 1, -- 1-5 scale
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    published_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}'
);

-- Course modules for organizing course content
CREATE TABLE IF NOT EXISTS course_modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    module_order INTEGER NOT NULL,
    duration_minutes INTEGER DEFAULT 0,
    is_locked BOOLEAN DEFAULT false,
    prerequisites TEXT[], -- required module IDs
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(course_id, module_order)
);

-- Course lessons within modules
CREATE TABLE IF NOT EXISTS course_lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id UUID NOT NULL REFERENCES course_modules(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    lesson_type VARCHAR(20) NOT NULL, -- video, text, interactive, live, quiz
    content JSONB NOT NULL DEFAULT '{}', -- video_url, text_content, interactive_data
    duration_minutes INTEGER DEFAULT 0,
    lesson_order INTEGER NOT NULL,
    is_preview BOOLEAN DEFAULT false,
    transcript TEXT,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(module_id, lesson_order)
);

-- Course resources (PDFs, templates, tools)
CREATE TABLE IF NOT EXISTS course_resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    lesson_id UUID REFERENCES course_lessons(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    resource_type VARCHAR(20) NOT NULL, -- pdf, link, tool, template, video
    url VARCHAR(500) NOT NULL,
    description TEXT,
    file_size BIGINT DEFAULT 0,
    download_count INTEGER DEFAULT 0,
    is_free BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Course enrollments
CREATE TABLE IF NOT EXISTS course_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'enrolled', -- enrolled, in_progress, completed, dropped, refunded
    progress DECIMAL(5,2) DEFAULT 0.00, -- 0-100%
    start_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completion_date TIMESTAMP WITH TIME ZONE,
    last_accessed TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    time_spent_minutes INTEGER DEFAULT 0,
    payment_status VARCHAR(20) DEFAULT 'pending', -- pending, paid, failed, refunded
    payment_amount DECIMAL(10,2) DEFAULT 0.00,
    payment_method VARCHAR(20), -- stripe, crypto, free
    payment_reference VARCHAR(255),
    enrolled_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, course_id)
);

-- User lesson progress tracking
CREATE TABLE IF NOT EXISTS lesson_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    lesson_id UUID NOT NULL REFERENCES course_lessons(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'not_started', -- not_started, in_progress, completed
    progress DECIMAL(5,2) DEFAULT 0.00, -- 0-100%
    time_spent_minutes INTEGER DEFAULT 0,
    last_position INTEGER DEFAULT 0, -- for video lessons, position in seconds
    completed_at TIMESTAMP WITH TIME ZONE,
    first_accessed TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_accessed TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, lesson_id)
);

-- Course quizzes and assessments
CREATE TABLE IF NOT EXISTS course_quizzes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    module_id UUID REFERENCES course_modules(id) ON DELETE CASCADE,
    lesson_id UUID REFERENCES course_lessons(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    quiz_type VARCHAR(20) DEFAULT 'module_quiz', -- module_quiz, final_exam, practice_test
    time_limit_minutes INTEGER DEFAULT 0, -- 0 = no limit
    passing_score DECIMAL(5,2) DEFAULT 70.00,
    max_attempts INTEGER DEFAULT 3,
    is_required BOOLEAN DEFAULT false,
    randomize_questions BOOLEAN DEFAULT true,
    show_correct_answers BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Quiz questions
CREATE TABLE IF NOT EXISTS quiz_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_id UUID NOT NULL REFERENCES course_quizzes(id) ON DELETE CASCADE,
    question_type VARCHAR(20) NOT NULL, -- multiple_choice, true_false, fill_blank, essay, matching
    question_text TEXT NOT NULL,
    options JSONB DEFAULT '[]', -- for multiple choice questions
    correct_answer JSONB NOT NULL,
    explanation TEXT,
    points INTEGER DEFAULT 1,
    difficulty VARCHAR(10) DEFAULT 'medium', -- easy, medium, hard
    question_order INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Quiz attempts and results
CREATE TABLE IF NOT EXISTS quiz_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    quiz_id UUID NOT NULL REFERENCES course_quizzes(id) ON DELETE CASCADE,
    attempt_number INTEGER NOT NULL,
    score DECIMAL(5,2) DEFAULT 0.00,
    max_score INTEGER NOT NULL,
    percentage DECIMAL(5,2) DEFAULT 0.00,
    passed BOOLEAN DEFAULT false,
    time_taken_minutes INTEGER DEFAULT 0,
    answers JSONB DEFAULT '{}', -- question_id -> answer mapping
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(user_id, quiz_id, attempt_number)
);

-- Course assignments
CREATE TABLE IF NOT EXISTS course_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    module_id UUID REFERENCES course_modules(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    instructions TEXT NOT NULL,
    due_date TIMESTAMP WITH TIME ZONE,
    max_points INTEGER DEFAULT 100,
    rubric JSONB DEFAULT '[]', -- grading criteria
    submission_type VARCHAR(20) DEFAULT 'text', -- text, file, url, code
    is_required BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Assignment submissions
CREATE TABLE IF NOT EXISTS assignment_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id UUID NOT NULL REFERENCES course_assignments(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL,
    content TEXT,
    attachments TEXT[], -- file URLs
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    grade INTEGER,
    max_grade INTEGER,
    feedback TEXT,
    graded_by VARCHAR(255),
    graded_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'submitted', -- submitted, graded, returned
    UNIQUE(assignment_id, user_id)
);

-- Course certificates
CREATE TABLE IF NOT EXISTS course_certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    certificate_type VARCHAR(20) DEFAULT 'completion', -- completion, achievement, mastery
    title VARCHAR(255) NOT NULL,
    description TEXT,
    issued_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    verification_code VARCHAR(100) UNIQUE NOT NULL,
    digital_badge_url VARCHAR(500),
    pdf_url VARCHAR(500),
    blockchain_hash VARCHAR(255), -- for blockchain verification
    is_public BOOLEAN DEFAULT true,
    UNIQUE(user_id, course_id)
);

-- Course reviews and ratings
CREATE TABLE IF NOT EXISTS course_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    review_text TEXT,
    is_verified_purchase BOOLEAN DEFAULT false,
    helpful_votes INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, course_id)
);

-- Educational content library
CREATE TABLE IF NOT EXISTS content_library (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    content_type VARCHAR(20) NOT NULL, -- article, video, podcast, webinar, tool, template
    content JSONB NOT NULL DEFAULT '{}', -- type-specific content data
    author_id UUID REFERENCES instructors(id),
    category VARCHAR(50) NOT NULL,
    tags TEXT[],
    description TEXT,
    thumbnail VARCHAR(500),
    duration_minutes INTEGER DEFAULT 0,
    views INTEGER DEFAULT 0,
    likes INTEGER DEFAULT 0,
    downloads INTEGER DEFAULT 0,
    is_free BOOLEAN DEFAULT true,
    price DECIMAL(10,2) DEFAULT 0.00,
    published_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Live sessions and webinars
CREATE TABLE IF NOT EXISTS live_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    instructor_id UUID NOT NULL REFERENCES instructors(id),
    course_id UUID REFERENCES courses(id), -- optional, can be standalone
    session_type VARCHAR(20) DEFAULT 'webinar', -- webinar, office_hours, workshop, q_and_a
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    duration_minutes INTEGER NOT NULL,
    meeting_url VARCHAR(500),
    max_attendees INTEGER DEFAULT 100,
    current_attendees INTEGER DEFAULT 0,
    price DECIMAL(10,2) DEFAULT 0.00,
    status VARCHAR(20) DEFAULT 'scheduled', -- scheduled, live, completed, cancelled
    recording_url VARCHAR(500),
    chat_enabled BOOLEAN DEFAULT true,
    q_and_a_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Live session registrations
CREATE TABLE IF NOT EXISTS live_session_registrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    session_id UUID NOT NULL REFERENCES live_sessions(id) ON DELETE CASCADE,
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    attended BOOLEAN DEFAULT false,
    attendance_duration_minutes INTEGER DEFAULT 0,
    payment_status VARCHAR(20) DEFAULT 'free', -- free, paid, pending
    payment_amount DECIMAL(10,2) DEFAULT 0.00,
    UNIQUE(user_id, session_id)
);

-- Learning paths for structured education
CREATE TABLE IF NOT EXISTS learning_paths (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,
    difficulty_level INTEGER DEFAULT 1, -- 1-5 scale
    estimated_duration_hours INTEGER DEFAULT 0,
    course_sequence UUID[] NOT NULL, -- ordered array of course IDs
    prerequisites TEXT[],
    price DECIMAL(10,2) DEFAULT 0.00,
    enrollment_count INTEGER DEFAULT 0,
    completion_count INTEGER DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0.00,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Learning path enrollments
CREATE TABLE IF NOT EXISTS learning_path_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    path_id UUID NOT NULL REFERENCES learning_paths(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'enrolled', -- enrolled, in_progress, completed, dropped
    progress DECIMAL(5,2) DEFAULT 0.00, -- 0-100%
    current_course_index INTEGER DEFAULT 0,
    enrolled_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(user_id, path_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_courses_category ON courses(category);
CREATE INDEX IF NOT EXISTS idx_courses_level ON courses(level);
CREATE INDEX IF NOT EXISTS idx_courses_status ON courses(status);
CREATE INDEX IF NOT EXISTS idx_courses_instructor ON courses(instructor_id);
CREATE INDEX IF NOT EXISTS idx_courses_rating ON courses(rating DESC);
CREATE INDEX IF NOT EXISTS idx_courses_enrollment_count ON courses(enrollment_count DESC);

CREATE INDEX IF NOT EXISTS idx_enrollments_user ON course_enrollments(user_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_course ON course_enrollments(course_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_status ON course_enrollments(status);

CREATE INDEX IF NOT EXISTS idx_lesson_progress_user ON lesson_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_lesson_progress_lesson ON lesson_progress(lesson_id);

CREATE INDEX IF NOT EXISTS idx_quiz_attempts_user ON quiz_attempts(user_id);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_quiz ON quiz_attempts(quiz_id);

CREATE INDEX IF NOT EXISTS idx_content_library_type ON content_library(content_type);
CREATE INDEX IF NOT EXISTS idx_content_library_category ON content_library(category);
CREATE INDEX IF NOT EXISTS idx_content_library_published ON content_library(published_at);

-- Functions for automatic calculations
CREATE OR REPLACE FUNCTION update_course_stats()
RETURNS TRIGGER AS $$
BEGIN
    -- Update course enrollment count
    UPDATE courses 
    SET enrollment_count = (
        SELECT COUNT(*) FROM course_enrollments 
        WHERE course_id = NEW.course_id AND status != 'refunded'
    ),
    completion_count = (
        SELECT COUNT(*) FROM course_enrollments 
        WHERE course_id = NEW.course_id AND status = 'completed'
    ),
    updated_at = NOW()
    WHERE id = NEW.course_id;
    
    -- Update instructor stats
    UPDATE instructors 
    SET total_students = (
        SELECT COUNT(DISTINCT user_id) FROM course_enrollments ce
        JOIN courses c ON ce.course_id = c.id
        WHERE c.instructor_id = (SELECT instructor_id FROM courses WHERE id = NEW.course_id)
        AND ce.status != 'refunded'
    ),
    updated_at = NOW()
    WHERE id = (SELECT instructor_id FROM courses WHERE id = NEW.course_id);
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update course stats
CREATE TRIGGER trigger_update_course_stats
    AFTER INSERT OR UPDATE ON course_enrollments
    FOR EACH ROW
    EXECUTE FUNCTION update_course_stats();

-- Function to update course ratings
CREATE OR REPLACE FUNCTION update_course_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE courses 
    SET rating = (
        SELECT COALESCE(AVG(rating), 0) FROM course_reviews 
        WHERE course_id = NEW.course_id
    ),
    review_count = (
        SELECT COUNT(*) FROM course_reviews 
        WHERE course_id = NEW.course_id
    ),
    updated_at = NOW()
    WHERE id = NEW.course_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update course ratings
CREATE TRIGGER trigger_update_course_rating
    AFTER INSERT OR UPDATE OR DELETE ON course_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_course_rating();

-- Insert sample course categories and levels
INSERT INTO courses (id, title, description, category, level, price, instructor_id, status) VALUES
(gen_random_uuid(), 'Cryptocurrency Fundamentals', 'Learn the basics of cryptocurrency and blockchain technology', 'beginner', 'crypto_basics', 299.00, (SELECT id FROM instructors LIMIT 1), 'published'),
(gen_random_uuid(), 'AI-Powered Trading Strategies', 'Master advanced trading using artificial intelligence', 'advanced', 'ai_trading', 1999.00, (SELECT id FROM instructors LIMIT 1), 'published'),
(gen_random_uuid(), 'DeFi Mastery Course', 'Complete guide to decentralized finance protocols', 'intermediate', 'advanced_strategies', 799.00, (SELECT id FROM instructors LIMIT 1), 'published')
ON CONFLICT DO NOTHING;

-- Create views for analytics
CREATE OR REPLACE VIEW course_analytics AS
SELECT 
    c.id,
    c.title,
    c.category,
    c.level,
    c.price,
    c.enrollment_count,
    c.completion_count,
    CASE WHEN c.enrollment_count > 0 THEN 
        c.completion_count::decimal / c.enrollment_count::decimal * 100 
    ELSE 0 END as completion_rate,
    c.rating,
    c.review_count,
    c.total_revenue,
    i.name as instructor_name
FROM courses c
JOIN instructors i ON c.instructor_id = i.id
WHERE c.status = 'published';

-- Comments for documentation
COMMENT ON TABLE courses IS 'Educational courses and training programs';
COMMENT ON TABLE course_enrollments IS 'User enrollments in courses';
COMMENT ON TABLE lesson_progress IS 'Individual lesson completion tracking';
COMMENT ON TABLE course_certificates IS 'Course completion certificates';
COMMENT ON TABLE content_library IS 'Educational content repository';

COMMENT ON COLUMN courses.price IS 'Course price in specified currency';
COMMENT ON COLUMN course_enrollments.progress IS 'Course completion percentage 0-100%';
COMMENT ON COLUMN lesson_progress.progress IS 'Lesson completion percentage 0-100%';
