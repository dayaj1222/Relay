const API = import.meta.env.VITE_API_URL;

export async function getRooms() {
  const res = await fetch(`${API}/rooms`);
  return res.json();
}

export async function createRoom(name) {
  const res = await fetch(`${API}/rooms`, {
    method: "POST",
    headers: {
      "Content-type": "application/json",
    },
    body: JSON.stringify({ name }),
  });

  if (!res.ok) {
    throw new Error("failed to create new room.");
  }
  return res.json();
}
