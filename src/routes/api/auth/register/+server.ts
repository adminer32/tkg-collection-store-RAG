import { json } from "@sveltejs/kit";
import db from "$lib/server/db";
import { v4 as uuidv4 } from "uuid";
import bcrypt from "bcryptjs";

export async function POST({ request }) {
  const { username, password } = await request.json();

  if (!username || !password) {
    return json({ error: "Missing username or password" }, { status: 400 });
  }

  try {
    // Check if user exists
    const checkUser = db
      .prepare("SELECT id FROM users WHERE username = ?")
      .get(username);
    if (checkUser) {
      return json({ error: "Username already exists" }, { status: 409 });
    }

    // Hash password
    const passwordHash = await bcrypt.hash(password, 10);
    const id = uuidv4();

    const insertUser = db.prepare(
      "INSERT INTO users (id, username, password_hash) VALUES (?, ?, ?)",
    );
    const result = insertUser.run(id, username, passwordHash);

    return json(
      { message: "User created successfully", userId: id },
      { status: 201 },
    );
  } catch (error) {
    console.error("Registration error:", error);
    return json({ error: "Internal server error" }, { status: 500 });
  }
}
