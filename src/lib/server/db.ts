import Database from "better-sqlite3";
import { env } from "$env/dynamic/private";

const dbPath = env.DATABASE_PATH || "local.db";
const db = new Database(dbPath, { verbose: console.log });

// 初始化用户表
db.exec(`
  CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
  );

  CREATE TABLE IF NOT EXISTS inquiries (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    username TEXT,
    kind TEXT NOT NULL,
    name TEXT NOT NULL,
    contact TEXT NOT NULL,
    country TEXT,
    address TEXT,
    product_id TEXT,
    addon_id TEXT,
    quantity INTEGER NOT NULL DEFAULT 1,
    message TEXT,
    status TEXT NOT NULL DEFAULT 'received',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
  );

  CREATE TABLE IF NOT EXISTS rag_questions (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    question TEXT NOT NULL,
    answer TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
  );
`);

export default db;
