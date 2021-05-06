
-- +migrate Up
CREATE TABLE questions_answers(
    question    VARCHAR CONSTRAINT questions_answers_pk PRIMARY KEY NOT NULL,
    answer      VARCHAR NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    deleted_at  TIMESTAMP WITH TIME ZONE
);

COMMENT ON TABLE  questions_answers           IS 'Захардкоженые вход-выход данные';
COMMENT ON COLUMN questions_answers.question  IS 'входные данные';
COMMENT ON COLUMN questions_answers.answer    IS 'выходные данные';

INSERT INTO questions_answers (question, answer) VALUES ('ping', 'pong');

-- +migrate Down
DROP TABLE questions_answers;
