-- Create test Users
INSERT INTO public."user" (name, password, contact)
VALUES
       ('TestUser1', 'change me', '@TestUser1'),
       ('TestUser2', 'change me', '@TestUser2'),
       ('TestUser3', 'change me', '@TestUser3');

-- Create test Authors
INSERT INTO public.author (name, contact)
VALUES
       ('TestAuthor1', '@TestAuthor1_contact'),
       ('TestAuthor3', '@TestAuthor2_contact'),
       ('TestAuthor2', '@TestAuthor3_contact');

-- Create test Tasks
INSERT INTO public.task (title, description, category, complexity, points, hint, flag, is_active, is_disabled, author_id)
VALUES
       ('TestTask1', 'TestTask1_description', 'test', 'hard', 1337, 'Test Hint for task 1', 'flag{test_flag_for_task_1}', true, false, 1),
       ('TestTask1', 'TestTask1_description', 'test', 'easy', 1337, 'Test Hint for task 2', 'flag{test_flag_for_task_2}', false, false, 2),
       ('TestTask1', 'TestTask1_description', 'test', 'easy', 1337, 'Test Hint for task 3', 'flag{test_flag_for_task_3}', false, true, 2);

-- Create test Teams
INSERT INTO public.team (name, token, team_leader_id)
VALUES
       ('TestTeam1', '2b8dc03a-aa81-11ec-8928-04d9f50035c4', 1),
       ('TestTeam2', '2cbfd65a-aa81-11ec-b0f6-04d9f50035c4', 2),
       ('TestTeam3', '3858b180-aa81-11ec-96e5-04d9f50035c4', 3);

-- Create test Settings
INSERT INTO  public.settings (key, value) VALUES ('team-limit', '4')








