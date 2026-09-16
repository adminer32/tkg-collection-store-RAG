import { json } from "@sveltejs/kit";
import type { RequestHandler } from "@sveltejs/kit";

export const POST: RequestHandler = async ({ cookies }) => {
  cookies.delete("session_user", { path: "/" });
  cookies.delete("token", { path: "/" });
  return json({ message: "Logged out successfully" });
};
