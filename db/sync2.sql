-- Surveys management

CREATE TABLE IF NOT EXISTS surveys (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE question_type AS ENUM ('TEXT');

CREATE TABLE IF NOT EXISTS survey_questions (
    id SERIAL PRIMARY KEY,
    survey_id INT NOT NULL,
    question_text VARCHAR(255) NOT NULL,
    question_type question_type NOT NULL DEFAULT 'TEXT',
    FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS survey_respondents (
    id SERIAL PRIMARY KEY,
    survey_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(25) NOT NULL,
    rating INT NOT NULL,
    FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS question_responses (
    id SERIAL PRIMARY KEY,
    respondent_id INT NOT NULL,
    question_id INT NOT NULL,
    response_text VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (respondent_id) REFERENCES survey_respondents(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES survey_questions(id) ON DELETE CASCADE
);

-- Landing page management

CREATE TABLE IF NOT EXISTS landing (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    survey_id INT,
    FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS landing_images (
    landing_id INT NOT NULL,
    image_id INT NOT NULL,
    index INT NOT NULL DEFAULT 0,
    PRIMARY KEY (landing_id, image_id),
    FOREIGN KEY (landing_id) REFERENCES landing(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE
);
