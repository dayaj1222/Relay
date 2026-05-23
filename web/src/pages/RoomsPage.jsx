import { useEffect, useState } from "react";
import { createRoom, getRooms } from "../api/rooms";

export function RoomsPage() {
  const [rooms, setRooms] = useState([]);
  const [roomName, setRoomName] = useState("");

  async function handleCreate() {
    if (!roomName.trim()) return;

    await createRoom(roomName);

    const updated = await getRooms();
    setRooms(updated);

    setRoomName("");
  }

  useEffect(() => {
    getRooms().then(setRooms);
  }, []);

  return (
    <>
      <div>Rooms</div>

      <ul>
        {rooms.map((room) => (
          <li key={room.id}>{room.name}</li>
        ))}
      </ul>

      <div>
        <span>Create New Room: </span>
        <input
          className="border-2"
          value={roomName}
          onChange={(e) => setRoomName(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              handleCreate();
            }
          }}
        />
      </div>
    </>
  );
}
