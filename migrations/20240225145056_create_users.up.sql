CREATE TABLE users(
    user_id INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(60) NOT NULL UNIQUE,
    encrypted_password VARCHAR(60) NOT NULL
);


create table lk(
    lk_id int PRIMARY KEY auto_increment,
    email VARCHAR(60) NOT NULL UNIQUE,
    user_id int,
    nickname VARCHAR(60) NOT NULL UNIQUE,
    cards_count int DEFAULT 0,
    user_description VARCHAR(60) DEFAULT "",
    encrypted_password VARCHAR(60) NOT NULL,
    Foreign Key (user_id) REFERENCES users(user_id) on delete cascade
);

create table card_groups(
	group_id int primary key auto_increment,
    user_id int,
    group_name varchar(60),
    cards_count int default 0,
    foreign key (user_id) references users(user_id) on delete cascade
);
create table cards (
	card_id int primary key auto_increment,
    user_id int,
	front_side varchar(300),
    back_side varchar(300),
    card_time timestamp,
    time_flag time,
    group_id int,
    foreign key (user_id) references users(user_id) on delete cascade,
    foreign key (group_id) references card_groups(group_id) on delete cascade
);
delimiter //
CREATE TRIGGER increment_cards_count_lk
AFTER INSERT ON cards
FOR EACH ROW
BEGIN
    UPDATE lk
    SET cards_count = cards_count + 1
    WHERE lk.user_id = NEW.user_id;
END; //
CREATE TRIGGER increment_cards_count_group
AFTER INSERT ON cards
FOR EACH ROW
BEGIN
	UPDATE card_groups
    SET cards_count = cards_count + 1
    WHERE card_groups.user_id = NEW.user_id AND card_groups.group_id = NEW.group_id;
END; //
delimiter ;