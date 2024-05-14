CREATE TABLE users(
  id bigserial not null,
  name varchar(255) not null,
  age int not null,
  location varchar(255) not null,
  gender varchar(255) not null,
  description text not null,
  activities text[],
  primary key (id)
);

CREATE TABLE wishes(
  user_id bigserial unique not null,
  age int not null,
  gender varchar(255) not null,
  location varchar(255) not null,
  activities text[]
); 
