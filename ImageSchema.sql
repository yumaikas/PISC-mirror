-- This is the schemda for a PISC image


Create Table vocabulary (
	vocab_id INTEGER PRIMARY KEY,
	vocab_name text NOT NULL,
	-- TODO: figure out dependency tracking
);

Create Table word_definition (
	word_id INTEGER PRIMARY KEY,
	word_name text NOT NULL,
	stack_comment text NOT NULL,
	definition text NOT NULL,
	vocab_id INTEGER NOT NULL,
	created_date INTEGER NULL,
	modified_date INTEGER NULL,
	FOREIGN KEY(vocab_id) REFERENCES vocabulary(vocab_id)
);

Create Table word_documentation (
	doc_id INTEGER PRIMARY KEY,
	word_id INTEGER NOT NULL,
	word_documentation TEXT NULL,
	FOREIGN KEY(word_id) REFERENCES vocabulary(word_id)
);

Create Table kv_store (
);