-- Version: 1.1
-- Description: Create table et_users
create table et_users(
    user_id UUID ,
    first_name varchar(20) not null,
    last_name varchar(20) not null,
    email varchar(30) not null,
    password text not null,

    PRIMARY KEY(user_id)
);
-- Version: 1.2
-- Description: Create table et_categories
create table et_categories(
    category_id integer not null,
    user_id UUID ,
    title varchar(20) not null,
    description varchar(50) not null,

    PRIMARY KEY (category_id)
);

-- Version: 1.3
-- Description: Alter table et_categories
alter table et_categories add constraint cat_users_fk
foreign key (user_id) references et_users(user_id);

-- Version: 2.1
-- Description: Create table transactions
create table et_transactions(
    transaction_id integer not  null ,
    category_id integer not null,
    user_id UUID not null,
    amount numeric(10,2),
    note varchar(50) not null,
    transaction_date bigint not null,

    PRIMARY KEY(transaction_id)
);

-- Version: 2.2
-- Description:alter table et_transactions
alter table et_transactions add constraint trans_cat_fk
foreign  key (category_id) references  et_categories(category_id);
alter table et_transactions add constraint trans_users_fk
foreign key (user_id) references et_users(user_id);
create sequence et_categories_seq increment 1 start 1;
create sequence et_transactions_seq increment 1 start 1000;





