alter table `project__integration_extract_value` add `variable_type` varchar(255) not null;
update table `project__integration_extract_value` set `variable_type` = 'environment';