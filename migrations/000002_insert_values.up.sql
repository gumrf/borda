--Author
INSERT INTO author (name, contact) VALUES ('Test_Author_02', '@TestContact2');
--Task
INSERT INTO task (title, description, category, complexity, points, hint, flag, is_active, is_disabled, author_id) VALUES ('TestTask_02', 'Test_description', 'Test', 'Hard', 133, 'Test_hint', 'testFLAG', true, true, 1);
--User 8cb2237d0679ca88db6464eac60da96345513964 hash of 12345
INSERT INTO "user" (name, password, contact) VALUES ('TestUser02', '8cb2237d0679ca88db6464eac60da96345513964', '@UserContact1');
--Team
INSERT INTO team (name, token, team_leader_id) VALUES ('TestTeam02', '8888', 1);
