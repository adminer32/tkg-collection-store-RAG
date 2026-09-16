import { json } from '@sveltejs/kit';
import db from '$lib/server/db';
import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';
import { JWT_SECRET } from '$lib/server/auth';

type UserRow = {
  id: string;
  username: string;
  password_hash: string;
};

export async function POST({ request }) {
  const { username, password } = await request.json();

  if (!username || !password) {
    return json({ error: "Missing username or password" }, { status: 400 });
  }

  try {
    const user = db
      .prepare("SELECT * FROM users WHERE username = ?")
      .get(username) as UserRow | undefined;

    if (!user) {
      return json({ error: "Invalid username or password" }, { status: 401 });
    }

    const passwordMatch = await bcrypt.compare(password, user.password_hash);

    if (!passwordMatch) {
      return json({ error: "Invalid username or password" }, { status: 401 });
    }

    // Create JWT token
    const token = jwt.sign(
      { userId: user.id, username: user.username },
      JWT_SECRET,
      { expiresIn: "1h" },
    );

    return json(
      {
        message: "Login successful",
        token,
        user: { id: user.id, username: user.username },
      },
      { status: 200 },
    );
  } catch (error) {
    console.error("Login error:", error);
    return json({ error: "Internal server error" }, { status: 500 });
  }
}
