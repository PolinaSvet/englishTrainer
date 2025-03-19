/*
	script to create database structure
*/

DROP TABLE IF EXISTS users;


CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    mail TEXT UNIQUE NOT NULL,
    telega TEXT,
    password TEXT NOT NULL
);
COMMENT ON TABLE users IS 'Table for users studied';


DROP FUNCTION IF EXISTS f_users_insert, f_users_update, f_users_delete, f_users_view;
--=======================
--1::: table: letter
--=======================
--insert
CREATE FUNCTION f_users_insert(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	INSERT INTO users (
		name,
		mail,
		telega,
		password
	)    
	VALUES (
		(json_data ->> 'name')::TEXT,
		(json_data ->> 'mail')::TEXT,
		(json_data ->> 'telega')::TEXT,
		(json_data ->> 'password')::TEXT
	) 
	RETURNING id INTO new_id;
	
	IF new_id IS NULL THEN
		RAISE EXCEPTION 'Parameter value cannot be null. ';
	END IF;

	SELECT json_build_object('id',new_id,'err','') INTO json_result;
  	RETURN json_result;

EXCEPTION
    WHEN others THEN
		GET STACKED DIAGNOSTICS err_context = PG_EXCEPTION_CONTEXT;
    	GET STACKED DIAGNOSTICS err_mess = MESSAGE_TEXT;

        SELECT json_build_object('id',null,'err',err_mess||err_context) INTO json_result;
  		RETURN json_result;  
END;
$$ LANGUAGE plpgsql;

--update
CREATE FUNCTION f_users_update(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	UPDATE users SET
		name = (json_data ->> 'name')::TEXT,
		mail = (json_data ->> 'mail')::TEXT,
		telega = (json_data ->> 'telega')::TEXT,
		password = (json_data ->> 'password')::TEXT
	WHERE id = (json_data ->> 'id')::BIGINT 
	RETURNING id INTO new_id;
	
	IF new_id IS NULL THEN
		RAISE EXCEPTION 'Parameter value cannot be null. ';
	END IF;

	SELECT json_build_object('id',new_id,'err','') INTO json_result;
  	RETURN json_result;

EXCEPTION
    WHEN others THEN
		GET STACKED DIAGNOSTICS err_context = PG_EXCEPTION_CONTEXT;
    	GET STACKED DIAGNOSTICS err_mess = MESSAGE_TEXT;

        SELECT json_build_object('id',null,'err',err_mess||err_context) INTO json_result;
  		RETURN json_result;  
END;
$$ LANGUAGE plpgsql;


--delete
CREATE FUNCTION f_users_delete(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	DELETE FROM users WHERE id = (json_data ->> 'id')::BIGINT; 

	SELECT json_build_object('id',(json_data ->> 'id')::BIGINT,'err','') INTO json_result;
  	RETURN json_result;

EXCEPTION
    WHEN others THEN
		GET STACKED DIAGNOSTICS err_context = PG_EXCEPTION_CONTEXT;
    	GET STACKED DIAGNOSTICS err_mess = MESSAGE_TEXT;

        SELECT json_build_object('id',null,'err',err_mess||err_context) INTO json_result;
  		RETURN json_result;  
END;
$$ LANGUAGE plpgsql;

--select view
CREATE FUNCTION f_users_view(
		json_data jsonb
) 
RETURNS TABLE (
	id BIGINT,
    name TEXT,
	mail TEXT,
	telega TEXT,
	password TEXT
) AS $$
DECLARE
  	par_id BIGINT = 0;
	par_mail TEXT = '';

BEGIN

	IF (json_data ->> 'id') IS NOT NULL THEN
		par_id = (json_data ->> 'id')::BIGINT;
	END IF;

	IF (json_data ->> 'mail') IS NOT NULL THEN
		par_mail = (json_data ->> 'mail')::TEXT;
	END IF;

	RETURN QUERY
		SELECT users.id,
			   users.name,  
			    users.mail,
			    users.telega,
				users.password
		FROM users
		WHERE
			(par_id = 0 OR users.id = par_id) and
			(par_mail = '' OR users.mail = par_mail)
		ORDER BY users.id;
	
END;
$$ LANGUAGE plpgsql;


--users: insert
DO $$
DECLARE
BEGIN

	PERFORM * FROM f_users_insert(('{
				    "name": "user1",
			        "mail": "mail1",
			        "telega": "telega1",
			        "password": "1"
				}')::jsonb);

	PERFORM * FROM f_users_insert(('{
					    "name": "user2",
				        "mail": "mail2",
				        "telega": "telega2",
				        "password": "2"
					}')::jsonb);



END;
$$;

SELECT * FROM f_users_view('{
    "id": 0
     
}'::jsonb);
