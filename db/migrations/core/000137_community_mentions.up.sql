alter table mentions drop column contract_id;
alter table mentions add column community_id varchar(255) references communities(id);