CREATE TABLE projects(
    id serial NOT NULL UNIQUE,
    title varchar(255) NOT NULL,
    description varchar (255),
    done boolean NOT NULL DEFAULT false
);