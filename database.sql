CREATE TABLE blc_admin_module (
	admin_module_id int8 NOT NULL,
	display_order int4 NULL,
	icon varchar(255) NULL,
	module_key varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
	CONSTRAINT blc_admin_module_pkey PRIMARY KEY (admin_module_id)
);
CREATE INDEX adminmodule_name_index ON public.blc_admin_module USING btree (name);
CREATE TABLE blc_admin_password_token (
	password_token varchar(255) NOT NULL,
	admin_user_id int8 NOT NULL,
	create_date timestamp NOT NULL,
	token_used_date timestamp NULL,
	token_used_flag bool NOT NULL,
	CONSTRAINT blc_admin_password_token_pkey PRIMARY KEY (password_token)
);
CREATE TABLE blc_admin_permission (
	admin_permission_id int8 NOT NULL,
	description varchar(255) NOT NULL,
	is_friendly bool NULL,
	"name" varchar(255) NOT NULL,
	permission_type varchar(255) NOT NULL,
	CONSTRAINT blc_admin_permission_pkey PRIMARY KEY (admin_permission_id)
);
CREATE TABLE blc_admin_role (
	admin_role_id int8 NOT NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	description varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
	CONSTRAINT blc_admin_role_pkey PRIMARY KEY (admin_role_id)
);
CREATE TABLE blc_admin_user (
	admin_user_id int8 NOT NULL,
	active_status_flag bool NULL,
	archived bpchar(1) NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	email varchar(255) NOT NULL,
	login varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
	"password" varchar(255) NULL,
	phone_number varchar(255) NULL,
	CONSTRAINT blc_admin_user_pkey PRIMARY KEY (admin_user_id)
);
CREATE INDEX adminperm_email_index ON public.blc_admin_user USING btree (email);
CREATE INDEX adminuser_name_index ON public.blc_admin_user USING btree (name);
CREATE TABLE blc_bank_account_payment (
	payment_id int8 NOT NULL,
	account_number varchar(255) NOT NULL,
	reference_number varchar(255) NOT NULL,
	routing_number varchar(255) NOT NULL,
	CONSTRAINT blc_bank_account_payment_pkey PRIMARY KEY (payment_id)
);
CREATE INDEX bankaccount_index ON public.blc_bank_account_payment USING btree (reference_number);
CREATE TABLE blc_catalog (
	catalog_id int8 NOT NULL,
	archived bpchar(1) NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	"name" varchar(255) NULL,
	CONSTRAINT blc_catalog_pkey PRIMARY KEY (catalog_id)
);
CREATE TABLE blc_challenge_question (
	question_id int8 NOT NULL,
	question varchar(255) NOT NULL,
	CONSTRAINT blc_challenge_question_pkey PRIMARY KEY (question_id)
);
CREATE TABLE blc_cms_menu (
	menu_id int8 NOT NULL,
	"name" varchar(255) NOT NULL,
	CONSTRAINT blc_cms_menu_pkey PRIMARY KEY (menu_id)
);
CREATE INDEX idx_menu_name ON public.blc_cms_menu USING btree (name);
CREATE TABLE blc_code_types (
	code_id int8 NOT NULL,
	code_type varchar(255) NOT NULL,
	code_desc varchar(255) NULL,
	code_key varchar(255) NOT NULL,
	modifiable bpchar(1) NULL,
	CONSTRAINT blc_code_types_pkey PRIMARY KEY (code_id)
);
CREATE TABLE blc_country (
	abbreviation varchar(255) NOT NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	"name" varchar(255) NOT NULL,
	CONSTRAINT blc_country_pkey PRIMARY KEY (abbreviation)
);
CREATE TABLE blc_country_sub_cat (
	country_sub_cat_id int8 NOT NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	"name" varchar(255) NOT NULL,
	CONSTRAINT blc_country_sub_cat_pkey PRIMARY KEY (country_sub_cat_id)
);
CREATE TABLE blc_credit_card_payment (
	payment_id int8 NOT NULL,
	expiration_month int4 NOT NULL,
	expiration_year int4 NOT NULL,
	name_on_card varchar(255) NOT NULL,
	pan varchar(255) NOT NULL,
	reference_number varchar(255) NOT NULL,
	CONSTRAINT blc_credit_card_payment_pkey PRIMARY KEY (payment_id)
);
CREATE INDEX creditcard_index ON public.blc_credit_card_payment USING btree (reference_number);
CREATE TABLE blc_currency (
	currency_code varchar(255) NOT NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	default_flag bool NULL,
	friendly_name varchar(255) NULL,
	CONSTRAINT blc_currency_pkey PRIMARY KEY (currency_code)
);
CREATE TABLE blc_customer_password_token (
	password_token varchar(255) NOT NULL,
	create_date timestamp NOT NULL,
	customer_id int8 NOT NULL,
	token_used_date timestamp NULL,
	token_used_flag bool NOT NULL,
	CONSTRAINT blc_customer_password_token_pkey PRIMARY KEY (password_token)
);
CREATE TABLE blc_data_drvn_enum (
	enum_id int8 NOT NULL,
	enum_key varchar(255) NULL,
	modifiable bool NULL,
	CONSTRAINT blc_data_drvn_enum_pkey PRIMARY KEY (enum_id)
);
CREATE INDEX enum_key_index ON public.blc_data_drvn_enum USING btree (enum_key);
CREATE TABLE blc_email_tracking (
	email_tracking_id int8 NOT NULL,
	date_sent timestamp NULL,
	email_address varchar(255) NULL,
	"type" varchar(255) NULL,
	CONSTRAINT blc_email_tracking_pkey PRIMARY KEY (email_tracking_id)
);
CREATE INDEX datesent_index ON public.blc_email_tracking USING btree (date_sent);
CREATE INDEX emailtracking_index ON public.blc_email_tracking USING btree (email_address);
CREATE TABLE blc_field (
	field_id int8 NOT NULL,
	abbreviation varchar(255) NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	entity_type varchar(255) NOT NULL,
	friendly_name varchar(255) NULL,
	override_generated_prop_name bool NULL,
	property_name varchar(255) NOT NULL,
	translatable bool NULL,
	CONSTRAINT blc_field_pkey PRIMARY KEY (field_id)
);
CREATE INDEX entity_type_index ON public.blc_field USING btree (entity_type);
CREATE TABLE blc_fld_group (
	fld_group_id int8 NOT NULL,
	init_collapsed_flag bool NULL,
	is_master_field_group bool NULL,
	"name" varchar(255) NULL,
	CONSTRAINT blc_fld_group_pkey PRIMARY KEY (fld_group_id)
);
CREATE TABLE blc_fulfillment_option (
	fulfillment_option_id int8 NOT NULL,
	fulfillment_type varchar(255) NOT NULL,
	long_description text NULL,
	"name" varchar(255) NULL,
	tax_code varchar(255) NULL,
	taxable bool NULL,
	use_flat_rates bool NULL,
	CONSTRAINT blc_fulfillment_option_pkey PRIMARY KEY (fulfillment_option_id)
);
CREATE TABLE blc_gift_card_payment (
	payment_id int8 NOT NULL,
	pan varchar(255) NOT NULL,
	pin varchar(255) NULL,
	reference_number varchar(255) NOT NULL,
	CONSTRAINT blc_gift_card_payment_pkey PRIMARY KEY (payment_id)
);
CREATE INDEX giftcard_index ON public.blc_gift_card_payment USING btree (reference_number);
CREATE TABLE blc_id_generation (
	id_type varchar(255) NOT NULL,
	batch_size int8 NOT NULL,
	batch_start int8 NOT NULL,
	id_min int8 NULL,
	id_max int8 NULL,
	"version" int4 NULL,
	CONSTRAINT blc_id_generation_pkey PRIMARY KEY (id_type)
);
CREATE TABLE blc_iso_country (
	alpha_2 varchar(255) NOT NULL,
	alpha_3 varchar(255) NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	"name" varchar(255) NULL,
	numeric_code int4 NULL,
	status varchar(255) NULL,
	CONSTRAINT blc_iso_country_pkey PRIMARY KEY (alpha_2)
);
CREATE TABLE blc_media (
	media_id int8 NOT NULL,
	alt_text varchar(255) NULL,
	tags varchar(255) NULL,
	title varchar(255) NULL,
	url varchar(255) NOT NULL,
	CONSTRAINT blc_media_pkey PRIMARY KEY (media_id)
);
CREATE INDEX media_title_index ON public.blc_media USING btree (title);
CREATE INDEX media_url_index ON public.blc_media USING btree (url);
CREATE TABLE blc_module_configuration (
	module_config_id int8 NOT NULL,
	active_end_date timestamp NULL,
	active_start_date timestamp NULL,
	archived bpchar(1) NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	config_type varchar(255) NOT NULL,
	is_default bool NOT NULL,
	module_name varchar(255) NOT NULL,
	module_priority int4 NOT NULL,
	CONSTRAINT blc_module_configuration_pkey PRIMARY KEY (module_config_id)
);
CREATE TABLE blc_offer (
	offer_id int8 NOT NULL,
	offer_adjustment_type varchar(255) NULL,
	apply_to_child_items bool NULL,
	apply_to_sale_price bool NULL,
	archived bpchar(1) NULL,
	automatically_added bool NULL,
	combinable_with_other_offers bool NULL,
	offer_description varchar(255) NULL,
	offer_discount_type varchar(255) NULL,
	end_date timestamp NULL,
	marketing_message varchar(255) NULL,
	max_uses_per_customer int8 NULL,
	max_uses int4 NULL,
	max_uses_strategy varchar(255) NULL,
	minimum_days_per_usage int8 NULL,
	offer_name varchar(255) NOT NULL,
	offer_item_qualifier_rule varchar(255) NULL,
	offer_item_target_rule varchar(255) NULL,
	order_min_total numeric(19, 5) NULL,
	offer_priority int4 NULL,
	qualifying_item_min_total numeric(19, 5) NULL,
	requires_related_tar_qual bool NULL,
	start_date timestamp NULL,
	target_min_total numeric(19, 5) NULL,
	target_system varchar(255) NULL,
	totalitarian_offer bool NULL,
	offer_type varchar(255) NOT NULL,
	use_list_for_discounts bool NULL,
	offer_value numeric(19, 5) NOT NULL,
	CONSTRAINT blc_offer_pkey PRIMARY KEY (offer_id)
);
CREATE INDEX idx_blof_start_date ON public.blc_offer USING btree (start_date);
CREATE INDEX offer_automatically_added_index ON public.blc_offer USING btree (automatically_added);
CREATE INDEX offer_discount_index ON public.blc_offer USING btree (offer_discount_type);
CREATE INDEX offer_marketing_message_index ON public.blc_offer USING btree (marketing_message);
CREATE INDEX offer_name_index ON public.blc_offer USING btree (offer_name);
CREATE INDEX offer_type_index ON public.blc_offer USING btree (offer_type);
CREATE TABLE blc_offer_audit (
	offer_audit_id int8 NOT NULL,
	account_id int8 NULL,
	customer_id int8 NULL,
	offer_code_id int8 NULL,
	offer_id int8 NULL,
	order_id int8 NULL,
	redeemed_date timestamp NULL,
	CONSTRAINT blc_offer_audit_pkey PRIMARY KEY (offer_audit_id)
);
CREATE INDEX offeraudit_account_index ON public.blc_offer_audit USING btree (account_id, offer_id);
CREATE INDEX offeraudit_customer_index ON public.blc_offer_audit USING btree (customer_id, offer_id);
CREATE INDEX offeraudit_offer_code_index ON public.blc_offer_audit USING btree (offer_code_id);
CREATE INDEX offeraudit_offer_index ON public.blc_offer_audit USING btree (offer_id);
CREATE INDEX offeraudit_order_index ON public.blc_offer_audit USING btree (order_id);
CREATE TABLE blc_offer_info (
	offer_info_id int8 NOT NULL,
	CONSTRAINT blc_offer_info_pkey PRIMARY KEY (offer_info_id)
);
CREATE TABLE blc_offer_item_criteria (
	offer_item_criteria_id int8 NOT NULL,
	order_item_match_rule text NULL,
	quantity int4 NOT NULL,
	CONSTRAINT blc_offer_item_criteria_pkey PRIMARY KEY (offer_item_criteria_id)
);
CREATE TABLE blc_offer_rule (
	offer_rule_id int8 NOT NULL,
	match_rule text NULL,
	CONSTRAINT blc_offer_rule_pkey PRIMARY KEY (offer_rule_id)
);
CREATE TABLE blc_order_lock (
	lock_key varchar(255) NOT NULL,
	order_id int8 NOT NULL,
	last_updated int8 NULL,
	"locked" bpchar(1) NULL,
	CONSTRAINT blc_order_lock_pkey PRIMARY KEY (lock_key, order_id)
);
CREATE TABLE blc_page_item_criteria (
	page_item_criteria_id int8 NOT NULL,
	order_item_match_rule text NULL,
	quantity int4 NOT NULL,
	CONSTRAINT blc_page_item_criteria_pkey PRIMARY KEY (page_item_criteria_id)
);
CREATE TABLE blc_page_rule (
	page_rule_id int8 NOT NULL,
	match_rule text NULL,
	CONSTRAINT blc_page_rule_pkey PRIMARY KEY (page_rule_id)
);
CREATE TABLE blc_personal_message (
	personal_message_id int8 NOT NULL,
	message varchar(255) NULL,
	message_from varchar(255) NULL,
	message_to varchar(255) NULL,
	occasion varchar(255) NULL,
	CONSTRAINT blc_personal_message_pkey PRIMARY KEY (personal_message_id)
);
CREATE TABLE blc_phone (
	phone_id int8 NOT NULL,
	country_code varchar(255) NULL,
	"extension" varchar(255) NULL,
	is_active bool NULL,
	is_default bool NULL,
	phone_number varchar(255) NOT NULL,
	CONSTRAINT blc_phone_pkey PRIMARY KEY (phone_id)
);
CREATE TABLE blc_product_option (
	product_option_id int8 NOT NULL,
	attribute_name varchar(255) NULL,
	display_order int4 NULL,
	error_code varchar(255) NULL,
	error_message varchar(255) NULL,
	"label" varchar(255) NULL,
	long_description text NULL,
	"name" varchar(255) NULL,
	validation_strategy_type varchar(255) NULL,
	validation_type varchar(255) NULL,
	required bool NULL,
	option_type varchar(255) NULL,
	use_in_sku_generation bool NULL,
	validation_string varchar(255) NULL,
	CONSTRAINT blc_product_option_pkey PRIMARY KEY (product_option_id)
);
CREATE INDEX product_option_name_index ON public.blc_product_option USING btree (name);
CREATE TABLE blc_rating_summary (
	rating_summary_id int8 NOT NULL,
	average_rating float8 NOT NULL,
	item_id varchar(255) NOT NULL,
	rating_type varchar(255) NOT NULL,
	CONSTRAINT blc_rating_summary_pkey PRIMARY KEY (rating_summary_id)
);
CREATE INDEX ratingsumm_item_index ON public.blc_rating_summary USING btree (item_id);
CREATE INDEX ratingsumm_type_index ON public.blc_rating_summary USING btree (rating_type);
CREATE TABLE blc_role (
	role_id int8 NOT NULL,
	role_name varchar(255) NOT NULL,
	CONSTRAINT blc_role_pkey PRIMARY KEY (role_id)
);
CREATE INDEX role_name_index ON public.blc_role USING btree (role_name);
CREATE TABLE blc_sc_fld (
	sc_fld_id int8 NOT NULL,
	fld_key varchar(255) NULL,
	lob_value text NULL,
	value varchar(255) NULL,
	CONSTRAINT blc_sc_fld_pkey PRIMARY KEY (sc_fld_id)
);
CREATE TABLE blc_sc_fld_tmplt (
	sc_fld_tmplt_id int8 NOT NULL,
	"name" varchar(255) NULL,
	CONSTRAINT blc_sc_fld_tmplt_pkey PRIMARY KEY (sc_fld_tmplt_id)
);
CREATE TABLE blc_sc_rule (
	sc_rule_id int8 NOT NULL,
	match_rule text NULL,
	CONSTRAINT blc_sc_rule_pkey PRIMARY KEY (sc_rule_id)
);
CREATE TABLE blc_search_intercept (
	search_redirect_id int8 NOT NULL,
	active_end_date timestamp NULL,
	active_start_date timestamp NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	priority int4 NULL,
	search_term varchar(255) NOT NULL,
	url varchar(255) NOT NULL,
	CONSTRAINT blc_search_intercept_pkey PRIMARY KEY (search_redirect_id)
);
CREATE INDEX search_active_index ON public.blc_search_intercept USING btree (active_start_date, active_end_date);
CREATE TABLE blc_search_synonym (
	search_synonym_id int8 NOT NULL,
	synonyms varchar(255) NULL,
	term varchar(255) NULL,
	CONSTRAINT blc_search_synonym_pkey PRIMARY KEY (search_synonym_id)
);
CREATE INDEX searchsynonym_term_index ON public.blc_search_synonym USING btree (term);
CREATE TABLE blc_sku_availability (
	sku_availability_id int8 NOT NULL,
	availability_date timestamp NULL,
	availability_status varchar(255) NULL,
	location_id int8 NULL,
	qty_on_hand int4 NULL,
	reserve_qty int4 NULL,
	sku_id int8 NULL,
	CONSTRAINT blc_sku_availability_pkey PRIMARY KEY (sku_availability_id)
);
CREATE INDEX skuavail_location_index ON public.blc_sku_availability USING btree (location_id);
CREATE INDEX skuavail_sku_index ON public.blc_sku_availability USING btree (sku_id);
CREATE INDEX skuavail_status_index ON public.blc_sku_availability USING btree (availability_status);
CREATE TABLE blc_static_asset (
	static_asset_id int8 NOT NULL,
	alt_text varchar(255) NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	file_extension varchar(255) NULL,
	file_size int8 NULL,
	full_url varchar(255) NOT NULL,
	mime_type varchar(255) NULL,
	"name" varchar(255) NOT NULL,
	storage_type varchar(255) NULL,
	title varchar(255) NULL,
	CONSTRAINT blc_static_asset_pkey PRIMARY KEY (static_asset_id)
);
CREATE INDEX asst_full_url_indx ON public.blc_static_asset USING btree (full_url);
CREATE TABLE blc_static_asset_desc (
	static_asset_desc_id int8 NOT NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	description varchar(255) NULL,
	long_description varchar(255) NULL,
	CONSTRAINT blc_static_asset_desc_pkey PRIMARY KEY (static_asset_desc_id)
);
CREATE TABLE blc_static_asset_strg (
	static_asset_strg_id int8 NOT NULL,
	file_data oid NULL,
	static_asset_id int8 NOT NULL,
	CONSTRAINT blc_static_asset_strg_pkey PRIMARY KEY (static_asset_strg_id)
);
CREATE INDEX static_asset_id_index ON public.blc_static_asset_strg USING btree (static_asset_id);
CREATE TABLE blc_system_property (
	blc_system_property_id int8 NOT NULL,
	friendly_group varchar(255) NULL,
	friendly_name varchar(255) NULL,
	friendly_tab varchar(255) NULL,
	property_name varchar(255) NOT NULL,
	override_generated_prop_name bool NULL,
	property_type varchar(255) NULL,
	property_value varchar(255) NULL,
	CONSTRAINT blc_system_property_pkey PRIMARY KEY (blc_system_property_id)
);
CREATE INDEX idx_blsypr_property_name ON public.blc_system_property USING btree (property_name);
CREATE TABLE blc_translation (
	translation_id int8 NOT NULL,
	entity_id varchar(255) NULL,
	entity_type varchar(255) NULL,
	field_name varchar(255) NULL,
	locale_code varchar(255) NULL,
	translated_value text NULL,
	CONSTRAINT blc_translation_pkey PRIMARY KEY (translation_id)
);
CREATE INDEX translation_index ON public.blc_translation USING btree (entity_type, entity_id, field_name, locale_code);
CREATE TABLE blc_url_handler (
	url_handler_id int8 NOT NULL,
	incoming_url varchar(255) NOT NULL,
	is_regex bool NULL,
	new_url varchar(255) NOT NULL,
	url_redirect_type varchar(255) NULL,
	CONSTRAINT blc_url_handler_pkey PRIMARY KEY (url_handler_id)
);
CREATE INDEX incoming_url_index ON public.blc_url_handler USING btree (incoming_url);
CREATE INDEX is_regex_handler_index ON public.blc_url_handler USING btree (is_regex);
CREATE TABLE blc_userconnection (
	providerid varchar(255) NOT NULL,
	provideruserid varchar(255) NOT NULL,
	userid varchar(255) NOT NULL,
	accesstoken varchar(255) NOT NULL,
	displayname varchar(255) NULL,
	expiretime int8 NULL,
	imageurl varchar(255) NULL,
	profileurl varchar(255) NULL,
	"rank" int4 NOT NULL,
	refreshtoken varchar(255) NULL,
	secret varchar(255) NULL,
	CONSTRAINT blc_userconnection_pkey PRIMARY KEY (providerid, provideruserid, userid)
);
CREATE TABLE blc_zip_code (
	zip_code_id varchar(255) NOT NULL,
	zip_city varchar(255) NULL,
	zip_latitude float8 NULL,
	zip_longitude float8 NULL,
	zip_state varchar(255) NULL,
	zipcode int4 NULL,
	CONSTRAINT blc_zip_code_pkey PRIMARY KEY (zip_code_id)
);
CREATE INDEX zipcode_city_index ON public.blc_zip_code USING btree (zip_city);
CREATE INDEX zipcode_latitude_index ON public.blc_zip_code USING btree (zip_latitude);
CREATE INDEX zipcode_longitude_index ON public.blc_zip_code USING btree (zip_longitude);
CREATE INDEX zipcode_state_index ON public.blc_zip_code USING btree (zip_state);
CREATE INDEX zipcode_zip_index ON public.blc_zip_code USING btree (zipcode);
CREATE TABLE sequence_generator (
	id_name varchar(255) NOT NULL,
	id_val int8 NULL,
	CONSTRAINT sequence_generator_pkey PRIMARY KEY (id_name)
);
CREATE TABLE blc_address (
	address_id int8 NOT NULL,
	address_line1 varchar(255) NOT NULL,
	address_line2 varchar(255) NULL,
	address_line3 varchar(255) NULL,
	city varchar(255) NOT NULL,
	company_name varchar(255) NULL,
	county varchar(255) NULL,
	email_address varchar(255) NULL,
	fax varchar(255) NULL,
	first_name varchar(255) NULL,
	full_name varchar(255) NULL,
	is_active bool NULL,
	is_business bool NULL,
	is_default bool NULL,
	is_mailing bool NULL,
	is_street bool NULL,
	iso_country_sub varchar(255) NULL,
	last_name varchar(255) NULL,
	postal_code varchar(255) NULL,
	primary_phone varchar(255) NULL,
	secondary_phone varchar(255) NULL,
	standardized bool NULL,
	sub_state_prov_reg varchar(255) NULL,
	tokenized_address varchar(255) NULL,
	verification_level varchar(255) NULL,
	zip_four varchar(255) NULL,
	iso_country_alpha2 varchar(255) NULL,
	phone_fax_id int8 NULL,
	phone_primary_id int8 NULL,
	phone_secondary_id int8 NULL,
	CONSTRAINT blc_address_pkey PRIMARY KEY (address_id),
	CONSTRAINT fklafhchfputda32qhub54fa726 FOREIGN KEY (phone_primary_id) REFERENCES blc_phone(phone_id),
	CONSTRAINT fklbxqgy7cjnjn5ey2wqvpjnhe5 FOREIGN KEY (phone_secondary_id) REFERENCES blc_phone(phone_id),
	CONSTRAINT fkp37ru1cyeu6fq48ohmjmyvjej FOREIGN KEY (iso_country_alpha2) REFERENCES blc_iso_country(alpha_2),
	CONSTRAINT fkrgw6kfwuqepeo3u7i75t57l8w FOREIGN KEY (phone_fax_id) REFERENCES blc_phone(phone_id)
);
CREATE INDEX address_county_index ON public.blc_address USING btree (county);
CREATE INDEX address_iso_country_idx ON public.blc_address USING btree (iso_country_alpha2);
CREATE INDEX address_phone_fax_idx ON public.blc_address USING btree (phone_fax_id);
CREATE INDEX address_phone_pri_idx ON public.blc_address USING btree (phone_primary_id);
CREATE INDEX address_phone_sec_idx ON public.blc_address USING btree (phone_secondary_id);
CREATE TABLE blc_admin_permission_entity (
	admin_permission_entity_id int8 NOT NULL,
	ceiling_entity varchar(255) NOT NULL,
	admin_permission_id int8 NULL,
	CONSTRAINT blc_admin_permission_entity_pkey PRIMARY KEY (admin_permission_entity_id),
	CONSTRAINT fkr7lum3wwl9kacdlgw4cwdrsas FOREIGN KEY (admin_permission_id) REFERENCES blc_admin_permission(admin_permission_id)
);
CREATE TABLE blc_admin_permission_xref (
	child_permission_id int8 NOT NULL,
	admin_permission_id int8 NOT NULL,
	CONSTRAINT fk1m3h00oqtternnpeiupslooan FOREIGN KEY (admin_permission_id) REFERENCES blc_admin_permission(admin_permission_id),
	CONSTRAINT fk9gfarfrwe5wnew41w9oyd3j6y FOREIGN KEY (child_permission_id) REFERENCES blc_admin_permission(admin_permission_id)
);
CREATE TABLE blc_admin_role_permission_xref (
	admin_role_id int8 NOT NULL,
	admin_permission_id int8 NOT NULL,
	CONSTRAINT blc_admin_role_permission_xref_pkey PRIMARY KEY (admin_permission_id, admin_role_id),
	CONSTRAINT fkl1jm8qymrs3laxvyawcb7mlbt FOREIGN KEY (admin_role_id) REFERENCES blc_admin_role(admin_role_id),
	CONSTRAINT fkoj1ji2ummmtfdm0xb9jesi7g FOREIGN KEY (admin_permission_id) REFERENCES blc_admin_permission(admin_permission_id)
);
CREATE TABLE blc_admin_section (
	admin_section_id int8 NOT NULL,
	ceiling_entity varchar(255) NULL,
	display_controller varchar(255) NULL,
	display_order int4 NULL,
	folderable bool NULL,
	foldered_by_default bool NULL,
	"name" varchar(255) NOT NULL,
	section_key varchar(255) NOT NULL,
	url varchar(255) NULL,
	use_default_handler bool NULL,
	admin_module_id int8 NOT NULL,
	CONSTRAINT blc_admin_section_pkey PRIMARY KEY (admin_section_id),
	CONSTRAINT uk_2l8u0qyluf4fwp2iiqp3p4jrn UNIQUE (section_key),
	CONSTRAINT fk2gpd1e839i00bosr6e54mdnn2 FOREIGN KEY (admin_module_id) REFERENCES blc_admin_module(admin_module_id)
);
CREATE INDEX adminsection_module_index ON public.blc_admin_section USING btree (admin_module_id);
CREATE INDEX adminsection_name_index ON public.blc_admin_section USING btree (name);
CREATE TABLE blc_admin_user_addtl_fields (
	attribute_id int8 NOT NULL,
	field_name varchar(255) NOT NULL,
	field_value varchar(255) NULL,
	admin_user_id int8 NOT NULL,
	CONSTRAINT blc_admin_user_addtl_fields_pkey PRIMARY KEY (attribute_id),
	CONSTRAINT fkiiateds21bej9b6qvrom06ayr FOREIGN KEY (admin_user_id) REFERENCES blc_admin_user(admin_user_id)
);
CREATE INDEX adminuserattribute_index ON public.blc_admin_user_addtl_fields USING btree (admin_user_id);
CREATE INDEX adminuserattribute_name_index ON public.blc_admin_user_addtl_fields USING btree (field_name);
CREATE TABLE blc_admin_user_permission_xref (
	admin_user_id int8 NOT NULL,
	admin_permission_id int8 NOT NULL,
	CONSTRAINT blc_admin_user_permission_xref_pkey PRIMARY KEY (admin_permission_id, admin_user_id),
	CONSTRAINT fk8ia4c6mqqvm9pt1aghjbvdmtb FOREIGN KEY (admin_permission_id) REFERENCES blc_admin_permission(admin_permission_id),
	CONSTRAINT fkj7ms4sgplv582id7ftu4thyn3 FOREIGN KEY (admin_user_id) REFERENCES blc_admin_user(admin_user_id)
);
CREATE TABLE blc_admin_user_role_xref (
	admin_user_id int8 NOT NULL,
	admin_role_id int8 NOT NULL,
	CONSTRAINT blc_admin_user_role_xref_pkey PRIMARY KEY (admin_role_id, admin_user_id),
	CONSTRAINT fk4skhb24d5kego6i7iw4y1a448 FOREIGN KEY (admin_role_id) REFERENCES blc_admin_role(admin_role_id),
	CONSTRAINT fka48q2nut9wd1cktfjp3l2f3xv FOREIGN KEY (admin_user_id) REFERENCES blc_admin_user(admin_user_id)
);
CREATE TABLE blc_asset_desc_map (
	static_asset_id int8 NOT NULL,
	static_asset_desc_id int8 NOT NULL,
	map_key varchar(255) NOT NULL,
	CONSTRAINT blc_asset_desc_map_pkey PRIMARY KEY (static_asset_id, map_key),
	CONSTRAINT fkheybfelcjp3ave1pxgmjl78eo FOREIGN KEY (static_asset_desc_id) REFERENCES blc_static_asset_desc(static_asset_desc_id),
	CONSTRAINT fkhrepj8vehjv59lcn3xbiq7ays FOREIGN KEY (static_asset_id) REFERENCES blc_static_asset(static_asset_id)
);
CREATE TABLE blc_category (
	category_id int8 NOT NULL,
	active_end_date timestamp NULL,
	active_start_date timestamp NULL,
	archived bpchar(1) NULL,
	description varchar(255) NULL,
	display_template varchar(255) NULL,
	external_id varchar(255) NULL,
	fulfillment_type varchar(255) NULL,
	inventory_type varchar(255) NULL,
	long_description text NULL,
	meta_desc varchar(255) NULL,
	meta_title varchar(255) NULL,
	"name" varchar(255) NOT NULL,
	override_generated_url bool NULL,
	product_desc_pattern_override varchar(255) NULL,
	product_title_pattern_override varchar(255) NULL,
	root_display_order numeric(10, 6) NULL,
	tax_code varchar(255) NULL,
	url varchar(255) NULL,
	url_key varchar(255) NULL,
	default_parent_category_id int8 NULL,
	CONSTRAINT blc_category_pkey PRIMARY KEY (category_id),
	CONSTRAINT fk6lf7a3qgmh5m8aq8o8url408t FOREIGN KEY (default_parent_category_id) REFERENCES blc_category(category_id)
);
CREATE INDEX category_e_id_index ON public.blc_category USING btree (external_id);
CREATE INDEX category_name_index ON public.blc_category USING btree (name);
CREATE INDEX category_parent_index ON public.blc_category USING btree (default_parent_category_id);
CREATE INDEX category_url_index ON public.blc_category USING btree (url);
CREATE INDEX category_urlkey_index ON public.blc_category USING btree (url_key);
CREATE TABLE blc_category_attribute (
	category_attribute_id int8 NOT NULL,
	"name" varchar(255) NOT NULL,
	value varchar(255) NULL,
	category_id int8 NOT NULL,
	CONSTRAINT blc_category_attribute_pkey PRIMARY KEY (category_attribute_id),
	CONSTRAINT fkhkechh91jg8iog16ry7089anf FOREIGN KEY (category_id) REFERENCES blc_category(category_id)
);
CREATE INDEX categoryattribute_index ON public.blc_category_attribute USING btree (category_id);
CREATE INDEX categoryattribute_name_index ON public.blc_category_attribute USING btree (name);
CREATE TABLE blc_category_media_map (
	category_media_id int8 NOT NULL,
	map_key varchar(255) NOT NULL,
	blc_category_category_id int8 NOT NULL,
	media_id int8 NULL,
	CONSTRAINT blc_category_media_map_pkey PRIMARY KEY (category_media_id),
	CONSTRAINT fkel78nfydgqxta46k7uvsh5q3x FOREIGN KEY (media_id) REFERENCES blc_media(media_id),
	CONSTRAINT fkfi64uesmjfu96gc0i4urxyf12 FOREIGN KEY (blc_category_category_id) REFERENCES blc_category(category_id)
);
CREATE TABLE blc_category_xref (
	category_xref_id int8 NOT NULL,
	default_reference bool NULL,
	display_order numeric(10, 6) NULL,
	category_id int8 NOT NULL,
	sub_category_id int8 NOT NULL,
	CONSTRAINT blc_category_xref_pkey PRIMARY KEY (category_xref_id),
	CONSTRAINT fke9p1rqdyircs0atgu7e5xlwmx FOREIGN KEY (category_id) REFERENCES blc_category(category_id),
	CONSTRAINT fkgrlcy8qrn7lqyiou3vu1piuk1 FOREIGN KEY (sub_category_id) REFERENCES blc_category(category_id)
);
CREATE TABLE blc_country_sub (
	abbreviation varchar(255) NOT NULL,
	alt_abbreviation varchar(255) NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	"name" varchar(255) NOT NULL,
	country_sub_cat int8 NULL,
	country varchar(255) NOT NULL,
	CONSTRAINT blc_country_sub_pkey PRIMARY KEY (abbreviation),
	CONSTRAINT fkapicr4ced87ut6dfyh5fuway8 FOREIGN KEY (country) REFERENCES blc_country(abbreviation),
	CONSTRAINT fktjleoo0nukky2den29i7mlx0c FOREIGN KEY (country_sub_cat) REFERENCES blc_country_sub_cat(country_sub_cat_id)
);
CREATE INDEX country_sub_alt_abrv_idx ON public.blc_country_sub USING btree (alt_abbreviation);
CREATE INDEX country_sub_name_idx ON public.blc_country_sub USING btree (name);
CREATE TABLE blc_data_drvn_enum_val (
	enum_val_id int8 NOT NULL,
	display varchar(255) NULL,
	hidden bool NULL,
	enum_key varchar(255) NULL,
	enum_type int8 NULL,
	CONSTRAINT blc_data_drvn_enum_val_pkey PRIMARY KEY (enum_val_id),
	CONSTRAINT fkq180xbmiqw22rrc9kf0qokaea FOREIGN KEY (enum_type) REFERENCES blc_data_drvn_enum(enum_id)
);
CREATE INDEX enum_val_key_index ON public.blc_data_drvn_enum_val USING btree (enum_key);
CREATE INDEX hidden_index ON public.blc_data_drvn_enum_val USING btree (hidden);
CREATE TABLE blc_email_tracking_clicks (
	click_id int8 NOT NULL,
	customer_id varchar(255) NULL,
	date_clicked timestamp NOT NULL,
	destination_uri varchar(255) NULL,
	query_string varchar(255) NULL,
	email_tracking_id int8 NOT NULL,
	CONSTRAINT blc_email_tracking_clicks_pkey PRIMARY KEY (click_id),
	CONSTRAINT fk3jed270645ahspuibr8wau0po FOREIGN KEY (email_tracking_id) REFERENCES blc_email_tracking(email_tracking_id)
);
CREATE INDEX trackingclicks_customer_index ON public.blc_email_tracking_clicks USING btree (customer_id);
CREATE INDEX trackingclicks_tracking_index ON public.blc_email_tracking_clicks USING btree (email_tracking_id);
CREATE TABLE blc_email_tracking_opens (
	open_id int8 NOT NULL,
	date_opened timestamp NULL,
	user_agent varchar(255) NULL,
	email_tracking_id int8 NULL,
	CONSTRAINT blc_email_tracking_opens_pkey PRIMARY KEY (open_id),
	CONSTRAINT fkt6fi06g4y7riiqeyuhb0t659n FOREIGN KEY (email_tracking_id) REFERENCES blc_email_tracking(email_tracking_id)
);
CREATE INDEX trackingopen_tracking ON public.blc_email_tracking_opens USING btree (email_tracking_id);
CREATE TABLE blc_fld_def (
	fld_def_id int8 NOT NULL,
	allow_multiples bool NULL,
	column_width varchar(255) NULL,
	fld_order int4 NULL,
	fld_type varchar(255) NULL,
	friendly_name varchar(255) NULL,
	help_text varchar(255) NULL,
	hidden_flag bool NULL,
	hint varchar(255) NULL,
	max_length int4 NULL,
	"name" varchar(255) NULL,
	required_flag bool NULL,
	security_level varchar(255) NULL,
	text_area_flag bool NULL,
	tooltip varchar(255) NULL,
	vldtn_error_mssg_key varchar(255) NULL,
	vldtn_regex varchar(255) NULL,
	enum_id int8 NULL,
	fld_group_id int8 NULL,
	CONSTRAINT blc_fld_def_pkey PRIMARY KEY (fld_def_id),
	CONSTRAINT fk3p9sauu111ycv4gbk6tymcj9e FOREIGN KEY (enum_id) REFERENCES blc_data_drvn_enum(enum_id),
	CONSTRAINT fkcqfi7hxwka5y9sqqoiolsnssr FOREIGN KEY (fld_group_id) REFERENCES blc_fld_group(fld_group_id)
);
CREATE TABLE blc_fulfillment_opt_banded_prc (
	fulfillment_option_id int8 NOT NULL,
	CONSTRAINT blc_fulfillment_opt_banded_prc_pkey PRIMARY KEY (fulfillment_option_id),
	CONSTRAINT fksf9j5pdg9lo5e7xhasqn61m0y FOREIGN KEY (fulfillment_option_id) REFERENCES blc_fulfillment_option(fulfillment_option_id)
);
CREATE TABLE blc_fulfillment_opt_banded_wgt (
	fulfillment_option_id int8 NOT NULL,
	CONSTRAINT blc_fulfillment_opt_banded_wgt_pkey PRIMARY KEY (fulfillment_option_id),
	CONSTRAINT fksarbwhn57dgx7kt1es3ny384n FOREIGN KEY (fulfillment_option_id) REFERENCES blc_fulfillment_option(fulfillment_option_id)
);
CREATE TABLE blc_fulfillment_option_fixed (
	price numeric(19, 5) NOT NULL,
	fulfillment_option_id int8 NOT NULL,
	currency_code varchar(255) NULL,
	CONSTRAINT blc_fulfillment_option_fixed_pkey PRIMARY KEY (fulfillment_option_id),
	CONSTRAINT fkj5n6w6q7dk09k6ayif4g5t0t3 FOREIGN KEY (currency_code) REFERENCES blc_currency(currency_code),
	CONSTRAINT fkl96yhpl4w0989nil2g6v2t3kq FOREIGN KEY (fulfillment_option_id) REFERENCES blc_fulfillment_option(fulfillment_option_id)
);
CREATE TABLE blc_fulfillment_price_band (
	fulfillment_price_band_id int8 NOT NULL,
	result_amount numeric(19, 5) NOT NULL,
	result_amount_type varchar(255) NOT NULL,
	retail_price_minimum_amount numeric(19, 5) NOT NULL,
	fulfillment_option_id int8 NULL,
	CONSTRAINT blc_fulfillment_price_band_pkey PRIMARY KEY (fulfillment_price_band_id),
	CONSTRAINT fkh2i7xep5l3txpi65xpb3vxxdh FOREIGN KEY (fulfillment_option_id) REFERENCES blc_fulfillment_opt_banded_prc(fulfillment_option_id)
);
CREATE TABLE blc_fulfillment_weight_band (
	fulfillment_weight_band_id int8 NOT NULL,
	result_amount numeric(19, 5) NOT NULL,
	result_amount_type varchar(255) NOT NULL,
	minimum_weight numeric(19, 5) NULL,
	weight_unit_of_measure varchar(255) NULL,
	fulfillment_option_id int8 NULL,
	CONSTRAINT blc_fulfillment_weight_band_pkey PRIMARY KEY (fulfillment_weight_band_id),
	CONSTRAINT fkoij3p9iwe1856w6fd5283bpyl FOREIGN KEY (fulfillment_option_id) REFERENCES blc_fulfillment_opt_banded_wgt(fulfillment_option_id)
);
CREATE TABLE blc_img_static_asset (
	height int4 NULL,
	width int4 NULL,
	static_asset_id int8 NOT NULL,
	CONSTRAINT blc_img_static_asset_pkey PRIMARY KEY (static_asset_id),
	CONSTRAINT fk6pugoo2mcm6irchv42bui3tm6 FOREIGN KEY (static_asset_id) REFERENCES blc_static_asset(static_asset_id)
);
CREATE TABLE blc_index_field (
	index_field_id int8 NOT NULL,
	archived bpchar(1) NULL,
	searchable bool NULL,
	field_id int8 NOT NULL,
	CONSTRAINT blc_index_field_pkey PRIMARY KEY (index_field_id),
	CONSTRAINT fkc1x5lu6romu8tsjrlpjmsqqxy FOREIGN KEY (field_id) REFERENCES blc_field(field_id)
);
CREATE INDEX index_field_searchable_index ON public.blc_index_field USING btree (searchable);
CREATE TABLE blc_index_field_type (
	index_field_type_id int8 NOT NULL,
	archived bpchar(1) NULL,
	field_type varchar(255) NULL,
	index_field_id int8 NOT NULL,
	CONSTRAINT blc_index_field_type_pkey PRIMARY KEY (index_field_type_id),
	CONSTRAINT fkmv0l0yt2099ffo2pjdrdbbj9h FOREIGN KEY (index_field_id) REFERENCES blc_index_field(index_field_id)
);
CREATE TABLE blc_locale (
	locale_code varchar(255) NOT NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	default_flag bool NULL,
	friendly_name varchar(255) NULL,
	use_in_search_index bool NULL,
	currency_code varchar(255) NULL,
	CONSTRAINT blc_locale_pkey PRIMARY KEY (locale_code),
	CONSTRAINT fk6gs37rhrtyd5ei2oqspxxrc3x FOREIGN KEY (currency_code) REFERENCES blc_currency(currency_code)
);
CREATE TABLE blc_offer_code (
	offer_code_id int8 NOT NULL,
	archived bpchar(1) NULL,
	email_address varchar(255) NULL,
	max_uses int4 NULL,
	offer_code varchar(255) NOT NULL,
	end_date timestamp NULL,
	start_date timestamp NULL,
	uses int4 NULL,
	offer_id int8 NOT NULL,
	CONSTRAINT blc_offer_code_pkey PRIMARY KEY (offer_code_id),
	CONSTRAINT fk4rcfx31u6n9esw1ob98u8o87o FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id)
);
CREATE INDEX offer_code_email_index ON public.blc_offer_code USING btree (email_address);
CREATE INDEX offercode_code_index ON public.blc_offer_code USING btree (offer_code);
CREATE INDEX offercode_offer_index ON public.blc_offer_code USING btree (offer_id);
CREATE TABLE blc_offer_info_fields (
	offer_info_fields_id int8 NOT NULL,
	field_value varchar(255) NULL,
	field_name varchar(255) NOT NULL,
	CONSTRAINT blc_offer_info_fields_pkey PRIMARY KEY (offer_info_fields_id, field_name),
	CONSTRAINT fkohr0h2751uuxgawkbkakbptqn FOREIGN KEY (offer_info_fields_id) REFERENCES blc_offer_info(offer_info_id)
);
CREATE TABLE blc_offer_price_data (
	offer_price_data_id int8 NOT NULL,
	end_date timestamp NULL,
	start_date timestamp NULL,
	amount numeric(19, 5) NOT NULL,
	archived bpchar(1) NULL,
	discount_type varchar(255) NULL,
	identifier_type varchar(255) NULL,
	identifier_value varchar(255) NULL,
	quantity int4 NOT NULL,
	offer_id int8 NOT NULL,
	CONSTRAINT blc_offer_price_data_pkey PRIMARY KEY (offer_price_data_id),
	CONSTRAINT fkkprkllbin16hh5ay89t8we431 FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id)
);
CREATE INDEX offer_price_data_offer_index ON public.blc_offer_price_data USING btree (offer_id);
CREATE TABLE blc_offer_rule_map (
	offer_offer_rule_id int8 NOT NULL,
	map_key varchar(255) NOT NULL,
	blc_offer_offer_id int8 NOT NULL,
	offer_rule_id int8 NULL,
	CONSTRAINT blc_offer_rule_map_pkey PRIMARY KEY (offer_offer_rule_id),
	CONSTRAINT fk8ndq3dtgs1cr4ds9eil3sxcti FOREIGN KEY (offer_rule_id) REFERENCES blc_offer_rule(offer_rule_id),
	CONSTRAINT fkmwr0lw44aa4hulm8c9qg39x9x FOREIGN KEY (blc_offer_offer_id) REFERENCES blc_offer(offer_id)
);
CREATE TABLE blc_page_tmplt (
	page_tmplt_id int8 NOT NULL,
	tmplt_descr varchar(255) NULL,
	tmplt_name varchar(255) NULL,
	tmplt_path varchar(255) NULL,
	locale_code varchar(255) NULL,
	CONSTRAINT blc_page_tmplt_pkey PRIMARY KEY (page_tmplt_id),
	CONSTRAINT fk19poavwqssando5ll1qid9kmf FOREIGN KEY (locale_code) REFERENCES blc_locale(locale_code)
);
CREATE TABLE blc_pgtmplt_fldgrp_xref (
	pg_tmplt_fld_grp_id int8 NOT NULL,
	group_order numeric(10, 6) NULL,
	fld_group_id int8 NULL,
	page_tmplt_id int8 NULL,
	CONSTRAINT blc_pgtmplt_fldgrp_xref_pkey PRIMARY KEY (pg_tmplt_fld_grp_id),
	CONSTRAINT fkoy5hlxlq3pii0gj8yalskxv88 FOREIGN KEY (fld_group_id) REFERENCES blc_fld_group(fld_group_id),
	CONSTRAINT fkr3xcn67le94r6oxnaebm5ebnk FOREIGN KEY (page_tmplt_id) REFERENCES blc_page_tmplt(page_tmplt_id)
);
CREATE TABLE blc_product (
	product_id int8 NOT NULL,
	archived bpchar(1) NULL,
	can_sell_without_options bool NULL,
	canonical_url varchar(255) NULL,
	display_template varchar(255) NULL,
	enable_default_sku_in_inventory bool NULL,
	manufacture varchar(255) NULL,
	meta_desc varchar(255) NULL,
	meta_title varchar(255) NULL,
	model varchar(255) NULL,
	override_generated_url bool NULL,
	url varchar(255) NULL,
	url_key varchar(255) NULL,
	default_category_id int8 NULL,
	default_sku_id int8 NULL,
	CONSTRAINT blc_product_pkey PRIMARY KEY (product_id),
	CONSTRAINT fk57aoxhpvwg389v7sx4m153mde FOREIGN KEY (default_category_id) REFERENCES blc_category(category_id)
);
CREATE INDEX product_url_index ON public.blc_product USING btree (url, url_key);
CREATE INDEX product_url_key_index ON public.blc_product USING btree (url_key);
CREATE TABLE blc_product_attribute (
	product_attribute_id int8 NOT NULL,
	"name" varchar(255) NOT NULL,
	value varchar(255) NULL,
	product_id int8 NOT NULL,
	CONSTRAINT blc_product_attribute_pkey PRIMARY KEY (product_attribute_id),
	CONSTRAINT fk5rahmy0l6hsgbvgb37ojlkx09 FOREIGN KEY (product_id) REFERENCES blc_product(product_id)
);
CREATE INDEX productattribute_index ON public.blc_product_attribute USING btree (product_id);
CREATE INDEX productattribute_name_index ON public.blc_product_attribute USING btree (name);
CREATE TABLE blc_product_bundle (
	auto_bundle bool NULL,
	bundle_promotable bool NULL,
	items_promotable bool NULL,
	pricing_model varchar(255) NULL,
	bundle_priority int4 NULL,
	product_id int8 NOT NULL,
	CONSTRAINT blc_product_bundle_pkey PRIMARY KEY (product_id),
	CONSTRAINT fk2hern8ie7vx4k6cawbryglb9g FOREIGN KEY (product_id) REFERENCES blc_product(product_id)
);
CREATE TABLE blc_product_cross_sale (
	cross_sale_product_id int8 NOT NULL,
	promotion_message varchar(255) NULL,
	"sequence" numeric(10, 6) NULL,
	category_id int8 NULL,
	product_id int8 NULL,
	related_sale_product_id int8 NOT NULL,
	CONSTRAINT blc_product_cross_sale_pkey PRIMARY KEY (cross_sale_product_id),
	CONSTRAINT fkak6hk19vp8ioy27lrt7x9be7w FOREIGN KEY (category_id) REFERENCES blc_category(category_id),
	CONSTRAINT fkeq0i4yj6td2qxh0tnekeomrxk FOREIGN KEY (product_id) REFERENCES blc_product(product_id),
	CONSTRAINT fkovg4s9i9ejesgcygfpyjq7eqa FOREIGN KEY (related_sale_product_id) REFERENCES blc_product(product_id)
);
CREATE INDEX crosssale_category_index ON public.blc_product_cross_sale USING btree (category_id);
CREATE INDEX crosssale_index ON public.blc_product_cross_sale USING btree (product_id);
CREATE INDEX crosssale_related_index ON public.blc_product_cross_sale USING btree (related_sale_product_id);
CREATE TABLE blc_product_featured (
	featured_product_id int8 NOT NULL,
	promotion_message varchar(255) NULL,
	"sequence" numeric(10, 6) NULL,
	category_id int8 NULL,
	product_id int8 NULL,
	CONSTRAINT blc_product_featured_pkey PRIMARY KEY (featured_product_id),
	CONSTRAINT fk83xhh0734p8wo8w64di8qwd9o FOREIGN KEY (category_id) REFERENCES blc_category(category_id),
	CONSTRAINT fkr4v6adrqqmd1qe09i6mb8fj8p FOREIGN KEY (product_id) REFERENCES blc_product(product_id)
);
CREATE INDEX prodfeatured_category_index ON public.blc_product_featured USING btree (category_id);
CREATE INDEX prodfeatured_product_index ON public.blc_product_featured USING btree (product_id);
CREATE TABLE blc_product_option_value (
	product_option_value_id int8 NOT NULL,
	attribute_value varchar(255) NULL,
	display_order int8 NULL,
	price_adjustment numeric(19, 5) NULL,
	product_option_id int8 NULL,
	CONSTRAINT blc_product_option_value_pkey PRIMARY KEY (product_option_value_id),
	CONSTRAINT fk9ixc1gbymkin77d6krmc3mub7 FOREIGN KEY (product_option_id) REFERENCES blc_product_option(product_option_id)
);
CREATE TABLE blc_product_option_xref (
	product_option_xref_id int8 NOT NULL,
	product_id int8 NOT NULL,
	product_option_id int8 NOT NULL,
	CONSTRAINT blc_product_option_xref_pkey PRIMARY KEY (product_option_xref_id),
	CONSTRAINT fkhqikdn2uw75plhcwn4cmjtt4m FOREIGN KEY (product_id) REFERENCES blc_product(product_id),
	CONSTRAINT fkswm8iotfkm6a9iyj6ru3rmpv8 FOREIGN KEY (product_option_id) REFERENCES blc_product_option(product_option_id)
);
CREATE TABLE blc_product_up_sale (
	up_sale_product_id int8 NOT NULL,
	promotion_message varchar(255) NULL,
	"sequence" numeric(10, 6) NULL,
	category_id int8 NULL,
	product_id int8 NULL,
	related_sale_product_id int8 NULL,
	CONSTRAINT blc_product_up_sale_pkey PRIMARY KEY (up_sale_product_id),
	CONSTRAINT fkgefhcqet8553xhh9bdjb1jbjo FOREIGN KEY (product_id) REFERENCES blc_product(product_id),
	CONSTRAINT fkkcx4cl73kctdinewei1fk2vvl FOREIGN KEY (category_id) REFERENCES blc_category(category_id),
	CONSTRAINT fkm1r8s9j593gpcgluy5uyfa862 FOREIGN KEY (related_sale_product_id) REFERENCES blc_product(product_id)
);
CREATE INDEX upsale_category_index ON public.blc_product_up_sale USING btree (category_id);
CREATE INDEX upsale_product_index ON public.blc_product_up_sale USING btree (product_id);
CREATE INDEX upsale_related_index ON public.blc_product_up_sale USING btree (related_sale_product_id);
CREATE TABLE blc_promotion_message (
	promotion_message_id int8 NOT NULL,
	archived bpchar(1) NULL,
	end_date timestamp NULL,
	promotion_messasge varchar(255) NULL,
	message_placement varchar(255) NULL,
	"name" varchar(255) NULL,
	promotion_message_priority int4 NULL,
	start_date timestamp NULL,
	locale_code varchar(255) NULL,
	media_id int8 NULL,
	CONSTRAINT blc_promotion_message_pkey PRIMARY KEY (promotion_message_id),
	CONSTRAINT fk3dgs3j2b8mshpafd25qvtftgv FOREIGN KEY (locale_code) REFERENCES blc_locale(locale_code),
	CONSTRAINT fk59dkr5skbs8ve2v27truld8kc FOREIGN KEY (media_id) REFERENCES blc_media(media_id)
);
CREATE INDEX promotion_message_name_index ON public.blc_promotion_message USING btree (name);
CREATE TABLE blc_qual_crit_offer_xref (
	offer_qual_crit_id int8 NOT NULL,
	offer_id int8 NOT NULL,
	offer_item_criteria_id int8 NULL,
	CONSTRAINT blc_qual_crit_offer_xref_pkey PRIMARY KEY (offer_qual_crit_id),
	CONSTRAINT fk6e8y3yk68wvw90gtsesduqbrb FOREIGN KEY (offer_item_criteria_id) REFERENCES blc_offer_item_criteria(offer_item_criteria_id),
	CONSTRAINT fkmmxl8rjhiuu6hc7qhyy85pjov FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id)
);
CREATE TABLE blc_sandbox (
	sandbox_id int8 NOT NULL,
	archived bpchar(1) NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	author int8 NULL,
	color varchar(255) NULL,
	description varchar(255) NULL,
	go_live_date timestamp NULL,
	sandbox_name varchar(255) NULL,
	sandbox_type varchar(255) NULL,
	parent_sandbox_id int8 NULL,
	CONSTRAINT blc_sandbox_pkey PRIMARY KEY (sandbox_id),
	CONSTRAINT fk5e7j7mfpr1en8q48yxkbjmduw FOREIGN KEY (parent_sandbox_id) REFERENCES blc_sandbox(sandbox_id)
);
CREATE INDEX sandbox_name_index ON public.blc_sandbox USING btree (sandbox_name);
CREATE TABLE blc_sandbox_mgmt (
	sandbox_mgmt_id int8 NOT NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	sandbox_id int8 NOT NULL,
	CONSTRAINT blc_sandbox_mgmt_pkey PRIMARY KEY (sandbox_mgmt_id),
	CONSTRAINT uk_owins1o4pyal0j5pdlqfkd4b7 UNIQUE (sandbox_id),
	CONSTRAINT fkri581qivns8jshddbsl6m83hr FOREIGN KEY (sandbox_id) REFERENCES blc_sandbox(sandbox_id)
);
CREATE TABLE blc_sc_fldgrp_xref (
	blc_sc_fldgrp_xref_id int8 NOT NULL,
	group_order int4 NULL,
	fld_group_id int8 NULL,
	sc_fld_tmplt_id int8 NULL,
	CONSTRAINT blc_sc_fldgrp_xref_pkey PRIMARY KEY (blc_sc_fldgrp_xref_id),
	CONSTRAINT fkotfd5rhlje73tldskasabxd7k FOREIGN KEY (fld_group_id) REFERENCES blc_fld_group(fld_group_id),
	CONSTRAINT fktqvhk2j6dxo8kruvflvpstgf FOREIGN KEY (sc_fld_tmplt_id) REFERENCES blc_sc_fld_tmplt(sc_fld_tmplt_id)
);
CREATE TABLE blc_sc_type (
	sc_type_id int8 NOT NULL,
	description varchar(255) NULL,
	"name" varchar(255) NULL,
	sc_fld_tmplt_id int8 NULL,
	CONSTRAINT blc_sc_type_pkey PRIMARY KEY (sc_type_id),
	CONSTRAINT fkh7idqa2kh7vusepjor3bc80b3 FOREIGN KEY (sc_fld_tmplt_id) REFERENCES blc_sc_fld_tmplt(sc_fld_tmplt_id)
);
CREATE INDEX sc_type_name_index ON public.blc_sc_type USING btree (name);
CREATE TABLE blc_search_facet (
	search_facet_id int8 NOT NULL,
	archived bpchar(1) NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	multiselect bool NULL,
	"label" varchar(255) NULL,
	"name" varchar(255) NULL,
	requires_all_dependent bool NULL,
	search_display_priority int4 NULL,
	show_on_search bool NULL,
	use_facet_ranges bool NULL,
	index_field_type_id int8 NULL,
	CONSTRAINT blc_search_facet_pkey PRIMARY KEY (search_facet_id),
	CONSTRAINT fkrrhp7pwx3bjh2rhadrtv2ro81 FOREIGN KEY (index_field_type_id) REFERENCES blc_index_field_type(index_field_type_id)
);
CREATE TABLE blc_search_facet_range (
	search_facet_range_id int8 NOT NULL,
	archived bpchar(1) NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	max_value numeric(19, 5) NULL,
	min_value numeric(19, 5) NOT NULL,
	search_facet_id int8 NULL,
	CONSTRAINT blc_search_facet_range_pkey PRIMARY KEY (search_facet_range_id),
	CONSTRAINT fkm1k6kfkc59c8jdx51qym3tcai FOREIGN KEY (search_facet_id) REFERENCES blc_search_facet(search_facet_id)
);
CREATE INDEX search_facet_index ON public.blc_search_facet_range USING btree (search_facet_id);
CREATE TABLE blc_search_facet_xref (
	id int8 NOT NULL,
	required_facet_id int8 NOT NULL,
	search_facet_id int8 NOT NULL,
	CONSTRAINT blc_search_facet_xref_pkey PRIMARY KEY (id),
	CONSTRAINT fk4xpicfgot9h1utp316cbi8268 FOREIGN KEY (required_facet_id) REFERENCES blc_search_facet(search_facet_id),
	CONSTRAINT fktdvsplwk8dl6mnb0p7fdyte13 FOREIGN KEY (search_facet_id) REFERENCES blc_search_facet(search_facet_id)
);
CREATE TABLE blc_site (
	site_id int8 NOT NULL,
	archived bpchar(1) NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	deactivated bool NULL,
	"name" varchar(255) NULL,
	site_identifier_type varchar(255) NULL,
	site_identifier_value varchar(255) NULL,
	default_locale varchar(255) NULL,
	CONSTRAINT blc_site_pkey PRIMARY KEY (site_id),
	CONSTRAINT fka6umcqko5gobxtd4hvv8g2d4o FOREIGN KEY (default_locale) REFERENCES blc_locale(locale_code)
);
CREATE INDEX blc_site_id_val_index ON public.blc_site USING btree (site_identifier_value);
CREATE TABLE blc_site_catalog (
	site_catalog_xref_id int8 NOT NULL,
	catalog_id int8 NOT NULL,
	site_id int8 NOT NULL,
	CONSTRAINT blc_site_catalog_pkey PRIMARY KEY (site_catalog_xref_id),
	CONSTRAINT fkho5bxxfvt21ijan47er38vnyu FOREIGN KEY (catalog_id) REFERENCES blc_catalog(catalog_id),
	CONSTRAINT fkmktxeb1okchyhs2mxat1nk6s5 FOREIGN KEY (site_id) REFERENCES blc_site(site_id)
);
CREATE TABLE blc_site_map_cfg (
	indexed_site_map_file_name varchar(255) NULL,
	indexed_site_map_file_pattern varchar(255) NULL,
	max_url_entries_per_file int4 NULL,
	site_map_file_name varchar(255) NULL,
	module_config_id int8 NOT NULL,
	CONSTRAINT blc_site_map_cfg_pkey PRIMARY KEY (module_config_id),
	CONSTRAINT fkdskgdyr42vk7c8g92bxir3wej FOREIGN KEY (module_config_id) REFERENCES blc_module_configuration(module_config_id)
);
CREATE TABLE blc_site_map_gen_cfg (
	gen_config_id int8 NOT NULL,
	change_freq varchar(255) NOT NULL,
	disabled bool NOT NULL,
	generator_type varchar(255) NOT NULL,
	priority varchar(255) NULL,
	module_config_id int8 NOT NULL,
	CONSTRAINT blc_site_map_gen_cfg_pkey PRIMARY KEY (gen_config_id),
	CONSTRAINT fkmw4fb38sdenx8kjrxg5v8mjei FOREIGN KEY (module_config_id) REFERENCES blc_site_map_cfg(module_config_id)
);
CREATE TABLE blc_sku (
	sku_id int8 NOT NULL,
	active_end_date timestamp NULL,
	active_start_date timestamp NULL,
	available_flag bpchar(1) NULL,
	"cost" numeric(19, 5) NULL,
	description varchar(255) NULL,
	container_shape varchar(255) NULL,
	"depth" numeric(19, 2) NULL,
	dimension_unit_of_measure varchar(255) NULL,
	girth numeric(19, 2) NULL,
	height numeric(19, 2) NULL,
	container_size varchar(255) NULL,
	width numeric(19, 2) NULL,
	discountable_flag bpchar(1) NULL,
	display_template varchar(255) NULL,
	external_id varchar(255) NULL,
	fulfillment_type varchar(255) NULL,
	inventory_type varchar(255) NULL,
	is_machine_sortable bool NULL,
	long_description text NULL,
	"name" varchar(255) NULL,
	quantity_available int4 NULL,
	retail_price numeric(19, 5) NULL,
	sale_price numeric(19, 5) NULL,
	tax_code varchar(255) NULL,
	taxable_flag bpchar(1) NULL,
	upc varchar(255) NULL,
	url_key varchar(255) NULL,
	weight numeric(19, 2) NULL,
	weight_unit_of_measure varchar(255) NULL,
	currency_code varchar(255) NULL,
	default_product_id int8 NULL,
	addl_product_id int8 NULL,
	CONSTRAINT blc_sku_pkey PRIMARY KEY (sku_id),
	CONSTRAINT fkdowfc15iv11csxhs4itbfbowh FOREIGN KEY (currency_code) REFERENCES blc_currency(currency_code),
	CONSTRAINT fkseqmjfh1kdphq3eeplsuf6nj4 FOREIGN KEY (addl_product_id) REFERENCES blc_product(product_id)
);
CREATE INDEX sku_active_end_index ON public.blc_sku USING btree (active_end_date);
CREATE INDEX sku_active_start_index ON public.blc_sku USING btree (active_start_date);
CREATE INDEX sku_available_index ON public.blc_sku USING btree (available_flag);
CREATE INDEX sku_discountable_index ON public.blc_sku USING btree (discountable_flag);
CREATE INDEX sku_external_id_index ON public.blc_sku USING btree (external_id);
CREATE INDEX sku_name_index ON public.blc_sku USING btree (name);
CREATE INDEX sku_taxable_index ON public.blc_sku USING btree (taxable_flag);
CREATE INDEX sku_upc_index ON public.blc_sku USING btree (upc);
CREATE INDEX sku_url_key_index ON public.blc_sku USING btree (url_key);
CREATE TABLE blc_sku_attribute (
	sku_attr_id int8 NOT NULL,
	"name" varchar(255) NOT NULL,
	value varchar(255) NOT NULL,
	sku_id int8 NOT NULL,
	CONSTRAINT blc_sku_attribute_pkey PRIMARY KEY (sku_attr_id),
	CONSTRAINT fk6w8gul2489hdbmxha9ftu6qbq FOREIGN KEY (sku_id) REFERENCES blc_sku(sku_id)
);
CREATE INDEX skuattr_name_index ON public.blc_sku_attribute USING btree (name);
CREATE INDEX skuattr_sku_index ON public.blc_sku_attribute USING btree (sku_id);
CREATE TABLE blc_sku_bundle_item (
	sku_bundle_item_id int8 NOT NULL,
	item_sale_price numeric(19, 5) NULL,
	quantity int4 NOT NULL,
	"sequence" numeric(10, 6) NULL,
	product_bundle_id int8 NOT NULL,
	sku_id int8 NOT NULL,
	CONSTRAINT blc_sku_bundle_item_pkey PRIMARY KEY (sku_bundle_item_id),
	CONSTRAINT fk78yrrdqjalrqrei5kfnybkrs8 FOREIGN KEY (product_bundle_id) REFERENCES blc_product_bundle(product_id),
	CONSTRAINT fkbhe93esmsur5uyhl0v6dj392t FOREIGN KEY (sku_id) REFERENCES blc_sku(sku_id)
);
CREATE TABLE blc_sku_fee (
	sku_fee_id int8 NOT NULL,
	amount numeric(19, 5) NOT NULL,
	description varchar(255) NULL,
	"expression" text NULL,
	fee_type varchar(255) NULL,
	"name" varchar(255) NULL,
	taxable bool NULL,
	currency_code varchar(255) NULL,
	CONSTRAINT blc_sku_fee_pkey PRIMARY KEY (sku_fee_id),
	CONSTRAINT fkm9vf5c5ktjqu4wilne2f6926k FOREIGN KEY (currency_code) REFERENCES blc_currency(currency_code)
);
CREATE TABLE blc_sku_fee_xref (
	sku_fee_id int8 NOT NULL,
	sku_id int8 NOT NULL,
	CONSTRAINT fk3vmkt7ojjlpk2fle4cp8eq55f FOREIGN KEY (sku_id) REFERENCES blc_sku(sku_id),
	CONSTRAINT fkky8dxmeg4o49qyc7kiwojnuek FOREIGN KEY (sku_fee_id) REFERENCES blc_sku_fee(sku_fee_id)
);
CREATE TABLE blc_sku_fulfillment_excluded (
	sku_id int8 NOT NULL,
	fulfillment_option_id int8 NOT NULL,
	CONSTRAINT fkbf81qj596ta3xs2erw4o7m1ft FOREIGN KEY (fulfillment_option_id) REFERENCES blc_fulfillment_option(fulfillment_option_id),
	CONSTRAINT fks0toanodthismt1hugerli3e8 FOREIGN KEY (sku_id) REFERENCES blc_sku(sku_id)
);
CREATE TABLE blc_sku_fulfillment_flat_rates (
	sku_id int8 NOT NULL,
	rate numeric(19, 5) NULL,
	fulfillment_option_id int8 NOT NULL,
	CONSTRAINT blc_sku_fulfillment_flat_rates_pkey PRIMARY KEY (sku_id, fulfillment_option_id),
	CONSTRAINT fk1ruxosbsx27uicd9dny1ls9td FOREIGN KEY (sku_id) REFERENCES blc_sku(sku_id),
	CONSTRAINT fkklcbu8knfitgnhlj1sa2vyd30 FOREIGN KEY (fulfillment_option_id) REFERENCES blc_fulfillment_option(fulfillment_option_id)
);
CREATE TABLE blc_sku_media_map (
	sku_media_id int8 NOT NULL,
	map_key varchar(255) NOT NULL,
	media_id int8 NULL,
	blc_sku_sku_id int8 NOT NULL,
	CONSTRAINT blc_sku_media_map_pkey PRIMARY KEY (sku_media_id),
	CONSTRAINT fkc3mu07614ovbqwbnd1lxdg2ac FOREIGN KEY (blc_sku_sku_id) REFERENCES blc_sku(sku_id),
	CONSTRAINT fkl3netvy66i56mjj6bo43mjmn2 FOREIGN KEY (media_id) REFERENCES blc_media(media_id)
);
CREATE TABLE blc_sku_option_value_xref (
	sku_option_value_xref_id int8 NOT NULL,
	product_option_value_id int8 NOT NULL,
	sku_id int8 NOT NULL,
	CONSTRAINT blc_sku_option_value_xref_pkey PRIMARY KEY (sku_option_value_xref_id),
	CONSTRAINT fkc9e8sa4v1mqlbhd9hjp6bxujh FOREIGN KEY (product_option_value_id) REFERENCES blc_product_option_value(product_option_value_id),
	CONSTRAINT fkci6tv12pfsh2srrakx3ridy8v FOREIGN KEY (sku_id) REFERENCES blc_sku(sku_id)
);
CREATE TABLE blc_store (
	store_id int8 NOT NULL,
	archived bpchar(1) NULL,
	latitude float8 NULL,
	longitude float8 NULL,
	store_name varchar(255) NOT NULL,
	store_open bool NULL,
	store_hours varchar(255) NULL,
	store_number varchar(255) NULL,
	address_id int8 NULL,
	CONSTRAINT blc_store_pkey PRIMARY KEY (store_id),
	CONSTRAINT fkg65fln1wkn5rai85klf8ei1uy FOREIGN KEY (address_id) REFERENCES blc_address(address_id)
);
CREATE TABLE blc_tar_crit_offer_xref (
	offer_tar_crit_id int8 NOT NULL,
	offer_id int8 NOT NULL,
	offer_item_criteria_id int8 NULL,
	CONSTRAINT blc_tar_crit_offer_xref_pkey PRIMARY KEY (offer_tar_crit_id),
	CONSTRAINT fk5n28fyhs3hvqqn38rap5yns9i FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id),
	CONSTRAINT fkj44eau35bu6hfq5w53civq01y FOREIGN KEY (offer_item_criteria_id) REFERENCES blc_offer_item_criteria(offer_item_criteria_id)
);
CREATE TABLE blc_tax_detail (
	tax_detail_id int8 NOT NULL,
	amount numeric(19, 5) NULL,
	tax_country varchar(255) NULL,
	jurisdiction_name varchar(255) NULL,
	rate numeric(19, 5) NULL,
	tax_region varchar(255) NULL,
	tax_name varchar(255) NULL,
	"type" varchar(255) NULL,
	currency_code varchar(255) NULL,
	module_config_id int8 NULL,
	CONSTRAINT blc_tax_detail_pkey PRIMARY KEY (tax_detail_id),
	CONSTRAINT fk53heksajqlpbnfd8yrbudum8a FOREIGN KEY (module_config_id) REFERENCES blc_module_configuration(module_config_id),
	CONSTRAINT fk7rwcm52210yymslbjj8m25cvi FOREIGN KEY (currency_code) REFERENCES blc_currency(currency_code)
);
CREATE TABLE blc_admin_sec_perm_xref (
	admin_section_id int8 NOT NULL,
	admin_permission_id int8 NOT NULL,
	CONSTRAINT fk3k1buujb5let82ixj1k9nha3r FOREIGN KEY (admin_section_id) REFERENCES blc_admin_section(admin_section_id),
	CONSTRAINT fkns2d7kvauk8kgskridssn1gcv FOREIGN KEY (admin_permission_id) REFERENCES blc_admin_permission(admin_permission_id)
);
CREATE TABLE blc_admin_user_sandbox (
	sandbox_id int8 NULL,
	admin_user_id int8 NOT NULL,
	CONSTRAINT blc_admin_user_sandbox_pkey PRIMARY KEY (admin_user_id),
	CONSTRAINT fkay43c311x89bqu7lbswc7xan6 FOREIGN KEY (admin_user_id) REFERENCES blc_admin_user(admin_user_id),
	CONSTRAINT fkehxq8fct257ml7j0rbya7ripb FOREIGN KEY (sandbox_id) REFERENCES blc_sandbox(sandbox_id)
);
CREATE TABLE blc_cat_search_facet_excl_xref (
	cat_excl_search_facet_id int8 NOT NULL,
	"sequence" numeric(19, 2) NULL,
	category_id int8 NULL,
	search_facet_id int8 NULL,
	CONSTRAINT blc_cat_search_facet_excl_xref_pkey PRIMARY KEY (cat_excl_search_facet_id),
	CONSTRAINT fk66xu32canhiu19e6or98vufcw FOREIGN KEY (category_id) REFERENCES blc_category(category_id),
	CONSTRAINT fkmmy51xuqakfxoflomh4dgl7on FOREIGN KEY (search_facet_id) REFERENCES blc_search_facet(search_facet_id)
);
CREATE TABLE blc_cat_search_facet_xref (
	category_search_facet_id int8 NOT NULL,
	"sequence" numeric(19, 2) NULL,
	category_id int8 NULL,
	search_facet_id int8 NULL,
	CONSTRAINT blc_cat_search_facet_xref_pkey PRIMARY KEY (category_search_facet_id),
	CONSTRAINT fk15e8rvxxafd6h16c1ul3pynqh FOREIGN KEY (category_id) REFERENCES blc_category(category_id),
	CONSTRAINT fk68dqudo00pmvd760r53lmcq1q FOREIGN KEY (search_facet_id) REFERENCES blc_search_facet(search_facet_id)
);
CREATE TABLE blc_cat_site_map_gen_cfg (
	ending_depth int4 NOT NULL,
	starting_depth int4 NOT NULL,
	gen_config_id int8 NOT NULL,
	root_category_id int8 NOT NULL,
	CONSTRAINT blc_cat_site_map_gen_cfg_pkey PRIMARY KEY (gen_config_id),
	CONSTRAINT fkerl6k0300vd4y8haxljr92rmo FOREIGN KEY (root_category_id) REFERENCES blc_category(category_id),
	CONSTRAINT fkn5liq0ue5rtn6h7bmsv7q85nn FOREIGN KEY (gen_config_id) REFERENCES blc_site_map_gen_cfg(gen_config_id)
);
CREATE TABLE blc_category_product_xref (
	category_product_id int8 NOT NULL,
	default_reference bool NULL,
	display_order numeric(10, 6) NULL,
	category_id int8 NOT NULL,
	product_id int8 NOT NULL,
	CONSTRAINT blc_category_product_xref_pkey PRIMARY KEY (category_product_id),
	CONSTRAINT fkj8gn00lwi7fih0ueqwdat589e FOREIGN KEY (category_id) REFERENCES blc_category(category_id),
	CONSTRAINT fknwoet42m887na9hjfvqqgr58v FOREIGN KEY (product_id) REFERENCES blc_product(product_id)
);
CREATE TABLE blc_cust_site_map_gen_cfg (
	gen_config_id int8 NOT NULL,
	CONSTRAINT blc_cust_site_map_gen_cfg_pkey PRIMARY KEY (gen_config_id),
	CONSTRAINT fks5s4vmpbxh4edqjtbted9gxmw FOREIGN KEY (gen_config_id) REFERENCES blc_site_map_gen_cfg(gen_config_id)
);
CREATE TABLE blc_customer (
	customer_id int8 NOT NULL,
	archived bpchar(1) NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	challenge_answer varchar(255) NULL,
	deactivated bool NULL,
	email_address varchar(255) NULL,
	external_id varchar(255) NULL,
	first_name varchar(255) NULL,
	is_tax_exempt bool NULL,
	last_name varchar(255) NULL,
	"password" varchar(255) NULL,
	password_change_required bool NULL,
	is_preview bool NULL,
	receive_email bool NULL,
	is_registered bool NULL,
	tax_exemption_code varchar(255) NULL,
	user_name varchar(255) NULL,
	challenge_question_id int8 NULL,
	locale_code varchar(255) NULL,
	CONSTRAINT blc_customer_pkey PRIMARY KEY (customer_id),
	CONSTRAINT fk4utjhbg9600iwhb05m40wspj1 FOREIGN KEY (locale_code) REFERENCES blc_locale(locale_code),
	CONSTRAINT fksgsex6rdheq2nm6pl23gggtqs FOREIGN KEY (challenge_question_id) REFERENCES blc_challenge_question(question_id)
);
CREATE INDEX customer_challenge_index ON public.blc_customer USING btree (challenge_question_id);
CREATE INDEX customer_email_index ON public.blc_customer USING btree (email_address);
CREATE TABLE blc_customer_address (
	customer_address_id int8 NOT NULL,
	address_name varchar(255) NULL,
	archived bpchar(1) NULL,
	address_id int8 NOT NULL,
	customer_id int8 NOT NULL,
	CONSTRAINT blc_customer_address_pkey PRIMARY KEY (customer_address_id),
	CONSTRAINT fkn79uhm41n1b23e6brajb4ggpw FOREIGN KEY (address_id) REFERENCES blc_address(address_id),
	CONSTRAINT fkrpdahw86mewf46g63nitq0w9p FOREIGN KEY (customer_id) REFERENCES blc_customer(customer_id)
);
CREATE INDEX customeraddress_address_index ON public.blc_customer_address USING btree (address_id);
CREATE TABLE blc_customer_attribute (
	customer_attr_id int8 NOT NULL,
	"name" varchar(255) NOT NULL,
	value varchar(255) NULL,
	customer_id int8 NOT NULL,
	CONSTRAINT blc_customer_attribute_pkey PRIMARY KEY (customer_attr_id),
	CONSTRAINT fko7j035lp80xu9wncbv96a1ry0 FOREIGN KEY (customer_id) REFERENCES blc_customer(customer_id)
);
CREATE TABLE blc_customer_offer_xref (
	customer_offer_id int8 NOT NULL,
	customer_id int8 NOT NULL,
	offer_id int8 NOT NULL,
	CONSTRAINT blc_customer_offer_xref_pkey PRIMARY KEY (customer_offer_id),
	CONSTRAINT fkg81dq5yxrtsy6cjivd0afkxcj FOREIGN KEY (customer_id) REFERENCES blc_customer(customer_id),
	CONSTRAINT fkrks1nkejqmm3n7y4xo5rs6wuk FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id)
);
CREATE INDEX custoffer_customer_index ON public.blc_customer_offer_xref USING btree (customer_id);
CREATE INDEX custoffer_offer_index ON public.blc_customer_offer_xref USING btree (offer_id);
CREATE TABLE blc_customer_payment (
	customer_payment_id int8 NOT NULL,
	is_default bool NULL,
	gateway_type varchar(255) NULL,
	payment_token varchar(255) NULL,
	payment_type varchar(255) NULL,
	address_id int8 NULL,
	customer_id int8 NOT NULL,
	CONSTRAINT blc_customer_payment_pkey PRIMARY KEY (customer_payment_id),
	CONSTRAINT cstmr_pay_unique_cnstrnt UNIQUE (customer_id, payment_token),
	CONSTRAINT fkhd53v4ilet9h8jxbjh7m2k7yj FOREIGN KEY (customer_id) REFERENCES blc_customer(customer_id),
	CONSTRAINT fkouuqjxsn30esr7ftg7i5mmr4p FOREIGN KEY (address_id) REFERENCES blc_address(address_id)
);
CREATE INDEX customerpayment_type_index ON public.blc_customer_payment USING btree (payment_type);
CREATE TABLE blc_customer_payment_fields (
	customer_payment_id int8 NOT NULL,
	field_value text NULL,
	field_name varchar(255) NOT NULL,
	CONSTRAINT blc_customer_payment_fields_pkey PRIMARY KEY (customer_payment_id, field_name),
	CONSTRAINT fkpwpbmvuo4pd8y76dswmq4cr00 FOREIGN KEY (customer_payment_id) REFERENCES blc_customer_payment(customer_payment_id)
);
CREATE TABLE blc_customer_phone (
	customer_phone_id int8 NOT NULL,
	phone_name varchar(255) NULL,
	customer_id int8 NOT NULL,
	phone_id int8 NOT NULL,
	CONSTRAINT blc_customer_phone_pkey PRIMARY KEY (customer_phone_id),
	CONSTRAINT cstmr_phone_unique_cnstrnt UNIQUE (customer_id, phone_name),
	CONSTRAINT fk1uy5spxqx6kxiqnvg5la7bjbb FOREIGN KEY (phone_id) REFERENCES blc_phone(phone_id),
	CONSTRAINT fk4sg479sv9t1dj7b1pso158tr8 FOREIGN KEY (customer_id) REFERENCES blc_customer(customer_id)
);
CREATE INDEX custphone_phone_index ON public.blc_customer_phone USING btree (phone_id);
CREATE TABLE blc_customer_role (
	customer_role_id int8 NOT NULL,
	customer_id int8 NOT NULL,
	role_id int8 NOT NULL,
	CONSTRAINT blc_customer_role_pkey PRIMARY KEY (customer_role_id),
	CONSTRAINT fkqcnikrg70t86oju6xf6622f5x FOREIGN KEY (role_id) REFERENCES blc_role(role_id),
	CONSTRAINT fksqxeay9un5o0l77mrtdgjxps4 FOREIGN KEY (customer_id) REFERENCES blc_customer(customer_id)
);
CREATE INDEX custrole_customer_index ON public.blc_customer_role USING btree (customer_id);
CREATE INDEX custrole_role_index ON public.blc_customer_role USING btree (role_id);
CREATE TABLE blc_order (
	order_id int8 NOT NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	email_address varchar(255) NULL,
	"name" varchar(255) NULL,
	order_number varchar(255) NULL,
	is_preview bool NULL,
	order_status varchar(255) NULL,
	order_subtotal numeric(19, 5) NULL,
	submit_date timestamp NULL,
	tax_override bool NULL,
	order_total numeric(19, 5) NULL,
	total_shipping numeric(19, 5) NULL,
	total_tax numeric(19, 5) NULL,
	currency_code varchar(255) NULL,
	customer_id int8 NOT NULL,
	locale_code varchar(255) NULL,
	CONSTRAINT blc_order_pkey PRIMARY KEY (order_id),
	CONSTRAINT fkc90jmu6i66weyh7o0u5q94jtv FOREIGN KEY (locale_code) REFERENCES blc_locale(locale_code),
	CONSTRAINT fkh0ouiaamkm2k7qfgc6cjacukg FOREIGN KEY (customer_id) REFERENCES blc_customer(customer_id),
	CONSTRAINT fkl58agohje8ndhoow8s8hlday1 FOREIGN KEY (currency_code) REFERENCES blc_currency(currency_code)
);
CREATE INDEX order_customer_index ON public.blc_order USING btree (customer_id);
CREATE INDEX order_email_index ON public.blc_order USING btree (email_address);
CREATE INDEX order_name_index ON public.blc_order USING btree (name);
CREATE INDEX order_number_index ON public.blc_order USING btree (order_number);
CREATE INDEX order_status_index ON public.blc_order USING btree (order_status);
CREATE TABLE blc_order_adjustment (
	order_adjustment_id int8 NOT NULL,
	is_future_credit bool NULL,
	adjustment_reason varchar(255) NOT NULL,
	adjustment_value numeric(19, 5) NOT NULL,
	offer_id int8 NOT NULL,
	order_id int8 NULL,
	CONSTRAINT blc_order_adjustment_pkey PRIMARY KEY (order_adjustment_id),
	CONSTRAINT fkh9agwlogxxgfxbxe7rcgrwv4u FOREIGN KEY (order_id) REFERENCES blc_order(order_id),
	CONSTRAINT fkmlymutb81ohtx11e2u64tjqyu FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id)
);
CREATE INDEX orderadjust_offer_index ON public.blc_order_adjustment USING btree (offer_id);
CREATE INDEX orderadjust_order_index ON public.blc_order_adjustment USING btree (order_id);
CREATE TABLE blc_order_attribute (
	order_attribute_id int8 NOT NULL,
	"name" varchar(255) NOT NULL,
	value varchar(255) NULL,
	order_id int8 NOT NULL,
	CONSTRAINT attr_name_order_id UNIQUE (name, order_id),
	CONSTRAINT blc_order_attribute_pkey PRIMARY KEY (order_attribute_id),
	CONSTRAINT fka5k0dl8lmasauj4cmems5e16s FOREIGN KEY (order_id) REFERENCES blc_order(order_id)
);
CREATE TABLE blc_order_offer_code_xref (
	order_id int8 NOT NULL,
	offer_code_id int8 NOT NULL,
	CONSTRAINT fkdtidg8l9a5wpcuwpnqbwwhuve FOREIGN KEY (offer_code_id) REFERENCES blc_offer_code(offer_code_id),
	CONSTRAINT fkljh9nh9lgxkgnebscn8u8sbgf FOREIGN KEY (order_id) REFERENCES blc_order(order_id)
);
CREATE TABLE blc_order_payment (
	order_payment_id int8 NOT NULL,
	amount numeric(19, 5) NULL,
	archived bpchar(1) NULL,
	gateway_type varchar(255) NULL,
	reference_number varchar(255) NULL,
	payment_type varchar(255) NOT NULL,
	address_id int8 NULL,
	order_id int8 NULL,
	CONSTRAINT blc_order_payment_pkey PRIMARY KEY (order_payment_id),
	CONSTRAINT fk7k9dsqtdku90rongw4f2xsgg5 FOREIGN KEY (address_id) REFERENCES blc_address(address_id),
	CONSTRAINT fkh0n8n75hk2l646hsxyyqrwgpx FOREIGN KEY (order_id) REFERENCES blc_order(order_id)
);
CREATE INDEX orderpayment_address_index ON public.blc_order_payment USING btree (address_id);
CREATE INDEX orderpayment_order_index ON public.blc_order_payment USING btree (order_id);
CREATE INDEX orderpayment_reference_index ON public.blc_order_payment USING btree (reference_number);
CREATE INDEX orderpayment_type_index ON public.blc_order_payment USING btree (payment_type);
CREATE TABLE blc_order_payment_transaction (
	payment_transaction_id int8 NOT NULL,
	transaction_amount numeric(19, 2) NULL,
	archived bpchar(1) NULL,
	customer_ip_address varchar(255) NULL,
	date_recorded timestamp NULL,
	raw_response text NULL,
	save_token bool NULL,
	success bool NULL,
	transaction_type varchar(255) NULL,
	order_payment int8 NOT NULL,
	parent_transaction int8 NULL,
	CONSTRAINT blc_order_payment_transaction_pkey PRIMARY KEY (payment_transaction_id),
	CONSTRAINT fkjg77vkh5u48omyy8uhagkswxs FOREIGN KEY (parent_transaction) REFERENCES blc_order_payment_transaction(payment_transaction_id),
	CONSTRAINT fkq3hymseoakriel7rw57g3vh5n FOREIGN KEY (order_payment) REFERENCES blc_order_payment(order_payment_id)
);
CREATE TABLE blc_page (
	page_id int8 NOT NULL,
	active_end_date timestamp NULL,
	active_start_date timestamp NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	description varchar(255) NULL,
	exclude_from_site_map bool NULL,
	full_url varchar(255) NULL,
	meta_description varchar(255) NULL,
	meta_title varchar(255) NULL,
	offline_flag bool NULL,
	priority int4 NULL,
	page_tmplt_id int8 NULL,
	CONSTRAINT blc_page_pkey PRIMARY KEY (page_id),
	CONSTRAINT fko95c7m41ycmhf4dwpebvemasl FOREIGN KEY (page_tmplt_id) REFERENCES blc_page_tmplt(page_tmplt_id)
);
CREATE INDEX page_full_url_index ON public.blc_page USING btree (full_url);
CREATE TABLE blc_page_attributes (
	attribute_id int8 NOT NULL,
	field_name varchar(255) NOT NULL,
	field_value varchar(255) NULL,
	page_id int8 NOT NULL,
	CONSTRAINT blc_page_attributes_pkey PRIMARY KEY (attribute_id),
	CONSTRAINT fk94785hg4iuw1k22qh6b8hysxe FOREIGN KEY (page_id) REFERENCES blc_page(page_id)
);
CREATE INDEX pageattribute_index ON public.blc_page_attributes USING btree (page_id);
CREATE INDEX pageattribute_name_index ON public.blc_page_attributes USING btree (field_name);
CREATE TABLE blc_page_fld (
	page_fld_id int8 NOT NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	fld_key varchar(255) NULL,
	lob_value text NULL,
	value varchar(255) NULL,
	page_id int8 NOT NULL,
	CONSTRAINT blc_page_fld_pkey PRIMARY KEY (page_fld_id),
	CONSTRAINT fk8t4im2p53x0mfufl92k87tsnx FOREIGN KEY (page_id) REFERENCES blc_page(page_id)
);
CREATE TABLE blc_page_rule_map (
	blc_page_page_id int8 NOT NULL,
	page_rule_id int8 NOT NULL,
	map_key varchar(255) NOT NULL,
	CONSTRAINT blc_page_rule_map_pkey PRIMARY KEY (blc_page_page_id, map_key),
	CONSTRAINT fkhj9uu6o7fb0n81g5yvk48skem FOREIGN KEY (page_rule_id) REFERENCES blc_page_rule(page_rule_id),
	CONSTRAINT fktqx8xsmgx0hkrvery3ipqwwi0 FOREIGN KEY (blc_page_page_id) REFERENCES blc_page(page_id)
);
CREATE TABLE blc_qual_crit_page_xref (
	page_id int8 NOT NULL,
	page_item_criteria_id int8 NOT NULL,
	CONSTRAINT blc_qual_crit_page_xref_pkey PRIMARY KEY (page_id, page_item_criteria_id),
	CONSTRAINT uk_dg6txhn3dw4k680q2sjyeumml UNIQUE (page_item_criteria_id),
	CONSTRAINT fkm0ov6kstmsab8gy93m53c05tg FOREIGN KEY (page_id) REFERENCES blc_page(page_id),
	CONSTRAINT fkpe7oenmm4t3g8ypvo5j2yjrd7 FOREIGN KEY (page_item_criteria_id) REFERENCES blc_page_item_criteria(page_item_criteria_id)
);
CREATE TABLE blc_rating_detail (
	rating_detail_id int8 NOT NULL,
	rating float8 NOT NULL,
	rating_submitted_date timestamp NOT NULL,
	customer_id int8 NOT NULL,
	rating_summary_id int8 NOT NULL,
	CONSTRAINT blc_rating_detail_pkey PRIMARY KEY (rating_detail_id),
	CONSTRAINT fkjjjou706ellmb65wmy7gpv07s FOREIGN KEY (customer_id) REFERENCES blc_customer(customer_id),
	CONSTRAINT fkorvii76lm0rnac92c67q1oles FOREIGN KEY (rating_summary_id) REFERENCES blc_rating_summary(rating_summary_id)
);
CREATE INDEX rating_customer_index ON public.blc_rating_detail USING btree (customer_id);
CREATE TABLE blc_review_detail (
	review_detail_id int8 NOT NULL,
	helpful_count int4 NOT NULL,
	not_helpful_count int4 NOT NULL,
	review_submitted_date timestamp NOT NULL,
	review_status varchar(255) NOT NULL,
	review_text varchar(255) NOT NULL,
	customer_id int8 NOT NULL,
	rating_detail_id int8 NULL,
	rating_summary_id int8 NOT NULL,
	CONSTRAINT blc_review_detail_pkey PRIMARY KEY (review_detail_id),
	CONSTRAINT fkdc0r2t22u3ghe0ihma8dcd5y6 FOREIGN KEY (customer_id) REFERENCES blc_customer(customer_id),
	CONSTRAINT fkhp1g51hyv3y8gr2tedxm0pgyl FOREIGN KEY (rating_detail_id) REFERENCES blc_rating_detail(rating_detail_id),
	CONSTRAINT fkn9hvs3m8fhodmipm3fvmwhw74 FOREIGN KEY (rating_summary_id) REFERENCES blc_rating_summary(rating_summary_id)
);
CREATE INDEX reviewdetail_customer_index ON public.blc_review_detail USING btree (customer_id);
CREATE INDEX reviewdetail_rating_index ON public.blc_review_detail USING btree (rating_detail_id);
CREATE INDEX reviewdetail_status_index ON public.blc_review_detail USING btree (review_status);
CREATE INDEX reviewdetail_summary_index ON public.blc_review_detail USING btree (rating_summary_id);
CREATE TABLE blc_review_feedback (
	review_feedback_id int8 NOT NULL,
	is_helpful bool NOT NULL,
	customer_id int8 NOT NULL,
	review_detail_id int8 NOT NULL,
	CONSTRAINT blc_review_feedback_pkey PRIMARY KEY (review_feedback_id),
	CONSTRAINT fkaledeh8wwn4ykopccqh3to8k5 FOREIGN KEY (review_detail_id) REFERENCES blc_review_detail(review_detail_id),
	CONSTRAINT fkmppbg4bf4h8v1m9efgm10ty4b FOREIGN KEY (customer_id) REFERENCES blc_customer(customer_id)
);
CREATE INDEX reviewfeed_customer_index ON public.blc_review_feedback USING btree (customer_id);
CREATE INDEX reviewfeed_detail_index ON public.blc_review_feedback USING btree (review_detail_id);
CREATE TABLE blc_sc (
	sc_id int8 NOT NULL,
	content_name varchar(255) NOT NULL,
	offline_flag bool NULL,
	priority int4 NOT NULL,
	locale_code varchar(255) NOT NULL,
	sc_type_id int8 NULL,
	CONSTRAINT blc_sc_pkey PRIMARY KEY (sc_id),
	CONSTRAINT fk13qnfvvq355l9cckfxkqqh59 FOREIGN KEY (locale_code) REFERENCES blc_locale(locale_code),
	CONSTRAINT fkp9be5g25yydwn151wnwvbj9hu FOREIGN KEY (sc_type_id) REFERENCES blc_sc_type(sc_type_id)
);
CREATE INDEX content_name_index_archived ON public.blc_sc USING btree (content_name, sc_type_id);
CREATE INDEX content_priority_index ON public.blc_sc USING btree (priority);
CREATE INDEX sc_offln_flg_indx ON public.blc_sc USING btree (offline_flag);
CREATE TABLE blc_sc_fld_map (
	blc_sc_sc_field_id int8 NOT NULL,
	map_key varchar(255) NOT NULL,
	sc_id int8 NOT NULL,
	sc_fld_id int8 NULL,
	CONSTRAINT blc_sc_fld_map_pkey PRIMARY KEY (blc_sc_sc_field_id),
	CONSTRAINT fkh54fvkukkun10yu69gr7neyj9 FOREIGN KEY (sc_id) REFERENCES blc_sc(sc_id),
	CONSTRAINT fkrwpb5a8l5uoeu4u046uihdx1g FOREIGN KEY (sc_fld_id) REFERENCES blc_sc_fld(sc_fld_id)
);
CREATE TABLE blc_sc_item_criteria (
	sc_item_criteria_id int8 NOT NULL,
	order_item_match_rule text NULL,
	quantity int4 NOT NULL,
	sc_id int8 NULL,
	CONSTRAINT blc_sc_item_criteria_pkey PRIMARY KEY (sc_item_criteria_id),
	CONSTRAINT fki62rdb9fuxn6lfdo7d4q9haow FOREIGN KEY (sc_id) REFERENCES blc_sc(sc_id)
);
CREATE TABLE blc_sc_rule_map (
	blc_sc_sc_id int8 NOT NULL,
	sc_rule_id int8 NOT NULL,
	map_key varchar(255) NOT NULL,
	CONSTRAINT blc_sc_rule_map_pkey PRIMARY KEY (blc_sc_sc_id, map_key),
	CONSTRAINT fk31d3qpemphv6qdbha0cl1x361 FOREIGN KEY (sc_rule_id) REFERENCES blc_sc_rule(sc_rule_id),
	CONSTRAINT fko4q8t9hx8iprk9bc9tllwhdhk FOREIGN KEY (blc_sc_sc_id) REFERENCES blc_sc(sc_id)
);
CREATE TABLE blc_site_map_url_entry (
	url_entry_id int8 NOT NULL,
	change_freq varchar(255) NOT NULL,
	last_modified timestamp NOT NULL,
	"location" varchar(255) NOT NULL,
	priority varchar(255) NOT NULL,
	gen_config_id int8 NOT NULL,
	CONSTRAINT blc_site_map_url_entry_pkey PRIMARY KEY (url_entry_id),
	CONSTRAINT fkrvkeinfysjshg9ulmxno9rhla FOREIGN KEY (gen_config_id) REFERENCES blc_cust_site_map_gen_cfg(gen_config_id)
);
CREATE TABLE blc_trans_additnl_fields (
	payment_transaction_id int8 NOT NULL,
	field_value text NULL,
	field_name varchar(255) NOT NULL,
	CONSTRAINT blc_trans_additnl_fields_pkey PRIMARY KEY (payment_transaction_id, field_name),
	CONSTRAINT fkdmq1toto9pwrhw5uife2ssq45 FOREIGN KEY (payment_transaction_id) REFERENCES blc_order_payment_transaction(payment_transaction_id)
);
CREATE TABLE blc_additional_offer_info (
	blc_order_order_id int8 NOT NULL,
	offer_info_id int8 NOT NULL,
	offer_id int8 NOT NULL,
	CONSTRAINT blc_additional_offer_info_pkey PRIMARY KEY (blc_order_order_id, offer_id),
	CONSTRAINT fk40nm1ylfeiv2t6ojcmehu8gcr FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id),
	CONSTRAINT fklkk2kdjpv0v0ybxnc11p7hg4e FOREIGN KEY (offer_info_id) REFERENCES blc_offer_info(offer_info_id),
	CONSTRAINT fkrfc8a02u7yp8qqk206ug62lnb FOREIGN KEY (blc_order_order_id) REFERENCES blc_order(order_id)
);
CREATE TABLE blc_candidate_order_offer (
	candidate_order_offer_id int8 NOT NULL,
	discounted_price numeric(19, 5) NULL,
	offer_id int8 NOT NULL,
	order_id int8 NULL,
	CONSTRAINT blc_candidate_order_offer_pkey PRIMARY KEY (candidate_order_offer_id),
	CONSTRAINT fk59se4s0394sw56c1rvdw5aepf FOREIGN KEY (order_id) REFERENCES blc_order(order_id),
	CONSTRAINT fk7ckpnvor07qankv258p1vxwww FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id)
);
CREATE INDEX candidate_order_index ON public.blc_candidate_order_offer USING btree (order_id);
CREATE INDEX candidate_orderoffer_index ON public.blc_candidate_order_offer USING btree (offer_id);
CREATE TABLE blc_cms_menu_item (
	menu_item_id int8 NOT NULL,
	action_url varchar(255) NULL,
	alt_text varchar(255) NULL,
	custom_html text NULL,
	image_url varchar(255) NULL,
	"label" varchar(255) NULL,
	"sequence" numeric(10, 6) NULL,
	menu_item_type varchar(255) NULL,
	linked_menu_id int8 NULL,
	linked_page_id int8 NULL,
	parent_menu_id int8 NULL,
	CONSTRAINT blc_cms_menu_item_pkey PRIMARY KEY (menu_item_id),
	CONSTRAINT fka7ogt4huutal1mirsufnmy9lr FOREIGN KEY (parent_menu_id) REFERENCES blc_cms_menu(menu_id),
	CONSTRAINT fkbgy5higr7beta0sxdqvkm9k7r FOREIGN KEY (linked_page_id) REFERENCES blc_page(page_id),
	CONSTRAINT fksfd7p9istk4908bchapktbnr0 FOREIGN KEY (linked_menu_id) REFERENCES blc_cms_menu(menu_id)
);
CREATE TABLE blc_fulfillment_group (
	fulfillment_group_id int8 NOT NULL,
	delivery_instruction varchar(255) NULL,
	price numeric(19, 5) NULL,
	shipping_price_taxable bool NULL,
	merchandise_total numeric(19, 5) NULL,
	"method" varchar(255) NULL,
	is_primary bool NULL,
	reference_number varchar(255) NULL,
	retail_price numeric(19, 5) NULL,
	sale_price numeric(19, 5) NULL,
	fulfillment_group_sequnce int4 NULL,
	service varchar(255) NULL,
	shipping_override bool NULL,
	status varchar(255) NULL,
	total numeric(19, 5) NULL,
	total_fee_tax numeric(19, 5) NULL,
	total_fg_tax numeric(19, 5) NULL,
	total_item_tax numeric(19, 5) NULL,
	total_tax numeric(19, 5) NULL,
	"type" varchar(255) NULL,
	address_id int8 NULL,
	fulfillment_option_id int8 NULL,
	order_id int8 NOT NULL,
	personal_message_id int8 NULL,
	phone_id int8 NULL,
	CONSTRAINT blc_fulfillment_group_pkey PRIMARY KEY (fulfillment_group_id),
	CONSTRAINT fk336lsxych2j78fsd12dxacn7n FOREIGN KEY (fulfillment_option_id) REFERENCES blc_fulfillment_option(fulfillment_option_id),
	CONSTRAINT fk3lis2v456tgcmagt1tkdummdi FOREIGN KEY (phone_id) REFERENCES blc_phone(phone_id),
	CONSTRAINT fk44mielsxkxtt1ndfiat2wj9po FOREIGN KEY (personal_message_id) REFERENCES blc_personal_message(personal_message_id),
	CONSTRAINT fkavpobeg9yjr9k3wtycirv5i8a FOREIGN KEY (address_id) REFERENCES blc_address(address_id),
	CONSTRAINT fkbtadc11h6ysb0fbyq2bsegum7 FOREIGN KEY (order_id) REFERENCES blc_order(order_id)
);
CREATE INDEX fg_address_index ON public.blc_fulfillment_group USING btree (address_id);
CREATE INDEX fg_message_index ON public.blc_fulfillment_group USING btree (personal_message_id);
CREATE INDEX fg_method_index ON public.blc_fulfillment_group USING btree (method);
CREATE INDEX fg_order_index ON public.blc_fulfillment_group USING btree (order_id);
CREATE INDEX fg_phone_index ON public.blc_fulfillment_group USING btree (phone_id);
CREATE INDEX fg_primary_index ON public.blc_fulfillment_group USING btree (is_primary);
CREATE INDEX fg_reference_index ON public.blc_fulfillment_group USING btree (reference_number);
CREATE INDEX fg_service_index ON public.blc_fulfillment_group USING btree (service);
CREATE INDEX fg_status_index ON public.blc_fulfillment_group USING btree (status);
CREATE TABLE blc_fulfillment_group_fee (
	fulfillment_group_fee_id int8 NOT NULL,
	amount numeric(19, 5) NULL,
	fee_taxable_flag bool NULL,
	"name" varchar(255) NULL,
	reporting_code varchar(255) NULL,
	total_fee_tax numeric(19, 5) NULL,
	fulfillment_group_id int8 NOT NULL,
	CONSTRAINT blc_fulfillment_group_fee_pkey PRIMARY KEY (fulfillment_group_fee_id),
	CONSTRAINT fkss79152pprx7xdwjkmelwf4xo FOREIGN KEY (fulfillment_group_id) REFERENCES blc_fulfillment_group(fulfillment_group_id)
);
CREATE TABLE blc_qual_crit_sc_xref (
	sc_id int8 NOT NULL,
	sc_item_criteria_id int8 NOT NULL,
	CONSTRAINT blc_qual_crit_sc_xref_pkey PRIMARY KEY (sc_id, sc_item_criteria_id),
	CONSTRAINT uk_afqd4tvahqdouwkfb55xjuycm UNIQUE (sc_item_criteria_id),
	CONSTRAINT fk6v9jfn06vkk5kpio9jdu08t30 FOREIGN KEY (sc_item_criteria_id) REFERENCES blc_sc_item_criteria(sc_item_criteria_id),
	CONSTRAINT fkq0wnd7j8o8ss4umkpdjn38ota FOREIGN KEY (sc_id) REFERENCES blc_sc(sc_id)
);
CREATE TABLE blc_candidate_fg_offer (
	candidate_fg_offer_id int8 NOT NULL,
	discounted_price numeric(19, 5) NULL,
	fulfillment_group_id int8 NULL,
	offer_id int8 NOT NULL,
	CONSTRAINT blc_candidate_fg_offer_pkey PRIMARY KEY (candidate_fg_offer_id),
	CONSTRAINT fkg5qmns7vl5e1pu96axl8uknal FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id),
	CONSTRAINT fkh9csft0rxeopd0s4in7qp15am FOREIGN KEY (fulfillment_group_id) REFERENCES blc_fulfillment_group(fulfillment_group_id)
);
CREATE INDEX candidate_fg_index ON public.blc_candidate_fg_offer USING btree (fulfillment_group_id);
CREATE INDEX candidate_fgoffer_index ON public.blc_candidate_fg_offer USING btree (offer_id);
CREATE TABLE blc_fg_adjustment (
	fg_adjustment_id int8 NOT NULL,
	is_future_credit bool NULL,
	adjustment_reason varchar(255) NOT NULL,
	adjustment_value numeric(19, 5) NOT NULL,
	fulfillment_group_id int8 NULL,
	offer_id int8 NOT NULL,
	CONSTRAINT blc_fg_adjustment_pkey PRIMARY KEY (fg_adjustment_id),
	CONSTRAINT fk2ceuqqy88te84f61f0n7kvaw1 FOREIGN KEY (fulfillment_group_id) REFERENCES blc_fulfillment_group(fulfillment_group_id),
	CONSTRAINT fkt0l1mgyccsuq76n8b0b6pc9a9 FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id)
);
CREATE INDEX fgadjustment_index ON public.blc_fg_adjustment USING btree (fulfillment_group_id);
CREATE INDEX fgadjustment_offer_index ON public.blc_fg_adjustment USING btree (offer_id);
CREATE TABLE blc_fg_fee_tax_xref (
	fulfillment_group_fee_id int8 NOT NULL,
	tax_detail_id int8 NOT NULL,
	CONSTRAINT uk_59ow3plvbkxjfs57k92ahf3eg UNIQUE (tax_detail_id),
	CONSTRAINT fk1aueplsngm018mlqqq9yhgrn6 FOREIGN KEY (tax_detail_id) REFERENCES blc_tax_detail(tax_detail_id),
	CONSTRAINT fk2t3oa9322dqgya6r27pb2bcsd FOREIGN KEY (fulfillment_group_fee_id) REFERENCES blc_fulfillment_group_fee(fulfillment_group_fee_id)
);
CREATE TABLE blc_fg_fg_tax_xref (
	fulfillment_group_id int8 NOT NULL,
	tax_detail_id int8 NOT NULL,
	CONSTRAINT uk_57834q276cjrrnwjj1ilnj6ve UNIQUE (tax_detail_id),
	CONSTRAINT fkla7cgvy244irmood3xt8rpsjb FOREIGN KEY (tax_detail_id) REFERENCES blc_tax_detail(tax_detail_id),
	CONSTRAINT fknah3gdurbtogb0s9sf3humt14 FOREIGN KEY (fulfillment_group_id) REFERENCES blc_fulfillment_group(fulfillment_group_id)
);
CREATE TABLE blc_bund_item_fee_price (
	bund_item_fee_price_id int8 NOT NULL,
	amount numeric(19, 5) NULL,
	is_taxable bool NULL,
	"name" varchar(255) NULL,
	reporting_code varchar(255) NULL,
	bund_order_item_id int8 NOT NULL,
	CONSTRAINT blc_bund_item_fee_price_pkey PRIMARY KEY (bund_item_fee_price_id)
);
CREATE TABLE blc_bundle_order_item (
	base_retail_price numeric(19, 5) NULL,
	base_sale_price numeric(19, 5) NULL,
	order_item_id int8 NOT NULL,
	product_bundle_id int8 NULL,
	sku_id int8 NULL,
	CONSTRAINT blc_bundle_order_item_pkey PRIMARY KEY (order_item_id)
);
CREATE TABLE blc_candidate_item_offer (
	candidate_item_offer_id int8 NOT NULL,
	discounted_price numeric(19, 5) NULL,
	offer_id int8 NOT NULL,
	order_item_id int8 NULL,
	CONSTRAINT blc_candidate_item_offer_pkey PRIMARY KEY (candidate_item_offer_id)
);
CREATE INDEX candidate_item_index ON public.blc_candidate_item_offer USING btree (order_item_id);
CREATE INDEX candidate_itemoffer_index ON public.blc_candidate_item_offer USING btree (offer_id);
CREATE TABLE blc_disc_item_fee_price (
	disc_item_fee_price_id int8 NOT NULL,
	amount numeric(19, 5) NULL,
	"name" varchar(255) NULL,
	reporting_code varchar(255) NULL,
	order_item_id int8 NOT NULL,
	CONSTRAINT blc_disc_item_fee_price_pkey PRIMARY KEY (disc_item_fee_price_id)
);
CREATE TABLE blc_discrete_order_item (
	base_retail_price numeric(19, 5) NULL,
	base_sale_price numeric(19, 5) NULL,
	order_item_id int8 NOT NULL,
	bundle_order_item_id int8 NULL,
	product_id int8 NULL,
	sku_id int8 NOT NULL,
	sku_bundle_item_id int8 NULL,
	CONSTRAINT blc_discrete_order_item_pkey PRIMARY KEY (order_item_id)
);
CREATE INDEX discrete_product_index ON public.blc_discrete_order_item USING btree (product_id);
CREATE INDEX discrete_sku_index ON public.blc_discrete_order_item USING btree (sku_id);
CREATE TABLE blc_dyn_discrete_order_item (
	order_item_id int8 NOT NULL,
	CONSTRAINT blc_dyn_discrete_order_item_pkey PRIMARY KEY (order_item_id)
);
CREATE TABLE blc_fg_item_tax_xref (
	fulfillment_group_item_id int8 NOT NULL,
	tax_detail_id int8 NOT NULL,
	CONSTRAINT uk_hs9yvwvlwdy668hf186rgfyvq UNIQUE (tax_detail_id)
);
CREATE TABLE blc_fulfillment_group_item (
	fulfillment_group_item_id int8 NOT NULL,
	prorated_order_adj numeric(19, 2) NULL,
	quantity int4 NOT NULL,
	status varchar(255) NULL,
	total_item_amount numeric(19, 5) NULL,
	total_item_taxable_amount numeric(19, 5) NULL,
	total_item_tax numeric(19, 5) NULL,
	fulfillment_group_id int8 NOT NULL,
	order_item_id int8 NOT NULL,
	CONSTRAINT blc_fulfillment_group_item_pkey PRIMARY KEY (fulfillment_group_item_id)
);
CREATE INDEX fgitem_fg_index ON public.blc_fulfillment_group_item USING btree (fulfillment_group_id);
CREATE INDEX fgitem_status_index ON public.blc_fulfillment_group_item USING btree (status);
CREATE TABLE blc_giftwrap_order_item (
	order_item_id int8 NOT NULL,
	CONSTRAINT blc_giftwrap_order_item_pkey PRIMARY KEY (order_item_id)
);
CREATE TABLE blc_item_offer_qualifier (
	item_offer_qualifier_id int8 NOT NULL,
	quantity int8 NULL,
	offer_id int8 NOT NULL,
	order_item_id int8 NULL,
	CONSTRAINT blc_item_offer_qualifier_pkey PRIMARY KEY (item_offer_qualifier_id)
);
CREATE TABLE blc_order_item (
	order_item_id int8 NOT NULL,
	created_by int8 NULL,
	date_created timestamp NULL,
	date_updated timestamp NULL,
	updated_by int8 NULL,
	discounts_allowed bool NULL,
	has_validation_errors bool NULL,
	item_taxable_flag bool NULL,
	"name" varchar(255) NULL,
	order_item_type varchar(255) NULL,
	price numeric(19, 5) NULL,
	quantity int4 NOT NULL,
	retail_price numeric(19, 5) NULL,
	retail_price_override bool NULL,
	sale_price numeric(19, 5) NULL,
	sale_price_override bool NULL,
	total_tax numeric(19, 2) NULL,
	category_id int8 NULL,
	gift_wrap_item_id int8 NULL,
	order_id int8 NULL,
	parent_order_item_id int8 NULL,
	personal_message_id int8 NULL,
	CONSTRAINT blc_order_item_pkey PRIMARY KEY (order_item_id)
);
CREATE INDEX orderitem_category_index ON public.blc_order_item USING btree (category_id);
CREATE INDEX orderitem_gift_index ON public.blc_order_item USING btree (gift_wrap_item_id);
CREATE INDEX orderitem_message_index ON public.blc_order_item USING btree (personal_message_id);
CREATE INDEX orderitem_order_index ON public.blc_order_item USING btree (order_id);
CREATE INDEX orderitem_parent_index ON public.blc_order_item USING btree (parent_order_item_id);
CREATE INDEX orderitem_type_index ON public.blc_order_item USING btree (order_item_type);
CREATE TABLE blc_order_item_add_attr (
	order_item_id int8 NOT NULL,
	value varchar(255) NULL,
	"name" varchar(255) NOT NULL,
	CONSTRAINT blc_order_item_add_attr_pkey PRIMARY KEY (order_item_id, name)
);
CREATE TABLE blc_order_item_adjustment (
	order_item_adjustment_id int8 NOT NULL,
	applied_to_sale_price bool NULL,
	adjustment_reason varchar(255) NOT NULL,
	adjustment_value numeric(19, 5) NOT NULL,
	offer_id int8 NOT NULL,
	order_item_id int8 NULL,
	CONSTRAINT blc_order_item_adjustment_pkey PRIMARY KEY (order_item_adjustment_id)
);
CREATE INDEX oiadjust_item_index ON public.blc_order_item_adjustment USING btree (order_item_id);
CREATE TABLE blc_order_item_attribute (
	order_item_attribute_id int8 NOT NULL,
	"name" varchar(255) NOT NULL,
	value varchar(255) NOT NULL,
	order_item_id int8 NOT NULL,
	CONSTRAINT attr_name_order_item_id UNIQUE (name, order_item_id),
	CONSTRAINT blc_order_item_attribute_pkey PRIMARY KEY (order_item_attribute_id)
);
CREATE TABLE blc_order_item_cart_message (
	order_item_id int8 NOT NULL,
	cart_message varchar(255) NULL
);
CREATE TABLE blc_order_item_dtl_adj (
	order_item_dtl_adj_id int8 NOT NULL,
	applied_to_sale_price bool NULL,
	is_future_credit bool NULL,
	offer_name varchar(255) NULL,
	adjustment_reason varchar(255) NOT NULL,
	adjustment_value numeric(19, 5) NOT NULL,
	offer_id int8 NOT NULL,
	order_item_price_dtl_id int8 NULL,
	CONSTRAINT blc_order_item_dtl_adj_pkey PRIMARY KEY (order_item_dtl_adj_id)
);
CREATE TABLE blc_order_item_price_dtl (
	order_item_price_dtl_id int8 NOT NULL,
	quantity int4 NOT NULL,
	use_sale_price bool NULL,
	order_item_id int8 NULL,
	CONSTRAINT blc_order_item_price_dtl_pkey PRIMARY KEY (order_item_price_dtl_id)
);
CREATE TABLE blc_order_multiship_option (
	order_multiship_option_id int8 NOT NULL,
	address_id int8 NULL,
	fulfillment_option_id int8 NULL,
	order_id int8 NULL,
	order_item_id int8 NULL,
	CONSTRAINT blc_order_multiship_option_pkey PRIMARY KEY (order_multiship_option_id)
);
CREATE INDEX multiship_option_order_index ON public.blc_order_multiship_option USING btree (order_id);
CREATE TABLE blc_prorated_order_item_adjust (
	prorated_order_item_adjust_id int8 NOT NULL,
	prorated_quantity int4 NOT NULL,
	adjustment_reason varchar(255) NOT NULL,
	prorated_adjustment_value numeric(19, 5) NOT NULL,
	offer_id int8 NOT NULL,
	order_item_id int8 NULL,
	CONSTRAINT blc_prorated_order_item_adjust_pkey PRIMARY KEY (prorated_order_item_adjust_id)
);
CREATE INDEX poiadjust_item_index ON public.blc_prorated_order_item_adjust USING btree (order_item_id);
ALTER TABLE public.blc_bund_item_fee_price ADD CONSTRAINT fkmlwh6qvntrxs81h26syrft6bj FOREIGN KEY (bund_order_item_id) REFERENCES blc_bundle_order_item(order_item_id);
ALTER TABLE public.blc_bundle_order_item ADD CONSTRAINT fkepomqf1wy0rsw6utuc5r7gflq FOREIGN KEY (sku_id) REFERENCES blc_sku(sku_id);
ALTER TABLE public.blc_bundle_order_item ADD CONSTRAINT fkhbcblyyh5lfrmrt1avy8wajgx FOREIGN KEY (order_item_id) REFERENCES blc_order_item(order_item_id);
ALTER TABLE public.blc_bundle_order_item ADD CONSTRAINT fko17ovssehxe4y3b38xjxodmrq FOREIGN KEY (product_bundle_id) REFERENCES blc_product_bundle(product_id);
ALTER TABLE public.blc_candidate_item_offer ADD CONSTRAINT fkb14jq3w7049s1h61pthy51m92 FOREIGN KEY (order_item_id) REFERENCES blc_order_item(order_item_id);
ALTER TABLE public.blc_candidate_item_offer ADD CONSTRAINT fkno8jmqw67ef9lpuwoumdyh7wj FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id);
ALTER TABLE public.blc_disc_item_fee_price ADD CONSTRAINT fk70ocmaswx7p3xymfvildubx5 FOREIGN KEY (order_item_id) REFERENCES blc_discrete_order_item(order_item_id);
ALTER TABLE public.blc_discrete_order_item ADD CONSTRAINT fk188b985egh16qfcjt8kv1asa4 FOREIGN KEY (order_item_id) REFERENCES blc_order_item(order_item_id);
ALTER TABLE public.blc_discrete_order_item ADD CONSTRAINT fk1micyx881c06d24amsg3sk2he FOREIGN KEY (sku_bundle_item_id) REFERENCES blc_sku_bundle_item(sku_bundle_item_id);
ALTER TABLE public.blc_discrete_order_item ADD CONSTRAINT fk2moe4tjwke365lo2s5qgmacx7 FOREIGN KEY (product_id) REFERENCES blc_product(product_id);
ALTER TABLE public.blc_discrete_order_item ADD CONSTRAINT fkmtcs7ax8jo2hy1ae4caafsfkp FOREIGN KEY (sku_id) REFERENCES blc_sku(sku_id);
ALTER TABLE public.blc_discrete_order_item ADD CONSTRAINT fkpu94j8xpk9uwpcgcy98ktle06 FOREIGN KEY (bundle_order_item_id) REFERENCES blc_bundle_order_item(order_item_id);
ALTER TABLE public.blc_dyn_discrete_order_item ADD CONSTRAINT fkhv263skp3pgb4wcxg44umwcjs FOREIGN KEY (order_item_id) REFERENCES blc_discrete_order_item(order_item_id);
ALTER TABLE public.blc_fg_item_tax_xref ADD CONSTRAINT fkb5rnxtly8pr3ihvlrxlovnjkb FOREIGN KEY (tax_detail_id) REFERENCES blc_tax_detail(tax_detail_id);
ALTER TABLE public.blc_fg_item_tax_xref ADD CONSTRAINT fkl5kovj2ayfp7idroml0qjwan3 FOREIGN KEY (fulfillment_group_item_id) REFERENCES blc_fulfillment_group_item(fulfillment_group_item_id);
ALTER TABLE public.blc_fulfillment_group_item ADD CONSTRAINT fkmra6tj092ugw58xhvvi43pdb2 FOREIGN KEY (order_item_id) REFERENCES blc_order_item(order_item_id);
ALTER TABLE public.blc_fulfillment_group_item ADD CONSTRAINT fkqfqxv2f0ita9ou48jpi7c3wi9 FOREIGN KEY (fulfillment_group_id) REFERENCES blc_fulfillment_group(fulfillment_group_id);
ALTER TABLE public.blc_giftwrap_order_item ADD CONSTRAINT fktq6vr571td9a8ihss8os1wtr8 FOREIGN KEY (order_item_id) REFERENCES blc_discrete_order_item(order_item_id);
ALTER TABLE public.blc_item_offer_qualifier ADD CONSTRAINT fk9fl5rced4g4u8sxh1j4mrwkto FOREIGN KEY (order_item_id) REFERENCES blc_order_item(order_item_id);
ALTER TABLE public.blc_item_offer_qualifier ADD CONSTRAINT fko9i9n1thqcqt9nu0fv2nlg1ec FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id);
ALTER TABLE public.blc_order_item ADD CONSTRAINT fk3553qqcmvw5i3durebksttod3 FOREIGN KEY (gift_wrap_item_id) REFERENCES blc_giftwrap_order_item(order_item_id);
ALTER TABLE public.blc_order_item ADD CONSTRAINT fk4vocoseu9tnln1vq4r2gygh3n FOREIGN KEY (parent_order_item_id) REFERENCES blc_order_item(order_item_id);
ALTER TABLE public.blc_order_item ADD CONSTRAINT fk737vx8aceqsa8wyb6hjt44x58 FOREIGN KEY (category_id) REFERENCES blc_category(category_id);
ALTER TABLE public.blc_order_item ADD CONSTRAINT fk8b71a8di9bu8jrssp98u8ka0s FOREIGN KEY (order_id) REFERENCES blc_order(order_id);
ALTER TABLE public.blc_order_item ADD CONSTRAINT fkccrkxx60l5x2q24dl97x9iu0a FOREIGN KEY (personal_message_id) REFERENCES blc_personal_message(personal_message_id);
ALTER TABLE public.blc_order_item_add_attr ADD CONSTRAINT fk2xfsv1rmg5hy926njxrgrkxja FOREIGN KEY (order_item_id) REFERENCES blc_discrete_order_item(order_item_id);
ALTER TABLE public.blc_order_item_adjustment ADD CONSTRAINT fkkw991n1so1bd0nvmtgky3d4rm FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id);
ALTER TABLE public.blc_order_item_adjustment ADD CONSTRAINT fkniw5eryl2ea895x5p3wh4sd0u FOREIGN KEY (order_item_id) REFERENCES blc_order_item(order_item_id);
ALTER TABLE public.blc_order_item_attribute ADD CONSTRAINT fk5f2l8atn9sh06yhbjx72i8tl1 FOREIGN KEY (order_item_id) REFERENCES blc_order_item(order_item_id);
ALTER TABLE public.blc_order_item_cart_message ADD CONSTRAINT fkpm9ip11x3m6rnkla1khhgqmfn FOREIGN KEY (order_item_id) REFERENCES blc_order_item(order_item_id);
ALTER TABLE public.blc_order_item_dtl_adj ADD CONSTRAINT fktaukfbw7rtre7kmvl6fla49t6 FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id);
ALTER TABLE public.blc_order_item_dtl_adj ADD CONSTRAINT fktecvjagoubp6v2337mqm5gnmg FOREIGN KEY (order_item_price_dtl_id) REFERENCES blc_order_item_price_dtl(order_item_price_dtl_id);
ALTER TABLE public.blc_order_item_price_dtl ADD CONSTRAINT fkerm8r2c1fj0vvd2rj0uxmavlj FOREIGN KEY (order_item_id) REFERENCES blc_order_item(order_item_id);
ALTER TABLE public.blc_order_multiship_option ADD CONSTRAINT fk1d1sd1fr2cdv0kvf2s5yclo55 FOREIGN KEY (fulfillment_option_id) REFERENCES blc_fulfillment_option(fulfillment_option_id);
ALTER TABLE public.blc_order_multiship_option ADD CONSTRAINT fk8poefpppeoej296cr2g1otfki FOREIGN KEY (order_item_id) REFERENCES blc_order_item(order_item_id);
ALTER TABLE public.blc_order_multiship_option ADD CONSTRAINT fkrpsf6ltf21ohrgimnktmq5dn3 FOREIGN KEY (order_id) REFERENCES blc_order(order_id);
ALTER TABLE public.blc_order_multiship_option ADD CONSTRAINT fkt77nf9y3nokcclqibjhihjily FOREIGN KEY (address_id) REFERENCES blc_address(address_id);
ALTER TABLE public.blc_prorated_order_item_adjust ADD CONSTRAINT fkfmq7jd0v7r11g8hlvyuw50u6q FOREIGN KEY (order_item_id) REFERENCES blc_order_item(order_item_id);
ALTER TABLE public.blc_prorated_order_item_adjust ADD CONSTRAINT fkkud7s4f923plknu6u4041v3oa FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id);
