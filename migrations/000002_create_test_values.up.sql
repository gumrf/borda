-- Create test Users
INSERT INTO public."user" (name, password, contact)
VALUES
       -- Password = User[i]Pass. Example: User1Pass
       ('User1', '70617373776f72645f73616c745f78363973069283af36faa8d331eb3ed5888e5ace9c4d9b6d', '@user1'),
       ('User2', '70617373776f72645f73616c745f7836397331cbb66d17819e36ef6cf1b1480526c68d9e22c5', '@user2'),
       ('User3', '70617373776f72645f73616c745f783639731b2cbeaaafc2f7dac2dd273b3c584da6e88122cf', '@user3');

-- Assign roles to users. TestUser1 has admin role, others have user role
INSERT INTO public.user_role (user_id, role_id) VALUES (1,1), (2,2), (3,2);

-- Create test Teams
INSERT INTO public.team (name, token, team_leader_id)
VALUES
    ('Team1', '7c8b4c73-9fdd-4fbb-b926-43de9aa6f24d', 1),
    ('Team2', 'e01e949e-9dd0-428e-96c0-28adebf4df3d', 2),
    ('Team3', 'bd58e756-7ef3-4043-bf4c-2c5ae9b9ad0b', 3);

-- Add users to team
INSERT INTO public.team_member(team_id, user_id) VALUES (1,1), (2,2), (3,3);

-- Create test Authors
INSERT INTO public.author (name, contact)
VALUES
       ('Author1', '@author1'),
       ('Author2', '@author2'),
       ('Author3', '@author3');

-- Create test Tasks
INSERT INTO public.task (title, description, category, complexity, points, hint, flag, is_active, is_disabled, author_id)
VALUES
       ('Task1', 'Task1 description', 'test', 'hard', 1337, 'Hint for task 1', 'flag{flag_for_task_1}', true, false, 1),
       ('Task2', 'Task2 description', 'test', 'easy', 1337, 'Hint for task 2', 'flag{flag_for_task_2}', false, false, 2),
       ('Task3', 'Task3 description', 'test', 'easy', 1337, 'Hint for task 3', 'flag{flag_for_task_3}', false, true, 2);

-- Create test Settings
INSERT INTO  public.settings (key, value) VALUES ('team_limit', '4')