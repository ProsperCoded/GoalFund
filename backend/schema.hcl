// Atlas schema definition for GoFund

schema "public" {
  comment = "GoFund database schema"
}

// Enums
enum "user_role" {
  schema = schema.public
  values = ["user", "admin"]
}

// Users Service Tables
table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "email" {
    null = false
    type = varchar(255)
  }
  column "username" {
    null = false
    type = varchar(100)
  }
  column "password_hash" {
    null = false
    type = varchar(255)
  }
  column "first_name" {
    null = true
    type = varchar(100)
  }
  column "last_name" {
    null = true
    type = varchar(100)
  }
  column "phone" {
    null = true
    type = varchar(20)
  }
  column "email_verified" {
    null = false
    type = boolean
    default = false
  }
  column "phone_verified" {
    null = false
    type = boolean
    default = false
  }
  column "role" {
    null = false
    type = enum("user_role")
    default = "user"
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_users_email" {
    unique = true
    columns = [column.email]
  }
  
  index "idx_users_username" {
    unique = true
    columns = [column.username]
  }
}


table "sessions" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "user_id" {
    null = false
    type = uuid
  }
  column "token_hash" {
    null = false
    type = varchar(255)
  }
  column "expires_at" {
    null = false
    type = timestamptz
  }
  column "metadata" {
    null = true
    type = jsonb
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_sessions_user_id" {
    columns = [column.user_id]
  }
  
  index "idx_sessions_token_hash" {
    unique = true
    columns = [column.token_hash]
  }
  
  foreign_key "fk_sessions_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }
}

table "password_reset_tokens" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "user_id" {
    null = false
    type = uuid
  }
  column "token_hash" {
    null = false
    type = varchar(255)
  }
  column "expires_at" {
    null = false
    type = timestamptz
  }
  column "used" {
    null = false
    type = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_password_reset_tokens_user_id" {
    columns = [column.user_id]
  }
  
  index "idx_password_reset_tokens_token_hash" {
    unique = true
    columns = [column.token_hash]
  }
  
  foreign_key "fk_password_reset_tokens_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }
}

// Goals Service Tables
table "goals" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "owner_id" {
    null = false
    type = uuid
  }
  column "title" {
    null = false
    type = varchar(255)
  }
  column "description" {
    null = true
    type = text
  }
  column "target_amount" {
    null = false
    type = bigint
  }
  column "currency" {
    null = false
    type = varchar(3)
    default = "NGN"
  }
  column "deadline" {
    null = true
    type = timestamptz
  }
  column "status" {
    null = false
    type = varchar(20)
    default = "OPEN"
  }
  column "is_public" {
    null = false
    type = boolean
    default = true
  }
  column "deposit_bank_name" {
    null = true
    type = varchar(100)
  }
  column "deposit_account_number" {
    null = true
    type = varchar(20)
  }
  column "deposit_account_name" {
    null = true
    type = varchar(255)
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_goals_owner_id" {
    columns = [column.owner_id]
  }
  
  index "idx_goals_status" {
    columns = [column.status]
  }
}

table "contributions" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "goal_id" {
    null = false
    type = uuid
  }
  column "user_id" {
    null = false
    type = uuid
  }
  column "payment_id" {
    null = true
    type = uuid
  }
  column "amount" {
    null = false
    type = bigint
  }
  column "currency" {
    null = false
    type = varchar(3)
    default = "NGN"
  }
  column "status" {
    null = false
    type = varchar(20)
    default = "PENDING"
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_contributions_goal_id" {
    columns = [column.goal_id]
  }
  
  index "idx_contributions_user_id" {
    columns = [column.user_id]
  }
  
  index "idx_contributions_payment_id" {
    columns = [column.payment_id]
  }
  
  foreign_key "fk_contributions_goal" {
    columns     = [column.goal_id]
    ref_columns = [table.goals.column.id]
    on_delete   = CASCADE
  }
}

table "proofs" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "goal_id" {
    null = false
    type = uuid
  }
  column "submitted_by" {
    null = false
    type = uuid
  }
  column "title" {
    null = false
    type = varchar(255)
  }
  column "description" {
    null = true
    type = text
  }
  column "media_urls" {
    null = true
    type = jsonb
  }
  column "status" {
    null = false
    type = varchar(20)
    default = "PENDING"
  }
  column "submitted_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "verified_at" {
    null = true
    type = timestamptz
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_proofs_goal_id" {
    columns = [column.goal_id]
  }
  
  index "idx_proofs_submitted_by" {
    columns = [column.submitted_by]
  }
  
  foreign_key "fk_proofs_goal" {
    columns     = [column.goal_id]
    ref_columns = [table.goals.column.id]
    on_delete   = CASCADE
  }
}

table "votes" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "proof_id" {
    null = false
    type = uuid
  }
  column "voter_id" {
    null = false
    type = uuid
  }
  column "is_approved" {
    null = false
    type = boolean
  }
  column "comment" {
    null = true
    type = text
  }
  column "voted_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_votes_proof_id" {
    columns = [column.proof_id]
  }
  
  index "idx_votes_voter_id" {
    columns = [column.voter_id]
  }
  
  index "idx_votes_proof_voter" {
    unique = true
    columns = [column.proof_id, column.voter_id]
  }
  
  foreign_key "fk_votes_proof" {
    columns     = [column.proof_id]
    ref_columns = [table.proofs.column.id]
    on_delete   = CASCADE
  }
}

// Ledger Service Tables
table "accounts" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "account_type" {
    null = false
    type = varchar(20)
  }
  column "entity_id" {
    null = false
    type = uuid
  }
  column "currency" {
    null = false
    type = varchar(3)
    default = "NGN"
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_accounts_type" {
    columns = [column.account_type]
  }
  
  index "idx_accounts_entity" {
    columns = [column.entity_id]
  }
  
  index "idx_accounts_type_entity" {
    unique = true
    columns = [column.account_type, column.entity_id, column.currency]
  }
}

table "transactions" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "type" {
    null = false
    type = varchar(50)
  }
  column "description" {
    null = false
    type = varchar(500)
  }
  column "amount" {
    null = false
    type = bigint
  }
  column "currency" {
    null = false
    type = varchar(3)
  }
  column "metadata" {
    null = true
    type = jsonb
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_transactions_type" {
    columns = [column.type]
  }
  
  index "idx_transactions_created_at" {
    columns = [column.created_at]
  }
}

table "ledger_entries" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "account_id" {
    null = false
    type = uuid
  }
  column "transaction_id" {
    null = false
    type = uuid
  }
  column "entry_type" {
    null = false
    type = varchar(10)
  }
  column "amount" {
    null = false
    type = bigint
  }
  column "currency" {
    null = false
    type = varchar(3)
    default = "NGN"
  }
  column "description" {
    null = false
    type = varchar(500)
  }
  column "metadata" {
    null = true
    type = jsonb
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_ledger_entries_account_id" {
    columns = [column.account_id]
  }
  
  index "idx_ledger_entries_transaction_id" {
    columns = [column.transaction_id]
  }
  
  index "idx_ledger_entries_created_at" {
    columns = [column.created_at]
  }
  
  foreign_key "fk_ledger_entries_account" {
    columns     = [column.account_id]
    ref_columns = [table.accounts.column.id]
    on_delete   = RESTRICT
  }
  
  foreign_key "fk_ledger_entries_transaction" {
    columns     = [column.transaction_id]
    ref_columns = [table.transactions.column.id]
    on_delete   = RESTRICT
  }
}

table "balance_snapshots" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "account_id" {
    null = false
    type = uuid
  }
  column "balance" {
    null = false
    type = bigint
  }
  column "currency" {
    null = false
    type = varchar(3)
  }
  column "updated_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_balance_snapshots_account" {
    unique = true
    columns = [column.account_id]
  }
  
  foreign_key "fk_balance_snapshots_account" {
    columns     = [column.account_id]
    ref_columns = [table.accounts.column.id]
    on_delete   = CASCADE
  }
}