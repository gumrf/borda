-- Create test Users
INSERT INTO public."user" (name, password, contact)
VALUES
       -- TestUser1Password
       ('TestUser1', '70617373776f72645f73616c745f783639735df7b147ad23e5ec0eb51eb6ef272e6a10f33d49', '@TestUser1'),
       -- TestUser2Password
       ('TestUser2', '70617373776f72645f73616c745f78363973f0e0b8ace07ed0c0c2c111f8a04a7b5bd99062de', '@TestUser2'),
       -- TestUser3Password
       ('TestUser3', '70617373776f72645f73616c745f78363973b56b6aa265cd6d4362fe387bc144deaf2991bf07', '@TestUser3');

-- Create test Authors
INSERT INTO public.author (name, contact)
VALUES
       ('TestAuthor1', '@TestAuthor1_contact'),
       ('TestAuthor2', '@TestAuthor2_contact'),
       ('TestAuthor3', '@TestAuthor3_contact');

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
INSERT INTO  public.settings (key, value) VALUES ('team_limit', '4')