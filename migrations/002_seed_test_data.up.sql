-- Create test Users
INSERT INTO public."user" (name, password, contact)
VALUES
       -- Password for test users: test_password
       ('TestUser1', '707377645f7365637265745f7068726173659fb7fe1217aed442b04c0f5e43b5d5a7d3287097', '@testuser1'),
       ('TestUser2', '707377645f7365637265745f7068726173659fb7fe1217aed442b04c0f5e43b5d5a7d3287097', '@testuser2'),
       ('TestUser3', '707377645f7365637265745f7068726173659fb7fe1217aed442b04c0f5e43b5d5a7d3287097', '@testuser3'),
       ('TestUser4', '707377645f7365637265745f7068726173659fb7fe1217aed442b04c0f5e43b5d5a7d3287097', '@testuser4'),
       ('TestUser5', '707377645f7365637265745f7068726173659fb7fe1217aed442b04c0f5e43b5d5a7d3287097', '@testuser5'),
       ('TestUser6', '707377645f7365637265745f7068726173659fb7fe1217aed442b04c0f5e43b5d5a7d3287097', '@testuser6'),
       ('TestUser7', '707377645f7365637265745f7068726173659fb7fe1217aed442b04c0f5e43b5d5a7d3287097', '@testuser7'),
       ('TestUser8', '707377645f7365637265745f7068726173659fb7fe1217aed442b04c0f5e43b5d5a7d3287097', '@testuser8');

-- Assign roles to users. TestUser1 has admin role, others have user role
INSERT INTO public.user_role (user_id, role_id) 
VALUES 
       (1,1), 
       (2,2), (3,2), (4,2), (5,2), (6,2), (7,2), (8,2);

-- Create test Teams
INSERT INTO public.team (name, token)
VALUES
    ('TestTeam1', '7c8b4c73-9fdd-4fbb-b926-43de9aa6f24d'),
    ('TestTeam2', 'e01e949e-9dd0-428e-96c0-28adebf4df3d'),
    ('TestTeam3', 'bd58e756-7ef3-4043-bf4c-2c5ae9b9ad0b');

-- Add users to team
INSERT INTO public.team_member(team_id, user_id, is_captain) 
VALUES 
       (1,1,true), (1,2,false), (1,3,false), (1,4,false),
       (2,5,true), (2,6,false),
       (3,7,true);

-- Create test Authors
INSERT INTO public.author (name, contact)
VALUES
       ('TestAuthor1', '@author1'),
       ('TestAuthor2', '@author2'),
       ('TestAuthor3', '@author3');

-- Create test Tasks
INSERT INTO public.task (title, description, category, complexity, points, hint, flag, is_active, is_disabled, author_id)
VALUES
       ('TestTask1', 'TestTask1 description', 'test', 'hard', 1337, 'Hint for test task 1', 'flag_for_task_1', true, false, 1),
       ('TestTask2', 'TestTask2 description', 'test', 'easy', 1337, 'Hint for test task 2', 'flag_for_task_2', false, false, 2),
       ('TestTask3', 'TestTask3 description', 'test', 'easy', 1337, 'Hint for test task 3', 'flag_for_task_3', false, true, 2);

-- Create test Settings
INSERT INTO  public.settings (key, value) VALUES ('team_limit', '4')