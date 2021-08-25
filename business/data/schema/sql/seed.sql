INSERT INTO ET_USERS(USER_ID ,FIRST_NAME, LAST_NAME,EMAIL , PASSWORD) VALUES
('5cf37266-3473-4006-984f-9325122678b7','mark', 'sarchi', 'sarchimark@example.com', '$2a$12$JOvqT7bBQlEZ7Cnkiur8teSRU8xwpzH4L9475X3zMpin/u7lBQueK'),
('45b5fbd3-755f-4379-8f07-a58d4a30fa2f','eva','max' ,'evamax@yahoo.com', '$2a$12$lcaFYoHrKTAFOtaaa5DgKO0b9GEULsAG23Z.q6cy91nFCujy91Z1y')
ON CONFLICT DO NOTHING;

INSERT INTO ET_CATEGORIES ( CATEGORY_ID, USER_ID, TITLE, DESCRIPTION) VALUES
( 1,'5cf37266-3473-4006-984f-9325122678b7' ,'Travel costs', 'All travel costs'),
( 2,'45b5fbd3-755f-4379-8f07-a58d4a30fa2f' ,'Car Costs', 'all car costs')
ON CONFLICT DO NOTHING;  

INSERT INTO ET_TRANSACTIONS (TRANSACTION_ID, CATEGORY_ID, USER_ID, AMOUNT, NOTE, TRANSACTION_DATE) VALUES
( 1001, 1, '5cf37266-3473-4006-984f-9325122678b7',2580, 'Travel upcountry', 1616014260159428300),
( 1002, 2 , '45b5fbd3-755f-4379-8f07-a58d4a30fa2f',33450, 'Repaired windscreen', 1616097061795640800),
( 1003, 2 ,'45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 33450, 'Bought new tyre', 1616097061795640800)
ON CONFLICT DO NOTHING;





