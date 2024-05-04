CREATE TABLE users(
    user_id INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(60) NOT NULL UNIQUE,
    encrypted_password VARCHAR(60) NOT NULL
);
create table cards (
	card_id int primary key auto_increment,
    user_id int,
	front_side varchar(300),
    back_side varchar(300),
    card_time timestamp,
    time_flag time,
    foreign key (user_id) references users(user_id) on delete cascade
);

create table lk(
    lk_id int PRIMARY KEY auto_increment,
    email VARCHAR(60) NOT NULL UNIQUE,
    user_id int,
    nickname VARCHAR(60) NOT NULL UNIQUE,
    cards_count int,
    user_description VARCHAR(60),
    encrypted_password VARCHAR(60) NOT NULL,
    Foreign Key (user_id) REFERENCES users(user_id) on delete cascade
);